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
	metric, ok := metrics["pulse"]

	if !ok {
		service := message.From
		log.Printf("service:%s does not have pulse or the service does not have the pulse metric. It needs to be removed", service)
		str := a.storage(REGISTRY_STORAGE)
		return str.RemoveValue(service.Service, storage.StringLine{Value: service.Address})
	}

	if metric.Error != nil {
		service := message.From
		log.Printf("service:%s does not have pulse. It needs to be removed", service)
		str := a.storage(REGISTRY_STORAGE)
		return str.RemoveValue(service.Service, storage.StringLine{Value: service.Address})
	}
	return nil
}
func CBHandler(a *AdminServer, message message.MetricsMessage) error {
	metrics := message.Metrics
	metric, ok := metrics["cb"]
	service := message.From
	str := a.storage(CIRCUIT_BREAKER_STORAGE)

	if !ok {
		return str.RemoveValue(service.Service, storage.CBLine{Address: service.Address})
	}

	active := true
	if metric.Error != nil || metric.Value != "true" {
		active = false
	}
	var ln storage.Line
	ln = storage.CBLine{Address: service.Address, Active: false}

	line, ok := str.GetValue(service.Service, &ln)
	if ok {
		cbLine := line.(storage.CBLine)
		if cbLine.Active == active {
			return nil
		}
	}

	return str.Put(service.Service, storage.CBLine{Address: service.Address, Active: active})
}
