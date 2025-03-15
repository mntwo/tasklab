package gen_event

import (
	"context"
	"fmt"
	"sync"
)

// Event is an empty interface representing an event.
type Event interface{}

// Handler is an interface that defines methods for handling events.
type Handler interface {
	Init()                              // Init initializes the handler.
	HandleEvent(context.Context, Event) // HandleEvent processes an event.
	Close() error                       // Close cleans up the handler.
}

// EventManager manages event handlers and dispatches events to them.
type EventManager struct {
	handlers map[Handler]struct{} // A set of handlers.
	mu       sync.RWMutex         // A read-write mutex to protect the handlers map.
	eventCh  chan Event           // A channel for incoming events.
	ctx      context.Context      // A context to manage the lifecycle.
	cancel   context.CancelFunc   // A function to cancel the context.
	wg       sync.WaitGroup       // A wait group to wait for all goroutines to finish.
}

// NewEventManager creates a new EventManager with a specified buffer size for the event channel.
func NewEventManager(bufferSize int) *EventManager {
	ctx, cancel := context.WithCancel(context.Background())
	em := &EventManager{
		handlers: make(map[Handler]struct{}),
		eventCh:  make(chan Event, bufferSize),
		ctx:      ctx,
		cancel:   cancel,
	}
	em.wg.Add(1)
	go em.dispatchLoop()
	return em
}

// dispatchLoop listens for events and dispatches them to handlers.
func (em *EventManager) dispatchLoop() {
	defer em.wg.Done()
	for {
		select {
		case e := <-em.eventCh:
			em.broadcast(e)
		case <-em.ctx.Done():
			em.cleanup()
			return
		}
	}
}

// broadcast sends an event to all registered handlers.
func (em *EventManager) broadcast(e Event) {
	em.mu.RLock()
	defer em.mu.RUnlock()

	for handler := range em.handlers {
		em.wg.Add(1)
		go func(h Handler) {
			defer em.wg.Done()
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("handler panic: %v\n", r)
					em.RemoveHandler(h)
				}
			}()
			h.HandleEvent(em.ctx, e)
		}(handler)
	}
}

// cleanup closes all handlers.
func (em *EventManager) cleanup() {
	em.mu.Lock()
	defer em.mu.Unlock()

	for handler := range em.handlers {
		handler.Close()
	}
}

// AddHandler adds a handler to the EventManager.
func (em *EventManager) AddHandler(h Handler) {
	em.mu.Lock()
	defer em.mu.Unlock()
	em.handlers[h] = struct{}{}
	h.Init()
}

// RemoveHandler removes a handler from the EventManager.
func (em *EventManager) RemoveHandler(h Handler) {
	em.mu.Lock()
	defer em.mu.Unlock()
	delete(em.handlers, h)
	h.Close()
}

// Notify sends an event to the event channel.
func (em *EventManager) Notify(e Event) {
	select {
	case em.eventCh <- e:
	case <-em.ctx.Done():
		fmt.Println("EventManager is shutting down, event ignored")
	}
}

// Close shuts down the EventManager and waits for all goroutines to finish.
func (em *EventManager) Close() {
	em.cancel()
	em.wg.Wait()
	close(em.eventCh)
}
