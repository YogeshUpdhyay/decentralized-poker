package ui

import (
	"context"
	"time"

	"fyne.io/fyne/v2"
	fyneApp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"

	"github.com/YogeshUpdhyay/ypoker/internal/constants"
	"github.com/YogeshUpdhyay/ypoker/internal/p2p"
	"github.com/YogeshUpdhyay/ypoker/internal/ui/pages"
	"github.com/YogeshUpdhyay/ypoker/internal/ui/router"
	"github.com/YogeshUpdhyay/ypoker/internal/utils"
)

type UI interface {
	StartUI() error
}

type DefaultUI struct{}

// initialize the UI and start the application
func (ui *DefaultUI) StartUI(ctx context.Context, isIdentityInitialized bool) error {
	app := fyneApp.New()

	app.Settings().SetTheme(theme.DarkTheme())
	window := app.NewWindow("yoker gamma")
	window.Resize(fyne.NewSize(constants.WindowWidth, constants.WindowHeight))
	rootCanvas := window.Canvas()

	// get router
	router := router.NewRouter(ctx, rootCanvas)

	// starting the server testing block
	appConfig := utils.GetAppConfig()
	serverConfig := p2p.ServerConfig{ListenAddr: appConfig.Port, Version: appConfig.Version, ServerName: appConfig.Name}
	server := p2p.NewServer(serverConfig)
	go server.Start("oggy@123")
	time.Sleep(2 * time.Second)

	// registering pages to router
	router.Register(ctx, constants.LoginRoute, &pages.Login{})
	router.Register(ctx, constants.ChatRoute, &pages.Chat{})
	router.Register(ctx, constants.RegisterRoute, &pages.Register{})
	router.Navigate(ctx, constants.RegisterRoute)
	if isIdentityInitialized {
		router.Navigate(ctx, constants.ChatRoute)
	}

	// content is set by the router
	window.ShowAndRun()

	return nil
}
