package p2p

import (
	"context"
	"encoding/json"
	"strings"

	p2pModels "github.com/YogeshUpdhyay/ypoker/internal/p2p/models"
	libp2pnetwork "github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	log "github.com/sirupsen/logrus"
)

type Peer struct {
	PeerID string
	conn   libp2pnetwork.Stream
	Status string
}

func (p *Peer) Send(msg *p2pModels.Envelope) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	// append newline so decoder on the other side can read messages cleanly
	data = append(data, '\n')

	_, err = p.conn.Write(data)
	return err
}

func (p *Peer) ReadLoop(ctx context.Context, msgCh chan *p2pModels.Envelope, delPeer chan peer.ID) {
	defer func() {
		p.conn.Close()
		delPeer <- p.conn.Conn().RemotePeer()
	}()

	decoder := json.NewDecoder(p.conn)
	for {
		var envelope p2pModels.Envelope
		if err := decoder.Decode(&envelope); err != nil {
			if strings.Contains(err.Error(), "EOF") {
				log.WithContext(ctx).Infof("connection closed")
				return
			}
			// log unexpected error and continue
			log.Printf("read error from peer %s: %v", p.conn.Conn().RemotePeer(), err)
			break
		}

		// Add peer ID from connection if missing
		if envelope.From == "" {
			envelope.From = p.conn.Conn().RemotePeer().String()
		}
		log.WithContext(ctx).Infof("new message received from %s of type %s", envelope.From, envelope.Type)
		msgCh <- &envelope
	}
}

func (p *Peer) GetMultiAddrs(ctx context.Context) multiaddr.Multiaddr {
	connMultiAddrs := p.conn.Conn().RemoteMultiaddr()
	return connMultiAddrs
}
