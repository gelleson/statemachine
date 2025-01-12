package statemachine

import (
	"fmt"
	"strings"
)

type OrphanFinder interface {
	FindOrphanTransitions() []TransitionRule
}

func FindOrphanTransitions(name string, sm *StateMachine) error {
	var stringBuilder strings.Builder

	orphans := sm.FindOrphanTransitions()
	if len(orphans) == 0 {
		return nil
	}

	stringBuilder.WriteString(name)
	stringBuilder.WriteString(" has orphan transitions:\n")

	for _, orphan := range orphans {
		stringBuilder.WriteString(orphan.EventType)
		stringBuilder.WriteString("\n")
	}

	return fmt.Errorf(stringBuilder.String())
}

func FindOrphansTransitions(m map[string]*StateMachine) error {
	for name, sm := range m {
		err := FindOrphanTransitions(name, sm)
		if err != nil {
			return err
		}
	}

	return nil
}
