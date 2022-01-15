package gah

import (
	"errors"
	"fmt"
)

var ErrExpectedSubcommand = errors.New("expected subcommand")

type ErrInvalidSubcommand struct {
	subcommand string
}

func (e *ErrInvalidSubcommand) Error() string {
	return fmt.Sprintf("invalid subcommand %s", e.subcommand)
}

type ErrUnexpectedFlag struct {
	flag string
}

func (e *ErrUnexpectedFlag) Error() string {
	return fmt.Sprintf("unexpected flag %s", e.flag)
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

type ErrUnmarshallingArgument struct {
	name  string
	value string
	error error
}

func (e *ErrUnmarshallingArgument) Error() string {
	return fmt.Sprintf("error unmarshalling argument for %s %s: %v",
		e.name, e.value, e.error)
}
