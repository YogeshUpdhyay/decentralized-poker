package p2p

import (
	"bytes"
	"encoding/gob"
	"net"

	"github.com/sirupsen/logrus"
)

type Server struct {
	ServerConfig
	handler    Handler
	transport  *TCPTransport
	peers      map[net.Addr]*Peer
	addPeer    chan net.Conn
	deletePeer chan net.Conn
	msgCh      chan *Message
}

type ServerConfig struct {
	ListenAddr string
	Version    string
}

type Handshake struct {
	Version string
}

func NewServer(cfg ServerConfig) *Server {
	server := &Server{
		ServerConfig: cfg,
		handler:      &DefaultHandler{},
		peers:        make(map[net.Addr]*Peer),
		addPeer:      make(chan net.Conn),
		deletePeer:   make(chan net.Conn),
		msgCh:        make(chan *Message),
	}

	server.transport = &TCPTransport{
		ListenAddr: server.ListenAddr,
		addPeer:    server.addPeer,
		delPeer:    server.deletePeer,
	}

	return server
}

func (s *Server) Start() {
	go s.loop()
	if err := s.transport.ListenAndAccept(); err != nil {
		logrus.Fatal("error starting the server")
	}
}

func (s *Server) loop() {
	for {
		select {
		case conn := <-s.addPeer:
			// validating peer using handshake
			peer := Peer{conn: conn}
			if err := s.handshake(&peer); err != nil {
				logrus.Infof("handshake failed with peer %s due to %s\n", peer.conn.RemoteAddr(), err)
				continue
			}

			// handshake was success adding it to the peers
			// starting the read loop
			s.peers[conn.RemoteAddr()] = &peer
			go peer.ReadLoop(s.msgCh, s.deletePeer)

			logrus.Infof("new player connected %s\n", peer.conn.RemoteAddr())
		case msg := <-s.msgCh:
			s.handler.HandleMessage(msg)
		case conn := <-s.deletePeer:
			// connection is closed removing from peers
			delete(s.peers, conn.RemoteAddr())
			logrus.Infof("player disconnected %s\n", conn.RemoteAddr())
		}
	}
}

func (s *Server) handshake(p *Peer) error {
	// read the handshake message and validate
	buff := make([]byte, 1024)
	n, err := p.conn.Read(buff)
	if err != nil {
		logrus.Errorf("error reading hanshake data %s\n", err)
		return err
	}

	hs := Handshake{}
	if err := gob.NewDecoder(bytes.NewReader(buff[:n])).Decode(&hs); err != nil {
		logrus.Errorf("error decoding the handshake data %s\n", err)
	}

	if s.Version != hs.Version {

	}

	return nil
}
