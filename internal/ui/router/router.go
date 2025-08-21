package router

import (
	"context"

	"fyne.io/fyne/v2"
	models "github.com/YogeshUpdhyay/ypoker/internal/ui/commons"
	log "github.com/sirupsen/logrus"
)

type Router struct {
	rootCanvas fyne.Canvas
	pages      map[string]models.Page
	current    string
}

var router *Router

func NewRouter(ctx context.Context, rootCanvas fyne.Canvas) *Router {
	router = &Router{
		pages:      map[string]models.Page{},
		rootCanvas: rootCanvas,
	}
	log.Info("Router initialized")
	return router
}

func GetRouter() *Router {
	return router
}

func (r *Router) Container(ctx context.Context) {
	fyne.Do(func() {
		r.rootCanvas.SetContent(r.pages[r.current].Content(ctx))
	})
}

func (r *Router) Register(_ context.Context, name string, p models.Page) {
	r.pages[name] = p // default hidden
}

func (r *Router) Navigate(ctx context.Context, name string) {
	if r.current == name {
		return
	}
	// hide current
	if cur, ok := r.pages[r.current]; ok {
		cur.OnHide(ctx)
	}

	// show next
	if nxt, ok := r.pages[name]; ok {
		nxt.OnShow(ctx)
		r.current = name
	}

	r.Container(ctx)
}
