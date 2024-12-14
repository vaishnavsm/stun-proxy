package proxy

import (
	"io"
	"log"
	"net"
	"os"
	"sync"

	"github.com/pkg/errors"
	brokerconn "github.com/vaishnavsm/stun-proxy/proxy/src/pkgs/broker_conn"
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
		conn, ok := key.(*Connection)
		if !ok {
			return true
		}
		conn.backend.Close()
		conn.frontend.Close()
		return true
	})
}

func proxy(dst, src net.Conn, srcClosed chan struct{}, logdata any) {
	_, err := io.Copy(dst, src)

	if err != nil {
		log.Println("proxy error during copy", err, logdata)
	}
	if err := src.Close(); err != nil {
		log.Println("proxy error during connection close", err, logdata)
	}
	srcClosed <- struct{}{}
}

func (a *Application) createConnection(req brokerconn.ConnectRequest) {
	log.Printf("creating a connection: %v\n", req)

	localAddr, err := net.ResolveTCPAddr("tcp", req.LocalAddr)
	if err != nil {
		log.Println("could not resolve tcp local address", err, req)
	}
	remoteAddr, err := net.ResolveTCPAddr("tcp", req.RemoteAddr)
	if err != nil {
		log.Println("could not resolve tcp remote address", err, req)
	}
	frontendConn, err := net.DialTCP("tcp", localAddr, remoteAddr)
	if err != nil {
		log.Printf("failed to establish connection to remote. this could mean STUN failed! %+v", req)
		return
	}
	backendConn, err := net.DialTCP("tcp", nil, a.Addr)
	if err != nil {
		log.Printf("error connecting to upstream %s: %s", a.Addr, err)
		a.broker.FailConnection(req, "failed to connect to upstream server")
		return
		// close the connection
	}

	log.Printf("established connection: %v\n", req)

	bClosed := make(chan struct{}, 1)
	fClosed := make(chan struct{}, 1)
	go proxy(backendConn, frontendConn, fClosed, req)
	go proxy(frontendConn, backendConn, bClosed, req)

	conn := &Connection{
		frontend: frontendConn,
		backend:  backendConn,
	}

	a.connections.Store(conn, true)

	var wait chan struct{}
	select {
	case <-bClosed:
		frontendConn.SetLinger(0)
		frontendConn.CloseRead()
		wait = fClosed
	case <-fClosed:
		backendConn.SetLinger(0)
		backendConn.CloseRead()
		wait = bClosed
	}
	<-wait
	log.Printf("closed connection: %v\n", req)
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
