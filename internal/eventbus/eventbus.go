package eventbus

import (
	"sync"

	log "github.com/sirupsen/logrus"
)

// Event represents a generic event with a name and any payload.
type Event struct {
	Name    string
	Payload interface{}
}

// Handler is a function that processes an event.
type Handler func(Event)

// EventBus is a simple in-memory pub/sub bus.
type EventBus struct {
	subscribers map[string][]Handler
	mu          sync.RWMutex
}

var eventBus *EventBus

// New creates a new EventBus instance.
func New() *EventBus {
	eventBus = &EventBus{
		subscribers: make(map[string][]Handler),
	}
	return eventBus
}

// Get the event bus client
func Get() *EventBus {
	if eventBus == nil {
		log.Info("event bus not initialized")
		eventBus = &EventBus{subscribers: make(map[string][]Handler)}
	}
	return eventBus
}

// Subscribe adds a new handler for a specific event name.
func (b *EventBus) Subscribe(eventName string, handler Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subscribers[eventName] = append(b.subscribers[eventName], handler)
}

// Publish sends an event to all its subscribers asynchronously.
func (b *EventBus) Publish(eventName string, payload interface{}) {
	b.mu.RLock()
	handlers, ok := b.subscribers[eventName]
	b.mu.RUnlock()
	if !ok {
		return
	}

	event := Event{Name: eventName, Payload: payload}
	for _, handler := range handlers {
		// Dispatch asynchronously
		go handler(event)
	}
}
