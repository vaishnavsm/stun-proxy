package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/vaishnavsm/stun-proxy/proxy/src/pkgs/connector"
	"github.com/vaishnavsm/stun-proxy/proxy/src/pkgs/proxy"
)

var broker = flag.String("broker", ":8080", "address of the broker")
var mode = flag.String("mode", "connector", "mode of operation: `connector` allows you to connect to it to reach an application, `proxy` allows you to register applications to it and allows others to connect to it")
var connectorAddr = flag.String("connectorAddr", ":8081", "[connector] the address to bind to")
var proxyApplication = flag.String("proxyConfig", "{}", "[proxy] JSON specifying the application this proxy exposes. set to default to see schema.")

func main() {
	flag.Parse()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	if mode == nil {
		log.Fatal("operation mode was nil")
	}
	if *mode == "proxy" {
		proxy.Run(proxyApplication, broker, interrupt)
		return
	}
	if *mode == "connector" {
		connector.Run(connectorAddr, broker, interrupt)
		return
	}
	log.Fatalf("unknown mode: %s - expected proxy or connector", *mode)
}
