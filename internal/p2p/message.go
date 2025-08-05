package p2p

import (
	"io"

	"github.com/libp2p/go-libp2p/core/peer"
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
	From    peer.ID
	Payload io.Reader
}
