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

func ServiceDiscoveryPulse(a *AdminServer, message message.MetricsMessage) error {
	metrics := message.Metrics
	metric := metrics["pulse"]
	if metric.Error != nil {
		service := message.From
		log.Printf("service:%s does not have pulse. It needs to be removed", service)
		str := a.storages[REGISTRY_STORAGE]
		return str.RemoveValue(service.Service,
			storage.StringLine{Value: service.Address})
	}
	return nil
}
