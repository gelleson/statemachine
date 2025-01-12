// File: state_test.go
package statemachine

import (
	"context"
	"fmt"
	"strings"
	"testing"
)

// ----------------------------
// Define Events and Handlers
// ----------------------------

// EventA is a sample event with two integer fields.
type EventA struct {
	A int
	B int
}

func (e *EventA) EventType() string {
	return "EventA"
}

// EventB is another sample event with a single integer field.
type EventB struct {
	X int
}

func (e *EventB) EventType() string {
	return "EventB"
}

// StartEvent is used to initiate a transition.
type StartEvent struct{}

func (e *StartEvent) EventType() string {
	return "StartEvent"
}

// StopEvent is used to terminate a transition.
type StopEvent struct{}

func (e *StopEvent) EventType() string {
	return "StopEvent"
}

// UnlockEvent is used with a guard condition.
type UnlockEvent struct {
	Code int
}

func (e *UnlockEvent) EventType() string {
	return "UnlockEvent"
}

// UndefinedEvent has no registered handler to simulate orphan transitions.
type UndefinedEvent struct{}

func (e *UndefinedEvent) EventType() string {
	return "UndefinedEvent"
}

// NoTransitionEvent has a handler but no transition defined from the current state.
type NoTransitionEvent struct{}

func (e *NoTransitionEvent) EventType() string {
	return "NoTransitionEvent"
}

// ComplexEvent is used to test multiple guard conditions.
type ComplexEvent struct{}

func (e *ComplexEvent) EventType() string {
	return "ComplexEvent"
}

// MultiHandlerEvent is used to test multiple handlers for the same event type.
type MultiHandlerEvent struct{}

func (e *MultiHandlerEvent) EventType() string {
	return "MultiHandlerEvent"
}

// FaultyEvent is used to test handler execution failure.
type FaultyEvent struct{}

func (e *FaultyEvent) EventType() string {
	return "FaultyEvent"
}

// PingEvent is used to test transitions to the same state.
type PingEvent struct{}

func (e *PingEvent) EventType() string {
	return "PingEvent"
}

// CorrectEvent is used to test valid handler registrations.
type CorrectEvent struct{}

func (e *CorrectEvent) EventType() string {
	return "CorrectEvent"
}

// EventWithHandler is used to test orphan transitions.
type EventWithHandler struct{}

func (e *EventWithHandler) EventType() string {
	return "EventWithHandler"
}

// DefinedEvent is used in error handling tests.
type DefinedEvent struct{}

func (e *DefinedEvent) EventType() string {
	return "DefinedEvent"
}

// OrderEvent is used to test handler execution order.
type OrderEvent struct{}

func (e *OrderEvent) EventType() string {
	return "OrderEvent"
}

// IncrementEvent is used in concurrency tests.
type IncrementEvent struct {
	Value int
}

func (e *IncrementEvent) EventType() string {
	return "IncrementEvent"
}

// ResetEvent is used to reset the state machine.
type ResetEvent struct{}

func (e *ResetEvent) EventType() string {
	return "ResetEvent"
}

// Event1 and Event2 are used in DOT generation tests.
type Event1 struct{}

func (e *Event1) EventType() string {
	return "Event1"
}

type Event2 struct{}

func (e *Event2) EventType() string {
	return "Event2"
}

// ----------------------------
// Define Handlers
// ----------------------------

// HandlerEventA handles EventA.
type HandlerEventA struct{}

func (h *HandlerEventA) Execute(ctx context.Context, event *EventA) error {
	// Handle EventA
	return nil
}

// HandlerEventB handles EventB.
type HandlerEventB struct{}

func (h *HandlerEventB) Execute(ctx context.Context, event *EventB) error {
	// Handle EventB
	fmt.Printf("Handling EventB: X=%d\n", event.X)
	return nil
}

// StartHandler handles StartEvent.
type StartHandler struct{}

func (h *StartHandler) Execute(ctx context.Context, event *StartEvent) error {
	// Handle StartEvent
	return nil
}

// StopHandler handles StopEvent.
type StopHandler struct{}

