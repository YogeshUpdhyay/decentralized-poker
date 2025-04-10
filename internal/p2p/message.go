package p2p

import (
	"io"
	"net"
)

type Message struct {
	From    net.Addr
	Payload io.Reader
}
