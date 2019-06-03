package server

import (
	"encoding/json"
	"fmt"
	"msa-nd-box-go/message"
	"net/http"
)

func AdminServer()  {
	registerService()
	_= http.ListenAndServe(":9000", nil)
}

func registerService() {
	http.HandleFunc("/register", registerServiceHandler)
}

func registerServiceHandler(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var sm message.ServerMessage
	err := decoder.Decode(&sm)
	if err != nil {
		fmt.Printf(" error %s while parsing json %s \n", err, request.Body)
		return
	}
	fmt.Printf("got message from server: %s and address %s \n", sm.Service.Service, sm.Service.Address)

}
