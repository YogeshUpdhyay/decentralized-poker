package p2p

import (
	"github.com/YogeshUpdhyay/ypoker/internal/database"
	"github.com/YogeshUpdhyay/ypoker/internal/utils"
	log "github.com/sirupsen/logrus"
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
	metadata := database.UserMetadata{}
	if err := metadata.GetFirst(); err != nil {
		log.Errorf("error getting user metadata: %v", err)
		return nil
	}

	return &HandShake{
		Version: utils.GetAppConfig().Version,
		UserInfo: UserInfo{
			Name:      metadata.Username,
			AvatarUrl: metadata.AvatarUrl,
		},
	}
}
