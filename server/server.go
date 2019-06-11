package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"msa-nd-box-go/message"
	"net"
	"net/http"
)

func StartAndRegisterItself(service string) {

	port := fmt.Sprintf("localhost:%d", findNextPort())
	sm := message.ServerMessage{Service: message.Service{Address: port, Service: service}}
	buffer := new(bytes.Buffer)
	_ = json.NewEncoder(buffer).Encode(sm)
	b, err := http.Post("http://localhost:9000/register", "application/json; charset=utf-8", buffer)
	if err != nil {
		log.Printf("service %s can't start at %s, because error: %s ",service,port,err)
		panic(err)
	}
	if b.StatusCode != 200 {
		log.Printf("service %s can't start at %s, because status:%s, code:%d",service,port,b.Status,b.StatusCode)
		panic(err)
	}
	log.Printf("service %s is starting at %s \n", service,port)
	log.Println(http.ListenAndServe(sm.Service.Address, nil))
}

func findNextPort() int {
	port := 30000
	for {
		port++
		prt := fmt.Sprintf(":%d", port)
		_, err := net.Listen("tcp", prt)
		if err == nil {
			return port
		}
	}
}
