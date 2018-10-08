package main

import (
	"log"
	"octopus/config"
	"octopus/services/stream"
	"os"
)

func main() {
	var con config.Connection
	if err := con.GetCon("/home/octopus/config.yaml"); err != nil {
		log.Printf("%+v", err)
		os.Exit(1)
	}

	s := stream.NewStreamerService()
	if err := s.Process(con); err != nil {
		log.Printf("%+v", err)
		os.Exit(1)
	}
}
