package main

import (
	"github.com/mdesson/ew/echoserver"
	"github.com/mdesson/ew/frontendserver"
	"log/slog"
)

func main() {
	reqChan := make(chan echoserver.RequestDetails)
	echo := echoserver.New(3000, slog.Default(), reqChan)
	frontend := frontendserver.New(3001, slog.Default(), reqChan)
	go echo.Start()
	frontend.Start()
}
