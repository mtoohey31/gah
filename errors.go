package gah

import "fmt"

type ErrExpectedSubcommand struct{}

func (e *ErrExpectedSubcommand) Error() string {
	return "expected subcommand"
}

func (e *ErrExpectedSubcommand) Is(target error) bool {
	var t interface{} = target
	_, ok := t.(*ErrExpectedSubcommand)
	return ok
}

type ErrInvalidSubcommand struct {
	subcommand string
}

func (e *ErrInvalidSubcommand) Error() string {
	return fmt.Sprintf("invalid subcommand %s", e.subcommand)
}

func (e *ErrInvalidSubcommand) Is(target error) bool {
	var t interface{} = target
	_, ok := t.(*ErrInvalidSubcommand)
	return ok
}

type ErrUnexpectedFlag struct {
	flag string
}

func (e *ErrUnexpectedFlag) Error() string {
	return fmt.Sprintf("unexpected flag %s", e.flag)
}

func (e *ErrUnexpectedFlag) Is(target error) bool {
	var t interface{} = target
	_, ok := t.(*ErrUnexpectedFlag)
	return ok
}

func unexpectedShort(f rune) error {
	return &ErrUnexpectedFlag{flag: string([]rune{'-', f})}
}

func unexpectedLong(f string) error {
	return &ErrUnexpectedFlag{flag: "--" + f}
}

type ErrExpectedFlagValue struct {
	flag string
}

func (e *ErrExpectedFlagValue) Error() string {
	return fmt.Sprintf("expected value for flag %s", e.flag)
}

func (e *ErrExpectedFlagValue) Is(target error) bool {
	var t interface{} = target
	_, ok := t.(*ErrExpectedFlagValue)
	return ok
}

func expectedFlagValueShort(f rune) error {
	return &ErrExpectedFlagValue{flag: string([]rune{'-', f})}
}

func expectedFlagValueLong(f string) error {
	return &ErrExpectedFlagValue{flag: "--" + f}
}

type ErrUnexpectedFlagValue struct {
	flag  string
	value string
}

func (e *ErrUnexpectedFlagValue) Error() string {
	return fmt.Sprintf("unexpected value for flag %s: %s", e.flag, e.value)
}

func (e *ErrUnexpectedFlagValue) Is(target error) bool {
	var t interface{} = target
	_, ok := t.(*ErrUnexpectedFlagValue)
	return ok
}

func unexpectedFlagValueShort(f rune, v string) error {
	return &ErrUnexpectedFlagValue{flag: string([]rune{'-', f}), value: v}
}

func unexpectedFlagValueLong(f string, v string) error {
	return &ErrUnexpectedFlagValue{flag: "--" + f, value: v}
}

type ErrUnmarshallingFlagValue struct {
	flag  string
	error error
}

func (e *ErrUnmarshallingFlagValue) Error() string {
	return fmt.Sprintf("error unmarshalling flag %s: %v", e.flag, e.error)
}

func (e *ErrUnmarshallingFlagValue) Is(target error) bool {
	var t interface{} = target
	_, ok := t.(*ErrUnmarshallingFlagValue)
	return ok
}

func unmarshallingFlagShort(f rune, e error) error {
	return &ErrUnmarshallingFlagValue{flag: string([]rune{'-', f}), error: e}
}

func unmarshallingFlagLong(f string, e error) error {
	return &ErrUnmarshallingFlagValue{flag: "--" + f, error: e}
}

type ErrUnexpectedArgument struct {
	argument string
}

func (e *ErrUnexpectedArgument) Error() string {
	return fmt.Sprintf("unexpected argument %s", e.argument)
}

func (e *ErrUnexpectedArgument) Is(target error) bool {
	var t interface{} = target
	_, ok := t.(*ErrUnexpectedArgument)
	return ok
}

type ErrUnmarshallingArgument struct {
	name  string
	value string
	error error
}

func (e *ErrUnmarshallingArgument) Error() string {
	return fmt.Sprintf("error unmarshalling argument for %s %s: %v",
		e.name, e.value, e.error)
}

func (e *ErrUnmarshallingArgument) Is(target error) bool {
	var t interface{} = target
	_, ok := t.(*ErrUnmarshallingArgument)
	return ok
}

type ErrExpectedArgumentValue struct {
	name string
}

func (e *ErrExpectedArgumentValue) Error() string {
	return fmt.Sprintf("expected value for argument %s", e.name)
}

func (e *ErrExpectedArgumentValue) Is(target error) bool {
	var t interface{} = target
	_, ok := t.(*ErrExpectedArgumentValue)
	return ok
}
