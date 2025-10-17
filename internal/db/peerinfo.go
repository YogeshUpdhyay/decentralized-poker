package db

import "gorm.io/gorm"

type PeerInfo struct {
	gorm.Model
	PeerID    string `gorm:"primaryKey"`
	Username  string
	AvatarUrl string
	Status    string
	Address   string
}
