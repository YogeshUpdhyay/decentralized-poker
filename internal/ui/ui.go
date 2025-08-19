package ui

import (
	"context"

	"fyne.io/fyne/v2"
	fyneApp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"

	"github.com/YogeshUpdhyay/ypoker/internal/constants"
	"github.com/YogeshUpdhyay/ypoker/internal/ui/pages"
	"github.com/YogeshUpdhyay/ypoker/internal/ui/router"
)

type UI interface {
	StartUI() error
}

type DefaultUI struct{}

// initialize the UI and start the application
func (ui *DefaultUI) StartUI(ctx context.Context) error {
	app := fyneApp.New()
	app.Settings().SetTheme(theme.DarkTheme())
	window := app.NewWindow("yoker")
	window.Resize(fyne.NewSize(constants.WindowWidth, constants.WindowHeight))
	rootCanvas := window.Canvas()

	// get router
	router := router.NewRouter(rootCanvas)

	// registering pages to router
	router.Register(constants.LoginRoute, &pages.Login{})
	router.Register(constants.ChatRoute, &pages.Chat{})
	router.Navigate(constants.LoginRoute)

	// content is set by the router
	window.ShowAndRun()

	return nil
}
