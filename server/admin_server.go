package server

import (
	"encoding/json"
	"log"
	"msa-nd-box-go/message"
	"msa-nd-box-go/storage"
	"net/http"
	"strings"
)

type AdminServer struct {
	*storage.Storage
	serverMux *http.ServeMux
}

func CreateAdminServer(serviceRegistryStorage string, listeners ...storage.Listener) *AdminServer {
	str, err := storage.CreateStorage(serviceRegistryStorage, "service_registry_storage",
		storage.CreateStringLines, listeners)
	if err != nil {
		panic(err)
	}

	server := AdminServer{str, http.NewServeMux()}
	server.serverMux.HandleFunc("/register", server.registerServiceHandler)
	server.serverMux.HandleFunc("/service/", server.getServiceList)
	return &server
}

func (a *AdminServer) Start() {
	log.Println("start the admin server ")
	log.Println(storage.Snapshot(a.Storage))

	_ = http.ListenAndServe(":9000", a.serverMux)
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
	err = a.Put(sm.Service.Service, storage.StringLine{Value: sm.Service.Address})
	if err != nil {
		log.Fatalf(" error:%s, saving at storage", err)
	}
}
func (a *AdminServer) getServiceList(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	serviceName := strings.TrimPrefix(request.URL.Path, "/service/")

	if strings.Contains(serviceName, "/all") {
		serviceName = strings.TrimSuffix(serviceName, "/all")
		lines, ok := a.Get(serviceName)
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
		lines, ok := a.Get(serviceName)
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