func (h *StopHandler) Execute(ctx context.Context, event *StopEvent) error {
	// Handle StopEvent
	return nil
}

// UnlockHandler handles UnlockEvent.
type UnlockHandler struct{}

func (h *UnlockHandler) Execute(ctx context.Context, event *UnlockEvent) error {
	// Handle UnlockEvent
	return nil
}

// DefinedHandler handles DefinedEvent.
type DefinedHandler struct{}

func (h *DefinedHandler) Execute(ctx context.Context, event *DefinedEvent) error {
	return nil
}

// NoTransitionHandler handles NoTransitionEvent.
type NoTransitionHandler struct{}

func (h *NoTransitionHandler) Execute(ctx context.Context, event *NoTransitionEvent) error {
	return nil
}

// HandlerWith handles EventWithHandler.
type HandlerWith struct{}

func (h *HandlerWith) Execute(ctx context.Context, event *EventWithHandler) error {
	return nil
}

// HandlerOne handles MultiHandlerEvent.
type HandlerOne struct {
	ExecutionOrder *[]string
}

func (h *HandlerOne) Execute(ctx context.Context, event *MultiHandlerEvent) error {
	*h.ExecutionOrder = append(*h.ExecutionOrder, "HandlerOne")
	return nil
}

// HandlerTwo handles MultiHandlerEvent.
type HandlerTwo struct {
	ExecutionOrder *[]string
}

func (h *HandlerTwo) Execute(ctx context.Context, event *MultiHandlerEvent) error {
	*h.ExecutionOrder = append(*h.ExecutionOrder, "HandlerTwo")
	return nil
}

// FaultyHandler handles FaultyEvent and returns an error.
type FaultyHandler struct{}

func (h *FaultyHandler) Execute(ctx context.Context, event *FaultyEvent) error {
	return fmt.Errorf("handler failed")
}

// CorrectHandler handles CorrectEvent.
type CorrectHandler struct{}

func (h *CorrectHandler) Execute(ctx context.Context, event *CorrectEvent) error {
	return nil
}

// IncorrectHandler incorrectly handles CorrectEvent.
type IncorrectHandler struct{}

func (h *IncorrectHandler) Execute(ctx context.Context, event *CorrectEvent) error {
	return nil
}

// HandlerAOrder records execution order for OrderEvent.
type HandlerAOrder struct {
	ExecutionOrder *[]string
}

func (h *HandlerAOrder) Execute(ctx context.Context, event *OrderEvent) error {
	*h.ExecutionOrder = append(*h.ExecutionOrder, "HandlerA")
	return nil
}

// HandlerBOrder records execution order for OrderEvent.
type HandlerBOrder struct {
	ExecutionOrder *[]string
}

func (h *HandlerBOrder) Execute(ctx context.Context, event *OrderEvent) error {
	*h.ExecutionOrder = append(*h.ExecutionOrder, "HandlerB")
	return nil
}

// Handler1 handles Event1.
type Handler1 struct{}

func (h *Handler1) Execute(ctx context.Context, event *Event1) error {
	return nil
}

// Handler2 handles Event2.
type Handler2 struct{}

func (h *Handler2) Execute(ctx context.Context, event *Event2) error {
	return nil
}

// ResetHandler handles ResetEvent.
type ResetHandler struct{}

func (h *ResetHandler) Execute(ctx context.Context, event *ResetEvent) error {
	return nil
}

// ComplexHandler handles ComplexEvent.
type ComplexHandler struct{}

func (h *ComplexHandler) Execute(ctx context.Context, event *ComplexEvent) error {
	return nil
}

// ----------------------------
// Unit Tests
// ----------------------------

