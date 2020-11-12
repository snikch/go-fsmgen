# go-fsmgen
[![Documentation](https://godoc.org/github.com/snikch/go-fsmgen?status.svg)](http://godoc.org/github.com/snikch/go-fsmgen)
[![Go Report Card](https://goreportcard.com/badge/github.com/snikch/go-fsmgen)](https://goreportcard.com/report/github.com/snikch/go-fsmgen)
[![GitHub issues](https://img.shields.io/github/issues/snikch/go-fsmgen.svg)](https://github.com/snikch/go-fsmgen/issues)
[![license](https://img.shields.io/github/license/snikch/go-fsmgen.svg?maxAge=2592000)](https://github.com/snikch/go-fsmgen/blob/main/LICENSE)


`go-fsmgen` is a Finite State Machine generator. It provides a strongly typed implementation of a FSM through code generation.

```
go get github.com/snikch/go-fsmgen
```

## Generation

For a full example, see the [examples](./examples) directory.

```go
func main() {
	gen := fsmgen.New("audio_player", AudioPlayerState{}, AudioPlayerEnvironment{}, "init", "loading", "playing", "paused")
	gen.PackageName = "main"
	gen.AddEvent(fsmgen.NewEvent("load", EventLoad{}).FromAny().To("loading"))
	gen.AddEvent(fsmgen.NewEvent("play", EventPlay{}).From("paused", "loading").To("playing"))
	gen.AddEvent(fsmgen.NewEvent("pause", EventPause{}).From("playing").To("paused"))
	gen.AddEvent(fsmgen.NewEvent("error", EventError{}).FromAny().To("init"))
	err := gen.Write()
	if err != nil {
		log.Panic(err)
	}
}
```

## Usage

```go
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
```

```
$ ./audioplayer
fsm Did enter State Loading
2020/11/13 06:35:38 Loading
fsm Did enter State Playing
2020/11/13 06:35:38 Play
fsm Did enter State Paused
2020/11/13 06:35:40 Pause
2020/11/13 06:35:40 error: invalid transition: no transition target from paused via pause
```
## See Also

Design influenced by the following projects:

* [go-statemachine](https://github.com/filecoin-project/go-statemachine/tree/master/fsm)
* [xstate](https://xstate.js.org/)
