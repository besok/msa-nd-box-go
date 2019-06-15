package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"msa-nd-box-go/message"
	"net/http"
	"reflect"
	"strconv"
)


const (
	DISCOVERY Param = "discovery"
	CIRCUIT_BREAKER Param ="circuit_breaker"
)

type Param string
type Params map[Param]string
type Config struct {
	Params
}

var defaultConfig = Config{make(Params)}

func (c *Config) addParam(param Param, value string) {
	log.Printf("param:%s, value:%s", param, value)
	c.Params[Param(param)] = value
}

func (c *Config) Bool(param Param) bool {
	el, ok := c.Params[param]
	if !ok {
		return false
	}
	if el == "true" {
		return true
	}
	return false
}
func (c *Config) Int(param Param) (int, bool) {
	el, ok := c.Params[param]
	if !ok {
		return 0, ok
	}

	intV, err := strconv.Atoi(el)
	if err != nil {
		return 0, false
	}
	return intV, true
}

func (c *Config) String(param Param) (string, bool) {
	r, ok := c.Params[param]
	return r, ok
}

type InitOperator func(server *Server) error

type InitHandler struct {
	operators []InitOperator
}

var defHandler = InitHandler{operators: make([]InitOperator, 0)}

func defaultInitHandler() *InitHandler {
	defHandler.operators = append(defHandler.operators, discovery)
	return &defHandler
}

func (h *InitHandler) handle(server *Server) {
	for _, op := range h.operators {
		log.Printf("init operator %s is starting ", reflect.TypeOf(op))
		err := op(server)
		if err != nil {
			log.Fatalf("init operator failed, error:%s ", err)
		}
	}
}

func discovery(server *Server) error {
	s, ok := server.config.String(DISCOVERY)
	if !ok {
		log.Printf("discovery does not need")
		return nil
	}
	sm := message.ServerMessage{Service: server.service}
	buffer := new(bytes.Buffer)
	_ = json.NewEncoder(buffer).Encode(sm)
	b, err := http.Post(fmt.Sprintf("http://%s/register", s), "application/json; charset=utf-8", buffer)
	if err != nil {
		log.Fatalf("service %s can't start, because the admin server does not response,\n error: %s ", server.service, err)
	} else {
		if b.StatusCode != 200 {
			log.Fatalf("service %s can't start , because the admin server does not response, status:%s, code:%d", server.service, b.Status, b.StatusCode)
		}
	}

	return err
}
