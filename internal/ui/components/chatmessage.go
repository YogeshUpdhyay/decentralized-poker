package components

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/YogeshUpdhyay/ypoker/internal/ui/models"
)

type ChatMessage struct {
	widget.BaseWidget
	message *models.Message
}

func NewChatMessage(message *models.Message) *ChatMessage {
	chatMessage := &ChatMessage{message: message}
	chatMessage.ExtendBaseWidget(chatMessage)
	return chatMessage
}

func messageBoxColor(isSelf bool) color.Color {
	if isSelf {
		return color.RGBA{0x20, 0x90, 0xD0, 0xFF} // blue for self
	}
	return color.RGBA{0xE0, 0xE0, 0xE0, 0xFF} // light gray for others
}

func (c *ChatMessage) CreateRenderer() fyne.WidgetRenderer {
	messageBox := container.NewHBox(layout.NewSpacer(), container.NewPadded(
		container.NewStack(
			canvas.NewRectangle(messageBoxColor(c.message.IsSelf)),
			widget.NewLabel(c.message.Message),
		),
	))

	return widget.NewSimpleRenderer(messageBox)
}
