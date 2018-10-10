package main

import (
	stdlog "log"
	"octopus/config"
	"octopus/services/stream"
	"os"

	"github.com/go-kit/kit/log"
)

func main() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))

	var con config.Connection
	if err := con.GetCon("/Users/aputra/go/src/octopus/config/default.yaml"); err != nil {
		stdlog.Printf("%+v", err)
		os.Exit(1)
	}

	s := stream.LogMiddleware{
		Logger: logger,
		Next:   stream.NewStreamerService(),
	}
	if err := s.Process(con); err != nil {
		stdlog.Printf("%+v", err)
		os.Exit(1)
	}
}
