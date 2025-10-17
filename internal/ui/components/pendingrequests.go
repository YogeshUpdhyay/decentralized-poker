package components

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/YogeshUpdhyay/ypoker/internal/ui/models"
)

type PendingRequestsCTA struct {
	widget.BaseWidget
	PendingRequests binding.UntypedList
	OnAccept        func(peerID string)
	OnReject        func(peerID string)
}

func NewPendingRequestsCTA(pendingRequests binding.UntypedList, onAccept func(peerID string), onReject func(peerID string)) *PendingRequestsCTA {
	prc := &PendingRequestsCTA{PendingRequests: pendingRequests, OnAccept: onAccept, OnReject: onReject}
	prc.ExtendBaseWidget(prc)
	pendingRequests.AddListener(binding.NewDataListener(func() {
		prc.Refresh()
	}))
	return prc
}

func (p *PendingRequestsCTA) CreateRenderer() fyne.WidgetRenderer {
	content := container.NewStack()
	renderer := &pendingRequestsCTARenderer{
		pendingRequestsCTA: p,
		container:          content,
		currentIndex:       0,
	}
	return renderer
}

type pendingRequestsCTARenderer struct {
	pendingRequestsCTA *PendingRequestsCTA
	acceptButton       *widget.Button
	rejectButton       *widget.Button
	container          *fyne.Container
	currentIndex       int
}

func (r *pendingRequestsCTARenderer) Layout(size fyne.Size) {
	r.container.Resize(size)
}

func (r *pendingRequestsCTARenderer) MinSize() fyne.Size {
	return r.container.MinSize()
}

func (r *pendingRequestsCTARenderer) Refresh() {
	r.container.Objects = nil
	requests, err := r.pendingRequestsCTA.PendingRequests.Get()
	if err == nil {
		if len(requests) > 0 {
			pendingRequest, ok := requests[r.currentIndex].(models.PeerData)
			if ok {
				r.container.RemoveAll()
				r.container.Add(
					container.NewHBox(
						container.NewPadded(NewIconButton(theme.NavigateBackIcon(), func() { r.currentIndex = (r.currentIndex - 1 + len(requests)) % len(requests); r.Refresh() })),
						container.NewPadded(NewAvatar(pendingRequest.Avatar)),
						container.NewPadded(container.NewVBox(
							canvas.NewText(pendingRequest.Username, color.White),
							canvas.NewText(fmt.Sprintf("#%s", pendingRequest.PeerID), color.White),
						)),
						layout.NewSpacer(),
						container.NewPadded(NewIconButton(theme.ConfirmIcon(), func() { r.pendingRequestsCTA.OnAccept(pendingRequest.PeerID) })),
						container.NewPadded(NewIconButton(theme.CancelIcon(), func() { r.pendingRequestsCTA.OnReject(pendingRequest.PeerID) })),
						container.NewPadded(NewIconButton(theme.NavigateNextIcon(), func() { r.currentIndex = (r.currentIndex + 1 + len(requests)) % len(requests); r.Refresh() })),
					),
				)
			}
		}
	}
	r.container.Refresh()
}

func (r *pendingRequestsCTARenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.container}
}

func (r *pendingRequestsCTARenderer) Destroy() {}
