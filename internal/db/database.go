package db

import (
	"context"

	"github.com/YogeshUpdhyay/ypoker/internal/utils"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func Initialize(ctx context.Context) error {
	appConfig := utils.GetAppConfig()
	conn, err := gorm.Open(
		sqlite.Open(appConfig.DatabasePath),
		&gorm.Config{},
	)
	if err != nil {
		log.WithContext(ctx).Dup().Infof("error initializing local db: %v", err)
		return err
	}

	// automigrate the database schema
	err = conn.AutoMigrate(
		&UserMetadata{},
		&ConnectionRequests{},
		&PeerInfo{},
		&ChatMessage{},
	)

	db = conn

	return nil
}

func Get() *gorm.DB {
	return db
}
