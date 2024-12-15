package client

import "log"

type ConnectRequest struct {
	ConnectionId string `json:"connectionId"`
	LocalAddr    string `json:"localAddr"`
	RemoteAddr   string `json:"remoteAddr"`
	Error        string `json:"error"`
}

func (c *Client) SendMsgConnectionRequest(payload ConnectRequest) error {
	err := c.conn.WriteJSON(map[string]interface{}{
		"kind":    MsgKindConnectionRequest,
		"payload": payload,
	})
	if err != nil {
		log.Printf("error sending connection request to client %s: %s\n%v\n", c.Id, err, payload)
	}
	return err
}
func (c *Client) SendMsgConnectApplicationResponse(payload ConnectRequest) error {
	err := c.conn.WriteJSON(map[string]interface{}{
		"kind":    MsgKindConnectionRequest,
		"payload": payload,
	})
	if err != nil {
		log.Printf("error sending connect application response to client %s: %s\n%v\n", c.Id, err, payload)
	}
	return err
}
