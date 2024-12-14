package connector

import (
	"log"
	"net"

	"github.com/pkg/errors"
	brokerconn "github.com/vaishnavsm/stun-proxy/proxy/src/pkgs/broker_conn"
	"github.com/vaishnavsm/stun-proxy/proxy/src/pkgs/stunnel"
)

type AppConnect struct {
	frontConn  net.Conn
	remoteConn *net.TCPConn
	broker     *brokerconn.BrokerConn
	appName    string
	conn       *stunnel.Stunnel
}

func NewAppConnect(conn net.Conn, broker *brokerconn.BrokerConn, appName string) (*AppConnect, error) {
	ac := &AppConnect{
		frontConn: conn,
		broker:    broker,
		appName:   appName,
	}

	return ac, nil
}

func (a *AppConnect) Serve() error {
	req, err := a.broker.ConnectApplication(a.appName)
	if err != nil {
		log.Println("failed to connect to app via broker", err, a.appName)
		return errors.Wrap(err, "failed to connect to app via broker")
	}

	conn, err := stunnel.New(req, a.remoteConn)
	if err != nil {
		log.Println("failed to create tunnel connection with peer", err)
		return errors.Wrap(err, "failed to connect to peer")
	}

	defer conn.Close()

	a.conn = conn

	conn.Serve()
	return nil
}

func (a *AppConnect) Close() {
	a.conn.Close()
}
