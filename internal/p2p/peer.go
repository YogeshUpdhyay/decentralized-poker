package p2p

import (
	"bytes"

	libp2pnetwork "github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
)

type Peer struct {
	Username string
	Avatar   string
	conn     libp2pnetwork.Stream
	Status   string
}

func (p *Peer) Send(b []byte) error {
	_, err := p.conn.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func (p *Peer) ReadLoop(msgCh chan *Message, delPeer chan peer.ID) {
	// read loop for a stream
	buff := make([]byte, 4096)
	for {
		n, err := p.conn.Read(buff)
		if err != nil {
			break
		}

		msg := Message{
			Type:    GameMessage,
			From:    p.conn.Conn().RemotePeer(),
			Payload: bytes.NewReader(buff[:n]),
		}

		msgCh <- &msg
	}

	// close connection and mark peer as deleted in the list
	p.conn.Close()
	delPeer <- p.conn.Conn().RemotePeer()
}
