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
	data.AddListener(binding.NewDataListener(func() {
		log.Info("chat threads data changed, refreshing the list")
		ctl.Refresh()
	}))
	return ctl
}

// func (c *ChatThreadsList) CreateRenderer() fyne.WidgetRenderer {
// 	log.Info("creating chat threads list renderer")
// 	content := container.NewVBox()
// 	threads, err := c.ThreadsData.Get()
// 	if err != nil {
// 		log.Errorf("error getting chat threads data: %s", err.Error())
// 		return widget.NewSimpleRenderer(content)
// 	}

// 	for _, t := range threads {
// 		peerData, ok := t.(models.PeerData)
// 		if !ok {
// 			log.Error("error asserting peer data")
// 			continue
// 		}

// 		content.Add(NewChatThread(
// 			peerData.Username,
// 			peerData.Avatar,
// 			peerData.LastMessage,
// 			c.OnItemTapped,
// 		))
// 	}

// 	chatThreadListRenderer := widget.NewSimpleRenderer(content)
// 	return chatThreadListRenderer
// }

type chatThreadsListRenderer struct {
	list    *ChatThreadsList
	content *fyne.Container
}

func (r *chatThreadsListRenderer) Layout(size fyne.Size) {
	r.content.Resize(size)
}

func (r *chatThreadsListRenderer) MinSize() fyne.Size {
	return r.content.MinSize()
}

func (r *chatThreadsListRenderer) Refresh() {
	r.content.Objects = nil
	threads, err := r.list.ThreadsData.Get()
	if err == nil {
		for _, t := range threads {
			peerData, ok := t.(models.PeerData)
			if ok {
				r.content.Add(NewChatThread(
					peerData.Username,
					peerData.Avatar,
					peerData.LastMessage,
					r.list.OnItemTapped,
				))
			}
		}
	}
	r.content.Refresh()
}

func (r *chatThreadsListRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.content}
}

func (r *chatThreadsListRenderer) Destroy() {}

func (c *ChatThreadsList) CreateRenderer() fyne.WidgetRenderer {
	content := container.NewVBox()
	return &chatThreadsListRenderer{list: c, content: content}
}
