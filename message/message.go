package message

import (
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
	Params  map[string]string
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

type Chapter struct {
	Service  string
	Chapter  string
	Rollback string
	Input    string
}
type ChapterResult struct {
	Chapter Chapter
	State   ChapterState
}

func NewChapterResult(chapter Chapter, state ChapterState) *ChapterResult {
	return &ChapterResult{Chapter: chapter, State: state}
}

type ChapterState string

const (
	Abort    ChapterState = "failed"
	Start    ChapterState = "start"
	Rollback ChapterState = "rollback"
	Success  ChapterState = "success"
	Finish   ChapterState = "finish"
)

func NewChapter(service string, chapter string, rollback string, input string) *Chapter {
	return &Chapter{Service: service, Chapter: chapter, Rollback: rollback, Input: input}
}

type SagaMessage struct {
	Chapters []Chapter
}

func NewSagaMessage(chapters ...Chapter) *SagaMessage {
	return &SagaMessage{Chapters: chapters}
}

func CreateGetServiceAllMessage(service string, lines storage.Lines) GetServiceAllMessage {
	services := make([]Service, lines.Size())
	records := lines.ToString()

	for i, v := range records {
		services[i] = Service{v, service}
	}

	return GetServiceAllMessage{
		Message{
			Status: Ready, From: Service{":9000", "admin-Service"},
		}, services,
	}
}
func CreateGetServiceMessage(service string, addr string) GetServiceMessage {
	return GetServiceMessage{
		Message{
			Status: Ready, From: Service{":9000", "admin-Service"},
		}, Service{Service: service, Address: addr},
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
