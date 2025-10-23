package components

import (
	"context"
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
	PeerID      string
	Ctx         context.Context
	OnTapped    func(ctx context.Context, username, avatarUrl, peerID string)
}

func NewChatThread(ctx context.Context, username, avatarURL, lastMessage, peerID string, onTap func(ctx context.Context, username, avatarUrl, peerID string)) fyne.CanvasObject {
	// Build and return a container with avatar, username, and last message
	chatThread := &ChatThread{
		Ctx:         ctx,
		Username:    username,
		AvatarURL:   avatarURL,
		LastMessage: lastMessage,
		PeerID:      peerID,
		OnTapped:    onTap,
	}
	chatThread.ExtendBaseWidget(chatThread)
	return chatThread
}

func (c *ChatThread) CreateRenderer() fyne.WidgetRenderer {
	// truncate last message if too long
	truncatedMessage := c.truncateWithEllipsis(c.LastMessage, 30)
	subTitleText := canvas.NewText(truncatedMessage, color.Gray{Y: 200})

	usernameText := canvas.NewText(c.Username, color.White)
	usernameText.TextStyle = fyne.TextStyle{Bold: true}
	content := container.NewBorder(
		nil, canvas.NewLine(color.White), nil, nil,
		container.NewHBox(
			container.NewPadded(NewAvatar(c.AvatarURL)),
			container.NewPadded(
				container.NewCenter(
					container.NewVBox(
						usernameText,
						subTitleText,
					),
				),
			),
		),
	)

	return widget.NewSimpleRenderer(content)
}

func (c *ChatThread) Tapped(_ *fyne.PointEvent) {
	if c.OnTapped != nil {
		c.OnTapped(c.Ctx, c.Username, c.AvatarURL, c.PeerID)
	}
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
