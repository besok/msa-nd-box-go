package server

import (
	"fmt"
	"msa-nd-box-go/message"
	"net/http"
	"testing"
	"time"
)

func TestGaugeStore_Take(t *testing.T) {
	store := createGaugeStore()

	store.AddGauge(Pulse)
	mes := store.Take(message.Service{
		Service: "Test", Address: "1",
	})
	fmt.Println(mes)
}

func TestCircuitBreaker(t *testing.T) {
	serv := CreateServer("test-server")
	serv.AddGauge(Pulse)
	serv.AddGauge(CircuitBreaker)
	serv.AddParam(CIRCUIT_BREAKER, "true")

	serv.AddHandlerWithCircuitBreaker("/long-op", func(writer http.ResponseWriter, request *http.Request) {
		time.Sleep(2 * time.Second)
	}, 1)

	go serv.Start()

	time.Sleep(1 * time.Second)

	metric := serv.takeMetrics().Metrics["cb"]

	if metric.Value != "false" {
		t.Fatalf(" should be false")
	}
	url := fmt.Sprintf("http://localhost%s/long-op", serv.service.Address)

	_, err := http.Get(url)
	if err != nil {
		t.Fatalf(" should here")
	}

	metric = serv.takeMetrics().Metrics["cb"]

	if metric.Value != "true" {
		t.Fatalf(" should be true")
	}

}