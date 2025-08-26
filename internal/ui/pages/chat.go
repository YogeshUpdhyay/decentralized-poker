package pages

import (
	"context"
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/YogeshUpdhyay/ypoker/internal/constants"
	"github.com/YogeshUpdhyay/ypoker/internal/ui/components"
	"github.com/YogeshUpdhyay/ypoker/internal/ui/models"
	log "github.com/sirupsen/logrus"
)

type Chat struct {
}

func (c *Chat) OnShow(_ context.Context) {
	// Logic to execute when the chat page is shown
}

func (c *Chat) OnHide(_ context.Context) {
	// Logic to execute when the chat page is hidden
}

func (c *Chat) Content(ctx context.Context) fyne.CanvasObject {
	chat := container.NewHSplit(
		getSideNav(ctx),
		getChatLayout(ctx),
	)
	chat.SetOffset(0.25)

	return chat
}

func getSideNav(ctx context.Context) fyne.CanvasObject {
	// logo
	logo := canvas.NewImageFromFile(constants.LogoPath)
	logo.FillMode = canvas.ImageFillContain
	logo.SetMinSize(fyne.NewSize(40, 40))

	// account info button
	accountInfoIcon := components.NewIconButton(theme.AccountIcon(), func() {
		log.WithContext(ctx).Info("account info button tapped")
		// server := p2p.GetServer()
		// myAddresses := server.GetMyFullAddr()
		// _, err := utils.GetAddressQRS(myAddresses)
		// if err != nil {
		// 	log.WithContext(ctx).Infof("error creating address qrs %s", err.Error())
		// }
	})

	return container.NewBorder(
		nil, nil, nil, canvas.NewLine(color.White),
		container.NewVBox(
			container.NewPadded(
				container.NewHBox(
					logo,
					layout.NewSpacer(),
					accountInfoIcon,
				),
			),
			layout.NewSpacer(),
			container.NewPadded(
				widget.NewButtonWithIcon(
					"New Chat",
					theme.ContentAddIcon(),
					func() {
						log.WithContext(ctx).Info("New chat button tapped")
					},
				),
			),
		),
	)
}

func getChatLayout(ctx context.Context) fyne.CanvasObject {
	avatarUsername := canvas.NewText("MSo17", color.White)
	avatarUsername.TextSize = 16
	avatarUsername.TextStyle.Bold = true

	textBox := canvas.NewRectangle(color.Transparent)
	textBox.StrokeColor = color.White
	textBox.StrokeWidth = 2
	textBox.CornerRadius = 16

	// messages model
	// msgs := binding.NewStringList()
	messages := binding.NewUntypedList()
	list := container.NewVBox()

	entry := widget.NewMultiLineEntry()
	entry.SetPlaceHolder("Type a message…")
	send := widget.NewButtonWithIcon("", theme.MailSendIcon(), func() {
		text := strings.TrimSpace(entry.Text)
		if text == "" {
			return
		}

		newMessage := models.GetSelfMessage(text)
		_ = messages.Append(newMessage)
		chatMessage := components.NewChatMessage(newMessage)
		list.Add(chatMessage)

		entry.SetText("")
		entry.FocusLost()
	})
	entry.OnSubmitted = func(s string) { send.OnTapped() }

	return container.NewBorder(
		container.NewPadded(
			container.NewBorder(
				nil, canvas.NewLine(color.White), nil, nil,
				container.NewHBox(
					container.NewPadded(components.NewAvatar("https://avatar.iran.liara.run/public/boy?username=Ash")),
					container.NewPadded(
						container.NewVBox(
							avatarUsername,
							canvas.NewText("online", color.White),
						),
					),
				),
			),
		),
		container.NewPadded(
			container.NewBorder(
				canvas.NewLine(color.White),
				nil, nil,
				send,
				entry,
			),
		),
		nil, nil,
		list,
	)
}
