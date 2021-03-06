package events

import (
	"sort"
	"time"
)

// The Event type represents unique events generated by the program. They carry
// context about how they were triggered and information to pass to handlers.
type Event struct {
	// Message carries information about the event in a human-readable format.
	Message string

	// Source represents the location where this event was generated from.
	Source string

	// Args is the list of arguments of the event, it is intended to give
	// context about the information carried by the even in a format that can
	// be processed by a program.
	Args Args

	// Time is the time at which the event was generated.
	Time time.Time

	// Debug is set to true if this is a debugging event.
	Debug bool
}

// Clone makes a deep copy of the event, the returned value doesn't shared any
// pointer with the original.
func (e *Event) Clone() *Event {
	var a Args
	var m []byte
	var s []byte

	if n := len(e.Args); n != 0 {
		a = make(Args, n)
		for i := range a {
			a[i].Name = e.Args[i].Name
			a[i].Value = cloneValue(e.Args[i].Value)
		}
	}

	if n := len(e.Message); n != 0 {
		m = make([]byte, n)
		copy(m, e.Message)
	}

	if n := len(e.Source); n != 0 {
		s = make([]byte, n)
		copy(s, e.Source)
	}

	return &Event{
		Message: string(m),
		Source:  string(s),
		Args:    a,
		Time:    e.Time,
		Debug:   e.Debug,
	}
}

// Args reprsents a list of event arguments.
type Args []Arg

// Arg represents a single event argument.
type Arg struct {
	Name  string
	Value interface{}
}

// Get returns the value of the argument with name within args.
func (args Args) Get(name string) (v interface{}, ok bool) {
	for _, arg := range args {
		if arg.Name == name {
			v, ok = arg.Value, true
			break
		}
	}
	return
}

// Map converts an argument list to a map representation. In cases where the
// list contains multiple arguments with the same name the value of the last
// one will be seen in the map.
func (args Args) Map() map[string]interface{} {
	m := make(map[string]interface{})
	for _, arg := range args {
		m[arg.Name] = arg.Value
	}
	return m
}

// A constructs an argument list from a map.
func A(m map[string]interface{}) Args {
	args := make(Args, 0, len(m))
	for name, value := range m {
		args = append(args, Arg{name, value})
	}
	return args
}

// SortArgs sorts a list of argument by their argument names.
//
// This is not a stable sorting operation, elements with equal values may not be
// in the same order they were originally after the function returns.
func SortArgs(args Args) {
	sort.Sort(byArgName(args))
}

type byArgName []Arg

func (a byArgName) Len() int {
	return len(a)
}

func (a byArgName) Less(i int, j int) bool {
	return a[i].Name < a[j].Name
}

func (a byArgName) Swap(i int, j int) {
	a[i], a[j] = a[j], a[i]
}
