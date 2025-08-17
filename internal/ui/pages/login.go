package pages

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/YogeshUpdhyay/ypoker/internal/constants"
)

type Login struct {
}

func (l *Login) OnShow() {
	// Logic to execute when the login page is shown
}

func (l *Login) OnHide() {
	// Logic to execute when the login page is hidden
}

func (l *Login) Content() fyne.CanvasObject {
	// underlay container with a border
	underLayContainer := canvas.NewRectangle(color.Transparent)
	underLayContainer.SetMinSize(fyne.NewSize(320, 220))
	underLayContainer.CornerRadius = 10
	underLayContainer.StrokeColor = color.White
	underLayContainer.StrokeWidth = 2

	// logo
	logo := canvas.NewImageFromFile(constants.LogoPath)
	logo.FillMode = canvas.ImageFillOriginal
	logo.Resize(fyne.NewSize(100, 100))

	// form elements
	username := widget.NewEntry()
	username.SetPlaceHolder("Username")

	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("Password")

	submit := widget.NewButton("Submit", func() {})
	submit.Resize(fyne.NewSize(300*0.75, 40))

	// vbox size
	vBoxsize := canvas.NewRectangle(color.Transparent)
	vBoxsize.SetMinSize(fyne.NewSize(300, 200))

	login := container.NewCenter(
		container.NewStack(
			underLayContainer,
			container.NewCenter(
				container.NewPadded(
					container.NewVBox(
						logo,
						username,
						password,
						container.NewWithoutLayout(submit),
						vBoxsize,
					),
				),
			),
		),
	)

	return container.NewCenter(login)
}
