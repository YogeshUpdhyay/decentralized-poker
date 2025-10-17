package p2p

import (
	"context"
	"encoding/json"
	"errors"
	"path/filepath"

	"github.com/YogeshUpdhyay/ypoker/internal/constants"
	"github.com/YogeshUpdhyay/ypoker/internal/db"
	"github.com/YogeshUpdhyay/ypoker/internal/utils"
	libp2pnetwork "github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	ma "github.com/multiformats/go-multiaddr"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	ServerConfig
	handler       Handler
	transport     *P2PTransport
	peers         map[string]*Peer
	addPeer       chan libp2pnetwork.Stream
	deletePeer    chan peer.ID
	msgCh         chan *Message
	awaitingPeers chan *Peer
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
		ServerConfig:  cfg,
		handler:       &DefaultHandler{},
		peers:         make(map[string]*Peer),
		addPeer:       make(chan libp2pnetwork.Stream),
		deletePeer:    make(chan peer.ID),
		msgCh:         make(chan *Message),
		awaitingPeers: make(chan *Peer),
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
			log.Info("adding new peer connection")

			// validating peer using handshake connection remains in the pending state until handshake is complete
			peerId := conn.Conn().RemotePeer().String()
			peer := Peer{
				conn:   conn,
				Status: constants.ConnectionStateActive,
				PeerID: string(conn.Conn().RemotePeer().ShortString()),
			}

			if err := s.handshake(&peer); err != nil {
				log.Infof("handshake failed with peer %s due to %s", peerId, err)
				break
			}

			// handshake was success adding it to the peers
			// starting the read loop
			// s.peers[conn.Conn().RemotePeer().String()] = &peer
			// go peer.ReadLoop(s.msgCh, s.deletePeer)

			log.WithField(constants.ServerName, s.ServerName).Infof("new player connected %s", peerId)
		case msg := <-s.msgCh:
			s.handler.HandleMessage(msg)
		case peerId := <-s.deletePeer:
			// connection is closed removing from peers
			delete(s.peers, peerId.String())
			log.Infof("player disconnected %s", peerId.String())
		}
	}
}

func (s *Server) Connect(ctx context.Context, remoteAddr string) (*Peer, error) {
	// parse the multiaddress
	maddr, err := ma.NewMultiaddr(remoteAddr)
	if err != nil {
		log.Errorf("invalid multiaddress: %s", err)
		return nil, err
	}

	// extract peer info from multiaddr (contains ID + address)
	peerInfo, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		log.Errorf("invalid peer info from multiaddr: %s", err)
		return nil, err
	}

	// add peer to peerstore so we can dial it
	if err := s.transport.host.Connect(ctx, *peerInfo); err != nil {
		log.Errorf("error connecting to peer: %s", err)
		return nil, err
	}

	// open a stream to the peer using your protocol
	appConfig := utils.GetAppConfig()
	stream, err := s.transport.host.NewStream(ctx, peerInfo.ID, protocol.ID(appConfig.StreamProtocol))
	if err != nil {
		log.Errorf("error opening stream: %s", err)
		return nil, err
	}

	peer := &Peer{
		conn:   stream,
		Status: constants.ConnectionStateActive,
		PeerID: stream.Conn().RemotePeer().ShortString(),
	}

	// send handshake message
	if err := s.sendHandshake(peer); err != nil {
		log.Errorf("error sending handshake: %s", err)
		return nil, err
	}
	log.WithContext(ctx).Infof("connection request sucess %s at %s adding to peers", peerInfo.ID, remoteAddr)

	// adding this pending connection to db
	connectionRequest := db.ConnectionRequests{
		PeerID: stream.Conn().RemotePeer().ShortString(),
		Status: constants.RequestStatusSent,
	}
	db.Get().Create(&connectionRequest)

	return peer, nil
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

func (s *Server) GetPeers() map[string]*Peer {
	return s.peers
}

func (s *Server) GetPeerByUsername(username string) *Peer {
	return nil
}

func (s *Server) handshake(p *Peer) error {
	peerHandshake, err := p.ReadHandshakeMessage()
	if err != nil {
		return err
	}

	if peerHandshake.Version != s.Version {
		return errors.New("version mismatch")
	}

	log.Infof("handshake received from peer %s with version %s", p.PeerID, peerHandshake.Version)

	// add as pending connection awaiting approval
	peerInfo := db.ConnectionRequests{
		PeerID:    p.PeerID,
		Username:  peerHandshake.UserInfo.Name,
		AvatarUrl: peerHandshake.UserInfo.AvatarUrl,
		Status:    constants.RequestStatusAwaitingDecision,
	}
	db.Get().Create(&peerInfo)

	return nil
}

func (s *Server) sendHandshake(p *Peer) error {
	hs := GetSelfHandshakeMessage()
	hsData, err := json.Marshal(hs)
	if err != nil {
		return err
	}

	if err := p.Send(hsData); err != nil {
		return err
	}

	return nil
}
