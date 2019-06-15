package main

import (
	"msa-nd-box-go/server"
	"net/http"
	"time"
)

func main() {
	serv := server.CreateServer("test-server")
	serv.AddGauge(server.Pulse)
	serv.AddGauge(server.CircuitBreaker)
	serv.AddParam(server.DISCOVERY, "localhost:9000")
	serv.AddParam(server.CIRCUIT_BREAKER, "true")

	serv.AddHandlerWithCircuitBreaker("/long-op", longFunc(), 1)
	serv.Start()
}

func longFunc() func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		time.Sleep(2 * time.Second)
	}
}
