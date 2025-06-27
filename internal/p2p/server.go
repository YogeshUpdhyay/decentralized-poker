package p2p

import (
	"bytes"
	"encoding/gob"
	"errors"
	"net"

	"github.com/YogeshUpdhyay/ypoker/internal/constants"
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
	ServerName string
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
	if err := s.transport.ListenAndAccept(s.ServerName); err != nil {
		logrus.Fatal("error starting the server")
	}
}

func (s *Server) loop() {
	for {
		select {
		case conn := <-s.addPeer:
			// validating peer using handshake
			peer := Peer{conn: conn}
			logrus.WithField(constants.ServerName, s.ServerName).Infof("awaiting handshake from %s", peer.conn.RemoteAddr())
			if err := s.handshake(&peer); err != nil {
				logrus.Infof("handshake failed with peer %s due to %s", peer.conn.RemoteAddr(), err)
				break
			}

			// handshake was success adding it to the peers
			// starting the read loop
			s.peers[conn.RemoteAddr()] = &peer
			go peer.ReadLoop(s.msgCh, s.deletePeer)

			logrus.WithField(constants.ServerName, s.ServerName).Infof("new player connected %s", peer.conn.RemoteAddr())
		case msg := <-s.msgCh:
			s.handler.HandleMessage(msg)
		case conn := <-s.deletePeer:
			// connection is closed removing from peers
			delete(s.peers, conn.RemoteAddr())
			logrus.Infof("player disconnected %s", conn.RemoteAddr())
		}
	}
}

func (s *Server) handshake(p *Peer) error {
	// read the handshake message and validate
	buff := make([]byte, 1024)
	n, err := p.conn.Read(buff)
	if err != nil {
		logrus.Errorf("error reading hanshake data %s", err)
		return err
	}

	hs := Handshake{}
	if err := gob.NewDecoder(bytes.NewReader(buff[:n])).Decode(&hs); err != nil {
		logrus.Errorf("error decoding the handshake data %s", err)
		return err
	}

	if s.Version != hs.Version {
		logrus.Error("version mismatch")
		return errors.New("invalid peer: version mismatch")
	}

	logrus.WithField(constants.ServerName, s.ServerName).Info("handshake recieved and version matched.")
	return nil
}

func (s *Server) Connect(addr string) error {
	// connecting to the peer
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		logrus.Errorf("error connecting to peer")
		return err
	}

	// validating the peer connection using the handshake
	logrus.WithField(constants.ServerName, s.ServerName).Info("sending handshake message")
	err = s.sendHandshake(conn)
	if err != nil {
		logrus.Errorf("error sending handshake %s", err)
		return err
	}

	// handshake was success adding it to the peers
	// starting the read loop
	peer := Peer{conn: conn}
	s.peers[conn.RemoteAddr()] = &peer
	go peer.ReadLoop(s.msgCh, s.deletePeer)

	logrus.WithField(constants.ServerName, s.ServerName).Infof("new player connected %s", peer.conn.RemoteAddr())

	return nil
}

func (s *Server) sendHandshake(conn net.Conn) error {
	// preparing the handshake message
	var buf bytes.Buffer

	hs := Handshake{Version: s.Version}

	if err := gob.NewEncoder(&buf).Encode(hs); err != nil {
		logrus.Errorf("error creating the handshake message: %s", err)
		return err
	}

	// Writing handshake to the connection
	_, err := conn.Write(buf.Bytes())
	if err != nil {
		logrus.Errorf("error writing handshake to the connection %s", err)
		return err
	}

	logrus.WithField(constants.ServerName, s.ServerName).Info("hanshake message sent successfully")
	return nil
}
