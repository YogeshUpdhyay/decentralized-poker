package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/YogeshUpdhyay/ypoker/internal/constants"
	"github.com/YogeshUpdhyay/ypoker/internal/ui"
	log "github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()

	// initialize application
	initializeApp(ctx)
	uiImpl := ui.DefaultUI{}

	if !identityFileExists() {
		log.WithContext(ctx).Infof("identity not initialized, starting initialization flow")
		uiImpl.StartUI(ctx, false)
		return
	}
	log.WithContext(ctx).Info("identity already present, skipping initialization, starting UI")
	uiImpl.StartUI(ctx, true)

	// // start the peer
	// peer1Cfg := p2p.ServerConfig{ListenAddr: "3000", Version: "ypoker v0.1-alpha", ServerName: "yoker alpha"}
	// peer1 := p2p.NewServer(peer1Cfg)
	// go peer1.Start()

	// time.Sleep(1 * time.Second)
	// select {}

}

// check if the identity file exists
func identityFileExists() bool {
	_, err := os.Stat(fmt.Sprintf("%s/%s", constants.ApplicationDataDir, constants.ApplicationIdentityFileName))
	return !errors.Is(err, os.ErrNotExist)
}

func initializeApp(_ context.Context) {
	log.SetFormatter(&log.JSONFormatter{})
	// creating default config yml file
}
