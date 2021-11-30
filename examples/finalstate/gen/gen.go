//go:build ignore

package main

import (
	"log"

	"github.com/snikch/go-fsmgen"
	"github.com/snikch/go-fsmgen/examples/finalstate"
)

func main() {
	gen := fsmgen.New("init_final", finalstate.State{}, finalstate.Environment{}, finalstate.StateInit, finalstate.StateRunning, finalstate.StateFinal)
	gen.PackageName = "finalstate"
	gen.AddEvent(fsmgen.NewEvent("run", finalstate.EventRun{}).From(finalstate.StateInit).To(finalstate.StateRunning))
	gen.AddEvent(fsmgen.NewEvent("finish", finalstate.EventFinish{}).From(finalstate.StateRunning).To(finalstate.StateFinal))
	err := gen.Write()
	if err != nil {
		log.Panic(err)
	}
}
