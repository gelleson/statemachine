# State Machine Library in Go

This repository provides a simple yet powerful state machine implementation in Go. The library allows you to define states, transitions, and handlers for different events, making it easy to manage complex state transitions in your applications.

## Features

- **Simple API**: Define states, transitions, and event handlers with an easy-to-use API.
- **Event Handling**: Register handlers for specific events and execute them when the event occurs.
- **Guard Conditions**: Add conditions to transitions to ensure they only occur under specific circumstances.
- **Orphan Transition Detection**: Identify transitions that lack registered handlers.
- **Concurrency Safe**: The state machine is safe to use in concurrent environments.
- **DOT Representation**: Generate a DOT representation of the state machine for visualization.

## Installation

To install the library, use `go get`:

```bash
go get github.com/gelleson/statemachine
```

## Usage

### Defining Events and Handlers

First, define the events and their corresponding handlers.

```go
package statemachine

import (
	"context"
	"fmt"
)

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
```

### Setting Up the State Machine

Next, initialize the state machine and define the transitions.

```go
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
```

### Adding Guard Conditions

You can add guard conditions to transitions to ensure they only occur under specific circumstances.

```go
sm.From("Locked").On(&UnlockEvent{}).When(func(ctx context.Context, e Event) bool {
	ue, ok := e.(*UnlockEvent)
	return ok && ue.Code == 1234
}).To("Unlocked")
```

### Detecting Orphan Transitions

Identify transitions that do not have registered handlers.

```go
orphans := sm.FindOrphanTransitions()
if len(orphans) > 0 {
	fmt.Println("Orphan transitions found:", orphans)
}
```

### Generating DOT Representation

Generate a DOT representation of the state machine for visualization.

```go
dot := sm.GenerateDOT()
fmt.Println(dot)
```

## Examples

Check out the [examples](examples) directory for more usage examples.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any bugs or feature requests.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by various state machine implementations in different programming languages.
- Special thanks to the Go community for their support and contributions.

---

Feel free to explore and use this library to manage state transitions in your Go applications. Happy coding! 🚀
