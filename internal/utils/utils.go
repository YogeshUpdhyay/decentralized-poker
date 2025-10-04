package utils

import (
	"context"

	"fyne.io/fyne/v2"
)

func UpdateParentInContext(ctx context.Context, parent interface{}) context.Context {
	return context.WithValue(ctx, "parent", parent)
}

func GetParentSize(ctx context.Context) fyne.Size {
	parentSize := fyne.CurrentApp().Driver().CanvasForObject(ctx.Value("parent").(fyne.CanvasObject)).Size()
	return parentSize
}
