package brokerconn

import (
	"encoding/json"
	"log"
)

type MessageKind struct {
	Kind MsgKind `json:"kind"`
}

type MsgConnect struct {
	Kind    MsgKind        `json:"kind"`
	Payload ConnectRequest `json:"payload"`
}

func (b *BrokerConn) handleMsgConnectionRequest(msg []byte) {
	d := &MsgConnect{}
	err := json.Unmarshal(msg, d)
	if err != nil {
		log.Printf("could not parse connection message %s: ERROR: %s\n", string(msg), err)
		return
	}

	b.ConnectRequests <- d.Payload
}

func (b *BrokerConn) handleMsgConnectApplicationResponse(msg []byte) {
	d := &MsgConnect{}
	err := json.Unmarshal(msg, d)
	if err != nil {
		log.Printf("could not parse connection message %s: ERROR: %s\n", string(msg), err)
		return
	}
	ch, ok := b.connectionQueue[d.Payload.ConnectionId]
	if !ok {
		log.Printf("tried to create connection to unknown connection id %s\n", d.Payload.ConnectionId)
		return
	}
	ch <- d.Payload
}

func (b *BrokerConn) handleMsgDiagnostic(msg []byte) {
	log.Printf("received diagnostic message: %v", string(msg))
}

func (b *BrokerConn) handleMessage(msg []byte) {
	k := &MessageKind{}
	err := json.Unmarshal(msg, k)
	if err != nil {
		log.Println("error parsing received message: ", err, string(msg))
		return
	}

	switch k.Kind {
	case MsgKindConnectionRequest:
		{
			b.handleMsgConnectionRequest(msg)
		}
	case MsgKindFailConnection:
		{
			b.handleMsgDiagnostic(msg)
		}
	case MsgKindDisconnect:
		{
			b.handleMsgDiagnostic(msg)
		}
	case MsgKindConnectApplicationResponse:
		{
			b.handleMsgConnectApplicationResponse(msg)
		}
	}
}
