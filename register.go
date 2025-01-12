package statemachine

// RegisterHandler registers a handler for a specific event type
func RegisterHandler[T Event, H Handler[T]](sm *StateMachine, handler H) {
	var event T
	eventType := event.EventType()

	sm.mu.Lock()
	sm.handlers[eventType] = handler
	sm.mu.Unlock()
}
