package main

import (
	"github.com/mdesson/ew/echoserver"
	"log/slog"
)

func main() {
	s := echoserver.New(3000, slog.Default())
	s.Start()
}