func TestBasicTransitions(t *testing.T) {
	sm := NewStateMachine("Idle")

	// Register handlers
	RegisterHandler[*EventA, *HandlerEventA](sm, &HandlerEventA{})
	RegisterHandler[*EventB, *HandlerEventB](sm, &HandlerEventB{})

	// Define transitions
	sm.From("Idle").On(&EventA{}).To("Idle")
	sm.From("Idle").On(&EventB{}).To("Idle")

	// Execute EventA
	err := Execute(context.Background(), sm, &EventA{A: 1, B: 2})
	if err != nil {
		t.Fatalf("failed to execute EventA: %v", err)
	}
	if sm.GetCurrentState() != "Idle" {
		t.Errorf("expected state 'Idle', got '%s'", sm.GetCurrentState())
	}

	// Execute EventB
	err = Execute(context.Background(), sm, &EventB{X: 3})
	if err != nil {
		t.Fatalf("failed to execute EventB: %v", err)
	}
	if sm.GetCurrentState() != "Idle" {
		t.Errorf("expected state 'Idle', got '%s'", sm.GetCurrentState())
	}
}

func TestTransitionWithGuard(t *testing.T) {
	sm := NewStateMachine("Locked")

	// Register handler
	RegisterHandler[*UnlockEvent, *UnlockHandler](sm, &UnlockHandler{})

	// Define transition with guard
	sm.From("Locked").On(&UnlockEvent{}).When(func(ctx context.Context, e Event) bool {
		ue, ok := e.(*UnlockEvent)
		return ok && ue.Code == 1234
	}).To("Unlocked")

	// Attempt to unlock with incorrect code
	err := Execute(context.Background(), sm, &UnlockEvent{Code: 1111})
	if err == nil || !strings.Contains(err.Error(), "no valid transition") {
		t.Errorf("expected guard to prevent transition, got error: %v", err)
	}
	if sm.GetCurrentState() != "Locked" {
		t.Errorf("expected state 'Locked', got '%s'", sm.GetCurrentState())
	}

	// Unlock with correct code
	err = Execute(context.Background(), sm, &UnlockEvent{Code: 1234})
	if err != nil {
		t.Fatalf("failed to execute UnlockEvent with correct code: %v", err)
	}
	if sm.GetCurrentState() != "Unlocked" {
		t.Errorf("expected state 'Unlocked', got '%s'", sm.GetCurrentState())
	}
}

func TestErrorHandling(t *testing.T) {
	sm := NewStateMachine("Initial")

	// Register handler for DefinedEvent
	RegisterHandler[*DefinedEvent, *DefinedHandler](sm, &DefinedHandler{})

	// Define transitions
	sm.From("Initial").On(&DefinedEvent{}).To("DefinedState")
	sm.From("Initial").On(&UndefinedEvent{}).To("UndefinedState")

	// Attempt to execute UndefinedEvent which has no handler
	err := Execute(context.Background(), sm, &UndefinedEvent{})
	if err == nil || !strings.Contains(err.Error(), "no handler registered") {
		t.Errorf("expected error for missing handler, got: %v", err)
	}

	// Execute DefinedEvent
	err = Execute(context.Background(), sm, &DefinedEvent{})
	if err != nil {
		t.Errorf("unexpected error for DefinedEvent: %v", err)
	}
	if sm.GetCurrentState() != "DefinedState" {
		t.Errorf("expected state 'DefinedState', got '%s'", sm.GetCurrentState())
	}

	// Attempt to execute NoTransitionEvent which has no transition from "DefinedState"
	RegisterHandler[*NoTransitionEvent, *NoTransitionHandler](sm, &NoTransitionHandler{})
	err = Execute(context.Background(), sm, &NoTransitionEvent{})
	if err == nil || !strings.Contains(err.Error(), "no transitions defined") {
		t.Errorf("expected error for no transitions, got: %v", err)
	}
}

func TestOrphanTransitions(t *testing.T) {
	sm := NewStateMachine("Start")

	// Register handler only for EventWithHandler
	RegisterHandler[*EventWithHandler, *HandlerWith](sm, &HandlerWith{})

	// Define transitions
	sm.From("Start").On(&EventWithHandler{}).To("HandledState")
	sm.From("Start").On(&UndefinedEvent{}).To("OrphanState")

	// Find orphan transitions
	orphans := sm.FindOrphanTransitions()
	if len(orphans) != 1 {
		t.Fatalf("expected 1 orphan transition, got %d", len(orphans))
	}
	if orphans[0].EventType != "UndefinedEvent" {
		t.Errorf("expected orphan transition for 'UndefinedEvent', got '%s'", orphans[0].EventType)
	}
}

