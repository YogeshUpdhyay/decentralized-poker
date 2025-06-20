package main

import (
	"time"

	"github.com/YogeshUpdhyay/ypoker/internal/p2p"
	"github.com/sirupsen/logrus"
)

func main() {
	// initialize application utils
	initializeApp()

	peer1Cfg := p2p.ServerConfig{ListenAddr: ":3000", Version: "ypoker v0.1-alpha", ServerName: "yoker alpha"}
	peer1 := p2p.NewServer(peer1Cfg)
	go peer1.Start()
	time.Sleep(1 * time.Second)

	peer2Cfg := p2p.ServerConfig{ListenAddr: ":4000", Version: "ypoker v0.1-alpha", ServerName: "yoker beta"}
	peer2 := p2p.NewServer(peer2Cfg)
	go peer2.Start()
	time.Sleep(1 * time.Second)

	err := peer2.Connect(":3000")
	if err != nil {
		logrus.Info("peer not connected")
	}
	select {}
}

func initializeApp() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
}
