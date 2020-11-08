// +build ignore

package main

import (
	"log"
)

func main() {
	gen := go_fsmgen.New("audio_player", "init", "loading", "playing", "paused")
	gen.PackageName = "main"
	gen.AddEvent(go_fsmgen.NewEvent("load", EventLoad{}).FromAny().To("loading"))
	gen.AddEvent(go_fsmgen.NewEvent("play", EventPlay{}).From("paused", "loading").To("playing"))
	gen.AddEvent(go_fsmgen.NewEvent("pause", EventPause{}).From("playing").To("paused"))
	gen.AddEvent(go_fsmgen.NewEvent("error", EventError{}).FromAny().To("init"))
	err := gen.Write()
	if err != nil {
		log.Panic(err)
	}
}
