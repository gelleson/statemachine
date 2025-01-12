package statemachine

func Pipe(stateMachine *StateMachine, middlewares ...func(machine *StateMachine) *StateMachine) *StateMachine {
	for _, middleware := range middlewares {
		stateMachine = middleware(stateMachine)
	}
	return stateMachine
}
