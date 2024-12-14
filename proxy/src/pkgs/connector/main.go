package connector

import (
	"log"
	"net"

	brokerconn "github.com/vaishnavsm/stun-proxy/proxy/src/pkgs/broker_conn"
	"github.com/vaishnavsm/stun-proxy/proxy/src/pkgs/config"
)

func Run(c *config.Config) {

	tcpaddr, err := net.ResolveTCPAddr("tcp", c.Addr)
	if err != nil {
		log.Fatal("could not resolve tcp address", c.Addr, err)
	}

	serv, err := net.ListenTCP("tcp", tcpaddr)
	if err != nil {
		log.Fatal("could not start tcp server at address", c.Addr, err)
	}

	b, err := brokerconn.New(c.BrokerAddr, c.Interrupt)
	if err != nil {
		log.Fatal("could not connect to broker", err)
	}

	defer serv.Close()
	log.Println("started connector server at", c.Addr)

	connector, err := NewConnector(serv, b, c)
	if err != nil {
		log.Fatal("could not start connector", err)
	}
	connector.Start()
	defer connector.Close()
	select {
	case <-connector.Done:
	case <-c.Interrupt:
	}
}
