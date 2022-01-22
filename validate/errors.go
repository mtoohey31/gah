package validate

import (
	"fmt"
	"reflect"
	"strings"
)

type ErrFunctionIsNotFunction struct {
	functionKind reflect.Kind
}

func (e *ErrFunctionIsNotFunction) Error() string {
	return fmt.Sprintf("provided function is not a function, found kind %v",
		e.functionKind)
}

func (e *ErrFunctionIsNotFunction) Is(target error) bool {
	var t interface{} = target
	_, ok := t.(*ErrFunctionIsNotFunction)
	return ok
}

type ErrFunctionTakesNonTwoArgs struct {
	numFunctionArgs int
}

func (e *ErrFunctionTakesNonTwoArgs) Error() string {
	return fmt.Sprintf("provided function takes the wrong number of args: %d, should take 2",
		e.numFunctionArgs)
}

func (e *ErrFunctionTakesNonTwoArgs) Is(target error) bool {
	var t interface{} = target
	_, ok := t.(*ErrFunctionTakesNonTwoArgs)
	return ok
}

type ErrFunctionTakesNonStructArg struct {
	argumentIndex int
	argumentKind  reflect.Kind
}

func (e *ErrFunctionTakesNonStructArg) Error() string {
	return fmt.Sprintf("function argument %d is not of kind struct, found kind: %v",
		e.argumentIndex, e.argumentKind)
}

func (e *ErrFunctionTakesNonStructArg) Is(target error) bool {
	var t interface{} = target
	_, ok := t.(*ErrFunctionTakesNonStructArg)
	return ok
}

type ErrMissingValueUnmarshaller struct {
	valueType reflect.Type
}

func (e *ErrMissingValueUnmarshaller) Error() string {
	return fmt.Sprintf("missing value unmarshaller for type %s, add to CustomValueUnmarshallers",
		e.valueType)
}

func (e *ErrMissingValueUnmarshaller) Is(target error) bool {
	var t interface{} = target
	_, ok := t.(*ErrMissingValueUnmarshaller)
	return ok
}

type ErrMissingValuelessUnmarshaller struct {
	valueType reflect.Type
}

func (e *ErrMissingValuelessUnmarshaller) Error() string {
	return fmt.Sprintf("missing valueless unmarshaller for type %s, add to CustomValuelessUnmarshallers",
		e.valueType)
}

func (e *ErrMissingValuelessUnmarshaller) Is(target error) bool {
	var t interface{} = target
	_, ok := t.(*ErrMissingValuelessUnmarshaller)
	return ok
}

type ErrSubcommandArgsOnIncorrectType struct {
	argType reflect.Type
}

func (e *ErrSubcommandArgsOnIncorrectType) Error() string {
	return fmt.Sprintf("subcommandArgs tag on incorrect type %s, type should be []string",
		e.argType)
}

func (e *ErrSubcommandArgsOnIncorrectType) Is(target error) bool {
	var t interface{} = target
	_, ok := t.(*ErrSubcommandArgsOnIncorrectType)
	return ok
}

type ErrEmptyShortFlag struct {
	flagName string
}

func (e *ErrEmptyShortFlag) Error() string {
	return fmt.Sprintf("empty short flag declared for flag %s", e.flagName)
}

func (e *ErrEmptyShortFlag) Is(target error) bool {
	var t interface{} = target
	_, ok := t.(*ErrEmptyShortFlag)
	return ok
}

type ErrEmptyLongFlag struct {
	flagName string
}

func (e *ErrEmptyLongFlag) Error() string {
	return fmt.Sprintf("empty long flag declared for flag %s", e.flagName)
}

func (e *ErrEmptyLongFlag) Is(target error) bool {
	var t interface{} = target
	_, ok := t.(*ErrEmptyLongFlag)
	return ok
}

type ErrMultiRuneShortFlag struct {
	flagName  string
	shortFlag string
}

func (e *ErrMultiRuneShortFlag) Error() string {
	return fmt.Sprintf("multi rune short flag for %s: %s, should be a single rune",
		e.flagName, e.shortFlag)
}

func (e *ErrMultiRuneShortFlag) Is(target error) bool {
	var t interface{} = target
	_, ok := t.(*ErrMultiRuneShortFlag)
	return ok
}

type ErrConflictingShortFlags struct {
	flagNames []string
}

func (e *ErrConflictingShortFlags) Error() string {
	return fmt.Sprintf("conflicting short flags: %s",
		strings.Join(e.flagNames, ", "))
}

func (e *ErrConflictingShortFlags) Is(target error) bool {
	var t interface{} = target
	_, ok := t.(*ErrConflictingShortFlags)
	return ok
}

type ErrConflictingLongFlags struct {
	flagNames []string
}

func (e *ErrConflictingLongFlags) Error() string {
	return fmt.Sprintf("conflicting long flags: %s",
		strings.Join(e.flagNames, ", "))
}

func (e *ErrConflictingLongFlags) Is(target error) bool {
	var t interface{} = target
	_, ok := t.(*ErrConflictingLongFlags)
	return ok
}

type ErrConflictingSubcommands struct {
	subcommandNames []string
	aliasOrName     string
}

func (e *ErrConflictingSubcommands) Error() string {
	return fmt.Sprintf("conflicting subcommands or aliases %s with: %s",
		strings.Join(e.subcommandNames, ", "), e.aliasOrName)
}

func (e *ErrConflictingSubcommands) Is(target error) bool {
	var t interface{} = target
	_, ok := t.(*ErrConflictingSubcommands)
	return ok
}

type ErrFailingDefault struct {
	defaultString string
	flagName      string
	error         error
}

func (e *ErrFailingDefault) Error() string {
	return fmt.Sprintf("failing default value %s for flag %s with error: %v",
		e.defaultString, e.flagName, e.error)
}

func (e *ErrFailingDefault) Is(target error) bool {
	var t interface{} = target
	_, ok := t.(*ErrFailingDefault)
	return ok
}

// NOTE: this can't validate params in custom unmarshallers, users are
// responsible for that, the only way to test that here would be to convert
// unmarshallers to interfaces or structs, which is an overkill solution

type ErrFailingParam struct {
	paramName   string
	paramString string
	flagName    string
	error       error
}

func (e *ErrFailingParam) Error() string {
	return fmt.Sprintf("failing param %s:\"%s\" for flag %s with error: %v",
		e.paramName, e.paramString, e.flagName, e.error)
}

func (e *ErrFailingParam) Is(target error) bool {
	var t interface{} = target
	_, ok := t.(*ErrFailingParam)
	return ok
}

type ErrMultipleVariableArguments struct {
	argumentNames []string
}

func (e *ErrMultipleVariableArguments) Error() string {
	return fmt.Sprintf("multiple variable arguments: %s, deciding which argument belongs where is ambiguous",
		strings.Join(e.argumentNames, ", "))
}

func (e *ErrMultipleVariableArguments) Is(target error) bool {
	var t interface{} = target
	_, ok := t.(*ErrMultipleVariableArguments)
	return ok
}

type ErrArgsAndSubcommands struct{}

func (e *ErrArgsAndSubcommands) Error() string {
	return "command contains both arguments and subcommands, these are mutually exclusive"
}

func (e *ErrArgsAndSubcommands) Is(target error) bool {
	var t interface{} = target
	_, ok := t.(*ErrArgsAndSubcommands)
	return ok
}
