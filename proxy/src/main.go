package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/vaishnavsm/stun-proxy/proxy/src/pkgs/config"
	"github.com/vaishnavsm/stun-proxy/proxy/src/pkgs/connector"
	"github.com/vaishnavsm/stun-proxy/proxy/src/pkgs/proxy"
)

var broker = flag.String("broker", ":8080", "address of the broker")
var mode = flag.String("mode", "connector", "mode of operation: `connector` allows you to connect to it to reach an application, `proxy` allows you to register applications to it and allows others to connect to it")
var addr = flag.String("addr", ":8081", "connector: the address for the connector to bind to, proxy: the address of the upstream server")
var app = flag.String("name", "", "name of the application")

func main() {
	flag.Parse()

	if *app == "" {
		log.Fatal("the app name cannot be blank. please set the `name` flag.")
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	c := &config.Config{
		Name:       *app,
		Addr:       *addr,
		BrokerAddr: *broker,
		Interrupt:  interrupt,
	}

	if mode == nil {
		log.Fatal("operation mode was nil")
	}
	if *mode == "proxy" {
		proxy.Run(c)
		return
	}
	if *mode == "connector" {
		connector.Run(c)
		return
	}
	log.Fatalf("unknown mode: %s - expected proxy or connector", *mode)
}
