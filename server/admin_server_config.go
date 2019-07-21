package server

import (
	"log"
	"msa-nd-box-go/message"
	"msa-nd-box-go/storage"
)

type ParamHandler func(*AdminServer, message.Service, string, string) error
type ParamHandlers struct {
	handlers []ParamHandler
}

var paramHandlers = new(ParamHandlers)
var defaultAdminConfig = Config{make(Params)}

func AddParamHandler(p ParamHandler) {
	paramHandlers.handlers = append(paramHandlers.handlers, p)
}

func CircuitBreakerParam(a *AdminServer, s message.Service, k string, v string) error {
	str := a.Storage(CircuitBreakerStorage)
	var ln storage.Line
	ln = storage.CBLine{Address: s.Address, Active: false}

	servName := s.Service
	if k == string(CIRCUIT_BREAKER) && v == "true" {
		log.Println(" create circuit breaker for service", servName, ", ", s.Address)
		return str.Put(servName, ln)
	} else {
		return str.RemoveValueIfExist(servName,&ln)
	}
}

func ReloadFunc(a *AdminServer, s message.Service, k string, v string) error {
	if k == string(RESTART) {
		rStr := a.Storage(ReloadStorage)
		var ln storage.Line
		ln = storage.ReloadLine{
			Service: s.Service, Address: s.Address,
			Path:  v,
			Limit: 10,
			Count: 0,
		}

		if exL, ok := rStr.GetValue("services", &ln); ok {
			ln = storage.ReloadLine{
				Service: s.Service, Address: s.Address,
				Path:  v,
				Limit: 10,
				Count: exL.(storage.ReloadLine).Count,
			}
		}

		err := rStr.Put("services", ln)
		if err != nil {
			log.Println("error while put new valu to reload Storage: ", err)
		}

	} else if k == string(RESTART_KEEP_PARAMS) {
	}
	return nil
}
