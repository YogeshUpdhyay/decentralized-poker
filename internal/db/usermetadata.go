package db

import (
	"time"
)

type UserMetadata struct {
	PeerID      string `gorm:"uniqueIndex"`
	Username    string `gorm:"primaryKey"`
	AvatarUrl   string
	LastLoginTs time.Time
	CreateTs    time.Time
	UpdateTs    time.Time
}
