package main

//go:generate go run gen.go types.go
type InitFinalState struct {
}

type InitFinalEnvironment struct{}

const (
	StateInit    = "init"
	StateRunning = "running"
	StateFinal   = "final"
)

type EventInit struct{}
type EventRun struct{}
type EventFinish struct{}
