package p2p

import (
	"fmt"

	"github.com/YogeshUpdhyay/ypoker/internal/constants"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	libp2pnetwork "github.com/libp2p/go-libp2p/core/network"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"
)

type P2PTransport struct {
	host             host.Host
	ListenAddr       string
	addPeer          chan libp2pnetwork.Stream
	identityFilePath string
}

func (t *P2PTransport) ListenAndAccept(serverName string) error {
	encryption := DefaultEncryption{}
	edPrivKey, err := encryption.LoadAndDecryptKey(constants.Empty, t.identityFilePath)

	// start listening for peers
	h, err := libp2p.New(
		libp2p.Identity(edPrivKey),
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", t.ListenAddr)),
	)
	if err != nil {
		return err
	}
	t.host = h

	logrus.WithField(constants.ServerName, fmt.Sprintf("%s@%s", serverName, h.ID().String())).Infof("server listening at %v", t.GetMyFullAddr())

	t.host.SetStreamHandler("/ypoker/1.0.0", func(s libp2pnetwork.Stream) {
		logrus.WithField(constants.ServerName, serverName).Infof("incoming stream from %s adding to peer list", s.Conn().RemotePeer())
		t.addPeer <- s
	})

	return nil
}

func (t *P2PTransport) GetMyFullAddr() []string {
	var addrs []string
	for _, addr := range t.host.Addrs() {
		// Append /p2p/<peer-id> to each base address
		fullAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/p2p/%s", t.host.ID().String()))
		addrWithID := addr.Encapsulate(fullAddr)
		addrs = append(addrs, addrWithID.String())
	}
	return addrs
}
