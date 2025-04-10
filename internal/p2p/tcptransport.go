package p2p

import (
	"net"

	"github.com/sirupsen/logrus"
)

type TCPTransport struct {
	ListenAddr string
	listener   net.Listener
	addPeer    chan net.Conn
	delPeer    chan net.Conn
}

func (t *TCPTransport) ListenAndAccept() error {
	// start listening for peers
	listener, err := net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}

	logrus.Infof("server listening at %s\n", t.ListenAddr)
	t.listener = listener

	// accepting connections
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			logrus.Infof("error accepting connection %s\n", err)
		}

		// all accepted connections need to pass the handshake protocol
		// if not the connection needs to be terminated
		// adding the connection to the add peer for handshake and peer creation
		t.addPeer <- conn
	}
}
