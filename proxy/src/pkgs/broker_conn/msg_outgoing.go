package brokerconn

import (
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (b *BrokerConn) RegisterApplication(name string) error {
	return b.conn.WriteJSON(map[string]interface{}{
		"kind": MsgKindRegisterApplication,
		"payload": map[string]interface{}{
			"name": name,
		},
	})
}

func (b *BrokerConn) Disconnect() error {
	b.Close()
	return nil
}

func (b *BrokerConn) FailConnection(req ConnectRequest, msg string) error {
	return b.conn.WriteJSON(map[string]interface{}{
		"kind":    MsgKindFailConnection,
		"payload": msg,
	})
}

func (b *BrokerConn) ConnectApplication(appName string) (ConnectRequest, error) {
	connectionIdFull, err := uuid.NewUUID()
	if err != nil {
		return ConnectRequest{}, errors.Wrap(err, "error creating connection id")
	}
	connectionId := connectionIdFull.String()
	ch := make(chan ConnectRequest, 1)
	b.connectionQueue[connectionId] = ch
	err = b.conn.WriteJSON(map[string]interface{}{
		"kind": MsgKindConnectApplication,
		"payload": map[string]interface{}{
			"name": appName,
		},
	})
	if err != nil {
		delete(b.connectionQueue, connectionId)
		return ConnectRequest{}, errors.Wrap(err, "error sending message to broker")
	}

	select {
	case req := <-ch:
		{
			delete(b.connectionQueue, connectionId)
			if req.Error != "" {
				return req, errors.Errorf("error present in response body: %s", req.Error)
			}
			return req, nil
		}
	case <-time.After(2 * time.Minute):
		{
			return ConnectRequest{}, errors.Errorf("request timed out for connection %s", connectionId)
		}
	}
}
