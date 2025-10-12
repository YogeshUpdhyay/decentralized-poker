package p2p

import (
	"github.com/YogeshUpdhyay/ypoker/internal/utils"
)

type UserInfo struct {
	Name      string
	AvatarUrl string
}

type HandShake struct {
	Version  string
	UserInfo UserInfo
}

func (hs *HandShake) GetSelfHandshakeMessage() *HandShake {
	return &HandShake{
		Version: utils.GetAppConfig().Version,
		// UserInfo: utils,
	}
}
