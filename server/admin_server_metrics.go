package server

import (
	"log"
	"msa-nd-box-go/message"
	"msa-nd-box-go/storage"
)

type MetricHandler func(a *AdminServer, message message.MetricsMessage) error

type MetricsProcessor struct {
	handlers []MetricHandler
}

var defaultMetricProcessor = MetricsProcessor{make([]MetricHandler, 0)}

func initDefaultMetrics() {
	NewMetricHandler(PulseMetricHandler)
	NewMetricHandler(CBMetricHandler)
}

func NewMetricHandler(handler MetricHandler) {
	defaultMetricProcessor.handlers = append(defaultMetricProcessor.handlers, handler)
}
func HandleMetrics(a *AdminServer, metricsMessage message.MetricsMessage) {
	for _, h := range defaultMetricProcessor.handlers {
		err := h(a, metricsMessage)
		if err != nil {
			log.Fatalf("error in metric handler , message: %s, error: %s", metricsMessage, err)
		}
	}
}

func PulseMetricHandler(a *AdminServer, message message.MetricsMessage) error {
	metrics := message.Metrics
	metric, ok := metrics["pulse"]

	if !ok {
		service := message.From
		log.Printf("service:%s does not have pulse or the service does not have the pulse metric. It needs to be removed", service)
		str := a.Storage(RegistryStorage)
		return str.RemoveValue(service.Service, storage.StringLine{Value: service.Address})
	}

	if metric.Error != nil {
		service := message.From
		log.Printf("service:%s does not have pulse. It needs to be removed", service)
		str := a.Storage(RegistryStorage)
		return str.RemoveValue(service.Service, storage.StringLine{Value: service.Address})
	}
	return nil
}
func CBMetricHandler(a *AdminServer, message message.MetricsMessage) error {
	if metric, ok := message.Metrics["cb"]; ok {
		service := message.From
		str := a.Storage(CircuitBreakerStorage)
		active := true
		if metric.Error != nil || metric.Value != "true" {
			active = false
		}
		return str.Put(service.Service, storage.CBLine{Address: service.Address, Active: active})
	}

	return nil
}
