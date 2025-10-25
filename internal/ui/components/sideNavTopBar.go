package components

import (
	"context"
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/YogeshUpdhyay/ypoker/internal/p2p"
	"github.com/YogeshUpdhyay/ypoker/internal/ui/assets"
	"github.com/YogeshUpdhyay/ypoker/internal/utils"
	log "github.com/sirupsen/logrus"
)

func NewSideNavTopBar(ctx context.Context) *fyne.Container {
	// logo
	logo := canvas.NewImageFromResource(assets.ResourceYChatPng)
	logo.FillMode = canvas.ImageFillContain
	logo.SetMinSize(fyne.NewSize(40, 40))

	// account info button
	accountInfoIcon := NewIconButton(theme.AccountIcon(), func() {
		log.WithContext(ctx).Info("account info button tapped")

		// fetch address and generate qr code
		server := p2p.GetServer()
		myAddresses := server.GetMyFullAddr()
		qrPath, err := utils.GetAddressQRS(myAddresses)
		if err != nil {
			log.WithContext(ctx).Infof("error creating address qrs %s", err.Error())
		}

		// display a share link or qr code in the dailog
		getConnectionPopUp(ctx, qrPath, myAddresses)
	})

	return container.NewPadded(
		container.NewHBox(
			logo,
			layout.NewSpacer(),
			accountInfoIcon,
		),
	)
}

func getConnectionPopUp(ctx context.Context, qrPath string, myAddresses []string) {
	log.WithContext(ctx).Info("generating connection pop up")

	// modal title
	modalTitle := canvas.NewText("Share your address", color.White)
	modalTitle.TextSize = 16
	modalTitle.TextStyle.Bold = true
	modalTitleContainer := container.NewPadded(container.NewCenter(modalTitle))

	//connection addresses
	connectionEntry := widget.NewMultiLineEntry()
	connectionEntry.SetText(strings.Join(myAddresses, "\n"))
	connectionEntry.Disable()
	connectionAddresses := container.NewPadded(connectionEntry)

	qrImage := canvas.NewImageFromFile(qrPath)
	qrImage.FillMode = canvas.ImageFillContain
	qrImage.SetMinSize(fyne.NewSize(200, 200))
	qrImageContainer := container.NewPadded(container.NewCenter(qrImage))

	widthContainer := canvas.NewRectangle(color.Transparent)
	widthContainer.SetMinSize(fyne.NewSize(450, 1))

	// Create the content for the pop-up
	content := container.NewVBox(
		widthContainer,
		modalTitleContainer,
		connectionAddresses,
		qrImageContainer,
	)

	// Create the pop-up
	popUp := widget.NewModalPopUp(
		content,
		fyne.CurrentApp().Driver().AllWindows()[0].Canvas(),
	)

	content.Add(widget.NewButton("Close", func() { popUp.Hide() }))

	popUp.Show()
}
