package stunnel

import (
	"io"
	"log"
	"net"

	"github.com/pkg/errors"
	brokerconn "github.com/vaishnavsm/stun-proxy/proxy/src/pkgs/broker_conn"
)

type Stunnel struct {
	req        brokerconn.ConnectRequest
	oConn      *net.TCPConn
	sConn      *net.TCPConn
	localAddr  *net.TCPAddr
	remoteAddr *net.TCPAddr
}

func New(req brokerconn.ConnectRequest, otherConnection *net.TCPConn) (*Stunnel, error) {
	localAddr, err := net.ResolveTCPAddr("tcp", req.LocalAddr)
	if err != nil {
		return nil, errors.Wrapf(err, "could not resolve tcp local address: %+v", req)
	}
	remoteAddr, err := net.ResolveTCPAddr("tcp", req.RemoteAddr)
	if err != nil {
		return nil, errors.Wrapf(err, "could not resolve tcp remote address: %+v", req)
	}
	s := &Stunnel{
		req:        req,
		oConn:      otherConnection,
		localAddr:  localAddr,
		remoteAddr: remoteAddr,
	}
	return s, nil
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

func (s *Stunnel) Close() {
	s.oConn.Close()
	s.sConn.Close()
}

func (s *Stunnel) Serve() {
	sConn, err := net.DialTCP("tcp", s.localAddr, s.remoteAddr)
	if err != nil {
		log.Printf("failed to establish connection to remote. this could mean STUN failed! %+v\n", s.req)
		return
	}
	s.sConn = sConn

	defer s.Close()

	log.Printf("established connection: %v\n", s.req)

	sClosed := make(chan struct{}, 1)
	oClosed := make(chan struct{}, 1)
	go proxy(s.oConn, sConn, sClosed, s.req)
	go proxy(sConn, s.oConn, oClosed, s.req)

	var wait chan struct{}
	select {
	case <-sClosed:
		s.oConn.SetLinger(0)
		s.oConn.CloseRead()
		wait = oClosed
	case <-oClosed:
		sConn.SetLinger(0)
		sConn.CloseRead()
		wait = sClosed
	}
	<-wait
	log.Printf("closed connection: %v\n", s.req)
}
