package brokerconn

import (
	"log"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

type ConnectRequest struct {
	ConnectionId string `json:"connectionId"`
	LocalAddr    string `json:"localAddr"`
	RemoteAddr   string `json:"remoteAddr"`
}

type BrokerConn struct {
	interrupt       chan os.Signal
	addr            string
	ConnectRequests chan ConnectRequest
	connectionQueue map[string]chan ConnectRequest
	conn            *websocket.Conn
	readerClosed    chan bool
}

func New(brokerAddr string, interrupt chan os.Signal) (*BrokerConn, error) {
	b := &BrokerConn{
		addr:            brokerAddr,
		interrupt:       interrupt,
		ConnectRequests: make(chan ConnectRequest, 512),
	}
	err := b.connectToBroker()
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (b *BrokerConn) Close() {
	err := b.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("write close:", err)
		return
	}
	select {
	case <-b.readerClosed:
	case <-time.After(time.Second):
	}
	b.conn.Close()
}

func (b *BrokerConn) serve() {
	defer b.conn.Close()

	b.readerClosed = make(chan bool)

	go func() {
		defer close(b.readerClosed)
		for {
			_, message, err := b.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
					log.Println("read error:", err)
				}
				return
			}
			b.handleMessage(message)
		}
	}()

	for {
		select {
		case <-b.readerClosed:
			return
		case <-b.interrupt:
			log.Println("interrupted")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			b.Close()
			return
		}
	}
}

func (b *BrokerConn) connectToBroker() error {
	u := url.URL{Scheme: "ws", Host: b.addr, Path: "/ws"}
	log.Printf("connecting to broker %s\n", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("could not connect to broker", err)
		return err
	}

	b.conn = c
	go b.serve()
	return nil
}
