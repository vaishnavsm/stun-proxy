package broker

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/vaishnavsm/stun-proxy/broker/src/pkgs/client"
)

type Broker struct {
	clients map[*client.Client]bool
}

func New() (*Broker, error) {
	broker := &Broker{
		clients: make(map[*client.Client]bool),
	}
	return broker, nil
}

func (b *Broker) CreateClient(conn *websocket.Conn) error {
	c, err := client.New(conn, b)
	if err != nil {
		return errors.Wrap(err, "error creating a client")
	}
	log.Printf("connected to a new client: %s\n", c.Id)

	b.clients[c] = true
	return nil
}

func (b *Broker) RegisterApplication() {

}

func (b *Broker) ConnectToApplication() {

}

func (b *Broker) Disconnect(c *client.Client) error {
	b.clients[c] = false
	// Remove applications associated with c
	// for active application connections to c, tell the clients that the server died
	// note that this doesn't mean that the connection died for the client, just that if it did, then the proxy won't be able to reconnect
	return nil
}
