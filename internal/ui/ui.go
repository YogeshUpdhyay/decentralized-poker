package ui

import (
	"context"

	"fyne.io/fyne/v2"
	fyneApp "fyne.io/fyne/v2/app"
	"github.com/YogeshUpdhyay/ypoker/internal/constants"
	"github.com/YogeshUpdhyay/ypoker/internal/ui/pages"
	"github.com/YogeshUpdhyay/ypoker/internal/ui/router"
	log "github.com/sirupsen/logrus"
)

type UI interface {
	StartUI() error
}

type DefaultUI struct{}

// initialize the UI and start the application
func (ui *DefaultUI) StartUI(ctx context.Context) error {
	app := fyneApp.New()
	window := app.NewWindow("yoker")
	window.Resize(fyne.NewSize(constants.WindowWidth, constants.WindowHeight))

	// get router
	router := router.NewRouter()
	log.WithContext(ctx).Info("router initialized")

	// registering pages to router
	router.Register(constants.LoginRoute, &pages.Login{})
	router.Navigate(constants.LoginRoute)

	window.SetContent(router.Container())
	window.ShowAndRun()

	return nil
}
