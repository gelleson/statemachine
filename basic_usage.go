package statemachine

import (
	"context"
	"fmt"
)

// Define Events

// StartEvent initiates the process.
type StartEvent struct{}

func (e *StartEvent) EventType() string {
	return "StartEvent"
}

// StopEvent terminates the process.
type StopEvent struct{}

func (e *StopEvent) EventType() string {
	return "StopEvent"
}

// Define Handlers

// StartHandler handles StartEvent.
type StartHandler struct{}

func (h *StartHandler) Execute(ctx context.Context, event *StartEvent) error {
	fmt.Println("Process started.")
	return nil
}

// StopHandler handles StopEvent.
type StopHandler struct{}

func (h *StopHandler) Execute(ctx context.Context, event *StopEvent) error {
	fmt.Println("Process stopped.")
	return nil
}

// ExampleBasicUsage demonstrates a simple state machine setup and event execution.
func ExampleBasicUsage() {
	// Initialize the state machine with the initial state "Idle".
	sm := NewStateMachine("Idle")

	// Register handlers for StartEvent and StopEvent.
	RegisterHandler[*StartEvent, *StartHandler](sm, &StartHandler{})
	RegisterHandler[*StopEvent, *StopHandler](sm, &StopHandler{})

	// Define transitions:
	// From "Idle" to "Running" on StartEvent.
	// From "Running" to "Idle" on StopEvent.
	sm.From("Idle").On(&StartEvent{}).To("Running")
	sm.From("Running").On(&StopEvent{}).To("Idle")

	// Execute StartEvent to transition from "Idle" to "Running".
	if err := Execute(context.Background(), sm, &StartEvent{}); err != nil {
		fmt.Println("Error:", err)
	}

	// Execute StopEvent to transition back to "Idle".
	if err := Execute(context.Background(), sm, &StopEvent{}); err != nil {
		fmt.Println("Error:", err)
	}

	// Output:
	// Process started.
	// Process stopped.
}
