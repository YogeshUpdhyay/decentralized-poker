package p2p

import (
	"github.com/YogeshUpdhyay/ypoker/internal/db"
	"github.com/YogeshUpdhyay/ypoker/internal/utils"
)

type UserInfo struct {
	Name      string
	AvatarUrl string
}

type HandShake struct {
	Version  string
	UserInfo UserInfo `json:"userInfo"`
}

func GetSelfHandshakeMessage() *HandShake {
	userMetadata := db.UserMetadata{}
	db.Get().First(&userMetadata)

	return &HandShake{
		Version: utils.GetAppConfig().Version,
		UserInfo: UserInfo{
			Name:      userMetadata.Username,
			AvatarUrl: userMetadata.AvatarUrl,
		},
	}
}
