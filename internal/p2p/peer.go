package p2p

import (
	"bytes"
	"net"
)

type Peer struct {
	conn net.Conn
}

func (p *Peer) Send(b []byte) error {
	_, err := p.conn.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func (p *Peer) ReadLoop(msgCh chan *Message, delPeer chan net.Conn) {
	buff := make([]byte, 1024)
	for {
		n, err := p.conn.Read(buff)
		if err != nil {
			break
		}

		msg := Message{
			Type:    GameMessage,
			From:    p.conn.RemoteAddr(),
			Payload: bytes.NewReader(buff[:n]),
		}
		msgCh <- &msg
	}
	p.conn.Close()
	delPeer <- p.conn
}
