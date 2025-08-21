package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/YogeshUpdhyay/ypoker/internal/constants"
	"github.com/YogeshUpdhyay/ypoker/internal/ui"
	"github.com/YogeshUpdhyay/ypoker/internal/utils"
	log "github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()

	// initialize application
	initializeApp(ctx)
	uiImpl := ui.DefaultUI{}

	if !identityFileExists() {
		// writing default config file
		log.WithContext(ctx).Info("first startup, creating thte application config file")
		if err := utils.WriteAppConfig(); err != nil {
			log.WithContext(ctx).WithError(err).Fatal("failed to write default config file")
		}

		log.WithContext(ctx).Infof("identity not initialized, starting initialization flow")
		uiImpl.StartUI(ctx, false)
		return
	}
	log.WithContext(ctx).Info("identity already present, skipping initialization, starting UI")
	uiImpl.StartUI(ctx, true)
}

// check if the identity file exists
func identityFileExists() bool {
	_, err := os.Stat(fmt.Sprintf("%s/%s", constants.ApplicationDataDir, constants.ApplicationIdentityFileName))
	return !errors.Is(err, os.ErrNotExist)
}

func initializeApp(_ context.Context) {
	log.SetFormatter(&log.JSONFormatter{})
}
