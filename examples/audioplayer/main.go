package main

import (
	"context"
	"log"
	"os"
	"time"
)

//go:generate go run gen.go types.go
func main() {
	ctx := context.Background()
	machine := NewMachine()
	err := machine.TriggerLoad(ctx, EventLoad{
		File: nil,
	})
	if err != nil {
		log.Panic(err)
	}
	time.Sleep(2 * time.Second)
	err = machine.TriggerPause(ctx, EventPause{})
	if err != nil {
		log.Panic(err)
	}
	err = machine.TriggerPause(ctx, EventPause{})
	if err != nil {
		log.Print("error: " + err.Error())
	}
}

func NewMachine() *AudioPlayerMachine {
	env := AudioPlayerEnvironment{Logger: log.New(os.Stdout, "fsm ", 0)}
	machine := NewAudioPlayerMachine(&AudioPlayerState{
		Player: &StringAudioPlayer{},
	}, env)
	machine.LoadAction = func(ctx AudioPlayerMachineContext, state *AudioPlayerState, ev EventLoad) error {
		state.file = ev.File
		return nil
	}
	machine.ErrorAction = func(ctx AudioPlayerMachineContext, state *AudioPlayerState, ev EventError) error {
		state.Message = ev.Message
		return nil
	}
	machine.OnStateLoading = func(ctx AudioPlayerMachineContext, env AudioPlayerEnvironment, state AudioPlayerState) error {
		env.Logger.Print("Did enter State Loading")
		err := state.Player.Load(state.file)
		if err != nil {
			return ctx.TriggerError(EventError{
				Message: err.Error(),
			})
		}
		return ctx.TriggerPlay(EventPlay{})
	}
	machine.OnStatePlaying = func(ctx AudioPlayerMachineContext, env AudioPlayerEnvironment, state AudioPlayerState) error {
		env.Logger.Print("Did enter State Playing")
		err := state.Player.Play()
		if err != nil {
			return ctx.TriggerError(EventError{
				Message: err.Error(),
			})
		}
		return nil
	}
	machine.OnStatePaused = func(ctx AudioPlayerMachineContext, env AudioPlayerEnvironment, state AudioPlayerState) error {
		env.Logger.Print("Did enter State Paused")
		err := state.Player.Pause()
		if err != nil {
			return ctx.TriggerError(EventError{
				Message: err.Error(),
			})
		}
		return nil
	}
	return machine
}
