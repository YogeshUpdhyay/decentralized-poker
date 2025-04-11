package p2p

import (
	"io"
	"net"
)

type MessageType uint8

const (
	GameMessage MessageType = iota
	PingMessage
	PongMessage
	HandShakeMessage
)

type Message struct {
	Type    MessageType
	From    net.Addr
	Payload io.Reader
}
