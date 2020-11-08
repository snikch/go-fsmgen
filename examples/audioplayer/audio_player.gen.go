
package main

import (
	"context"
	"errors"
)

type AudioPlayerMachine struct {
	CurrentState string
	State *AudioPlayerState

	transitions  map[string]map[string]string
	LoadAction func(ctx AudioPlayerMachineContext, state *AudioPlayerState, ev EventLoad) error
	PlayAction func(ctx AudioPlayerMachineContext, state *AudioPlayerState, ev EventPlay) error
	PauseAction func(ctx AudioPlayerMachineContext, state *AudioPlayerState, ev EventPause) error
	ErrorAction func(ctx AudioPlayerMachineContext, state *AudioPlayerState, ev EventError) error

	OnStateInit func(ctx AudioPlayerMachineContext, state AudioPlayerState) error
	OnStateLoading func(ctx AudioPlayerMachineContext, state AudioPlayerState) error
	OnStatePlaying func(ctx AudioPlayerMachineContext, state AudioPlayerState) error
	OnStatePaused func(ctx AudioPlayerMachineContext, state AudioPlayerState) error
}

type AudioPlayerMachineContext interface {
	Context() context.Context
	TriggerLoad(ev EventLoad) error
	TriggerPlay(ev EventPlay) error
	TriggerPause(ev EventPause) error
	TriggerError(ev EventError) error
}

type audioPlayerMachineContext struct {
	ctx context.Context
	machine *AudioPlayerMachine
}

func newAudioPlayerContext(ctx context.Context, machine *AudioPlayerMachine) AudioPlayerMachineContext {
	return &audioPlayerMachineContext{
		ctx: ctx,
		machine: machine,
	}
}

func (ctx audioPlayerMachineContext) Context() context.Context {
	return context.Background()
}

func (ctx audioPlayerMachineContext) TriggerLoad(ev EventLoad) error {
	return ctx.machine.TriggerLoad(ctx.ctx, ev)
}
func (ctx audioPlayerMachineContext) TriggerPlay(ev EventPlay) error {
	return ctx.machine.TriggerPlay(ctx.ctx, ev)
}
func (ctx audioPlayerMachineContext) TriggerPause(ev EventPause) error {
	return ctx.machine.TriggerPause(ctx.ctx, ev)
}
func (ctx audioPlayerMachineContext) TriggerError(ev EventError) error {
	return ctx.machine.TriggerError(ctx.ctx, ev)
}

func NewAudioPlayerMachine(state *AudioPlayerState) *AudioPlayerMachine{
	return &AudioPlayerMachine{
		State: state,
		CurrentState: "init",
		transitions: map[string]map[string]string{
				"": {
					"error": "init",
					"load": "loading",
				},
				"init": {
				},
				"loading": {
					"play": "playing",
				},
				"paused": {
					"play": "playing",
				},
				"playing": {
					"pause": "paused",
				},
		},
	}
}

func (machine *AudioPlayerMachine) getState(event string) (string, error) {
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

func (machine *AudioPlayerMachine) didEnterState(ctx context.Context) error {
	switch machine.CurrentState {
	case "init":
		if machine.OnStateInit == nil {
			break
		}
		return machine.OnStateInit(newAudioPlayerContext(ctx, machine), *machine.State)
	case "loading":
		if machine.OnStateLoading == nil {
			break
		}
		return machine.OnStateLoading(newAudioPlayerContext(ctx, machine), *machine.State)
	case "playing":
		if machine.OnStatePlaying == nil {
			break
		}
		return machine.OnStatePlaying(newAudioPlayerContext(ctx, machine), *machine.State)
	case "paused":
		if machine.OnStatePaused == nil {
			break
		}
		return machine.OnStatePaused(newAudioPlayerContext(ctx, machine), *machine.State)
	}
	return nil
}

func (machine *AudioPlayerMachine) TriggerLoad (ctx context.Context, ev EventLoad) error {
	target, err := machine.getState("load")
	if err != nil {
	return err
	}
	machine.CurrentState = target
	if machine.LoadAction != nil {
		err := machine.LoadAction(newAudioPlayerContext(ctx, machine), machine.State, ev)
		if err != nil {
			return err
		}
	}
	return machine.didEnterState(ctx)
}

func (machine *AudioPlayerMachine) TriggerPlay (ctx context.Context, ev EventPlay) error {
	target, err := machine.getState("play")
	if err != nil {
	return err
	}
	machine.CurrentState = target
	if machine.PlayAction != nil {
		err := machine.PlayAction(newAudioPlayerContext(ctx, machine), machine.State, ev)
		if err != nil {
			return err
		}
	}
	return machine.didEnterState(ctx)
}

func (machine *AudioPlayerMachine) TriggerPause (ctx context.Context, ev EventPause) error {
	target, err := machine.getState("pause")
	if err != nil {
	return err
	}
	machine.CurrentState = target
	if machine.PauseAction != nil {
		err := machine.PauseAction(newAudioPlayerContext(ctx, machine), machine.State, ev)
		if err != nil {
			return err
		}
	}
	return machine.didEnterState(ctx)
}

func (machine *AudioPlayerMachine) TriggerError (ctx context.Context, ev EventError) error {
	target, err := machine.getState("error")
	if err != nil {
	return err
	}
	machine.CurrentState = target
	if machine.ErrorAction != nil {
		err := machine.ErrorAction(newAudioPlayerContext(ctx, machine), machine.State, ev)
		if err != nil {
			return err
		}
	}
	return machine.didEnterState(ctx)
}

