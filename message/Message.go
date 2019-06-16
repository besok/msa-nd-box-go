package message

import (
	"math/rand"
	"msa-nd-box-go/storage"
)

type Status string

const (
	Failed Status = "Failed"
	Ready         = "Ready"
	Run           = "Run"
	Done          = "Done"
)

type ServerMessage struct {
	Service Service
}

type Service struct {
	Address string
	Service string
}

type Message struct {
	From   Service
	Status Status
}
type GetServiceMessage struct {
	Message
	Service Service
}
type GetServiceAllMessage struct {
	Message
	Services []Service
}

type Metrics map[string]Metric
type Metric struct {
	Value string
	Error error
}
type MetricsMessage struct {
	Message
	Metrics
}

func CreateGetServiceAllMessage(service string, lines storage.Lines) GetServiceAllMessage {
	services := make([]Service, lines.Size())
	records := lines.ToString()

	for i, v := range records {
		services[i] = Service{v, service}
	}

	return GetServiceAllMessage{
		Message{
			Status: Ready, From: Service{":9000", "admin-service"},
		}, services,
	}
}
func CreateGetServiceMessage(service string, lines storage.Lines) GetServiceMessage {
	records := lines.ToString()
	size := lines.Size()
	if size == 0{
		return GetServiceMessage{
			Message{
				Status: Failed, From: Service{":9000", "admin-service"},
			}, Service{},
		}
	}
	r := rand.Intn(size)

	return GetServiceMessage{
		Message{
			Status: Ready, From: Service{":9000", "admin-service"},
		}, Service{Service: service, Address: records[r]},
	}
}

func CreateMetricsMessage(service Service, status Status, metrics Metrics) MetricsMessage {
	return MetricsMessage{
		Message{
			Status: status, From: service,
		},
		metrics,
	}
}
func CreateMetricsMessageWithMetric(service Service, status Status, metricName string, metric Metric) MetricsMessage {
	metrics := make(Metrics)
	metrics[metricName] = metric
	return MetricsMessage{
		Message{
			Status: status, From: service,
		},
		metrics,
	}
}
