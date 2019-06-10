package integration

import (
	"msa-nd-box-go/server"
	"testing"
)

func TestServiceRegistry(t *testing.T){
	storagePath := "C:\\projects\\msa-nd-box-go\\file_storages"
	//adminStart := false


	go func() {
	server.CreateAdminServer(":",storagePath)

	}()
	go server.StartAndRegisterItself("service_test")
	go server.StartAndRegisterItself("service_test")
	go server.StartAndRegisterItself("service_test")


}
