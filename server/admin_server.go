package server

import (
	"encoding/json"
	"fmt"
	"log"
	"msa-nd-box-go/message"
	"msa-nd-box-go/storage"
	"net/http"
	"strings"
	"time"
)

type AdminParam struct {
	Param
}

type Storages map[string]*storage.Storage

const (
	REGISTRY_STORAGE = "service_registry_storage"
)

type AdminServer struct {
	storages  Storages
	serverMux *http.ServeMux
	config    Config
}

var defaultAdminConfig = Config{make(Params)}

func CreateAdminServer(serviceRegistryStorage string, listeners ...storage.Listener) *AdminServer {
	str, err := storage.CreateStorage(serviceRegistryStorage, REGISTRY_STORAGE,
		storage.CreateStringLines, listeners)

	if err != nil {
		panic(err)
	}

	strs := make(Storages)
	strs[REGISTRY_STORAGE] = str
	server := AdminServer{strs, http.NewServeMux(), defaultAdminConfig}
	server.serverMux.HandleFunc("/register", server.registerServiceHandler)
	server.serverMux.HandleFunc("/service/", server.getServiceList)
	return &server
}

func (a *AdminServer) Start() {
	log.Println("start the admin server ")
	log.Println(storage.Snapshot(a.storages[REGISTRY_STORAGE]))
	go a.fetchMetrics()
	_ = http.ListenAndServe(":9000", a.serverMux)
}

func (a *AdminServer) fetchMetrics() {
	NewMetricHandler(ServiceDiscoveryPulse)

	for {
		str := a.storage(REGISTRY_STORAGE)
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
					metricMessage = createFailedMetricMessage(k, addr, err)
				} else {
					decoder := json.NewDecoder(r.Body)
					err := decoder.Decode(&metricMessage)
					if err != nil {
						log.Printf("fetching metrics from the server:{%s,%s} has been finished with error: %s, \n", k, addr, err)
						metricMessage = createFailedMetricMessage(k, addr, err)
					}
				}
				HandleMetrics(a, metricMessage)
			}
		}

		time.Sleep(time.Second * 5)
	}
}

func createFailedMetricMessage(service string, addr string, err error) message.MetricsMessage {
	return message.CreateMetricsMessageWithMetric(
		message.Service{Service: service, Address: addr},
		message.Failed,
		"pulse", message.Metric{Value: "", Error: err})
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
	log.Printf("got message from server: %s and address %s \n", sm.Service.Service, sm.Service.Address)
	err = a.storage(REGISTRY_STORAGE).Put(sm.Service.Service, storage.StringLine{Value: sm.Service.Address})
	if err != nil {
		log.Fatalf(" error:%s, saving at storage", err)
	}
}
func (a *AdminServer) getServiceList(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	serviceName := strings.TrimPrefix(request.URL.Path, "/service/")

	if strings.Contains(serviceName, "/all") {
		serviceName = strings.TrimSuffix(serviceName, "/all")
		lines, ok := a.storage(REGISTRY_STORAGE).Get(serviceName)
		if !ok {
			writer.WriteHeader(404)
			return
		}
		js, e := json.Marshal(message.CreateGetServiceAllMessage(serviceName, lines))
		if e != nil {
			log.Fatalf("can't convert to json, %s", e)
		}
		_, e = writer.Write(js)
		if e != nil {
			log.Fatalf("can't send, %s", e)
		}

	} else {
		lines, ok := a.storage(REGISTRY_STORAGE).Get(serviceName)
		if !ok {
			writer.WriteHeader(404)
			return
		}
		js, e := json.Marshal(message.CreateGetServiceMessage(serviceName, lines))
		if e != nil {
			log.Fatalf("can't convert to json, %s", e)
		}
		_, e = writer.Write(js)
		if e != nil {
			log.Fatalf("can't send, %s", e)
		}
	}

}

func (a *AdminServer) storage(name string) *storage.Storage {
	return a.storages[name]
}
