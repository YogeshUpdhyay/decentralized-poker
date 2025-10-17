package db

import "gorm.io/gorm"

type ConnectionRequests struct {
	gorm.Model
	PeerID    string `gorm:"primaryKey"`
	Status    string // eg.sent pending accepted rejected
	Username  string
	AvatarUrl string
	Address   string
}
