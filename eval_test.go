package gah

import (
	"testing"

	"mtoohey.com/gah/test"
)

var simpleVersionedCmd = Cmd{
	Name:    "simple-test",
	Version: "v0.0.0",
	Content: func(f struct{}, a struct{}) {},
}

func TestNoArgs(t *testing.T) {
	err := simpleVersionedCmd.Eval([]string{""}, nil)
	test.AssertNil(err, t)
}

func TestHelp(t *testing.T) {
	test.AssertNil(simpleVersionedCmd.Eval([]string{"", "-h"}, nil), t)
	test.AssertNil(simpleVersionedCmd.Eval([]string{"", "-h", "extra", "ignored", "args"}, nil), t)
	test.AssertNil(simpleVersionedCmd.Eval([]string{"", "--help"}, nil), t)
	test.AssertNil(simpleVersionedCmd.Eval([]string{"", "--help", "extra", "ignored", "args"}, nil), t)
	test.AssertNil(simpleVersionedCmd.Eval([]string{"", "help"}, nil), t)
}

func TestVersionSuccess(t *testing.T) {
	test.AssertNil(simpleVersionedCmd.Eval([]string{"", "-v"}, nil), t)
	test.AssertNil(simpleVersionedCmd.Eval([]string{"", "-v", "extra", "ignored", "args"}, nil), t)
	test.AssertNil(simpleVersionedCmd.Eval([]string{"", "--version"}, nil), t)
	test.AssertNil(simpleVersionedCmd.Eval([]string{"", "--version", "extra", "ignored", "args"}, nil), t)
}

var simpleUnversionedCmd = Cmd{
	Name:    "simple-test",
	Content: func(f struct{}, a struct{}) {},
}

func TestVersionFailure(t *testing.T) {
	test.AssertErrIs(simpleUnversionedCmd.Eval([]string{"", "-v"}, nil),
		&ErrUnexpectedFlag{}, t)
	test.AssertErrIs(simpleUnversionedCmd.Eval([]string{"", "-v", "extra", "ignored",
		"args"}, nil), &ErrUnexpectedFlag{}, t)
	test.AssertErrIs(simpleUnversionedCmd.Eval([]string{"", "--version"}, nil),
		&ErrUnexpectedFlag{}, t)
	test.AssertErrIs(simpleUnversionedCmd.Eval([]string{"", "--version", "extra",
		"ignored", "args"}, nil), &ErrUnexpectedFlag{}, t)
}

func TestFlags(t *testing.T) {
	var test1 bool
	var test2 string

	cmd := Cmd{
		Name: "flag-only-test",
		Content: func(f struct {
			Test1 bool   `short:"1"`
			Test2 string `long:"test-two"`
		}, _ struct{}) {
			test1 = f.Test1
			test2 = f.Test2
		},
	}

	expected := "-test-value"

	test.AssertNil(cmd.Eval([]string{"", "-1", "--test-two", expected}, []string{}), t)
	test.Assert(test1, t)
	test.AssertEq(test2, expected, t)
	test.AssertNil(cmd.Eval([]string{"", "--test-1", "--test-two", expected}, []string{}), t)
	test.Assert(test1, t)
	test.AssertEq(test2, expected, t)
	test.AssertErrIs(cmd.Eval([]string{"", "--test-two"}, []string{}), &ErrExpectedFlagValue{}, t)
	test.AssertErrIs(cmd.Eval([]string{"", "--test-2", expected}, []string{}), &ErrUnexpectedFlag{}, t)
}

func TestArgs(t *testing.T) {
	var test1 string
	var test2 []int
	var test3 [3]string

	cmd := Cmd{
		Name: "args-only-test",
		Content: func(_ struct{},
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

	test.AssertNil(cmd.Eval([]string{"", "value1", "1", "2", "3", "4", "5", "6"},
		[]string{}), t)
	test.AssertEq(test1, "value1", t)
	test.AssertDeepEq(test2, []int{1, 2, 3}, t)
	test.AssertDeepEq(test3, [3]string{"4", "5", "6"}, t)
	test.AssertErrIs(cmd.Eval([]string{"", "value1", "1", "2", "3", "4", "5"},
		[]string{}), &ErrExpectedArgumentValue{}, t)
	test.AssertErrIs(cmd.Eval([]string{"", "value1", "-5", "a", "b", "c"},
		[]string{}), &ErrUnexpectedFlag{}, t)
	test.AssertNil(cmd.Eval([]string{"", "value1", "--", "-5", "a", "b", "c"},
		[]string{}), t)
	test.AssertEq(test1, "value1", t)
	test.AssertDeepEq(test2, []int{-5}, t)
	test.AssertDeepEq(test3, [3]string{"a", "b", "c"}, t)
	test.AssertErrIs(cmd.Eval([]string{"", "value1", "a", "b", "c"},
		[]string{}), &ErrUnmarshallingArgument{}, t)
}