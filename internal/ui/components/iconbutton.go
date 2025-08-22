package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type IconButton struct {
	widget.Icon
	OnTapped func()
}

func NewIconButton(res fyne.Resource, onTapped func()) *IconButton {
	icon := &IconButton{}
	icon.ExtendBaseWidget(icon)
	icon.SetResource(res)
	icon.OnTapped = onTapped

	return icon
}

func (t *IconButton) Tapped(_ *fyne.PointEvent) {
	t.OnTapped()
}
