
package main

// Code generated go-fsmgen DO NOT EDIT.

import (
	"context"
	"errors"
)

type InitFinalMachine struct {
	CurrentState string
	State *InitFinalState

	transitions  map[string]map[string]string
	RunAction func(ctx InitFinalMachineContext, state *InitFinalState, ev EventRun) error
	FinishAction func(ctx InitFinalMachineContext, state *InitFinalState, ev EventFinish) error

	OnStateInit func(ctx InitFinalMachineContext, state InitFinalState) error
	OnStateRunning func(ctx InitFinalMachineContext, state InitFinalState) error
	OnStateFinal func(ctx InitFinalMachineContext, state InitFinalState) error
}

type InitFinalMachineContext interface {
	Context() context.Context
	TriggerRun(ev EventRun) error
	TriggerFinish(ev EventFinish) error
}

type initFinalMachineContext struct {
	ctx context.Context
	machine *InitFinalMachine
}

func newInitFinalContext(ctx context.Context, machine *InitFinalMachine) InitFinalMachineContext {
	return &initFinalMachineContext{
		ctx: ctx,
		machine: machine,
	}
}

func (ctx initFinalMachineContext) Context() context.Context {
	return context.Background()
}

func (ctx initFinalMachineContext) TriggerRun(ev EventRun) error {
	return ctx.machine.TriggerRun(ctx.ctx, ev)
}
func (ctx initFinalMachineContext) TriggerFinish(ev EventFinish) error {
	return ctx.machine.TriggerFinish(ctx.ctx, ev)
}

func NewInitFinalMachine(state *InitFinalState) *InitFinalMachine{
	return &InitFinalMachine{
		State: state,
		CurrentState: "init",
		transitions: map[string]map[string]string{
				"": {
				},
				"final": {
				},
				"init": {
					"run": "running",
				},
				"running": {
					"finish": "final",
				},
		},
	}
}

func (machine *InitFinalMachine) Start(ctx context.Context) (error) {
	return machine.didEnterState(ctx)
}

func (machine *InitFinalMachine) getState(event string) (string, error) {
	target := machine.transitions[machine.CurrentState][event]
	if target != "" {
		return target, nil
	}
	target = machine.transitions[""][event]
	if target != "" {
		return target, nil
	}
	return "", errors.New("invalid transition: no transition target from " + machine.CurrentState + " via " + event)
}

func (machine *InitFinalMachine) didEnterState(ctx context.Context) error {
	switch machine.CurrentState {
	case "init":
		if machine.OnStateInit == nil {
			break
		}
		return machine.OnStateInit(newInitFinalContext(ctx, machine), *machine.State)
	case "running":
		if machine.OnStateRunning == nil {
			break
		}
		return machine.OnStateRunning(newInitFinalContext(ctx, machine), *machine.State)
	case "final":
		if machine.OnStateFinal == nil {
			break
		}
		return machine.OnStateFinal(newInitFinalContext(ctx, machine), *machine.State)
	}
	return nil
}

func (machine *InitFinalMachine) TriggerRun (ctx context.Context, ev EventRun) error {
	target, err := machine.getState("run")
	if err != nil {
	return err
	}
	machine.CurrentState = target
	if machine.RunAction != nil {
		err := machine.RunAction(newInitFinalContext(ctx, machine), machine.State, ev)
		if err != nil {
			return err
		}
	}
	return machine.didEnterState(ctx)
}

func (machine *InitFinalMachine) TriggerFinish (ctx context.Context, ev EventFinish) error {
	target, err := machine.getState("finish")
	if err != nil {
	return err
	}
	machine.CurrentState = target
	if machine.FinishAction != nil {
		err := machine.FinishAction(newInitFinalContext(ctx, machine), machine.State, ev)
		if err != nil {
			return err
		}
	}
	return machine.didEnterState(ctx)
}

