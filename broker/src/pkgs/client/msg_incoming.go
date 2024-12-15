package client

import (
	"encoding/json"
	"log"
)

type MessageKind struct {
	Kind MsgKind `json:"kind"`
}

type MsgFailConnection struct {
	Kind    MsgKind `json:"kind"`
	Payload string  `json:"payload"`
}

type MsgConnectApplicationPayload struct {
	Name         string `json:"name"`
	ConnectionId string `json:"connectionId"`
}

type MsgConnectApplication struct {
	Kind    MsgKind                      `json:"kind"`
	Payload MsgConnectApplicationPayload `json:"payload"`
}

type MsgRegisterApplicationPayload struct {
	Name string `json:"name"`
}

type MsgRegisterApplication struct {
	Kind    MsgKind                       `json:"kind"`
	Payload MsgRegisterApplicationPayload `json:"payload"`
}

func (c *Client) handleMsgFailConnection(msg []byte) {

}
func (c *Client) handleMsgConnectApplication(msg []byte) {
	d := &MsgConnectApplication{}
	err := json.Unmarshal(msg, d)
	if err != nil {
		log.Printf("could not parse connect message %s: ERROR: %s\n", string(msg), err)
		return
	}

	c.b.ConnectToApplication(c, d.Payload.Name, d.Payload.ConnectionId)
}
func (c *Client) handleMsgRegisterApplication(msg []byte) {
	d := &MsgRegisterApplication{}
	err := json.Unmarshal(msg, d)
	if err != nil {
		log.Printf("could not parse register message %s: ERROR: %s\n", string(msg), err)
		return
	}

	c.b.RegisterApplication(c, d.Payload.Name)
}

func (c *Client) handleMessage(msg []byte) {
	log.Printf("got message from client %s: %s\n", c.Id, string(msg))
	k := &MessageKind{}
	err := json.Unmarshal(msg, k)
	if err != nil {
		log.Println("error parsing received message: ", err, string(msg))
		return
	}

	switch k.Kind {
	case MsgKindFailConnection:
		{
			c.handleMsgFailConnection(msg)
		}
	case MsgKindConnectApplication:
		{
			c.handleMsgConnectApplication(msg)
		}
	case MsgKindRegisterApplication:
		{
			c.handleMsgRegisterApplication(msg)
		}
	}

}
