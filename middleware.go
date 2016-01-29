package main

import "github.com/rs/xhandler"

type middleware struct {
	xhandler.Chain
}

func (m *middleware) initCommonHandlers(handler ...func(xhandler.HandlerC) xhandler.HandlerC) {
	for _, v := range handler {
		m.UseC(v)
	}
}
