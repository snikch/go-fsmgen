package fsmgen

import (
	"bytes"
	"io/ioutil"
	"os"
	"reflect"
	"text/template"

	"github.com/iancoleman/strcase"
)

// Generator provides a type for defining a finite state machine's states and events.
type Generator struct {
	// Name is the name of the state machine. It defines the name of the types and file generated.
	Name string
	// PackageName defines the name of the package the generated file belongs to. Defaults to Name.
	PackageName string
	// Filename defines where the state machine file will be written to. Defaults to Name.gen.go
	Filename string
	// States contains all of the state names that the state machine may be in.
	States []string
	// Events is a slice of all possible events that can occur in the state machine.
	Events []*Event
}

// New returns a new Generator with the supplied name and states. The first supplied state is the initial state.
func New(name string, states ...string) *Generator {
	return &Generator{
		Name:        name,
		PackageName: name,
		Filename:    name + ".gen.go",
		States:      states,
	}
}

// AddEvent adds the supplied event to the generated state machine.
func (gen *Generator) AddEvent(ev *Event) {
	gen.Events = append(gen.Events, ev)
}

// Write will generator and output the state machine to file.
func (gen *Generator) Write() error {
	t, err := template.New("fsm").Parse(tmpl)
	if err != nil {
		return err
	}
	out := &bytes.Buffer{}
	err = t.Execute(out, &tmplGenerator{Generator: gen})
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(gen.Filename, out.Bytes(), os.FileMode(0644))
	return err
}

type tmplGenerator struct {
	*Generator
}

func (gen *tmplGenerator) ExportedName(str string) string {
	return strcase.ToCamel(str)
}

func (gen *tmplGenerator) UnexportedName(str string) string {
	return strcase.ToLowerCamel(str)
}

func (gen *tmplGenerator) TransitionMap() map[string]map[string]string {
	out := map[string]map[string]string{
		"": {},
	}
	for _, state := range gen.States {
		out[state] = make(map[string]string)
	}
	for _, event := range gen.Events {
		if len(event.FromStates) == 0 {
			out[""][event.Name] = event.ToState
			continue
		}
		for _, state := range event.FromStates {
			out[state][event.Name] = event.ToState
		}
	}
	return out
}

// Event defines an Event that transitions from a set of states to a new state.
type Event struct {
	Name       string
	FromStates []string
	ToState    string
	ObjName    reflect.Type
}

// NewEvent returns a new event with the supplied name and Event object.
func NewEvent(name string, obj interface{}) *Event {
	return &Event{
		Name:    name,
		ObjName: reflect.TypeOf(obj),
	}
}

// FromAny defines all states as valid source states for this event.
func (ev *Event) FromAny() *Event {
	ev.FromStates = nil
	return ev
}

// From defines the states that are valid as source states for this event.
func (ev *Event) From(from ...string) *Event {
	ev.FromStates = from
	return ev
}

// To defines the target state after this event.
func (ev *Event) To(to string) *Event {
	ev.ToState = to
	return ev
}

const tmpl = `
package {{ .PackageName }}

// Code generated go-fsmgen DO NOT EDIT.

import (
	"context"
	"errors"
)

type {{ .ExportedName .Name }}Machine struct {
	CurrentState string
	State *{{ .ExportedName .Name }}State

	transitions  map[string]map[string]string

{{- range $event := .Events }}
	{{ $.ExportedName $event.Name }}Action func(ctx {{ $.ExportedName $.Name }}MachineContext, state *{{ $.ExportedName $.Name }}State, ev {{ $event.ObjName.Name }}) error
{{- end }}
{{ range $state := .States }}
	OnState{{ $.ExportedName $state }} func(ctx {{ $.ExportedName $.Name }}MachineContext, state {{ $.ExportedName $.Name }}State) error
{{- end }}
}

type {{ .ExportedName .Name }}MachineContext interface {
	Context() context.Context
{{- range $event := .Events }}
	Trigger{{ $.ExportedName $event.Name }}(ev {{ $event.ObjName.Name }}) error
{{- end }}
}

type {{ .UnexportedName .Name }}MachineContext struct {
	ctx context.Context
	machine *{{ .ExportedName .Name }}Machine
}

func new{{ .ExportedName .Name }}Context(ctx context.Context, machine *{{ .ExportedName .Name }}Machine) {{ .ExportedName .Name }}MachineContext {
	return &{{ .UnexportedName .Name }}MachineContext{
		ctx: ctx,
		machine: machine,
	}
}

func (ctx {{ .UnexportedName .Name }}MachineContext) Context() context.Context {
	return context.Background()
}
{{ range $event := .Events }}
func (ctx {{ $.UnexportedName $.Name }}MachineContext) Trigger{{ $.ExportedName $event.Name }}(ev {{ $event.ObjName.Name }}) error {
	return ctx.machine.Trigger{{ $.ExportedName $event.Name }}(ctx.ctx, ev)
}
{{- end }}

func New{{ .ExportedName .Name }}Machine(state *{{ .ExportedName .Name }}State) *{{ .ExportedName  .Name}}Machine{
	return &{{ .ExportedName .Name }}Machine{
		State: state,
		CurrentState: "{{ (index .States 0) }}",
		transitions: map[string]map[string]string{
			{{- range $from, $events := .TransitionMap }}
				"{{ $from }}": {
				{{- range $event, $target := $events }}
					"{{ $event }}": "{{ $target }}",
				{{- end }}
				},
			{{- end }}
		},
	}
}

func (machine *{{ .ExportedName .Name }}Machine) getState(event string) (string, error) {
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

func (machine *{{ .ExportedName .Name }}Machine) didEnterState(ctx context.Context) error {
	switch machine.CurrentState {
	{{- range $state := .States }}
	case "{{ $state }}":
		if machine.OnState{{ $.ExportedName $state }} == nil {
			break
		}
		return machine.OnState{{ $.ExportedName $state }}(new{{ $.ExportedName $state }}Context(ctx, machine), *machine.State)
	{{- end }}
	}
	return nil
}
{{ range $event := .Events }}
func (machine *{{ $.ExportedName $.Name }}Machine) Trigger{{ $.ExportedName $event.Name }} (ctx context.Context, ev {{ $event.ObjName.Name }}) error {
	target, err := machine.getState("{{ $event.Name }}")
	if err != nil {
	return err
	}
	machine.CurrentState = target
	if machine.{{ $.ExportedName $event.Name }}Action != nil {
		err := machine.{{ $.ExportedName $event.Name }}Action(new{{ $.ExportedName $.Name }}Context(ctx, machine), machine.State, ev)
		if err != nil {
			return err
		}
	}
	return machine.didEnterState(ctx)
}
{{ end }}
`
