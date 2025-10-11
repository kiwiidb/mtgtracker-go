package events

import (
	"log"
	"sync"
)

// Event is the interface that all events must implement
type Event interface {
	EventName() string
}

// Handler is a function that processes an event
type Handler func(event Event) error

// EventBus manages event subscriptions and publishing
type EventBus struct {
	handlers map[string][]Handler
	mu       sync.RWMutex
}

// NewEventBus creates a new event bus instance
func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[string][]Handler),
	}
}

// Subscribe registers a handler for a specific event name
func (eb *EventBus) Subscribe(eventName string, handler Handler) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	eb.handlers[eventName] = append(eb.handlers[eventName], handler)
	log.Printf("Subscribed handler for event: %s", eventName)
}

// Publish sends an event to all registered handlers asynchronously
// Errors from handlers are logged but do not block other handlers
func (eb *EventBus) Publish(event Event) {
	eb.mu.RLock()
	handlers := eb.handlers[event.EventName()]
	eb.mu.RUnlock()

	if len(handlers) == 0 {
		log.Printf("No handlers registered for event: %s", event.EventName())
		return
	}

	// Execute handlers asynchronously
	for _, handler := range handlers {
		go func(h Handler) {
			if err := h(event); err != nil {
				log.Printf("Event handler error for %s: %v", event.EventName(), err)
			}
		}(handler)
	}
}

// PublishSync sends an event to all registered handlers synchronously
// Returns the first error encountered, but all handlers will be called
func (eb *EventBus) PublishSync(event Event) error {
	eb.mu.RLock()
	handlers := eb.handlers[event.EventName()]
	eb.mu.RUnlock()

	if len(handlers) == 0 {
		log.Printf("No handlers registered for event: %s", event.EventName())
		return nil
	}

	var firstError error
	for _, handler := range handlers {
		if err := handler(event); err != nil {
			log.Printf("Event handler error for %s: %v", event.EventName(), err)
			if firstError == nil {
				firstError = err
			}
		}
	}

	return firstError
}

// Unsubscribe removes all handlers for a specific event name
func (eb *EventBus) Unsubscribe(eventName string) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	delete(eb.handlers, eventName)
	log.Printf("Unsubscribed all handlers for event: %s", eventName)
}

// HandlerCount returns the number of handlers registered for an event
func (eb *EventBus) HandlerCount(eventName string) int {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	return len(eb.handlers[eventName])
}
