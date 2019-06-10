package integration

import (
	"msa-nd-box-go/server"
	"testing"
	"time"
)

func TestServiceRegistry(t *testing.T) {
	storagePath := "C:\\projects\\msa-nd-box-go\\file_storages"

	go server.CreateAdminServer(storagePath)
	second := time.Second
	time.Sleep(1 * second)
	go server.StartAndRegisterItself("service_test")
	go server.StartAndRegisterItself("service_test")
	go server.StartAndRegisterItself("service_test")
	time.Sleep(10 * second)


}
