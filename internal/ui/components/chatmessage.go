package components

import (
	"context"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/YogeshUpdhyay/ypoker/internal/ui/models"
	"github.com/YogeshUpdhyay/ypoker/internal/utils"
)

type ChatMessage struct {
	widget.BaseWidget
	message       *models.Message
	parentContext context.Context
}

func NewChatMessage(ctx context.Context, message *models.Message) *ChatMessage {
	chatMessage := &ChatMessage{message: message, parentContext: ctx}
	chatMessage.ExtendBaseWidget(chatMessage)
	return chatMessage
}

func messageBoxColor(isSelf bool) color.Color {
	if isSelf {
		return color.RGBA{0x20, 0x90, 0xD0, 0xFF} // blue for self
	}
	return color.RGBA{0x44, 0x44, 0x44, 0xFF} // dark gray for others
}

func (c *ChatMessage) CreateRenderer() fyne.WidgetRenderer {
	messageBox := container.NewHBox()

	if c.message.IsSelf {
		messageBox.Add(layout.NewSpacer())
	}

	textMessage := widget.NewLabel(c.message.Message)
	textMessage.Wrapping = fyne.TextWrapBreak

	messageWidth := canvas.NewRectangle(messageBoxColor(c.message.IsSelf))
	messageWidth.SetMinSize(fyne.NewSize(utils.GetParentSize(c.parentContext).Width*0.6, 40))
	messageWidth.CornerRadius = 10

	message := container.New(
		layout.NewStackLayout(),
		messageWidth,
		container.NewPadded(textMessage),
	)

	messageBox.Add(message)

	if !c.message.IsSelf {
		messageBox.Add(layout.NewSpacer())
	}

	return widget.NewSimpleRenderer(messageBox)
}
