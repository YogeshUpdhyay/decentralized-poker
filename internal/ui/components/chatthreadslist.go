package components

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/YogeshUpdhyay/ypoker/internal/ui/models"
	log "github.com/sirupsen/logrus"
)

type ChatThreadsList struct {
	widget.BaseWidget
	ctx          context.Context
	ThreadsData  binding.UntypedList
	OnItemTapped func(ctx context.Context, peerID, username, avatarUrl string)
}

func NewChatThreadListWithData(ctx context.Context, data binding.UntypedList, onItemTapped func(ctx context.Context, peerID, username, avatarUrl string)) *ChatThreadsList {
	ctl := &ChatThreadsList{ThreadsData: data, OnItemTapped: onItemTapped, ctx: ctx}
	ctl.ExtendBaseWidget(ctl)
	data.AddListener(binding.NewDataListener(func() {
		log.Info("chat threads data changed, refreshing the list")
		ctl.Refresh()
	}))
	return ctl
}

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
					r.list.ctx,
					peerData.Username,
					peerData.Avatar,
					peerData.LastMessage,
					peerData.PeerID,
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
