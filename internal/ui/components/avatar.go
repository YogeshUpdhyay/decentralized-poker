package components

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	log "github.com/sirupsen/logrus"
)

type Avatar struct {
	widget.BaseWidget
	URL string
}

func NewAvatar(url string) *Avatar {
	avatar := &Avatar{URL: url}
	avatar.ExtendBaseWidget(avatar)
	return avatar
}

func (a *Avatar) CreateRenderer() fyne.WidgetRenderer {
	imageResource, err := fyne.LoadResourceFromURLString(a.URL)
	if err != nil {
		log.Infof("error loading resource %s", err.Error())
		imageResource = theme.FyneLogo()
	}
	image := canvas.NewImageFromResource(imageResource)
	image.SetMinSize(fyne.NewSize(50, 50))

	backgroundCircle := canvas.NewCircle(color.White)
	backgroundCircle.StrokeColor = color.Gray{Y: 0x99}
	backgroundCircle.StrokeWidth = 2

	return widget.NewSimpleRenderer(
		container.NewStack(
			backgroundCircle,
			image,
		),
	)
}
