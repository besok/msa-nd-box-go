package main

import (
	. "msa-nd-box-go/server"
	"net/http"
	"time"
)

func main() {
	serv := CreateServer("test-server", CircuitBreaker)
	serv.AddGauge(Pulse)
	serv.AddParam(DISCOVERY, "localhost:9000")
	serv.AddParam(CIRCUIT_BREAKER, "true")
	serv.AddParam(PORT,"10000")
	serv.AddHandlerWithCircuitBreaker("/long-op", longFunc(), 1)
	serv.Start()
}

func longFunc() func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		time.Sleep(2 * time.Second)
	}
}
