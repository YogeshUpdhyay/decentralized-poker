package models

import (
	"time"

	"github.com/YogeshUpdhyay/ypoker/internal/constants"
)

type Message struct {
	To        string
	From      string
	Message   string
	CreateTs  time.Time
	Status    string
	DeliverTs time.Time
	IsSelf    bool
}

func GetSelfMessage(text string) *Message {
	return &Message{
		Message:  text,
		IsSelf:   true,
		CreateTs: time.Now(),
		Status:   constants.ToBeSent,
	}
}

func GetPeerMessage(from, text string) *Message {
	return &Message{
		From:     from,
		Message:  text,
		IsSelf:   false,
		CreateTs: time.Now(),
		Status:   constants.Received,
	}
}
