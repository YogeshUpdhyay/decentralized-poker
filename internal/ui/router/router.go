package router

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/YogeshUpdhyay/ypoker/internal/ui/pages"
)

type Router struct {
	stack   *fyne.Container
	pages   map[string]pages.Page
	current string
}

func NewRouter() *Router {
	// stack holds all pages on top of each other
	return &Router{
		stack: container.NewStack(), // we’ll add pages later
		pages: map[string]pages.Page{},
	}
}

func (r *Router) Container() *fyne.Container { return r.stack }

func (r *Router) Register(name string, p pages.Page) {
	r.pages[name] = p
	r.stack.Add(p.Content()) // all pages live in the stack
	p.Content().Hide()       // default hidden
}

func (r *Router) Navigate(name string) {
	if r.current == name {
		return
	}
	// hide current
	if cur, ok := r.pages[r.current]; ok {
		cur.Content().Hide()
		cur.OnHide()
	}

	// show next
	if nxt, ok := r.pages[name]; ok {
		nxt.Content().Show()
		nxt.OnShow()
		r.current = name
	}
}
