package integration

import (
	"fmt"
	"msa-nd-box-go/server"
	"msa-nd-box-go/storage"
	"testing"
	"time"
)

func TestServiceRegistry(t *testing.T) {
	storagePath := "C:\\projects\\msa-nd-box-go\\file_storages"

	adminServer := server.CreateAdminServer(storagePath)
	second := time.Second

	go adminServer.Start()
	time.Sleep(1 * second)
	serviceName := "service_test"
	go server.StartAndRegisterItself(serviceName)
	go server.StartAndRegisterItself(serviceName)
	go server.StartAndRegisterItself(serviceName)
	time.Sleep(3 * second)

	str := &adminServer.Storage
	snapshot := storage.Snapshot(str)
	fmt.Println(snapshot)
	if lines, ok := str.Get(serviceName); ok {
		if lines.Size() < 3 {
			t.Fatalf(" should be 3 services")
		}
	}

}
