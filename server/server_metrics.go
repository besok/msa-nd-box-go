package server

import (
	"log"
	"msa-nd-box-go/message"
	"reflect"
)

type Gauge interface {
	Take() (string, string, error)
}

type Pulse struct {}

func (Pulse) Take() (string, string, error) {
	return "pulse", "1", nil
}

type GaugeStore struct {
	Gauges []Gauge
	message.Service
}

func CreateGaugeStore(service message.Service, gauges ...Gauge) GaugeStore {
	if len(gauges) > 0 {
		return GaugeStore{Gauges: gauges, Service: service}
	}
	return GaugeStore{make([]Gauge, 0), service}
}

func (s *GaugeStore) Take() message.MetricsMessage {
	metricsMap := make(map[string]string)
	for _, g := range s.Gauges {
		k, v, e := g.Take()
		if e != nil {
			log.Fatalf("gauge has been finished with err : %s", e)
		}
		metricsMap[k] = v
	}
	return message.CreateMetricsMessage(s.Service.Service, s.Service.Address, message.Ready, metricsMap)
}

func (s *GaugeStore) AddGauge(gauge Gauge) {
	s.Gauges = append(s.Gauges, gauge)
	log.Printf("add new gauge %s", reflect.TypeOf(gauge))
}
