package p2p

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"net"
	"path/filepath"

	"github.com/YogeshUpdhyay/ypoker/internal/constants"
	libp2pnetwork "github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"
)

type Server struct {
	ServerConfig
	handler    Handler
	transport  *P2PTransport
	peers      map[peer.ID]*Peer
	addPeer    chan libp2pnetwork.Stream
	deletePeer chan peer.ID
	msgCh      chan *Message
}

type ServerConfig struct {
	ListenAddr       string
	Version          string
	ServerName       string
	IdentityFilePath string
}

func (cfg *ServerConfig) ApplyDefaults() {
	if cfg.IdentityFilePath == "" {
		cfg.IdentityFilePath = filepath.Join(constants.ApplicationDataDir, constants.ApplicationIdentityFileName)
	}

	if cfg.Version == "" {
		cfg.Version = "yoker v0.1-alpha"
	}

	if cfg.ServerName == "" {
		cfg.ServerName = "yoker"
	}

	if cfg.ListenAddr == "" {
		cfg.ListenAddr = ":3000"
	}
}

type Handshake struct {
	Version string
}

func NewServer(cfg ServerConfig) *Server {
	server := &Server{
		ServerConfig: cfg,
		handler:      &DefaultHandler{},
		peers:        make(map[peer.ID]*Peer),
		addPeer:      make(chan libp2pnetwork.Stream),
		deletePeer:   make(chan peer.ID),
		msgCh:        make(chan *Message),
	}

	server.transport = &P2PTransport{
		ListenAddr: server.ListenAddr,
		addPeer:    server.addPeer,

		// internals
		identityFilePath: server.IdentityFilePath,
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
			logrus.Info("addPeer called")
			// validating peer using handshake
			peerId := conn.Conn().RemotePeer().String()
			peer := Peer{conn: conn, Status: constants.ConnectionStateActive}

			// below is not required for libp2p implementation
			// if err := s.handshake(&peer); err != nil {
			// 	logrus.Infof("handshake failed with peer %s due to %s", peerId, err)
			// 	break
			// }

			// handshake was success adding it to the peers
			// starting the read loop
			s.peers[conn.Conn().RemotePeer()] = &peer
			go peer.ReadLoop(s.msgCh, s.deletePeer)

			logrus.WithField(constants.ServerName, s.ServerName).Infof("new player connected %s", peerId)
		case msg := <-s.msgCh:
			s.handler.HandleMessage(msg)
		case peerId := <-s.deletePeer:
			// connection is closed removing from peers
			delete(s.peers, peerId)
			logrus.Infof("player disconnected %s", peerId.String())
		}
	}
}

func (s *Server) Connect(remoteAddr string) error {
	// 1. Parse the multiaddress
	maddr, err := ma.NewMultiaddr(remoteAddr)
	if err != nil {
		logrus.Errorf("invalid multiaddress: %s", err)
		return err
	}

	// 2. Extract peer info from multiaddr (contains ID + address)
	peerInfo, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		logrus.Errorf("invalid peer info from multiaddr: %s", err)
		return err
	}

	// 3. Add peer to peerstore so we can dial it
	if err := s.transport.host.Connect(context.Background(), *peerInfo); err != nil {
		logrus.Errorf("error connecting to peer: %s", err)
		return err
	}

	// 4. Open a stream to the peer using your protocol
	stream, err := s.transport.host.NewStream(context.Background(), peerInfo.ID, "/ypoker/1.0.0")
	if err != nil {
		logrus.Errorf("error opening stream: %s", err)
		return err
	}

	logrus.WithField(constants.ServerName, s.ServerName).Infof("connection request sucess %s at %s adding to peers", peerInfo.ID, remoteAddr)

	s.addPeer <- stream

	return nil
}

func (s *Server) GetMyFullAddr() []string {
	return s.transport.GetMyFullAddr()
}

// ignored for libp2p
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

// ignored for libp2p
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
