package pages

import (
	"context"
	"fmt"
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
	eventModels "github.com/YogeshUpdhyay/ypoker/internal/eventbus/models"
	"github.com/YogeshUpdhyay/ypoker/internal/p2p"
	p2pModels "github.com/YogeshUpdhyay/ypoker/internal/p2p/models"
	"github.com/YogeshUpdhyay/ypoker/internal/ui/components"
	"github.com/YogeshUpdhyay/ypoker/internal/ui/models"
	"github.com/YogeshUpdhyay/ypoker/internal/utils"
	log "github.com/sirupsen/logrus"
)

type Chat struct {
	chatThreadsData    binding.UntypedList
	pendingRequestData binding.UntypedList
	currentChatData    binding.Untyped
}

func (c *Chat) OnShow(ctx context.Context) {
	// fetch chat threads data from the db
	c.chatThreadsData = binding.NewUntypedList()
	c.chatThreadsData.Set(getChatThreadsList(ctx))

	// fetch pending requests data from the db
	c.pendingRequestData = binding.NewUntypedList()
	c.pendingRequestData.Set(getPendingRequests(ctx))

	c.currentChatData = binding.NewUntyped()
}

func (c *Chat) OnHide(_ context.Context) {
	// Logic to execute when the chat page is hidden
}

func (c *Chat) Content(ctx context.Context) fyne.CanvasObject {
	// defining the main chat layout
	secondContent := c.getEmptyChatPlaceholder()
	currentChatData, _ := c.currentChatData.Get()
	if currentChatData != nil {
		currentChatPeerData := currentChatData.(models.PeerData)
		secondContent = c.getChatLayout(ctx, currentChatPeerData)
	}

	chat := container.NewHSplit(
		c.getSideNav(ctx),
		secondContent,
	)
	chat.SetOffset(0.25)

	c.currentChatData.AddListener(binding.NewDataListener(func() {
		log.WithContext(ctx).Info("open chat data changed, refreshing the layout")
		currentChatData, err := c.currentChatData.Get()
		if err != nil {
			log.WithContext(ctx).Infof("error getting current chat data")
		}

		if currentChatData == nil {
			return
		}

		currentChatPeerData := currentChatData.(models.PeerData)
		if currentChatPeerData.Username == constants.Empty {
			log.WithContext(ctx).Info("no chat selected, showing empty placeholder")
			return
		}
		chat.Trailing = c.getChatLayout(ctx, currentChatPeerData)
		chat.Refresh()
	}))

	return chat
}

func (c *Chat) getSideNav(ctx context.Context) fyne.CanvasObject {
	// chat threads list
	chatThreadsList := components.NewChatThreadListWithData(
		ctx,
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

func (c *Chat) getChatLayout(ctx context.Context, currentChatPeerData models.PeerData) fyne.CanvasObject {
	// username block with data binding
	avatarUsername := canvas.NewText(currentChatPeerData.Username, color.White)
	avatarUsername.TextSize = 16
	avatarUsername.TextStyle.Bold = true

	// avatar block with data binding
	avatar := components.NewAvatar(currentChatPeerData.Avatar)

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

		newMessage := models.GetSelfMessage(text)
		_ = messages.Append(newMessage)
		chatMessage := components.NewChatMessage(ctx, newMessage)
		list.Add(chatMessage)

		// sending message to the peer
		peer := p2p.GetServer().GetPeerFromPeerID(ctx, currentChatPeerData.PeerID)
		if peer != nil {
			msg := p2pModels.Envelope{
				Type: p2pModels.MsgTypeChat,
				Payload: utils.MarshalPayload(p2pModels.ChatPayload{
					Message: text,
				}),
			}
			err := peer.Send(&msg)
			if err != nil {
				log.WithContext(ctx).WithError(err).Errorf("error sending message to peer %s", currentChatPeerData.Username)
			} else {
				log.WithContext(ctx).Infof("message sent to peer %s", currentChatPeerData.Username)
			}
		}

		// adding message to db
		tx := db.Get().Create(&db.ChatMessage{
			From:    p2p.GetServer().GetNodeID(ctx),
			To:      currentChatPeerData.PeerID,
			Message: text,
		})
		if tx.Error != nil {
			log.WithContext(ctx).Infof("error writing chat to db: %s", tx.Error.Error())
		}

		entry.SetText("")
		entry.FocusLost()
	})
	entry.OnSubmitted = func(s string) { send.OnTapped() }

	eventbus.Get().Subscribe(constants.EventNewMessage, func(event eventbus.Event) {
		messageData := event.Payload.(eventModels.NewMessageEventData)
		if messageData.From != currentChatPeerData.PeerID {
			log.WithContext(ctx).Infof("current peer id is different")
			return
		}

		ctx = utils.UpdateParentInContext(ctx, list)

		newMessage := models.GetPeerMessage(messageData.From, messageData.Message)
		_ = messages.Append(newMessage)
		chatMessage := components.NewChatMessage(ctx, newMessage)
		list.Add(chatMessage)
		log.WithContext(ctx).Infof("message added to list")
	})

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

func (c *Chat) onChatThreadTap(ctx context.Context, username, avatarUrl, peerID string) {
	// check if the peer id in the active peers
	server := p2p.GetServer()
	peer := server.GetPeerFromPeerID(ctx, peerID)
	if peer == nil {
		// attempt connection
		peerInfo := db.PeerInfo{PeerID: peerID}
		tx := db.Get().First(&peerInfo)
		if tx.Error != nil {
			log.WithContext(ctx).Infof("error fetching the peer info to establish connection %s", tx.Error.Error())
			if err := c.currentChatData.Set(models.PeerData{
				Username:    username,
				PeerID:      peerID,
				Avatar:      avatarUrl,
				LastMessage: "...",
			}); err != nil {
				log.WithError(err).Info("error setting open chat data")
			}
			return
		}

		// build connection string
		connectStr := fmt.Sprintf("%s/p2p/%s", peerInfo.Address, peerInfo.PeerID)
		server.Connect(ctx, connectStr)
	}
	if err := c.currentChatData.Set(models.PeerData{
		Username:    username,
		PeerID:      peerID,
		Avatar:      avatarUrl,
		LastMessage: "...",
	}); err != nil {
		log.WithError(err).Info("error setting open chat data")
	}

}
