package statemachine

import "context"

type optionFn func(*StateMachine)

func applyOptions(sm *StateMachine, options []optionFn) *StateMachine {
	for _, o := range options {
		o(sm)
	}
	return sm
}

func WithPreTransitionMiddlewares(middlewares ...func(context.Context, Event) error) optionFn {
	return func(sm *StateMachine) {
		sm.mw.PreTransitionMiddlewares = middlewares
	}
}

func WithPostTransitionMiddlewares(middlewares ...func(context.Context, Event) error) optionFn {
	return func(sm *StateMachine) {
		sm.mw.PostTransitionMiddlewares = middlewares
	}
}
