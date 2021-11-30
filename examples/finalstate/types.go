package finalstate

//go:generate go run gen/gen.go
type State struct {
}

type Environment struct{}

const (
	StateInit    = "init"
	StateRunning = "running"
	StateFinal   = "final"
)

type EventInit struct{}
type EventRun struct{}
type EventFinish struct{}
