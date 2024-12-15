package broker

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/vaishnavsm/stun-proxy/broker/src/pkgs/client"
)

type Broker struct {
	clients      map[*client.Client]bool
	applications map[string]*client.Client
}

func New() (*Broker, error) {
	broker := &Broker{
		clients:      make(map[*client.Client]bool),
		applications: make(map[string]*client.Client),
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

func (b *Broker) RegisterApplication(c *client.Client, name string) {
	if ac, ok := b.applications[name]; ok {
		log.Printf("Tried to register application %s on client %s when client %s already has it. Ignoring.\n", name, c.Id, ac.Id)
		return
	}
	b.applications[name] = c
	log.Printf("Registered application %s on client %s\n", name, c.Id)
}

func (b *Broker) ConnectToApplication(c *client.Client, name string, connectionId string) {
	ac, ok := b.applications[name]
	if !ok {
		log.Printf("Client %s tried to connect to unknown application %s\n", c.Id, name)
		err := c.SendMsgConnectApplicationResponse(client.ConnectRequest{
			ConnectionId: connectionId,
			Error:        "unknown application " + name,
		})
		if err != nil {
			log.Printf("Failed responding to client %s: %v\n", c.Id, err)
		}
		return
	}

	log.Printf("Initiating connection between clients %s and %s for application %s\n", c.Id, ac.Id, name)

	serverRemote := ac.RandomRemoteAddr()
	clientRemote := c.RandomRemoteAddr()

	log.Printf("Allocating addresses:\n\tserver: %s\n\tclient: %s\nGodspeed!\n", serverRemote, clientRemote)

	c.SendMsgConnectApplicationResponse(client.ConnectRequest{
		ConnectionId: connectionId,
		LocalAddr:    clientRemote,
		RemoteAddr:   serverRemote,
	})
	ac.SendMsgConnectionRequest(client.ConnectRequest{
		ConnectionId: connectionId,
		LocalAddr:    serverRemote,
		RemoteAddr:   clientRemote,
	})
}

func (b *Broker) Disconnect(c *client.Client) error {
	b.clients[c] = false
	// Remove applications associated with c
	// for active application connections to c, tell the clients that the server died
	// note that this doesn't mean that the connection died for the client, just that if it did, then the proxy won't be able to reconnect
	return nil
}
