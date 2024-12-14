package client

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

type Broker interface {
	RegisterApplication()
	ConnectToApplication()
	Disconnect(c *Client) error
}

type Client struct {
	conn *websocket.Conn
	b    Broker
	Id   string
}

func (c *Client) close() {
	err := c.b.Disconnect(c)
	if err != nil {
		log.Printf("error disconnecting client from broker for %s: %v\n", c.Id, err)
	}
	err = c.conn.Close()
	if err != nil {
		log.Printf("error cleanly closing client connection for %s: %v\n", c.Id, err)
	}
}

var maxTimeout = 30 * time.Second

func (c *Client) handleMessage(msg []byte) {
	log.Printf("got message from client %s: %s\n", c.Id, string(msg))
}

func (c *Client) reader() {
	defer c.close()

	// Handle timeouts
	c.conn.SetReadDeadline(time.Now().Add(maxTimeout))
	c.conn.SetPongHandler(func(appData string) error {
		c.conn.SetReadDeadline(time.Now().Add(maxTimeout))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected websocket close for client %s, %v\n", c.Id, err)
			} else {
				log.Printf("websocket close for client %s, %v\n", c.Id, err)
			}
			break
		}
		c.handleMessage(message)
	}
}

func New(conn *websocket.Conn, b Broker) (*Client, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, errors.Wrap(err, "error creating a UUID for client")
	}
	client := &Client{
		conn: conn,
		b:    b,
		Id:   id.String(),
	}
	go client.reader()
	return client, nil
}
