package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/YogeshUpdhyay/ypoker/internal/constants"
	"github.com/YogeshUpdhyay/ypoker/internal/p2p"
	"github.com/sirupsen/logrus"
)

func main() {
	// initialize application utils
	initializeApp()

	peer1Cfg := p2p.ServerConfig{ListenAddr: "3000", Version: "ypoker v0.1-alpha", ServerName: "yoker alpha"}
	peer1 := p2p.NewServer(peer1Cfg)
	go peer1.Start()
	time.Sleep(1 * time.Second)
	select {}
}

func initializeApp() {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// check if the peer id and identity is established or not
	_, err := os.Stat(fmt.Sprintf("%s/%s", constants.ApplicationDataDir, constants.ApplicationIdentityFileName))
	if err != nil && errors.Is(err, os.ErrNotExist) {
		logrus.Infof("identity not initialized")
		err := initIdentityFlow()

		if err != nil {
			logrus.Fatalf("error initializing the identity flow %s", err.Error())
		}

		return
	}
	logrus.Infof("identity already present, skipping initialization")
}

func initIdentityFlow() error {
	password := constants.Empty

	encryption := p2p.DefaultEncryption{}
	idKey, err := encryption.GenerateIdentityKey()
	if err != nil {
		logrus.Errorf("error generating identity key %s", err.Error())
		return err
	}

	err = encryption.EncryptAndSaveIdentityKey(
		password,
		idKey,
		fmt.Sprintf("%s/%s", constants.ApplicationDataDir, constants.ApplicationIdentityFileName),
	)
	if err != nil {
		logrus.Errorf("error encrypting and saving the key %s", err.Error())
		return err
	}
	return nil
}
