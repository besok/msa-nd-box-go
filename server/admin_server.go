package server

import (
	"encoding/json"
	"log"
	"msa-nd-box-go/message"
	"msa-nd-box-go/storage"
	"net/http"
)

type AdminServer struct {
	storage.Storage
}

func CreateAdminServer(port string, serviceRegistryStorage string, listeners ...storage.Listener) {
	str, err := storage.CreateStorage(serviceRegistryStorage, "service_registry_storage",
		storage.CreateStringLines, listeners)
	if err != nil {
		panic(err)
	}

	server := AdminServer{str}
	http.HandleFunc("/register", server.registerServiceHandler)
	_ = http.ListenAndServe(port, nil)
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
