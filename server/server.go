package server

import (
	"encoding/json"
	"fmt"
	"log"
	"msa-nd-box-go/message"
	"net"
	"net/http"
	"time"
)

type Error string

func (e Error) Error() string {
	return string(e)
}

type CbValue struct {
	expected int
	actual   int
}
type CBProcessor struct {
	circuitBreakers map[string]CbValue
}

var defCBProcessor = CBProcessor{make(map[string]CbValue)}

type Server struct {
	gaugeStore GaugeStore
	service    message.Service
	mux        *http.ServeMux
	config     Config
	listener   *net.Listener
}

func(s *Server) Close() error {
	return (*s.listener).Close()
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
	defaultInitHandler().handle(s)

	addr := s.service.Address
	log.Printf("service %s is about to start \n", s.service)

	s.mux.HandleFunc("/metrics", s.processMetrics)
	s.mux.HandleFunc("/init", s.init)
	s.mux.HandleFunc("/close", s.close)
	srv := http.Server{Addr: addr, Handler: s.mux}
	li := *s.listener
	//to put it last
	AddCloseOperation(Close)
	log.Println(srv.Serve(li.(*net.TCPListener)))
}

//AddGauge add gauge to server
func (s *Server) AddGauge(gauge Gauge) *Server {
	s.gaugeStore.AddGauge(gauge)
	return s
}
func (s *Server) AddParam(param Param, value string) *Server {
	s.config.addParam(param, value)
	return s
}

func NewInitOperator(f func(server *Server) error) {
	defHandler.operators = append(defHandler.operators, f)
}

func (s *Server) AddHandler(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	log.Printf("server:%s, add new hanler: %s ", s.service, pattern)
	s.mux.HandleFunc(pattern, handler)
}
func (s *Server) AddHandlerWithCircuitBreaker(pattern string, handler func(http.ResponseWriter, *http.Request), expDuration int) {
	log.Printf("server:%s, add new hanler : %s with circuit breaker: %d ", s.service, pattern, expDuration)
	if expDuration < 0 {
		log.Println("should be more 0")
		panic(Error(""))
	}
	defCBProcessor.circuitBreakers[pattern] = CbValue{actual: 0, expected: expDuration}
	s.mux.HandleFunc(pattern, wrapCbHandler(pattern, handler))
}

func wrapCbHandler(pattern string, h func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		nowTime := time.Now()
		h(writer, request)
		since := time.Since(nowTime) / time.Second
		v := defCBProcessor.circuitBreakers[pattern]
		v.actual = int(since)
		defCBProcessor.circuitBreakers[pattern] = v
		log.Printf("handler duration: %s", since)
	}
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

func (s *Server) takeMetrics() message.MetricsMessage {
	return s.gaugeStore.Take(s.service)
}

func (s *Server) processMetrics(writer http.ResponseWriter, request *http.Request) {
	msg := s.takeMetrics()
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

type Operation func(s *Server) error
type OperationHandler struct {
	initOperations  []Operation
	closeOperations []Operation
}

var operationHandler = new(OperationHandler)

func AddInitOperation(op Operation) {
	operationHandler.initOperations = append(operationHandler.initOperations, op)
}
func AddCloseOperation(op Operation) {
	operationHandler.closeOperations = append(operationHandler.closeOperations, op)
}

func (s *Server) init(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := initOp(s)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
}
func (s *Server) close(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := closeOp(s)
	if err != nil {
		w.WriteHeader(500)
		log.Fatalf("can not marshal json %s", err)
		return
	}
	w.WriteHeader(200)
}
func initOp(s *Server) error {
	for _, op := range operationHandler.initOperations {
		e := op(s)
		if e != nil {
			log.Printf("error while invoke init operator: %s", e)
			return e
		}
	}
	return nil
}
func closeOp(s *Server) error {
	for _, op := range operationHandler.closeOperations {
		e := op(s)
		if e != nil {
			log.Printf("error while invoke close operator: %s", e)
			return e
		}
	}

	return nil
}
func (s *Server) GetService() message.Service {
	return s.service
}
func Close(s *Server) error {
	return s.Close()
}