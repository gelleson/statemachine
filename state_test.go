package statemachine

import (
	"context"
	"fmt"
	"testing"
)

type SubmitInitiateEvent struct{}
type SubmitCompleteEvent struct{}
type RequestEvent struct{}
type ApproveEvent struct{}

func (e *SubmitInitiateEvent) EventType() string { return "submit_initiate" }
func (e *SubmitCompleteEvent) EventType() string { return "submit_complete" }
func (e *RequestEvent) EventType() string        { return "request" }
func (e *ApproveEvent) EventType() string        { return "approve" }

// Define handlers
type GenericHandler[T Event] struct {
	message string
}

func (h *GenericHandler[T]) Execute(ctx context.Context, event T) error {
	return nil
}

func TestTransitionNotSubmittedToDraft(t *testing.T) {
	// Initialize the state machine
	sm := setupStateMachine()

	// Set the initial state
	sm.state = "not_submitted"

	// Execute the event
	err := Execute(context.Background(), sm, &SubmitInitiateEvent{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Assert the final state
	if sm.GetCurrentState() != "draft" {
		t.Errorf("Expected state 'draft', but got '%s'", sm.GetCurrentState())
	}
}

func TestTransitionNotSubmittedToRequested(t *testing.T) {
	sm := setupStateMachine()
	sm.state = "not_submitted"

	err := Execute(context.Background(), sm, &RequestEvent{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if sm.GetCurrentState() != "requested" {
		t.Errorf("Expected state 'requested', but got '%s'", sm.GetCurrentState())
	}
}

func TestTransitionDraftToSubmitted(t *testing.T) {
	sm := setupStateMachine()
	sm.state = "draft"

	err := Execute(context.Background(), sm, &SubmitCompleteEvent{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if sm.GetCurrentState() != "submitted" {
		t.Errorf("Expected state 'submitted', but got '%s'", sm.GetCurrentState())
	}

	fmt.Println(sm.GenerateDOT())
}

func TestTransitionSubmittedToRequested(t *testing.T) {
	sm := setupStateMachine()
	sm.state = "submitted"

	err := Execute(context.Background(), sm, &RequestEvent{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if sm.GetCurrentState() != "requested" {
		t.Errorf("Expected state 'requested', but got '%s'", sm.GetCurrentState())
	}
}

func TestTransitionRequestedToDraft(t *testing.T) {
	sm := setupStateMachine()
	sm.state = "requested"

	err := Execute(context.Background(), sm, &SubmitInitiateEvent{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if sm.GetCurrentState() != "draft" {
		t.Errorf("Expected state 'draft', but got '%s'", sm.GetCurrentState())
	}
}

func TestTransitionSubmittedToApproved(t *testing.T) {
	sm := setupStateMachine()
	sm.state = "submitted"

	err := Execute(context.Background(), sm, &ApproveEvent{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if sm.GetCurrentState() != "approved" {
		t.Errorf("Expected state 'approved', but got '%s'", sm.GetCurrentState())
	}
}

func TestInvalidTransition(t *testing.T) {
	sm := setupStateMachine()
	sm.state = "not_submitted"

	err := Execute(context.Background(), sm, &ApproveEvent{})
	if err == nil {
		t.Errorf("Expected error for invalid transition, but got none")
	}

	if sm.GetCurrentState() != "not_submitted" {
		t.Errorf("Expected state 'not_submitted', but got '%s'", sm.GetCurrentState())
	}
}

func setupStateMachine() *StateMachine {
	// Define events

	// Initialize the state machine
	sm := NewStateMachine("not_submitted")

	// Register handlers
	RegisterHandler[*SubmitInitiateEvent](sm, &GenericHandler[*SubmitInitiateEvent]{})
	RegisterHandler[*SubmitCompleteEvent](sm, &GenericHandler[*SubmitCompleteEvent]{})
	RegisterHandler[*RequestEvent](sm, &GenericHandler[*RequestEvent]{})
	RegisterHandler[*ApproveEvent](sm, &GenericHandler[*ApproveEvent]{})

	// Define transitions
	sm.From("not_submitted").On(&SubmitInitiateEvent{}).To("draft")
	sm.From("not_submitted").On(&RequestEvent{}).To("requested")
	sm.From("draft").On(&SubmitCompleteEvent{}).To("submitted")
	sm.From("submitted").On(&RequestEvent{}).To("requested")
	sm.From("submitted").On(&ApproveEvent{}).To("approved")
	sm.From("requested").On(&SubmitInitiateEvent{}).To("draft")

	return sm
}
