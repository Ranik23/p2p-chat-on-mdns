package main

import (
	"flag"
)

type config struct {
	rendevouz  string
	ProtocolID string
	listenHost string
	listenPort int
}

func parseFlags() *config {
	c := &config{}

	flag.StringVar(&c.rendevouz, "rendezvous", "meetme", "")
	flag.StringVar(&c.listenHost, "host", "0.0.0.0", "")
	flag.StringVar(&c.ProtocolID, "pid", "/chat/1.1.0", "")
	flag.IntVar(&c.listenPort, "port", 0, "")

	flag.Parse()
	return c
}
