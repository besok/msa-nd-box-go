package server

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"msa-nd-box-go/message"
	"msa-nd-box-go/storage"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

type AdminParam struct {
	Param
}

type Storages map[string]*storage.Storage

const (
	RegistryStorage       = "service_registry_storage"
	CircuitBreakerStorage = "circuit_breaker_storage"
	LoadBalanceStorage    = "load_balancer_storage"
	ReloadStorage         = "reload_storage"
)

type AdminServer struct {
	storages  Storages
	serverMux *http.ServeMux
	config    Config
}

func CreateAdminServer(serviceRegistryStorage string, listeners ...storage.Listener) *AdminServer {
	strs := createDefaultStorages(serviceRegistryStorage, listeners...)
	server := AdminServer{strs, http.NewServeMux(), defaultAdminConfig}
	server.serverMux.HandleFunc("/register", server.registerServiceHandler)
	server.serverMux.HandleFunc("/service/", server.getServiceList)
	server.serverMux.HandleFunc("/init/service/", server.initServices)
	server.serverMux.HandleFunc("/close/service/", server.closeServices)
	AddParamHandler(processLoadBalancer)
	AddParamHandler(ReloadFunc)
	AddParamHandler(CircuitBreakerParam)
	server.AddStorageListener(server.restartServer)
	server.AddStorageListener(server.removeUnusedValFromLBstr)
	server.AddStorageListener(server.removeCircuitBreaker)
	return &server
}

func (a *AdminServer) AddHandler(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	log.Printf("add new hanler: %s ", pattern)
	a.serverMux.HandleFunc(pattern, handler)
}
func (a *AdminServer) removeCircuitBreaker(event storage.Event, storageName storage.Name, key string, value storage.Line) {
	if event == storage.RemoveVal && storageName == RegistryStorage {
		str := a.Storage(CircuitBreakerStorage)
		addr := value.(storage.StringLine)
		var line storage.Line
		line = storage.CBLine{Address: addr.Value, Active: true}
		if err := str.RemoveValueIfExist(key, &line); err != nil {
			log.Printf("error while removing: %s  from cb str, error:%s ", line, err)
		}

	}
}

func (a *AdminServer) removeUnusedValFromLBstr(event storage.Event, storageName storage.Name, key string, value storage.Line) {
	if event == storage.RemoveKey && storageName == RegistryStorage {
		lbStr := a.Storage(LoadBalanceStorage)
		var l storage.Line
		l = storage.LBLine{Service: key}
		if line, ok := lbStr.GetValue("services", &l); ok {
			log.Printf("remove from loadbalancer Storage:%s, because there is not one running instance", line)
			if err := lbStr.RemoveValue("services", line); err != nil {
				log.Println("can not remove val from loadbalancer Storage, because ", err)
			}
		}
	}
}
func (a *AdminServer) restartServer(event storage.Event, storageName storage.Name, key string, value storage.Line) {
	if event == storage.RemoveVal && storageName == RegistryStorage {
		relStr := a.Storage(ReloadStorage)
		sl := value.(storage.StringLine)
		var l storage.Line
		l = storage.ReloadLine{Service: key, Address: sl.Value}
		if line, ok := relStr.GetValue("services", &l); ok {

			reloadLine := line.(storage.ReloadLine)
			if reloadLine.Count < reloadLine.Limit {
				err := a.StartNewCommand(reloadLine.Path)
				if err != nil {
					log.Printf("error to restart the server: %s, path:%s", err, reloadLine.Path)
				} else {
					log.Printf("server is restarted,path:%s", reloadLine.Path)
				}

				reloadLine.Count += 1
				_ = relStr.Put("services", reloadLine)
			} else {
				log.Printf("reloac limit: %d is reached for service:%s", reloadLine.Count, reloadLine.Service)
			}
		}
	}
}

func (a *AdminServer) initServices(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	serviceName := strings.TrimPrefix(request.URL.Path, "/init/service/")

	servers := findWorkingServers(serviceName, a).ToString()
	log.Printf("init request for service:%s, instances:%d", serviceName, len(servers))
	for _, addr := range servers {
		resp, err := http.Get(fmt.Sprintf("http://%s/init", addr))
		if err != nil {
			log.Printf("server with address[%s] is failed while initiation, error:%s", addr, err)
		}

		if resp.StatusCode == 500 {
			log.Printf("server with address[%s] is failed while initiation", addr)
		}
	}

}
func (a *AdminServer) closeServices(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	serviceName := strings.TrimPrefix(request.URL.Path, "/close/service/")
	if err := a.Storage(ReloadStorage).RemoveValue("services", storage.ReloadLine{Service: serviceName,}); err != nil {
		log.Printf(" error with clean the reload Storage up, service:%s, error:%s", serviceName, err)
	}
	servers := findWorkingServers(serviceName, a).ToString()

	for _, addr := range servers {
		resp, err := http.Get(fmt.Sprintf("http://%s/close", addr))
		if err != nil {
			log.Printf("server with address:%s is failed while termination, error:%s", addr, err)
		}
		if resp == nil {
			log.Printf("server with address:%s closed", addr)
		} else if resp.StatusCode == 500 {
			log.Printf("server with address:%s is failed while termination, error:%s", addr, err)
		}
	}

}

