package components

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type ChatThread struct {
	widget.BaseWidget
	Username    string
	AvatarURL   string
	LastMessage string
	OnTapped    func()
}

func NewChatThread(username, avatarURL, lastMessage string) fyne.CanvasObject {
	// Build and return a container with avatar, username, and last message
	chatThread := &ChatThread{
		Username:    username,
		AvatarURL:   avatarURL,
		LastMessage: lastMessage,
	}
	chatThread.ExtendBaseWidget(chatThread)
	return chatThread
}

func (c *ChatThread) CreateRenderer() fyne.WidgetRenderer {
	// Username label
	avatarUsername := widget.NewLabel(c.Username)
	avatarUsername.TextStyle.Bold = true
	avatarUsername.Truncation = fyne.TextTruncateEllipsis
	// avatarUsername.SizeName = theme.SizeNameText

	// Last message label with truncation
	lastMessage := widget.NewLabel(c.LastMessage)
	lastMessage.Truncation = fyne.TextTruncateEllipsis // nicer than Clip

	spacer := canvas.NewRectangle(color.Transparent)
	spacer.SetMinSize(fyne.NewSize(100, 0))

	// Layout: avatar fixed + right side expands
	// right := container.NewBorder(nil, nil, nil, nil)

	content := container.NewBorder(
		nil, canvas.NewLine(color.White), nil, nil,
		container.NewHBox(
			container.NewPadded(NewAvatar(c.AvatarURL)),
			container.NewPadded(container.NewVBox(
				spacer, // spacer
				avatarUsername,
				lastMessage,
			)), // this expands automatically
		),
	)

	return widget.NewSimpleRenderer(content)
}
