// +build ignore

package main

import (
	"log"

	"github.com/snikch/go-fsmgen"
)

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
