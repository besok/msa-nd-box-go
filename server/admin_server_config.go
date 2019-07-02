package server

import "msa-nd-box-go/message"

type ParamHandler func(*AdminServer, message.Service, string, string) error
type ParamHandlers struct {
	handlers []ParamHandler
}

var paramHandlers = new(ParamHandlers)
var defaultAdminConfig = Config{make(Params)}

func AddParamHandler(p ParamHandler) {
	paramHandlers.handlers = append(paramHandlers.handlers, p)
}