func TestMiddlewareApplication(t *testing.T) {
	sm := NewStateMachine("BaseState")

	// Define middleware
	addStateMiddleware := func(machine *StateMachine) *StateMachine {
		machine.From("BaseState").On(&EventA{}).To("AddedByMiddleware")
		return machine
	}

	logMiddleware := func(machine *StateMachine) *StateMachine {
		// Example middleware that logs transitions (dummy implementation)
		return machine
	}

	// Apply middleware using Pipe
	sm = Pipe(sm, addStateMiddleware, logMiddleware)

	// Register handler
	RegisterHandler[*EventA, *HandlerEventA](sm, &HandlerEventA{})

	// Execute EventA
	err := Execute(context.Background(), sm, &EventA{A: 10, B: 20})
	if err != nil {
		t.Fatalf("failed to execute EventA with middleware: %v", err)
	}
	if sm.GetCurrentState() != "AddedByMiddleware" {
		t.Errorf("expected state 'AddedByMiddleware', got '%s'", sm.GetCurrentState())
	}
}

func TestMultipleHandlersForSameEvent(t *testing.T) {
	sm := NewStateMachine("Start")

	// Track execution order
	var execOrder []string

	// Register multiple handlers for the same event type
	RegisterHandler[*MultiHandlerEvent, *HandlerOne](sm, &HandlerOne{ExecutionOrder: &execOrder})
	RegisterHandler[*MultiHandlerEvent, *HandlerTwo](sm, &HandlerTwo{ExecutionOrder: &execOrder})

	// Define transition
	sm.From("Start").On(&MultiHandlerEvent{}).To("Middle")

	// Execute event
	err := Execute(context.Background(), sm, &MultiHandlerEvent{})
	if err == nil || !strings.Contains(err.Error(), "invalid handler type") {
		t.Errorf("expected invalid handler type error due to multiple handlers, got: %v", err)
	}

	// Since multiple handlers are not supported, execution order should remain empty
	if len(execOrder) != 0 {
		t.Errorf("expected no handlers to execute, got: %v", execOrder)
	}
}

func TestStateConsistencyAfterFailedTransition(t *testing.T) {
	sm := NewStateMachine("Initial")

	// Register faulty handler
	RegisterHandler[*FaultyEvent, *FaultyHandler](sm, &FaultyHandler{})

	// Define transition
	sm.From("Initial").On(&FaultyEvent{}).To("NextState")

	// Execute FaultyEvent
	err := Execute(context.Background(), sm, &FaultyEvent{})
	if err == nil || !strings.Contains(err.Error(), "handler execution failed") {
		t.Errorf("expected handler execution failure, got: %v", err)
	}

	// Ensure state has not changed
	if sm.GetCurrentState() != "Initial" {
		t.Errorf("expected state to remain 'Initial' after failed transition, got '%s'", sm.GetCurrentState())
	}
}

func TestStateMachineInitialization(t *testing.T) {
	initialState := State("Initialized")
	sm := NewStateMachine(initialState)

	if sm.GetCurrentState() != initialState {
		t.Errorf("expected initial state '%s', got '%s'", initialState, sm.GetCurrentState())
	}

	if len(sm.handlers) != 0 {
		t.Errorf("expected no handlers initially, got %d", len(sm.handlers))
	}

	if len(sm.transitions) != 0 {
		t.Errorf("expected no transitions initially, got %d", len(sm.transitions))
	}
}

func TestHandlerExecutionOrder(t *testing.T) {
	sm := NewStateMachine("Start")

	// Track execution order
	var execOrder []string

	// Register handlers
	RegisterHandler[*OrderEvent, *HandlerAOrder](sm, &HandlerAOrder{ExecutionOrder: &execOrder})
	RegisterHandler[*OrderEvent, *HandlerBOrder](sm, &HandlerBOrder{ExecutionOrder: &execOrder})

	// Define transition
	sm.From("Start").On(&OrderEvent{}).To("Processed")

	// Execute event
	err := Execute(context.Background(), sm, &OrderEvent{})
	if err == nil || !strings.Contains(err.Error(), "invalid handler type") {
		t.Errorf("expected invalid handler type error due to multiple handlers, got: %v", err)
	}

	// Since multiple handlers are not supported, execution order should remain empty
	if len(execOrder) != 0 {
		t.Errorf("expected no handlers to execute, got: %v", execOrder)
	}
}

