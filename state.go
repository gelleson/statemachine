package statemachine

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

// Event is the interface that all events must implement
type Event interface {
	EventType() string
}

// State represents a state in the state machine
type State string

// Handler interface for handling events
type Handler[T Event] interface {
	Execute(ctx context.Context, event T) error
}

// TransitionRule defines when a transition can occur
type TransitionRule struct {
	From      State
	To        State
	EventType string
	Guard     func(context.Context, Event) bool // optional condition
}

// StateMachine handles different types of events with their corresponding handlers
type StateMachine struct {
	handlers    map[string]any
	transitions map[State][]TransitionRule
	state       State
	mu          sync.RWMutex
}

// StateMachineBuilder helps configure the state machine
type StateMachineBuilder struct {
	machine   *StateMachine
	fromState State
	eventType string
	guard     func(context.Context, Event) bool
}

// NewStateMachine creates a new instance of StateMachine
func NewStateMachine(initialState State) *StateMachine {
	return &StateMachine{
		handlers:    make(map[string]any),
		transitions: make(map[State][]TransitionRule),
		state:       initialState,
	}
}

// From starts defining a transition from a state
func (sm *StateMachine) From(state State) *StateMachineBuilder {
	return &StateMachineBuilder{
		machine:   sm,
		fromState: state,
	}
}

// On specifies the event type for the transition
func (b *StateMachineBuilder) On(eventType Event) *StateMachineBuilder {
	b.eventType = eventType.EventType()
	return b
}

// When adds a guard condition to the transition
func (b *StateMachineBuilder) When(guard func(context.Context, Event) bool) *StateMachineBuilder {
	b.guard = guard
	return b
}

// To completes the transition definition
func (b *StateMachineBuilder) To(toState State) *StateMachineBuilder {
	b.machine.mu.Lock()
	defer b.machine.mu.Unlock()

	rule := TransitionRule{
		From:      b.fromState,
		To:        toState,
		EventType: b.eventType,
		Guard:     b.guard,
	}

	b.machine.transitions[b.fromState] = append(b.machine.transitions[b.fromState], rule)
	return b
}

// findValidTransition finds a valid transition rule for the current state and event
func (sm *StateMachine) findValidTransition(ctx context.Context, event Event) (*TransitionRule, error) {
	rules, exists := sm.transitions[sm.state]
	if !exists {
		return nil, fmt.Errorf("no transitions defined for state: %s", sm.state)
	}

	for _, rule := range rules {
		if rule.EventType == event.EventType() {
			if rule.Guard == nil || rule.Guard(ctx, event) {
				return &rule, nil
			}
		}
	}

	return nil, fmt.Errorf("no valid transition found for event %s in state %s", event.EventType(), sm.state)
}

// GetCurrentState returns the current state
func (sm *StateMachine) GetCurrentState() State {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.state
}

// GenerateDOT generates a DOT representation of the state machine
func (sm *StateMachine) GenerateDOT() string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	var sb strings.Builder
	sb.WriteString("digraph StateMachine {\n")
	sb.WriteString("    rankdir=LR;\n")
	sb.WriteString("    node [shape = circle];\n")

	// Define all states
	for state := range sm.transitions {
		if state == sm.state {
			sb.WriteString(fmt.Sprintf("    \"%s\" [shape=doublecircle];\n", state))
		} else {
			sb.WriteString(fmt.Sprintf("    \"%s\";\n", state))
		}
	}

	// Define transitions
	for fromState, rules := range sm.transitions {
		for _, rule := range rules {
			label := rule.EventType
			if rule.Guard != nil {
				label += " [guard]"
			}
			sb.WriteString(fmt.Sprintf("    \"%s\" -> \"%s\" [ label = \"%s\" ];\n", fromState, rule.To, label))
		}
	}

	sb.WriteString("}\n")
	return sb.String()
}

// FindOrphanTransitions returns a list of TransitionRules that do not have a registered handler.
func (sm *StateMachine) FindOrphanTransitions() []TransitionRule {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	var orphans []TransitionRule
	for fromState, rules := range sm.transitions {
		for _, rule := range rules {
			if _, exists := sm.handlers[rule.EventType]; !exists {
				orphans = append(orphans, TransitionRule{
					From:      fromState,
					To:        rule.To,
					EventType: rule.EventType,
					Guard:     rule.Guard,
				})
			}
		}
	}
	return orphans
}
