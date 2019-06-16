package integration

import (
	"encoding/json"
	"msa-nd-box-go/message"
	. "msa-nd-box-go/server"
	"net/http"
	"testing"
	"time"
)

func TestPulseGroup(t *testing.T) {

	adminServer := CreateAdminServer("C:\\projects\\msa-nd-box-go\\file_storages_test")

	go adminServer.Start()
	time.Sleep(10 * time.Second)

	for i := 0; i < 10; i++ {
		go func() {
			srv := CreateServer("test_server", Pulse)
			srv.AddParam(DISCOVERY, "localhost:9000")
			srv.Start()
		}()
	}

	time.Sleep(10 * time.Second)
}

func TestCB(t *testing.T) {
	adminServer := CreateAdminServer("C:\\projects\\msa-nd-box-go\\file_storages_test")

	go adminServer.Start()
	time.Sleep(1 * time.Second)

	go func() {
		srv := CreateServer("test_server", Pulse, CircuitBreaker)
		srv.AddParam(DISCOVERY, "localhost:9000")
		srv.AddParam(CIRCUIT_BREAKER, "true")
		srv.AddParam(PORT, "10000")
		srv.AddHandlerWithCircuitBreaker("/long-op", longFunc(), 1)
		srv.Start()
	}()

	time.Sleep(1 * time.Second)

	resp, err := http.Get("http://localhost:9000/service/test_server/all")
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("something goes wrong")
	}

	decoder := json.NewDecoder(resp.Body)
	var sm message.GetServiceAllMessage
	err = decoder.Decode(&sm)
	if err != nil {
		t.Fatalf("something goes wrong")
	}

	ln := len(sm.Services)
	if ln < 1 {
		t.Fatalf("something goes wrong")
	}

	_, _ = http.Get("http://localhost:10000/long-op")

	time.Sleep(10 * time.Second)
	resp, err = http.Get("http://localhost:9000/service/test_server/all")
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("something goes wrong")
	}

	decoder = json.NewDecoder(resp.Body)
	err = decoder.Decode(&sm)
	if err != nil {
		t.Fatalf("something goes wrong")
	}

	ln = len(sm.Services)
	if ln != 0 {
		t.Fatalf("something goes wrong")
	}


}

func longFunc() func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		time.Sleep(2 * time.Second)
	}
}
