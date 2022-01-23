package gah

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"mtoohey.com/gah/unmarshal"
)

// TODO: test subcommand evaulation

var simpleVersionedCmd = Cmd{
	Version:  "v0.0.0",
	Function: func(f struct{}, a struct{}) {},
}

func TestNoArgs(t *testing.T) {
	assert.NoError(t, simpleUnversionedCmd.Eval([]string{""}, nil))
}

func TestHelp(t *testing.T) {
	assert.NoError(t, simpleVersionedCmd.Eval([]string{"", "-h"}, nil))
	assert.NoError(t, simpleVersionedCmd.Eval(
		[]string{"", "-h", "extra", "ignored", "args"}, nil))
	assert.NoError(t, simpleVersionedCmd.Eval([]string{"", "--help"}, nil))
	assert.NoError(t, simpleVersionedCmd.Eval(
		[]string{"", "--help", "extra", "ignored", "args"}, nil))
	assert.ErrorIs(t, simpleVersionedCmd.Eval([]string{"", "help"}, nil),
		&ErrUnexpectedArgument{})
}

func TestVersionSuccess(t *testing.T) {
	assert.NoError(t, simpleVersionedCmd.Eval([]string{"", "-v"}, nil))
	assert.NoError(t, simpleVersionedCmd.Eval(
		[]string{"", "-v", "extra", "ignored", "args"}, nil))
	assert.NoError(t, simpleVersionedCmd.Eval([]string{"", "--version"}, nil))
	assert.NoError(t, simpleVersionedCmd.Eval(
		[]string{"", "--version", "extra", "ignored", "args"}, nil))
}

var simpleUnversionedCmd = Cmd{
	Function: func(f struct{}, a struct{}) {},
}

func TestVersionFailure(t *testing.T) {
	assert.ErrorIs(t, simpleUnversionedCmd.Eval([]string{"", "-v"}, nil),
		&ErrUnexpectedFlag{})
	assert.ErrorIs(t, simpleUnversionedCmd.Eval([]string{"", "-v", "extra", "ignored",
		"args"}, nil), &ErrUnexpectedFlag{})
	assert.ErrorIs(t, simpleUnversionedCmd.Eval([]string{"", "--version"}, nil),
		&ErrUnexpectedFlag{})
	assert.ErrorIs(t, simpleUnversionedCmd.Eval([]string{"", "--version", "extra",
		"ignored", "args"}, nil), &ErrUnexpectedFlag{})
}

func TestFlags(t *testing.T) {
	var test1 bool
	var test2 string

	cmd := Cmd{
		Function: func(f struct {
			Test1 bool   `short:"1"`
			Test2 string `long:"test-two"`
		}, _ struct{}) {
			test1 = f.Test1
			test2 = f.Test2
		},
	}

	expected := "-test-value"

	assert.NoError(t, cmd.Eval([]string{"", "-1", "--test-two", expected}, []string{}))
	assert.True(t, test1)
	assert.Equal(t, test2, expected)
	assert.NoError(t, cmd.Eval([]string{"", "--test-1", "--test-two", expected}, []string{}))
	assert.True(t, test1)
	assert.Equal(t, test2, expected)
	assert.ErrorIs(t, cmd.Eval([]string{"", "--test-two"}, []string{}),
		&ErrExpectedFlagValue{})
	assert.ErrorIs(t, cmd.Eval([]string{"", "--test-2", expected}, []string{}),
		&ErrUnexpectedFlag{})
}

func TestDefaults(t *testing.T) {
	var test1 int
	var test2 string

	cmd := Cmd{
		Function: func(f struct {
			Test1 int    `default:"7"`
			Test2 string `default:"test2"`
		}, _ struct{}) {
			test1 = f.Test1
			test2 = f.Test2
		},
	}

	assert.NoError(t, cmd.Eval([]string{""}, []string{}))
	assert.Equal(t, test1, 7)
	assert.Equal(t, test2, "test2")
}

func TestArgs(t *testing.T) {
	var test1 string
	var test2 []int
	var test3 [3]string

	cmd := Cmd{
		Function: func(_ struct{},
			a struct {
				Test1 string
				Test2 []int `min:"1" max:"3"`
				Test3 [3]string
			}) {
			test1 = a.Test1
			test2 = a.Test2
			test3 = a.Test3
		},
	}

	assert.NoError(t, cmd.Eval([]string{"", "value1", "1", "2", "3", "4", "5", "6"},
		[]string{}))
	assert.Equal(t, test1, "value1")
	assert.Equal(t, test2, []int{1, 2, 3})
	assert.Equal(t, test3, [3]string{"4", "5", "6"})
	assert.NoError(t, cmd.Eval([]string{"", "value1", "1", "2", "3", "4", "5"},
		[]string{}))
	assert.ErrorIs(t, cmd.Eval([]string{"", "value1", "-5", "a", "b", "c"},
		[]string{}), &ErrUnexpectedFlag{})
	assert.NoError(t, cmd.Eval([]string{"", "value1", "--", "-5", "a", "b", "c"},
		[]string{}))
	assert.Equal(t, test1, "value1")
	assert.Equal(t, test2, []int{-5})
	assert.Equal(t, test3, [3]string{"a", "b", "c"})
	assert.ErrorIs(t, cmd.Eval([]string{"", "value1", "a", "b", "c"},
		[]string{}), &ErrUnmarshallingArgument{})
}

func TestCustomUnmarshallers(t *testing.T) {
	var b bool
	var test1 bool
	var test2 bool
	var test3 bool

	cmd := Cmd{
		Function: func(f struct {
			Test1 bool `takesVal:"true"`
			Test2 bool
		}, a struct {
			Test3 bool
		}) {
			test1 = f.Test1
			test2 = f.Test2
			test3 = a.Test3
		},
		CustomValueUnmarshallers: unmarshal.CustomValueUnmarshallers{
			reflect.TypeOf(b): func(s string, _ reflect.StructTag) (reflect.Value, error) {
				return reflect.ValueOf(len(s) == 0), nil
			},
		},
		CustomValuelessUnmarshallers: unmarshal.CustomValuelessUnmarshallers{
			reflect.TypeOf(b): func(_ reflect.Value, _ reflect.StructTag) (reflect.Value, error) {
				return reflect.ValueOf(false), nil
			},
		},
	}

	assert.NoError(t, cmd.Eval([]string{"", "--test-1", "", "asdf"}, nil))
	assert.True(t, test1)
	assert.True(t, !test2)
	assert.True(t, !test3)
}