func findWorkingServers(serviceName string, a *AdminServer) storage.Lines {
	var lines storage.Lines
	serviceName = strings.TrimSuffix(serviceName, "/")
	hasCB := a.Storage(CircuitBreakerStorage).Contains(serviceName)
	if hasCB {
		lines = *a.filterLines(CircuitBreakerStorage, serviceName, activeCBServices)
	} else {
		lines = *a.filterLines(RegistryStorage, serviceName, noFilter)
	}
	return lines
}

func createDefaultStorages(path string, listeners ...storage.Listener) Storages {
	strs := make(Storages)
	strs[RegistryStorage] = createStorage(path, RegistryStorage, storage.CreateStringLines, listeners...)
	strs[CircuitBreakerStorage] = createStorage(path, CircuitBreakerStorage, storage.CreateCBLines, listeners...)
	strs[LoadBalanceStorage] = createStorage(path, LoadBalanceStorage, storage.CreateLBLines, listeners...)
	strs[ReloadStorage] = createStorage(path, ReloadStorage, storage.CreateReloadLines, listeners...)
	return strs
}

func createStorage(path string, name string, f func() storage.Lines, listeners ...storage.Listener) *storage.Storage {
	str, err := storage.CreateStorage(path, name, f, listeners)
	if err != nil {
		log.Printf("can not create storga: %s", name)
		panic(err)
	}
	return str
}

func (a *AdminServer) AddStorageListener(listener storage.Listener) {
	for _, v := range a.storages {
		v.AddListener(listener)
	}
}

func (a *AdminServer) Start(port string) {
	log.Println("start the admin server at the port:", port)
	a.snapshot()
	go a.fetchMetrics()
	_ = http.ListenAndServe(fmt.Sprintf(":%s", port), a.serverMux)
}

func (a *AdminServer) snapshot() {
	for _, v := range a.storages {
		log.Println(storage.Snapshot(v))
	}
}

func (a *AdminServer) fetchMetrics() {
	initDefaultMetrics()
	failedMetricMes := func(service string, addr string, err error) message.MetricsMessage {
		return message.CreateMetricsMessageWithMetric(
			message.Service{Service: service, Address: addr},
			message.Failed,
			"pulse", message.Metric{Value: "", Error: err})
	}
	for {
		str := a.Storage(RegistryStorage)
		keys := str.Keys()

		for _, k := range keys {
			lines, ok := str.Get(k)
			if !ok {
				log.Printf("key %s has been removed\n", k)
			}
			addresses := lines.ToString()
			for _, addr := range addresses {
				r, err := http.Get(fmt.Sprintf("http://%s/metrics", addr))
				var metricMessage message.MetricsMessage
				if err != nil {
					log.Printf("fetching metrics from the server:{%s,%s} has been finished with error: %s, \n", k, addr, err)
					metricMessage = failedMetricMes(k, addr, err)
				} else {
					decoder := json.NewDecoder(r.Body)
					err := decoder.Decode(&metricMessage)
					if err != nil {
						log.Printf("fetching metrics from the server:{%s,%s} has been finished with error: %s, \n", k, addr, err)
						metricMessage = failedMetricMes(k, addr, err)
					}
				}
				HandleMetrics(a, metricMessage)
			}
		}

		time.Sleep(time.Second * 5)
	}
}

func (a *AdminServer) AddStorage(s *storage.Storage) {
	a.storages[s.Name] = s
}

func (a *AdminServer) registerServiceHandler(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var sm message.ServerMessage
	err := decoder.Decode(&sm)
	if err != nil {
		log.Fatalf(" error %s while parsing json %s \n", err, request.Body)
		return
	}
	service := sm.Service
	log.Printf("got message from server: %s and address %s \n", service.Service, service.Address)
	err = a.Storage(RegistryStorage).Put(service.Service, storage.StringLine{Value: service.Address})
	if err != nil {
		log.Fatalf(" error:%s, saving at Storage", err)
	}

	for _, h := range paramHandlers.handlers {
		ps := sm.Params
		for k, v := range ps {
			if err := h(a, service, k, v); err != nil {
				log.Printf("error processing param for param: %s, value :%s", k, v)
			}
		}
	}
}

