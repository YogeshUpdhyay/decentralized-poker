package main

import (
	"time"

	"github.com/YogeshUpdhyay/ypoker/internal/p2p"
)

func main() {
	peer1Cfg := p2p.ServerConfig{ListenAddr: ":3000", Version: "ypoker v0.1-alpha"}
	peer1 := p2p.NewServer(peer1Cfg)
	go peer1.Start()
	time.Sleep(1 * time.Second)

	peer2Cfg := p2p.ServerConfig{ListenAddr: ":4000", Version: "ypoker v0.1-alpha"}
	peer2 := p2p.NewServer(peer2Cfg)
	peer2.Connect(":3000")
}
