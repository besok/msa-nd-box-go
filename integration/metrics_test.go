package integration

import (
	"msa-nd-box-go/server"
	"testing"
	"time"
)

func TestPulseGroup(t *testing.T) {

	adminServer := server.CreateAdminServer("C:\\projects\\msa-nd-box-go\\file_storages_test")

	go adminServer.Start()
	time.Sleep(10 * time.Second)

	for i := 0; i < 10; i++ {
		go func() {
			srv := server.CreateServer("test_server", server.Pulse)
			srv.AddParam(server.DISCOVERY,"localhost:9000")
			srv.Start()
		}()
	}

	time.Sleep(10 * time.Second)

}
