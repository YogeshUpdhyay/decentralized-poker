package main

import "github.com/YogeshUpdhyay/ypoker/internal/p2p"

func main() {
	cfg := p2p.ServerConfig{ListenAddr: ":3000", Version: "ypoker v0.1-alpha\n"}
	server := p2p.NewServer(cfg)
	server.Start()
}
