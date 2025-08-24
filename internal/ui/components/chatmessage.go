package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/YogeshUpdhyay/ypoker/internal/ui/models"
)

type ChatMessage struct {
	widget.BaseWidget
	Message binding.String
}

func NewChatMessage() *ChatMessage {
	chatMessage := &ChatMessage{
		Message: binding.NewString(),
	}
	chatMessage.ExtendBaseWidget(chatMessage)
	return chatMessage
}

func (c *ChatMessage) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(
		widget.NewLabelWithData(c.Message),
	)
}

func (c *ChatMessage) SetData(data binding.Untyped) {
	unTypedMessage, _ := data.Get()
	message := unTypedMessage.(*models.Message)

	c.Message.Set(message.Message)

	c.Refresh()
}
