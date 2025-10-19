package pages

import (
	"context"
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/YogeshUpdhyay/ypoker/internal/constants"
	"github.com/YogeshUpdhyay/ypoker/internal/db"
	"github.com/YogeshUpdhyay/ypoker/internal/eventbus"
	"github.com/YogeshUpdhyay/ypoker/internal/ui/components"
	"github.com/YogeshUpdhyay/ypoker/internal/ui/models"
	"github.com/YogeshUpdhyay/ypoker/internal/utils"
	log "github.com/sirupsen/logrus"
)

type Chat struct {
	chatThreadsData    binding.UntypedList
	pendingRequestData binding.UntypedList
	openChatUsername   binding.String
	openChatAvatarUrl  binding.String
}

func (c *Chat) OnShow(ctx context.Context) {
	// fetch chat threads data from the db
	c.chatThreadsData = binding.NewUntypedList()
	c.chatThreadsData.Set(getChatThreadsList(ctx))

	// fetch pending requests data from the db
	c.pendingRequestData = binding.NewUntypedList()
	c.pendingRequestData.Set(getPendingRequests(ctx))

	c.openChatUsername = binding.NewString()
	c.openChatAvatarUrl = binding.NewString()
}

func (c *Chat) OnHide(_ context.Context) {
	// Logic to execute when the chat page is hidden
}

func (c *Chat) Content(ctx context.Context) fyne.CanvasObject {
	// defining the main chat layout
	secondContent := c.getEmptyChatPlaceholder()
	if username, _ := c.openChatUsername.Get(); username != constants.Empty {
		avatarUrl, _ := c.openChatAvatarUrl.Get()
		secondContent = c.getChatLayout(ctx, username, avatarUrl)
	}

	chat := container.NewHSplit(
		c.getSideNav(ctx),
		secondContent,
	)
	chat.SetOffset(0.25)

	c.openChatUsername.AddListener(binding.NewDataListener(func() {
		log.WithContext(ctx).Info("open chat username changed, refreshing the layout")
		username, _ := c.openChatUsername.Get()
		if username == constants.Empty {
			log.WithContext(ctx).Info("no chat selected, showing empty placeholder")
			return
		}
		avatarUrl, _ := c.openChatAvatarUrl.Get()
		chat.Trailing = c.getChatLayout(ctx, username, avatarUrl)
		chat.Refresh()
	}))

	return chat
}

func (c *Chat) getSideNav(ctx context.Context) fyne.CanvasObject {
	// chat threads list
	chatThreadsList := components.NewChatThreadListWithData(
		c.chatThreadsData,
		c.onChatThreadTap,
	)

	eventbus.Get().Subscribe(constants.EventThreadListUpdated, func(event eventbus.Event) {
		log.WithContext(ctx).Info("thread list updated")
		chatThreads := getChatThreadsList(ctx)
		c.chatThreadsData.Set(chatThreads)
	})

	return container.NewBorder(
		nil, nil, nil, canvas.NewLine(color.White),
		container.NewBorder(
			components.NewSideNavTopBar(ctx),
			components.NewSideNavBottomBar(ctx, c.pendingRequestData),
			nil, nil, chatThreadsList,
		),
	)
}

func (c *Chat) getEmptyChatPlaceholder() fyne.CanvasObject {
	return container.NewCenter(
		widget.NewLabel("Select a chat to start messaging"),
	)
}

func (c *Chat) getChatLayout(ctx context.Context, username, avatarUrl string) fyne.CanvasObject {
	// username block with data binding
	avatarUsername := canvas.NewText(username, color.White)
	avatarUsername.TextSize = 16
	avatarUsername.TextStyle.Bold = true

	// avatar block with data binding
	avatar := components.NewAvatar(avatarUrl)

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

		// sending message to the peer
		// peer := p2p.GetServer().GetPeerByUsername(username)
		// if peer != nil {
		// 	err := peer.Send([]byte(text))
		// 	if err != nil {
		// 		log.WithContext(ctx).WithError(err).Errorf("error sending message to peer %s", username)
		// 	} else {
		// 		log.WithContext(ctx).Infof("message sent to peer %s", username)
		// 	}
		// }

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
					container.NewPadded(avatar),
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

func getChatThreadsList(ctx context.Context) []any {
	peers := []db.PeerInfo{}
	_ = db.Get().Find(&peers)
	log.WithContext(ctx).Infof("found %d peers in the database", len(peers))

	chatThreads := []any{}
	for _, peer := range peers {
		peerData := models.PeerData{
			PeerID:      peer.PeerID,
			Username:    peer.Username,
			Avatar:      peer.AvatarUrl,
			LastMessage: peer.Status,
		}
		chatThreads = append(chatThreads, peerData)
	}

	return chatThreads
}

func getPendingRequests(ctx context.Context) []any {
	pendingRequests := []db.ConnectionRequests{}
	_ = db.Get().Where("status = ?", constants.RequestStatusAwaitingDecision).Find(&pendingRequests)

	pendingReqs := []any{}
	for _, pr := range pendingRequests {
		log.WithContext(ctx).Infof("found pending request from peer id %s", pr.PeerID)
		peerData := models.PeerData{
			PeerID:   pr.PeerID,
			Username: pr.Username,
			Avatar:   pr.AvatarUrl,
		}
		pendingReqs = append(pendingReqs, peerData)
	}

	return pendingReqs
}

func (c *Chat) onChatThreadTap(peerID, username, avatarUrl string) {
	err := c.openChatAvatarUrl.Set(avatarUrl)
	if err != nil {
		log.Infof("error setting open chat avatar url %s", err.Error())
	}

	err = c.openChatUsername.Set(username)
	if err != nil {
		log.Infof("error setting open chat username %s", err.Error())
	}
}
