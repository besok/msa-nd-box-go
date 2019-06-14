package integration

import (
	"encoding/json"
	"msa-nd-box-go/message"
	. "msa-nd-box-go/server"
	"msa-nd-box-go/storage"
	"net/http"
	"testing"
	"time"
)

func TestServiceRegistry(t *testing.T) {

	second := time.Second
	ch := make(chan int, 10)

	go func(chan int) {
		storagePath := "C:\\projects\\msa-nd-box-go\\file_storage"
		CreateAdminServer(storagePath, listenerPutOperation(ch)).Start()
	}(ch)

	time.Sleep(1 * second)

	service := "test_service"
	for i := 0; i < 10; i++ {
		go func() {
			server := CreateServer(service)
			server.AddGauge(Pulse)
			server.AddParam(DISCOVERY,"localhost:9000")
			server.Start()
		}()
	}
	time.Sleep(10 * second)
	el := 0
	for i := 0; i < 10; i++ {
		el += <-ch
	}

	if el < 10 {
		t.Fatal(" el must be = 10")
	}
}
func TestServiceDiscovery(t *testing.T) {
	second := time.Second

	go CreateAdminServer("C:\\projects\\msa-nd-box-go\\file_storage", ).Start()


	time.Sleep(1 * second)

	service := "test_service"
	for i := 0; i < 10; i++ {
		go func() {
			server := CreateServer(service)
			server.AddGauge(Pulse)
			server.AddParam(DISCOVERY,"localhost:9000")
			server.Start()
		}()
	}
	time.Sleep(10 * second)
	resp, err := http.Get("http://localhost:9000/service/test_service/all")
	if err != nil || resp.StatusCode != 200{
		t.Fatalf("something goes wrong")
	}

	decoder := json.NewDecoder(resp.Body)
	var sm message.GetServiceAllMessage
	err = decoder.Decode(&sm)
	if err != nil {
		t.Fatalf("something goes wrong")
	}

	if len(sm.Services) < 10 {
		t.Fatalf("something goes wrong")
	}


	resp, err = http.Get("http://localhost:9000/service/service_test1/all")
	if err != nil {
		t.Fatalf("something goes wrong")
	}
	if resp.StatusCode != 404{
		t.Fatalf("something goes wrong")
	}


}

func listenerPutOperation(ch chan int) func(event storage.StorageEvent, storageName storage.StorageName, key string, value storage.Line) {
	return func(event storage.StorageEvent, storageName storage.StorageName, key string, value storage.Line) {
		if event == storage.Put {
			ch <- 1
		}
	}
}