// todo	refactoring to chain or pipe
func (a *AdminServer) getServiceList(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	serviceName := strings.TrimPrefix(request.URL.Path, "/service/")

	var js []byte
	var e error

	if strings.Contains(serviceName, "/all") {
		var lines storage.Lines
		serviceName = strings.TrimSuffix(serviceName, "/all")
		hasCB := a.Storage(CircuitBreakerStorage).Contains(serviceName)
		if hasCB {
			lines = *a.filterLines(CircuitBreakerStorage, serviceName, activeCBServices)
		} else {
			lines = *a.filterLines(RegistryStorage, serviceName, noFilter)
		}
		js, e = json.Marshal(message.CreateGetServiceAllMessage(serviceName, lines))
	} else {
		serviceName = strings.TrimSuffix(serviceName, "/")
		addr := a.GetServiceByName(serviceName)

		js, e = json.Marshal(message.CreateGetServiceMessage(serviceName, addr))

	}
	if e != nil {
		log.Fatalf("can't convert to json, %s", e)
	}
	_, e = writer.Write(js)
	if e != nil {
		log.Fatalf("can't send, %s", e)
	}
}

func (a *AdminServer) GetServiceByName(serviceName string) string {
	var lines storage.Lines
	hasCB := a.Storage(CircuitBreakerStorage).Contains(serviceName)
	if hasCB {
		lines = *a.filterLines(CircuitBreakerStorage, serviceName, activeCBServices)
	} else {
		lines = *a.filterLines(RegistryStorage, serviceName, noFilter)
	}
	var addr = ""
	records := lines.ToString()
	rLn := len(records)
	lbStr := a.Storage(LoadBalanceStorage)
	var ln storage.Line = storage.LBLine{Service: serviceName}
	v, ok := lbStr.GetValue("services", &ln)
	if ok {
		lbLine := v.(storage.LBLine)
		idx := lbLine.Idx
		addr, idx = lbStrategyPicker(lbLine.Strategy, idx, records)
		lbLine.Idx = idx
		_ = lbStr.Put("services", lbLine)
	} else if rLn > 0 {
		addr = records[rand.Intn(rLn)]
	}
	return addr
}

func (a *AdminServer) filterLines(str string, key string, filter func(lines storage.Lines) *storage.Lines) *storage.Lines {
	lines, ok := a.Storage(str).Get(key)
	if !ok {
		return storage.CreateEmptyLines()
	}
	return filter(lines)
}

func noFilter(lines storage.Lines) *storage.Lines {
	return &lines
}

func activeCBServices(lines storage.Lines) *storage.Lines {
	cbLines := lines.(*storage.CBLines)

	tempLines := make([]storage.StringLine, 0)

	for _, v := range cbLines.Lines {
		if v.Active {
			tempLines = append(tempLines, storage.StringLine{Value: v.Address})
		}
	}

	var fLines storage.Lines
	fLines = &storage.StringLines{Lines: tempLines}
	return &fLines
}

func (a *AdminServer) Storage(name string) *storage.Storage {
	return a.storages[name]
}

func lbStrategyPicker(str storage.LBStrategy, idx int, records []string) (string, int) {
	rLn := len(records)
	addr := ""
	nextIdx := 0

	if rLn == 0 {
		return addr, nextIdx
	}

	switch str {
	case storage.Robin:
		if rLn > 0 {
			if idx < rLn {
				addr = records[idx]
			} else {
				addr = records[0]
			}
		}
		if idx == rLn-1 {
			nextIdx = 0
		} else {
			nextIdx = idx + 1
		}
	case storage.Random:
		addr = records[rand.Intn(rLn)]
	}
	return addr, nextIdx
}

func processLoadBalancer(a *AdminServer, service message.Service, p string, v string) error {
	if p == string(LOAD_BALANCER) {
		log.Printf("include load balancer:%s for %s", v, service)
		str := a.Storage(LoadBalanceStorage)
		var ln storage.Line
		ln = storage.LBLine{Service: service.Service, Strategy: storage.LBStrategy(v), Idx: 0}
		_, ok := str.GetValue("services", &ln)
		if !ok {
			if err := str.Put("services", ln); err != nil {
				log.Fatalf("error imposible to put to lb Storage:%s", err)
				return err
			}
		}
	}

	return nil
}

func (a *AdminServer) StartNewCommand(path string) error {
	log.Println("starting a new command for path: ", path)
	return exec.Command(path).Start()
}
