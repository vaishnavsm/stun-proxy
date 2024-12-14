package proxy

import (
	"log"

	brokerconn "github.com/vaishnavsm/stun-proxy/proxy/src/pkgs/broker_conn"
	"github.com/vaishnavsm/stun-proxy/proxy/src/pkgs/config"
)

func Run(c *config.Config) {

	b, err := brokerconn.New(c.BrokerAddr, c.Interrupt)
	if err != nil {
		log.Fatal("could not connect to broker", err)
	}

	appId, err := b.RegisterApplication(c.Name)
	if err != nil {
		log.Fatalf("Could not register application %s: %v", c.Name, err)
	}

	app, err := New(c.Name, c.Interrupt, b, c.Addr, appId)
	if err != nil {
		log.Fatalf("Could not create application %s: %v", c.Name, err)
	}

	app.Serve()
}
