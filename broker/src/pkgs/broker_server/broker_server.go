package broker_server

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/vaishnavsm/stun-proxy/broker/src/pkgs/broker"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type BrokerServer struct {
	b *broker.Broker
}

func New() (*BrokerServer, error) {
	broker, err := broker.New()
	if err != nil {
		return nil, err
	}
	return &BrokerServer{
		b: broker,
	}, nil
}

func (bs *BrokerServer) Serve(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection", err)
		http.Error(w, "Failed to upgrade connection to ws", http.StatusUpgradeRequired)
		return
	}

	err = bs.b.CreateClient(conn)
	if err != nil {
		log.Println("Failed to create client after upgrading connection", err)
		http.Error(w, "server failed to create client", http.StatusInternalServerError)
		return
	}
}
