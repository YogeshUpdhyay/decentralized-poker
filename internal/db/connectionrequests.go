package db

import "gorm.io/gorm"

type ConnectionRequests struct {
	gorm.Model
	PeerID    string
	Status    string
	Username  string
	AvatarUrl string
	Address   string
}
