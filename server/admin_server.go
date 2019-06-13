package server

import (
	"encoding/json"
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

	NewMetricHandler(ServiceDiscoveryPulse)

	storages := make(Storages)
	storages[REGISTRY_STORAGE] = str
	server := AdminServer{storages, http.NewServeMux(), defaultAdminConfig}
	server.serverMux.HandleFunc("/register", server.registerServiceHandler)
	server.serverMux.HandleFunc("/service/", server.getServiceList)
	return &server
}

func (a *AdminServer) Start() {
	log.Println("start the admin server ")
	log.Println(storage.Snapshot(a.storages[REGISTRY_STORAGE]))

	_ = http.ListenAndServe(":9000", a.serverMux)
}

func (a *AdminServer) fetchMetrics() {
	for {

		str := a.storage(REGISTRY_STORAGE)
		keys := str.Keys()
		for _, k := range keys {
			lines, ok := str.Get(k)
			if !ok {

			}
			addresses := lines.ToString()
			for _,adr := range addresses{

			}
		}

		time.Sleep(time.Second * 10)
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
