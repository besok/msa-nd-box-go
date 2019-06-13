package server

import (
	"encoding/json"
	"fmt"
	"log"
	"msa-nd-box-go/message"
	"net"
	"net/http"
)

type Server struct {
	gaugeStore GaugeStore
	service    message.Service
	mux        *http.ServeMux
	config     Config
	listener   *net.Listener
}

func (s *Server) TakeMetrics() message.MetricsMessage {
	return s.gaugeStore.Take(s.service)
}

func CreateServer(serviceName string, gauges ...Gauge) *Server {
	port, li := findNextPort()
	address := fmt.Sprintf(":%d", port)
	store := createGaugeStore(gauges[:]...)
	return &Server{
		mux:        http.NewServeMux(),
		service:    message.Service{Service: serviceName, Address: address},
		gaugeStore: store,
		config:     defaultConfig,
		listener:   &li,
	}
}

func (s *Server) Start() {
	defaultInitHandler().Handle(s)
	addr := s.service.Address
	log.Printf("service %s is about to start \n", s.service)
	s.mux.HandleFunc("/metrics", s.processMetrics)
	s.mux.HandleFunc("/h", h)

	srv := &http.Server{Addr: addr, Handler: s.mux}
	li := *s.listener
	_ = srv.Serve(li.(*net.TCPListener))
}

func (s *Server) AddGauge(gauge Gauge) *Server {
	s.gaugeStore.AddGauge(gauge)
	return s
}
func (s *Server) AddParam(param Param, value string) *Server {
	s.config.AddParam(param, value)
	return s
}

func findNextPort() (int, net.Listener) {
	port := 30000
	for {
		port++
		prt := fmt.Sprintf(":%d", port)
		c, err := net.Listen("tcp", prt)
		if err == nil {
			return port, c
		}
	}
}
func (s *Server) processMetrics(writer http.ResponseWriter, request *http.Request) {
	msg := s.TakeMetrics()
	writer.Header().Set("Content-Type", "application/json")
	js, e := json.Marshal(msg)
	if e != nil {
		writer.WriteHeader(500)
		log.Fatalf("can not marshal json %s", msg)
		return
	}
	writer.WriteHeader(200)
	_, _ = writer.Write(js)
	return
}

func h(writer http.ResponseWriter, r *http.Request) {
	_, _ = writer.Write([]byte("hello"))
}
