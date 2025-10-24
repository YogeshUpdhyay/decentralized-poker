package assets

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

//go:embed yChat.png
var resourceYChatPngData []byte
var ResourceYChatPng = &fyne.StaticResource{
	StaticName:    "yChat.png",
	StaticContent: resourceYChatPngData,
}
