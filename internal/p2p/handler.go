package p2p

import (
	"io"

	log "github.com/sirupsen/logrus"
)

type Handler interface {
	HandleMessage(*Message) error
}

type DefaultHandler struct{}

func (h *DefaultHandler) HandleMessage(msg *Message) error {
	b, err := io.ReadAll(msg.Payload)
	if err != nil {
		return err
	}
	log.Infof("handling message from %s: %s", msg.From, string(b))
	return nil
}
