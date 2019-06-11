package integration

import (
	"msa-nd-box-go/server"
	"msa-nd-box-go/storage"
	"testing"
	"time"
)

func TestServiceRegistry(t *testing.T) {

	second := time.Second
	ch := make(chan int, 10)

	go func(chan int) {
		storagePath := "C:\\projects\\msa-nd-box-go\\file_storages"
		server.CreateAdminServer(storagePath, putListener(ch)).Start()
	}(ch)

	time.Sleep(1 * second)

	service := "service_test"
	for i := 0; i < 10; i++ {
		go server.StartAndRegisterItself(service)
	}
	time.Sleep(5 * second)
	el := 0
	for i := 0; i < 10; i++ {
		el += <-ch
	}

	if el < 10 {
		t.Fatal(" el must be = 10")
	}
}

func putListener(ch chan int) func(event storage.StorageEvent, storageName storage.StorageName, key string, value storage.Line) {
	return func(event storage.StorageEvent, storageName storage.StorageName, key string, value storage.Line) {
		if event == storage.Put {
			ch <- 1
		}
	}
}
