package integration

import (
	"log"
	. "msa-nd-box-go/server"
	"testing"
	"time"
)

func TestPulse(t *testing.T) {
	srv := CreateServer("test-service")
	srv.AddGauge(Pulse)
	go srv.Start()

	time.Sleep(time.Second * 5)
	metrics := srv.TakeMetrics().Metrics

	for k, v := range metrics {
		if k == "pulse" {
			if v.Error != nil {
				log.Fatalf("%s , %s \n", k, v)
			}
		}
	}
}

func TestPulseGroup(t *testing.T) {

	adminServer := CreateAdminServer("C:\\projects\\msa-nd-box-go\\file_storages_test")

	go adminServer.Start()
	time.Sleep(10 * time.Second)

	for i := 0; i < 10; i++ {
		go func() {
			srv := CreateServer("test_server", Pulse)
			srv.AddParam(DISCOVERY,"localhost:9000")
			srv.Start()
		}()
	}

	time.Sleep(10 * time.Second)

}
