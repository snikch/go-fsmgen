package finalstate

import (
	"context"
	"testing"

	"gotest.tools/assert"
)

func TestFinalState(t *testing.T) {
	expectedTransitions := []string{"init", "running", "final"}
	expectedActions := []string{"run", "finish"}
	actions := make([]string, 0, 2)
	transitions := make([]string, 0, 3)
	machine := NewInitFinalMachine(&State{}, Environment{})
	machine.OnStateInit = func(ctx InitFinalMachineContext, env Environment, state State) error {
		transitions = append(transitions, "init")
		return ctx.TriggerRun(EventRun{})
	}
	machine.OnStateRunning = func(ctx InitFinalMachineContext, env Environment, state State) error {
		transitions = append(transitions, "running")
		return ctx.TriggerFinish(EventFinish{})
	}
	machine.OnStateFinal = func(ctx InitFinalMachineContext, env Environment, state State) error {
		transitions = append(transitions, "final")
		return nil
	}
	machine.RunAction = func(ctx InitFinalMachineContext, state *State, ev EventRun) error {
		actions = append(actions, "run")
		return nil
	}
	machine.FinishAction = func(ctx InitFinalMachineContext, state *State, ev EventFinish) error {
		actions = append(actions, "finish")
		return nil
	}

	assert.Equal(t, "init", machine.CurrentState)
	err := machine.Start(context.Background())
	assert.NilError(t, err)
	assert.Equal(t, "final", machine.CurrentState)
	assert.DeepEqual(t, expectedActions, actions)
	assert.DeepEqual(t, expectedTransitions, transitions)
}
