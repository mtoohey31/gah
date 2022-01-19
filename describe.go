/*
Package gah (Go Argument Handler) provides ergonomic argument parsing and
validation.
*/
package gah

import "mtoohey.com/gah/unmarshal"

// TODO: first class completion support
// TODO: provide validator function that people can run in their tests that
// will scan the cmd declaration for errors
// TODO: godocs!
// TODO: support default subcommands

type Cmd struct {
	Name        string
	Aliases     []string
	Author      string
	Version     string
	Description string
	// TODO: restrict the values of this as much as possible with some
	// modification of `interface{ []Cmd | interface{} }`
	Content                      interface{}
	CustomValueUnmarshallers     unmarshal.CustomValueUnmarshallers
	CustomValuelessUnmarshallers unmarshal.CustomValuelessUnmarshallers
}
