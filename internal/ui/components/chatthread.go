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
	// truncate last message if too long
	truncatedMessage := c.truncateWithEllipsis(c.LastMessage, 30)
	subTitleText := canvas.NewText(truncatedMessage, color.Gray{Y: 200})

	content := container.NewBorder(
		nil, canvas.NewLine(color.White), nil, nil,
		container.NewHBox(
			container.NewPadded(NewAvatar(c.AvatarURL)),
			container.NewPadded(
				container.NewCenter(
					container.NewVBox(
						canvas.NewText(c.Username, color.White),
						subTitleText,
					),
				),
			),
		),
	)

	return widget.NewSimpleRenderer(content)
}

func (c *ChatThread) truncateWithEllipsis(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}
