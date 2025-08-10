package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/YogeshUpdhyay/ypoker/internal/constants"
	"github.com/YogeshUpdhyay/ypoker/internal/p2p"
	"github.com/YogeshUpdhyay/ypoker/internal/ui"
	log "github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()

	// initialize application utils
	initializeApp(ctx)

	// // start the peer
	// peer1Cfg := p2p.ServerConfig{ListenAddr: "3000", Version: "ypoker v0.1-alpha", ServerName: "yoker alpha"}
	// peer1 := p2p.NewServer(peer1Cfg)
	// go peer1.Start()

	// time.Sleep(1 * time.Second)
	// select {}

	uiImpl := ui.DefaultUI{}
	uiImpl.StartUI(ctx)
}

func initializeApp(ctx context.Context) {
	log.SetFormatter(&log.JSONFormatter{})

	// check if the peer id and identity is established or not
	_, err := os.Stat(fmt.Sprintf("%s/%s", constants.ApplicationDataDir, constants.ApplicationIdentityFileName))
	if err != nil && errors.Is(err, os.ErrNotExist) {
		log.WithContext(ctx).Infof("identity not initialized")
		err := initIdentityFlow(ctx)

		if err != nil {
			log.WithContext(ctx).Fatalf("error initializing the identity flow %s", err.Error())
		}

		return
	}
	log.WithContext(ctx).Infof("identity already present, skipping initialization")
}

func initIdentityFlow(ctx context.Context) error {
	password := constants.Empty

	encryption := p2p.DefaultEncryption{}
	idKey, err := encryption.GenerateIdentityKey()
	if err != nil {
		log.WithContext(ctx).Errorf("error generating identity key %s", err.Error())
		return err
	}

	err = encryption.EncryptAndSaveIdentityKey(
		password,
		idKey,
		fmt.Sprintf("%s/%s", constants.ApplicationDataDir, constants.ApplicationIdentityFileName),
	)
	if err != nil {
		log.WithContext(ctx).Errorf("error encrypting and saving the key %s", err.Error())
		return err
	}
	return nil
}
