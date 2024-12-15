package client

type MsgKind string

const (
	MsgKindConnectionRequest          MsgKind = "connection_request"
	MsgKindConnectApplication         MsgKind = "connect_application"
	MsgKindFailConnection             MsgKind = "fail_connection"
	MsgKindDisconnect                 MsgKind = "disconnect"
	MsgKindRegisterApplication        MsgKind = "register_application"
	MsgKindConnectApplicationResponse MsgKind = "connect_application_response"
)
