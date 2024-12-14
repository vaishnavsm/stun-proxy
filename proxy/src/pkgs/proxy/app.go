package proxy

import (
	"log"
	"net"
	"os"
	"sync"

	"github.com/pkg/errors"
	brokerconn "github.com/vaishnavsm/stun-proxy/proxy/src/pkgs/broker_conn"
	"github.com/vaishnavsm/stun-proxy/proxy/src/pkgs/stunnel"
)

type Application struct {
	Name        string
	Id          string
	Addr        *net.TCPAddr
	broker      *brokerconn.BrokerConn
	interrupt   chan os.Signal
	connections sync.Map
}

type Connection struct {
	frontend net.Conn
	backend  net.Conn
}

func New(name string, interrupt chan os.Signal, broker *brokerconn.BrokerConn, addr string, appId string) (*Application, error) {
	tcpaddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, errors.Wrap(err, "error resolving tcp address")
	}
	app := &Application{
		Name:        name,
		interrupt:   interrupt,
		broker:      broker,
		Addr:        tcpaddr,
		Id:          appId,
		connections: sync.Map{},
	}
	return app, nil
}

func (a *Application) close() {
	a.broker.Disconnect()
	// close connections
	a.connections.Range(func(key any, value any) bool {
		conn, ok := key.(*stunnel.Stunnel)
		if !ok {
			return true
		}
		conn.Close()
		return true
	})
}

func (a *Application) createConnection(req brokerconn.ConnectRequest) {
	log.Printf("creating a connection: %v\n", req)

	backendConn, err := net.DialTCP("tcp", nil, a.Addr)
	if err != nil {
		log.Printf("error connecting to upstream %s: %s\n", a.Addr, err)
		a.broker.FailConnection(req, "failed to connect to upstream server")
		return
		// close the connection
	}

	conn, err := stunnel.New(req, backendConn)
	if err != nil {
		log.Printf("error connecting to peer: %s\n", err)
		a.broker.FailConnection(req, "failed to connect to peer")
		return
	}

	defer conn.Close()

	a.connections.Store(conn, true)

	conn.Serve()
}

func (a *Application) Serve() {
	log.Println("Proxy is ready and waiting for connections")
	defer a.close()
	for {
		select {
		case req := <-a.broker.ConnectRequests:
			{
				go a.createConnection(req)
			}
		case <-a.interrupt:
			{
				break
			}
		}
	}
}
