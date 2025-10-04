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
	"github.com/YogeshUpdhyay/ypoker/internal/p2p"
	"github.com/YogeshUpdhyay/ypoker/internal/ui/components"
	"github.com/YogeshUpdhyay/ypoker/internal/ui/models"
	"github.com/YogeshUpdhyay/ypoker/internal/utils"
	log "github.com/sirupsen/logrus"
)

type Chat struct {
}

func (c *Chat) OnShow(_ context.Context) {
	// Logic to execute when the chat page is shown
	// get the chat data stored locally or in the distributed storage over the peer-to-peer network
	log.Info("Chat page is now visible")
}

func (c *Chat) OnHide(_ context.Context) {
	// Logic to execute when the chat page is hidden
}

func (c *Chat) Content(ctx context.Context) fyne.CanvasObject {
	chat := container.NewHSplit(
		c.getSideNav(ctx),
		c.getChatLayout(ctx, nil),
	)
	chat.SetOffset(0.25)

	return chat
}

func (c *Chat) getSideNav(ctx context.Context) fyne.CanvasObject {
	// logo
	logo := canvas.NewImageFromFile(constants.LogoPath)
	logo.FillMode = canvas.ImageFillContain
	logo.SetMinSize(fyne.NewSize(40, 40))

	// account info button
	accountInfoIcon := components.NewIconButton(theme.AccountIcon(), func() {
		log.WithContext(ctx).Info("account info button tapped")

		// fetch address and generate qr code
		server := p2p.GetServer()
		myAddresses := server.GetMyFullAddr()
		qrPath, err := utils.GetAddressQRS(myAddresses)
		if err != nil {
			log.WithContext(ctx).Infof("error creating address qrs %s", err.Error())
		}

		// display a share link or qr code in the dailog
		c.getConnectionPopUp(ctx, qrPath, myAddresses)
	})

	// new chat button
	newChatButton := widget.NewButtonWithIcon(
		"New Chat",
		theme.ContentAddIcon(),
		func() {
			c.getAddNewChatPopUp(ctx)
		},
	)

	// chat threads list
	chatThreads := container.NewVBox()

	return container.NewBorder(
		nil, nil, nil, canvas.NewLine(color.White),
		container.NewBorder(
			container.NewPadded(
				container.NewHBox(
					logo,
					layout.NewSpacer(),
					accountInfoIcon,
				),
			),
			container.NewPadded(
				newChatButton,
			), nil, nil, chatThreads,
		),
		// chatThreads,
	)
}

