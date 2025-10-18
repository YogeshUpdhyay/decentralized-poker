package models

type UserInfo struct {
	Username  string `json:"username"`
	AvatarUrl string `json:"avatarUrl"`
}

type HandShake struct {
	Version  string   `json:"version"`
	UserInfo UserInfo `json:"userInfo"`
}
