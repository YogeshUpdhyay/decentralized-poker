package models

import (
	"encoding/json"
)

// MessageType defines the type of message being sent
type MessageType string

const (
	MsgTypeHandshake       MessageType = "handshake"
	MsgTypeHandshakeAck    MessageType = "handshake_ack"
	MsgTypeHandshakeReject MessageType = "handshake_reject"
	MsgTypeChat            MessageType = "chat" // or "game"
)

type Envelope struct {
	Type    MessageType     `json:"type"`
	From    string          `json:"from"` // optional, can also be set from conn
	Payload json.RawMessage `json:"payload"`
}
