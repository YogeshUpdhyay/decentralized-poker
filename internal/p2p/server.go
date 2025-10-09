package p2p

import (
	"context"
	"path/filepath"

	"github.com/YogeshUpdhyay/ypoker/internal/constants"
	libp2pnetwork "github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
	log "github.com/sirupsen/logrus"
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

var server *Server

func GetServer() *Server {
	return server
}

func NewServer(cfg ServerConfig) *Server {
	// apply default values if not set
	cfg.ApplyDefaults()

	server = &Server{
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

func (s *Server) Start(password string) {
	log.Info("starting server")
	go s.loop()
	if err := s.transport.ListenAndAccept(s.ServerName, password); err != nil {
		log.Fatal("error starting the server")
	}
}

func (s *Server) loop() {
	for {
		select {
		case conn := <-s.addPeer:
			log.Info("addPeer called")
			// validating peer using handshake
			peerId := conn.Conn().RemotePeer().String()
			peer := Peer{
				conn:     conn,
				Status:   constants.ConnectionStateActive,
				PeerID:   conn.Conn().RemotePeer().ShortString(),
				Username: "Oliver",
				Avatar:   "https://api.dicebear.com/9.x/adventurer/svg?seed=Oliver",
			}

			// below is not required for libp2p implementation
			// if err := s.handshake(&peer); err != nil {
			// 	log.Infof("handshake failed with peer %s due to %s", peerId, err)
			// 	break
			// }

			// handshake was success adding it to the peers
			// starting the read loop
			s.peers[conn.Conn().RemotePeer()] = &peer
			go peer.ReadLoop(s.msgCh, s.deletePeer)

			log.WithField(constants.ServerName, s.ServerName).Infof("new player connected %s", peerId)
		case msg := <-s.msgCh:
			s.handler.HandleMessage(msg)
		case peerId := <-s.deletePeer:
			// connection is closed removing from peers
			delete(s.peers, peerId)
			log.Infof("player disconnected %s", peerId.String())
		}
	}
}

func (s *Server) Connect(remoteAddr string) (*Peer, error) {
	// 1. Parse the multiaddress
	maddr, err := ma.NewMultiaddr(remoteAddr)
	if err != nil {
		log.Errorf("invalid multiaddress: %s", err)
		return nil, err
	}

	// 2. Extract peer info from multiaddr (contains ID + address)
	peerInfo, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		log.Errorf("invalid peer info from multiaddr: %s", err)
		return nil, err
	}

	// 3. Add peer to peerstore so we can dial it
	if err := s.transport.host.Connect(context.Background(), *peerInfo); err != nil {
		log.Errorf("error connecting to peer: %s", err)
		return nil, err
	}

	// 4. Open a stream to the peer using your protocol
	stream, err := s.transport.host.NewStream(context.Background(), peerInfo.ID, "/ypoker/1.0.0")
	if err != nil {
		log.Errorf("error opening stream: %s", err)
		return nil, err
	}

	log.WithField(constants.ServerName, s.ServerName).Infof("connection request sucess %s at %s adding to peers", peerInfo.ID, remoteAddr)

	s.addPeer <- stream

	return &Peer{
		conn:     stream,
		Status:   constants.ConnectionStateActive,
		PeerID:   stream.Conn().RemotePeer().ShortString(),
		Username: "Oliver",
		Avatar:   "https://api.dicebear.com/9.x/adventurer-neutral/svg?seed=Nolan&radius=50",
	}, nil
}

func (s *Server) GetMyFullAddr() []string {
	return s.transport.GetMyFullAddr()
}

func (s *Server) InitializeIdentityFlow(ctx context.Context, password string) error {
	encryption := DefaultEncryption{}
	idKey, err := encryption.GenerateIdentityKey()
	if err != nil {
		log.WithContext(ctx).Errorf("error generating identity key %s", err.Error())
		return err
	}

	err = encryption.EncryptAndSaveIdentityKey(
		password,
		idKey,
		s.IdentityFilePath,
	)
	if err != nil {
		log.WithContext(ctx).Errorf("error saving identity key %s", err.Error())
		return err
	}

	log.WithContext(ctx).Infof("identity key saved successfully at %s", s.IdentityFilePath)
	return nil
}

func (s *Server) GetPeers() map[peer.ID]*Peer {
	return s.peers
}
