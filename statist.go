// package statist provides interface Statist, which entities can implement in order
// to report their state.
//
// When implementing State() and SetState(), the time.Time
// parameter should be thought of as the time when the state was set.
//
// A caller to State() will interpret the output as
// the returned state 'since' the returned time
//
// Additionally, some functions allow creating a Lineup (a slice of Statists),
// pushing or popping Statists from the registry, and a function Muster which returns
// a string to be used (and embellished as needed) for reporting the state of
// all members of a Lineup
//
// A use-case for package is to periodically report sensor information on a schedule.
// By making your sensors Statists, just Enlist them into a Lineup and make them sound off over MQTT.
package statist

import (
	"strings"
)

type Statist interface {
	StateString() string
	Name() string
}

type Lineup []Statist

type Musterer interface {
	Muster() string
	MusterWithGreeting(string) string
}

// NewLineup creates a Statist slice (a Lineup) and returns it
func NewLineup() Lineup {
	statists := make([]Statist, 0, 10)
	return statists
}

// Enlist pushes a Statist into a Lineup and returns the new Lineup
func Enlist(s Statist, l Lineup) Lineup {
	l = append(l, s)
	return l
}

// Desert will remove a Statist (by 'Name()') from a Registry and returns the new registry, or return existing if no match
// Warning: Desert merely removes the first index matching s.Name() and does not check subsequent indicies
// so unique names are encouraged yet unenforced
func Desert(s Statist, l Lineup) Lineup {
	for i, v := range l {
		if v.Name() == s.Name() {
			return append(l[0:i], l[i+1:]...)
		}
	}
	return l
}

// MusterWithGreeting is a sample implementation that sets a greeting (eg, the date/time when the muster was called)
// and returns a multliline string of the StateString from each Statist in a Lineup;
// this is cleaner when StateString() is implemented with care
func (l Lineup) MusterWithGreeting(g string) string {
	s := strings.Builder{}
	s.Grow(1024)
	s.WriteString(g)
	s.WriteByte(NewLine())
	for _, v := range l {
		s.WriteString(v.StateString())
		s.WriteByte(NewLine())
	}
	return s.String()
}

// Muster does the same as MusterWithGreeting but sans greeting
func (l Lineup) Muster() string {
	s := strings.Builder{}
	s.Grow(1024)
	for _, v := range l {
		s.WriteString(v.StateString())
		s.WriteByte(NewLine())
	}
	return s.String()
}

/*
	Helpers for formatting within a Muster
*/

// newLine returns a line feed ascii byte
func NewLine() byte {
	return byte(10)
}

// tab returns a horizontal tab ascii byte
func Tab() byte {
	return byte(9)
}

// btc returns the emoji for the bitcoin logo lol
func Btc() rune {
	// return string("\uf15a")
	return rune('\u20bf')
}

// checkMark returns the emoji for a check mark
func CheckMark() rune {
	return rune('\u2713')
}

// xMark returns the emoji for an 'x' mark
func X() rune {
	return rune('\u2715')
}
