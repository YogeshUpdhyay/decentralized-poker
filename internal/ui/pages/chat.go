package pages

import (
	"context"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Chat struct {
}

func (c *Chat) OnShow(_ context.Context) {
	// Logic to execute when the chat page is shown
}

func (c *Chat) OnHide(_ context.Context) {
	// Logic to execute when the chat page is hidden
}

func (c *Chat) Content(_ context.Context) fyne.CanvasObject {
	// underlay container with a border
	underLayContainer := canvas.NewRectangle(color.Transparent)
	underLayContainer.SetMinSize(fyne.NewSize(350, 220))
	underLayContainer.CornerRadius = 10
	underLayContainer.StrokeColor = color.White
	underLayContainer.StrokeWidth = 2

	// chat content placeholder
	chatContent := widget.NewLabel("Chat content goes here")
	chatContent.Alignment = fyne.TextAlignCenter

	return container.NewBorder(nil, nil, nil, nil, underLayContainer, chatContent)
}
