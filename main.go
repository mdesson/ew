package main

import (
	"github.com/mdesson/ew/server"
	"log/slog"
)

func main() {
	s := server.New(3000, slog.Default())
	s.Start()
}
