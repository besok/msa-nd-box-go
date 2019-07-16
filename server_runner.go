package main

import (
	"log"
	. "msa-nd-box-go/server"
	"net/http"
	"time"
)

func main() {
	serv := CreateServer("test-server")

	serv.AddGauge(Pulse)
	serv.AddGauge(CircuitBreaker)


	serv.AddParam(DISCOVERY, "localhost:9000")
	serv.AddParam(LOAD_BALANCER, "robin")
	serv.AddParam(CIRCUIT_BREAKER, "true")
	serv.AddParam(PORT,"10000")
	serv.AddParam(RESTART,"C:\\projects\\msa-nd-box-go\\bin\\server_runner_go.exe")


	serv.AddHandlerWithCircuitBreaker("/long-op", longFunc(), 1)
	AddInitOperation(Hello)
	serv.Start()
}

func longFunc() func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		time.Sleep(2 * time.Second)
	}
}
func Hello(s *Server) error {
	log.Printf("init operation for server %s",s.GetService())
	return Error("error init ope")
}
