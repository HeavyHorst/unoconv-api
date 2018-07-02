package main

import (
	"os"
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
	uno *unoconv
)

func init() {
	// read config data
	cfg.initDefaultConfig()

	uno = initUnoconv()

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

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	router := xmux.New()
	addr := getEnv("LISTEN_ADDR", ":3000")

	router.GET("/unoconv/health", xhandler.HandlerFuncC(healthHandler))
	router.POST("/unoconv/:filetype", xhandler.HandlerFuncC(unoconvHandler))
	log.Fatal(http.ListenAndServe(addr, mw.Handler(router)))
}
