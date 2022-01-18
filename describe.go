/*
Package gah (Go Argument Handler) provides ergonomic argument parsing and
validation.
*/
package gah

// TODO: allow arguments before subcommands that can be passed to all
// subcommands
// TODO: first class completion support
// TODO: provide validator function that people can run in their tests that
// will scan the cmd declaration for errors

type Cmd struct {
	Name        string
	Aliases     []string
	Author      string
	Version     string
	Description string
	// TODO: restrict the values of this as much as possible with some
	// modification of `interface{ []Cmd | interface{} }`
	Content interface{}
}