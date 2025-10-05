package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/YogeshUpdhyay/ypoker/internal/ui/models"
	log "github.com/sirupsen/logrus"
)

type ChatThreadsList struct {
	widget.BaseWidget
	ThreadsData  binding.UntypedList
	OnItemTapped func(peerID, username, avatarUrl string)
}

func NewChatThreadListWithData(data binding.UntypedList, onItemTapped func(peerID, username, avatarUrl string)) *ChatThreadsList {
	ctl := &ChatThreadsList{ThreadsData: data, OnItemTapped: onItemTapped}
	ctl.ExtendBaseWidget(ctl)
	data.AddListener(binding.NewDataListener(ctl.Refresh))
	return ctl
}

func (c *ChatThreadsList) CreateRenderer() fyne.WidgetRenderer {
	content := container.NewVBox()

	threads, err := c.ThreadsData.Get()
	if err != nil {
		log.Errorf("error getting chat threads data: %s", err.Error())
		return widget.NewSimpleRenderer(content)
	}

	for _, t := range threads {
		peerData, ok := t.(models.PeerData)
		if !ok {
			log.Error("error asserting peer data")
			continue
		}

		content.Add(NewChatThread(
			peerData.Username,
			peerData.Avatar,
			peerData.LastMessage,
			c.OnItemTapped,
		))
	}

	return widget.NewSimpleRenderer(content)
}
