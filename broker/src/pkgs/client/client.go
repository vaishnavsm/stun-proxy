package client

import (
	"fmt"
	"log"
	"math/rand/v2"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

type Broker interface {
	RegisterApplication(c *Client, name string)
	ConnectToApplication(c *Client, name string, connectionId string)
	Disconnect(c *Client) error
}

type Client struct {
	conn *websocket.Conn
	b    Broker
	Id   string
}

func (c *Client) RemoteIp() string {
	s := c.conn.RemoteAddr().String()
	ipSegments := strings.Split(s, ":")
	if len(ipSegments) == 2 {
		// ipv4
		return ipSegments[0]
	}
	// ipv6, i don't really know if this works lol
	return strings.Join(ipSegments[:len(ipSegments)-1], ":")
}

func (c *Client) RandomRemoteAddr() string {
	ip := c.RemoteIp()
	port := rand.UintN(0000) + 40000 // random port between 40-50k
	return fmt.Sprintf("%s:%d", ip, port)
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
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("unexpected websocket close for client %s, %v\n", c.Id, err)
			} else {
				log.Printf("client disconnected %s, %v\n", c.Id, err)
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
