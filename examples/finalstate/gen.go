//go:build ignore

package main

import (
	"log"

	"github.com/snikch/go-fsmgen"
)

func main() {
	gen := fsmgen.New("init_final", StateInit, StateRunning, StateFinal)
	gen.PackageName = "main"
	gen.AddEvent(fsmgen.NewEvent("run", EventRun{}).From(StateInit).To(StateRunning))
	gen.AddEvent(fsmgen.NewEvent("finish", EventFinish{}).From(StateRunning).To(StateFinal))
	err := gen.Write()
	if err != nil {
		log.Panic(err)
	}
}
