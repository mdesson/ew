package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/mdesson/ew/Config"
	"github.com/mdesson/ew/echoserver"
	"github.com/mdesson/ew/frontendserver"
	"github.com/mdesson/ew/util"
)

func main() {
	// get config
	cfg, err := Config.New()
	if err != nil {
		log.Fatal(err)
	}

	// set up logger
	level, err := util.LogLevelInt(cfg.LogLevel)
	if err != nil {
		log.Fatal(err)
	}

	l := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
	}))

	reqChan := make(chan echoserver.RequestDetails)
	echo := echoserver.New(cfg.EchoPort, l, reqChan)
	frontend := frontendserver.New(cfg.FrontendPort, l, reqChan)
	go echo.Start()
	frontend.Start()
}
