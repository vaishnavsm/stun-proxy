package connector

import (
	"log"
	"net"
	"sync"

	brokerconn "github.com/vaishnavsm/stun-proxy/proxy/src/pkgs/broker_conn"
	"github.com/vaishnavsm/stun-proxy/proxy/src/pkgs/config"
)

type Connector struct {
	Done        chan bool
	broker      *brokerconn.BrokerConn
	cfg         *config.Config
	serv        net.Listener
	connections sync.Map
}

func NewConnector(serv net.Listener, broker *brokerconn.BrokerConn, cfg *config.Config) (*Connector, error) {
	c := &Connector{
		Done:        make(chan bool, 1),
		serv:        serv,
		broker:      broker,
		cfg:         cfg,
		connections: sync.Map{},
	}
	return c, nil
}

func (c *Connector) waitForConnections() {
	for {
		conn, err := c.serv.Accept()
		if err != nil {
			log.Println("failed to create connection", err)
			break
		}
		ac, err := NewAppConnect(conn, c.broker, c.cfg.Name)
		c.connections.Store(ac, true)
		go func() {
			ac.Serve()
			c.connections.Delete(ac)
			ac.Close()
		}()
	}
	c.Done <- true
}

func (c *Connector) Start() {
	go c.waitForConnections()
}

func (c *Connector) Close() {
	c.serv.Close()
	c.connections.Range(func(key any, value any) bool {
		ac, ok := key.(*AppConnect)
		if !ok {
			return true
		}
		ac.Close()
		return true
	})
}
