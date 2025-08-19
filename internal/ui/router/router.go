package router

import (
	"fyne.io/fyne/v2"
	"github.com/YogeshUpdhyay/ypoker/internal/ui/models"
)

type Router struct {
	rootCanvas fyne.Canvas
	pages      map[string]models.Page
	current    string
}

var router *Router

func NewRouter(rootCanvas fyne.Canvas) *Router {
	router = &Router{
		pages:      map[string]models.Page{},
		rootCanvas: rootCanvas,
	}
	return router
}

func GetRouter() *Router {
	return router
}

func (r *Router) Container() {
	fyne.Do(func() {
		r.rootCanvas.SetContent(r.pages[r.current].Content())
	})
}

func (r *Router) Register(name string, p models.Page) {
	r.pages[name] = p // default hidden
}

func (r *Router) Navigate(name string) {
	if r.current == name {
		return
	}
	// hide current
	if cur, ok := r.pages[r.current]; ok {
		cur.OnHide()
	}

	// show next
	if nxt, ok := r.pages[name]; ok {
		nxt.OnShow()
		r.current = name
	}

	r.Container()
}
