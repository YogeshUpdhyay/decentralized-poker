package pages

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Login struct{}

func (l *Login) OnShow() {
	// Logic to execute when the login page is shown
}

func (l *Login) OnHide() {
	// Logic to execute when the login page is hidden
}

func (l *Login) Content() fyne.CanvasObject {
	// form elements
	username := widget.NewEntry()
	username.SetPlaceHolder("Username")

	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("Password")

	submit := widget.NewButton("Submit", func() {
		// handle login
	})

	// VBox with form elements
	formBox := container.NewVBox(
		username,
		password,
		submit,
	)

	// add black border around VBox
	border := canvas.NewRectangle(color.White)
	border.SetMinSize(formBox.MinSize().Add(fyne.NewSize(20, 20))) // padding for border

	// overlay the border behind the formBox
	boxWithBorder := container.NewCenter(
		container.NewStack(
			border,
			container.NewPadded(formBox),
		),
	)

	return container.NewCenter(boxWithBorder)
}
