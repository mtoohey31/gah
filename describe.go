/*
Package gah (Go Argument Handler) provides ergonomic argument parsing and
validation.
*/
package gah

import "mtoohey.com/gah/unmarshal"

// TODO: first class completion support
// TODO: godocs!

type Cmd struct {
	Name        string
	Aliases     []string
	Author      string
	Version     string
	Description string
	// TODO: restrict the values of this as much as possible with some
	// modification of `interface{ []Cmd | interface{} }`
	Function                     interface{}
	Subcommands                  []Cmd
	DefaultFlags                 interface{}
	CustomValueUnmarshallers     unmarshal.CustomValueUnmarshallers
	CustomValuelessUnmarshallers unmarshal.CustomValuelessUnmarshallers
}
