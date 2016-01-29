package main

import (
	"os"

	"github.com/rs/xlog"
)

type config struct {
	loggerConfig xlog.Config
}

func (c *config) initDefaultConfig() {
	host, _ := os.Hostname()
	c.loggerConfig = xlog.Config{
		Level: xlog.LevelInfo,
		Fields: xlog.F{
			"role": "unoconv-api",
			"host": host,
		},
	}
	if os.Getenv("LOGFMT") == "json" {
		c.loggerConfig.Output = xlog.NewOutputChannel(xlog.NewJSONOutput(os.Stdout))
	} else {
		c.loggerConfig.Output = xlog.NewOutputChannel(xlog.NewConsoleOutput())
	}
}
