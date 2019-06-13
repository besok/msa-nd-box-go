package integration

import (
	"log"
	"msa-nd-box-go/server"
	"testing"
	"time"
)

func TestPulse(t *testing.T) {
	srv := server.CreateServer("test-service")
	srv.AddGauge(server.Pulse)
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
