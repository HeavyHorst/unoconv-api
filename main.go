package main

import (
	"log"
	"net/http"

	"github.com/rs/xaccess"
	"github.com/rs/xhandler"
	"github.com/rs/xlog"
	"github.com/rs/xmux"
)

var (
	cfg config
	mw  middleware
)

func init() {
	//plug the xlog handler's input to Go's default logger
	log.SetFlags(0)
	xlogger := xlog.New(cfg.loggerConfig)
	log.SetOutput(xlogger)

	//register some middleware handlers
	mw.initCommonHandlers(
		xlog.NewHandler(cfg.loggerConfig),
		xaccess.NewHandler(),
	)
}

func main() {
	router := xmux.New()

	router.POST("/unoconv/:filetype", xhandler.HandlerFuncC(unoconvHandler))
	log.Fatal(http.ListenAndServe(":3000", mw.Handler(router)))
}
