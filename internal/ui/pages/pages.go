package pages

import "fyne.io/fyne/v2"

type Page interface {
	Content() fyne.CanvasObject
	OnShow() // optional: refresh data when navigated to
	OnHide()
}
