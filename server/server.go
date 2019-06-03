package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"msa-nd-box-go/message"
	"net"
	"net/http"
	"os"
)

func StartAndRegisterItself(service string) {
	fmt.Printf("service %s is starting at ", service)

	port := fmt.Sprintf("localhost:%d", findNextPort())
	sm := message.ServerMessage{Service: message.Service{Address: port, Service: service}}
	buffer := new(bytes.Buffer)
	_ = json.NewEncoder(buffer).Encode(sm)
	_, err := http.Post("http://localhost:9000/register", "application/json; charset=utf-8", buffer)
	if err != nil {
		os.Exit(1)
	}
	fmt.Println(http.ListenAndServe(sm.Service.Address, nil))
}

func findNextPort() int {
	port := 30000
	for {
		port++
		prt := fmt.Sprintf(":%d", port)
		_, err := net.Listen("tcp", prt)
		if err == nil {
			fmt.Printf(" %s \n", prt)
			return port
		}
	}
}
