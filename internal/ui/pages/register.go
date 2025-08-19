package pages

import (
	"context"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/YogeshUpdhyay/ypoker/internal/constants"
	"github.com/YogeshUpdhyay/ypoker/internal/p2p"
	"github.com/YogeshUpdhyay/ypoker/internal/ui/router"
	log "github.com/sirupsen/logrus"
)

type Register struct {
}

func (l *Register) OnShow(_ context.Context) {
	// Logic to execute when the login page is shown
}

func (l *Register) OnHide(_ context.Context) {
	// Logic to execute when the login page is hidden
}

func (l *Register) Content(ctx context.Context) fyne.CanvasObject {
	// underlay container with a border
	underLayContainer := canvas.NewRectangle(color.Transparent)
	underLayContainer.SetMinSize(fyne.NewSize(350, 220))
	underLayContainer.CornerRadius = 10
	underLayContainer.StrokeColor = color.White
	underLayContainer.StrokeWidth = 2

	// logo
	logo := canvas.NewImageFromFile(constants.LogoPath)
	logo.FillMode = canvas.ImageFillContain
	logo.SetMinSize(fyne.NewSize(100, 100))

	// form elements and data bindings
	username := widget.NewEntry()
	username.SetPlaceHolder("Username")
	userNameString := binding.NewString()
	username.Bind(userNameString)

	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("Password")
	passwordString := binding.NewString()
	password.Bind(passwordString)

	confirmPassword := widget.NewPasswordEntry()
	confirmPassword.SetPlaceHolder("Confirm Password")
	confirmPasswordString := binding.NewString()
	confirmPassword.Bind(confirmPasswordString)

	submit := widget.NewButton("Submit", func() {
		_, err := userNameString.Get()
		if err != nil {
			log.Errorf("failed to get username: %v", err)
			return
		}

		passwordValue, err := passwordString.Get()
		if err != nil {
			log.Errorf("failed to get password: %v", err)
			return
		}

		confirmPasswordValue, err := confirmPasswordString.Get()
		if err != nil {
			log.Errorf("failed to get confirm password: %v", err)
			return
		}

		if passwordValue != confirmPasswordValue {
			log.Errorf("password and confirm password do not match")
			return
		}

		peer1Cfg := p2p.ServerConfig{ListenAddr: "3000", Version: "ypoker v0.1-alpha", ServerName: "yoker alpha"}
		peer1 := p2p.NewServer(peer1Cfg)
		peer1.InitializeIdentityFlow(ctx, passwordValue)

		router.GetRouter().Navigate(constants.ChatRoute)
	})
	submit.Alignment = widget.ButtonAlignCenter

	// vbox size
	vBoxsize := canvas.NewRectangle(color.Transparent)
	vBoxsize.SetMinSize(fyne.NewSize(300, 125))

	register := container.NewCenter(
		container.NewStack(
			underLayContainer,
			container.NewCenter(
				container.NewPadded(
					container.NewVBox(
						vBoxsize,
						logo,
						username,
						password,
						confirmPassword,
						submit,
						vBoxsize,
					),
				),
			),
		),
	)

	return register
}