func TestStateMachineReset(t *testing.T) {
	sm := NewStateMachine("Initial")

	// Register handler
	RegisterHandler[*ResetEvent, *ResetHandler](sm, &ResetHandler{})

	// Define transition
	sm.From("Initial").On(&ResetEvent{}).To("Initial")

	// Change state
	sm.mu.Lock()
	sm.state = "Changed"
	sm.mu.Unlock()

	// Execute ResetEvent
	err := Execute(context.Background(), sm, &ResetEvent{})
	if err != nil {
		t.Fatalf("failed to execute ResetEvent: %v", err)
	}

	if sm.GetCurrentState() != "Initial" {
		t.Errorf("expected state 'Initial' after reset, got '%s'", sm.GetCurrentState())
	}
}

func TestTransitionWithMultipleGuards(t *testing.T) {
	sm := NewStateMachine("Start")

	// Register handler
	RegisterHandler[*ComplexEvent, *ComplexHandler](sm, &ComplexHandler{})

	// Define multiple guards (Note: Current implementation supports only one guard per transition)
	sm.From("Start").On(&ComplexEvent{}).
		When(func(ctx context.Context, e Event) bool {
			return true
		}).
		When(func(ctx context.Context, e Event) bool {
			return false
		}).
		To("End")

	// Execute event
	err := Execute(context.Background(), sm, &ComplexEvent{})
	if err == nil || !strings.Contains(err.Error(), "no valid transition") {
		t.Errorf("expected no valid transition due to guards, got: %v", err)
	}

	if sm.GetCurrentState() != "Start" {
		t.Errorf("expected state to remain 'Start', got '%s'", sm.GetCurrentState())
	}
}

func TestMultipleHandlers(t *testing.T) {
	sm := NewStateMachine("Initial")

	// Register multiple handlers for the same event type
	RegisterHandler[*EventA, *HandlerEventA](sm, &HandlerEventA{})
	RegisterHandler[*EventA, *HandlerEventA](sm, &HandlerEventA{})

	// Define transition
	sm.From("Initial").On(&EventA{}).To("NextState")

	// Execute EventA
	err := Execute(context.Background(), sm, &EventA{A: 5, B: 10})
	if err == nil || !strings.Contains(err.Error(), "invalid handler type") {
		t.Errorf("expected invalid handler type error due to multiple handlers, got: %v", err)
	}
}

func TestFindOrphanTransitions(t *testing.T) {
	smMap := make(map[string]*StateMachine)

	// First state machine with no orphan transitions
	sm1 := NewStateMachine("State1")
	RegisterHandler[*EventA, *HandlerEventA](sm1, &HandlerEventA{})
	sm1.From("State1").On(&EventA{}).To("State2")
	smMap["SM1"] = sm1

	// Second state machine with one orphan transition
	sm2 := NewStateMachine("StateA")
	RegisterHandler[*EventB, *HandlerEventB](sm2, &HandlerEventB{})
	sm2.From("StateA").On(&EventB{}).To("StateB")
	sm2.From("StateA").On(&UndefinedEvent{}).To("StateC")
	smMap["SM2"] = sm2

	// Third state machine with multiple orphan transitions
	sm3 := NewStateMachine("Start")
	sm3.From("Start").On(&UndefinedEvent{}).To("End")
	sm3.From("Start").On(&NoTransitionEvent{}).To("NoEnd")
	smMap["SM3"] = sm3

	// Find orphan transitions
	err := FindOrphansTransitions(smMap)
	if err == nil {
		t.Errorf("expected orphan transitions error, got nil")
	} else {
		expectedError := "SM2 has orphan transitions:\nUndefinedEvent\nSM3 has orphan transitions:\nUndefinedEvent\nNoTransitionEvent\n"
		if err.Error() != expectedError {
			t.Errorf("expected error:\n%s\ngot:\n%s", expectedError, err.Error())
		}
	}
}