func (c *Chat) getChatLayout(ctx context.Context, peerData *models.PeerData) fyne.CanvasObject {
	// if no peer data is provided, show a placeholder
	if peerData == nil {
		return container.NewCenter(
			widget.NewLabel("Select a chat to start messaging"),
		)
	}

	avatarUsername := canvas.NewText(peerData.Username, color.White)
	avatarUsername.TextSize = 16
	avatarUsername.TextStyle.Bold = true

	// Create the username text for the chat header
	textBox := canvas.NewRectangle(color.Transparent)
	textBox.StrokeColor = color.White
	textBox.StrokeWidth = 2
	textBox.CornerRadius = 16

	// Create the messages model (data source for the chat list)
	messages := binding.NewUntypedList()
	list := container.NewVBox()

	// Create the entry box for typing new messages
	entry := widget.NewEntry()
	entry.SetPlaceHolder("Type a message…")
	send := widget.NewButtonWithIcon("", theme.MailSendIcon(), func() {
		text := strings.TrimSpace(entry.Text)
		if text == "" {
			return
		}

		ctx = utils.UpdateParentInContext(ctx, list)

		// testing  others message
		// newPeerMessage := models.GetPeerMessage("MSo17", text)
		// _ = messages.Append(newPeerMessage)
		// chatMessagePeer := components.NewChatMessage(ctx, newPeerMessage)
		// list.Add(chatMessagePeer)

		newMessage := models.GetSelfMessage(text)
		_ = messages.Append(newMessage)
		chatMessage := components.NewChatMessage(ctx, newMessage)
		list.Add(chatMessage)

		entry.SetText("")
		entry.FocusLost()
	})
	entry.OnSubmitted = func(s string) { send.OnTapped() }

	// Layout the chat UI using a border container
	// Top: chat header with avatar and username
	// Bottom: message entry and send button
	// Center: chat messages list
	return container.NewBorder(
		container.NewPadded(
			container.NewBorder(
				nil, canvas.NewLine(color.White), nil, nil,
				container.NewHBox(
					container.NewPadded(components.NewAvatar(peerData.Avatar)),
					container.NewPadded(
						container.NewVBox(
							avatarUsername,
							canvas.NewText("online", color.White), // get status func to be implemented
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

func (c *Chat) getConnectionPopUp(ctx context.Context, qrPath string, myAddresses []string) {
	log.WithContext(ctx).Info("generating connection pop up")

	// modal title
	modalTitle := canvas.NewText("Share your address", color.White)
	modalTitle.TextSize = 16
	modalTitle.TextStyle.Bold = true
	modalTitleContainer := container.NewPadded(container.NewCenter(modalTitle))

	//connection addresses
	connectionAddresses := container.NewPadded(widget.NewLabel(
		strings.Join(myAddresses, "\n"),
	))

	qrImage := canvas.NewImageFromFile(qrPath)
	qrImage.FillMode = canvas.ImageFillContain
	qrImage.SetMinSize(fyne.NewSize(200, 200))
	qrImageContainer := container.NewPadded(container.NewCenter(qrImage))

	// Create the content for the pop-up
	content := container.NewVBox(
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

func (c *Chat) getAddNewChatPopUp(ctx context.Context) {
	log.WithContext(ctx).Info("generating add new chat pop up")

	// modal title
	modalTitle := canvas.NewText("Start a new chat", color.White)
	modalTitle.TextSize = 16
	modalTitle.TextStyle.Bold = true
	modalTitleContainer := container.NewPadded(container.NewCenter(modalTitle))

	// input field for peer ID or address
	input := widget.NewEntry()
	input.SetPlaceHolder("Enter Peer ID or Address")
	inputFieldDimesion := canvas.NewRectangle(color.Transparent)
	inputFieldDimesion.SetMinSize(fyne.NewSize(300, 0))

	// cta row
	ctaRow := container.NewGridWithColumns(2)

	// underlay container with a border
	underLayContainer := canvas.NewRectangle(color.Transparent)
	underLayContainer.SetMinSize(fyne.NewSize(350, 220))

	// spacer rect
	spacerRect := canvas.NewRectangle(color.Transparent)
	spacerRect.SetMinSize(fyne.NewSize(0, 40))

	// create the content for the pop-up
	content := container.NewPadded(
		underLayContainer,
		container.NewCenter(container.NewVBox(
			modalTitleContainer,
			spacerRect,
			container.NewPadded(input, inputFieldDimesion),
			spacerRect,
			ctaRow,
		)),
	)

	// Create the pop-up
	popUp := widget.NewModalPopUp(
		content,
		fyne.CurrentApp().Driver().AllWindows()[0].Canvas(),
	)

	cancelBtn := widget.NewButton("Cancel", func() { popUp.Hide() })
	cancelBtn.Importance = widget.MediumImportance
	ctaRow.Add(cancelBtn)

	connectBtn := widget.NewButton("Connect", func() {
		address := strings.TrimSpace(input.Text)
		if address != "" {
			log.WithContext(ctx).Infof("Starting chat with peer ID: %s", address)
			// logic to initiate chat with the given peer ID goes here
			err := c.connectToPeer(ctx, address)
			if err != nil {
				log.WithContext(ctx).WithError(err).Errorf("error connecting to peer %s", address)
			}
			popUp.Hide()
		}
	})
	connectBtn.Importance = widget.HighImportance
	ctaRow.Add(connectBtn)

	popUp.Show()
}

func (c *Chat) connectToPeer(ctx context.Context, address string) error {
	server := p2p.GetServer()
	if server == nil {
		return nil
	}

	err := server.Connect(address)
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("error connecting to peer %s", address)
		return err
	}

	log.WithContext(ctx).Infof("successfully connected to peer %s", address)
	return nil
}
