package commons

import (
	"context"

	"fyne.io/fyne/v2"
)

type Page interface {
	Content(ctx context.Context) fyne.CanvasObject
	OnShow(ctx context.Context) // optional: refresh data when navigated to
	OnHide(ctx context.Context)
}
