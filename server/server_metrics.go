package server

import (
	"log"
	"msa-nd-box-go/message"
	"reflect"
)

type Gauge interface {
	Take() (string, string, error)
}

type Pulse struct{}

func (Pulse) Take() (string, string, error) {
	return "pulse", "1", nil
}

type GaugeStore struct {
	Gauges []Gauge
}

func CreateGaugeStore(gauges ...Gauge) GaugeStore {
	if len(gauges) > 0 {
		return GaugeStore{Gauges: gauges}
	}
	return GaugeStore{make([]Gauge, 0)}
}

func (s *GaugeStore) Take(service message.Service) message.MetricsMessage {
	metricsMap := make(map[string]message.Metric)
	for _, g := range s.Gauges {
		k, v, e := g.Take()
		metricsMap[k] = message.Metric{Value: v, Error:e}
	}
	return message.CreateMetricsMessage(service, message.Ready, metricsMap)
}

func (s *GaugeStore) AddGauge(gauge Gauge) {
	s.Gauges = append(s.Gauges, gauge)
	log.Printf("add new gauge %s", reflect.TypeOf(gauge))
}
