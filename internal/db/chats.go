package db

import "gorm.io/gorm"

type ChatMessage struct {
	gorm.Model
	From    string `gorm:"index"`
	To      string `gorm:"index"`
	Message string
}
