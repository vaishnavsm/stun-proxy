package brokerconn

import (
	"log"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

type ConnectRequest struct {
	LocalAddr  string
	RemoteAddr string
}

type BrokerConn struct {
	interrupt       chan os.Signal
	ConnectRequests chan ConnectRequest
}

func New(brokerAddr string, interrupt chan os.Signal) (*BrokerConn, error) {
	panic("todo")
}

func (b *BrokerConn) RegisterApplication(name string) (string, error) {
	panic("todo")
}

func (b *BrokerConn) Disconnect() error {
	panic("todo")
}

func (b *BrokerConn) FailConnection(req ConnectRequest, msg string) error {
	panic("todo")
}

func (b *BrokerConn) ConnectApplication(appName string) (ConnectRequest, error) {
	panic("todo")
}

func (b *BrokerConn) connectToBroker() {
	u := url.URL{Scheme: "ws", Host: *broker, Path: "/ws"}
	log.Printf("connecting to broker %s\n", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("could not connect to broker", err)
	}

	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
					log.Println("read error:", err)
				}
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println("write error:", err)
				return
			}
		case <-b.interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
