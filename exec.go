package statemachine

import (
	"context"
	"fmt"
)

// Execute processes an event with its registered handler
func Execute[T Event](ctx context.Context, sm *StateMachine, event T) error {
	sm.mu.Lock()
	transition, err := sm.findValidTransition(ctx, event)
	if err != nil {
		sm.mu.Unlock()
		return err
	}

	nextState := transition.To
	currentState := sm.state
	sm.mu.Unlock()

	// Get handler
	sm.mu.RLock()
	handler, exists := sm.handlers[event.EventType()]
	sm.mu.RUnlock()

	if !exists {
		return fmt.Errorf("no handler registered for event type: %s", event.EventType())
	}

	// Execute handler
	typedHandler, ok := handler.(Handler[T])
	if !ok {
		return fmt.Errorf("invalid handler type for event: %s", event.EventType())
	}

	for _, m := range sm.mw.PreTransitionMiddlewares {
		if err := m(ctx, event); err != nil {
			return err
		}
	}
	if err := typedHandler.Execute(ctx, event); err != nil {
		return fmt.Errorf("handler execution failed: %w", err)
	}

	for _, m := range sm.mw.PostTransitionMiddlewares {
		if err := m(ctx, event); err != nil {
			return err
		}
	}

	// Update state
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.state != currentState {
		return fmt.Errorf("concurrent state modification detected")
	}

	sm.state = nextState
	return nil
}
