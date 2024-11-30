package events

import (
	"sync"
)

// EventsHandler defines a function that handles an event
type EventsHandler func(event Event)

// EventsManager manages event subscriptions and dispatching
type EventsManager struct {
	subscribers map[string][]EventsHandler
	subMutex    sync.RWMutex
}

// NewEventsManager initializes and returns a new EventManager
func NewEventsManager() *EventsManager {
	return &EventsManager{
		subscribers: make(map[string][]EventsHandler),
	}
}

// Subscribe registers a handler for a specific event type
func (em *EventsManager) Subscribe(eventType string, handler EventsHandler) {
	em.subMutex.Lock()
	defer em.subMutex.Unlock()

	// If the event type doesn't exist in the map, create a new slice
	if _, exists := em.subscribers[eventType]; !exists {
		em.subscribers[eventType] = []EventsHandler{}
	}

	// Append the handler to the slice
	em.subscribers[eventType] = append(em.subscribers[eventType], handler)

}

// Dispatch sends an event to all registered handlers.
func (em *EventsManager) Dispatch(eventType string, event Event) {
	em.subMutex.RLock()
	defer em.subMutex.RUnlock()

	// Get the handlers for the event type
	handlers, handlersExists := em.subscribers[eventType]
	if !handlersExists {
		return
	}

	// Call each handler
	for _, handler := range handlers {
		handler(event)
	}
}
