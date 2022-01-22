package validate

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"mtoohey.com/gah"
)

func TestRecursiveValidation(t *testing.T) {
	cmd := gah.Cmd{
		Subcommands: []gah.Cmd{
			{
				Function: func() {},
			},
		},
	}
	assert.ErrorIs(t, Validate(cmd, true), &ErrFunctionTakesNonTwoArgs{})
}

func TestValidateFunctionIsFunction(t *testing.T) {
	cmd := gah.Cmd{
		Function: "asdf",
	}
	assert.ErrorIs(t, Validate(cmd, true), &ErrFunctionIsNotFunction{})
}

func TestValidateFunctionTakesTwoArgs(t *testing.T) {
	cmd := gah.Cmd{
		Function: func() {},
	}
	assert.ErrorIs(t, Validate(cmd, true), &ErrFunctionTakesNonTwoArgs{})
	cmd = gah.Cmd{
		Function: func(_ struct{}, _ struct{}, _ struct{}) {},
	}
	assert.ErrorIs(t, Validate(cmd, true), &ErrFunctionTakesNonTwoArgs{})
}

func TestValidateFunctionTakesStructArgs(t *testing.T) {
	cmd := gah.Cmd{
		Function: func(_ string, _ struct{}) {},
	}
	assert.ErrorIs(t, Validate(cmd, true), &ErrFunctionTakesNonStructArg{})
	cmd = gah.Cmd{
		Function: func(_ struct{}, _ string) {},
	}
	assert.ErrorIs(t, Validate(cmd, true), &ErrFunctionTakesNonStructArg{})
}

func TestValidateNoFailingParams(t *testing.T) {
	cmd := gah.Cmd{
		Function: func(f struct {
			Test bool `takesVal:"not a bool"`
		}, _ struct{}) {
		},
	}
	assert.ErrorIs(t, Validate(cmd, true), &ErrFailingParam{})
	cmd = gah.Cmd{
		Function: func(f struct {
			Test uint64 `minVal:"-32.5"`
		}, _ struct{}) {
		},
	}
	assert.ErrorIs(t, Validate(cmd, true), &ErrFailingParam{})
	cmd = gah.Cmd{
		Function: func(f struct {
			Test uint64 `maxVal:"-32.5"`
		}, _ struct{}) {
		},
	}
	assert.ErrorIs(t, Validate(cmd, true), &ErrFailingParam{})
	cmd = gah.Cmd{
		Function: func(_ struct{},
			a struct {
				Test []int `min:"1.1"`
			}) {
		},
	}
	assert.ErrorIs(t, Validate(cmd, true), &ErrFailingParam{})
	cmd = gah.Cmd{
		Function: func(_ struct{},
			a struct {
				Test []int `max:"-"`
			}) {
		},
	}
	assert.ErrorIs(t, Validate(cmd, true), &ErrFailingParam{})
	cmd = gah.Cmd{
		Function: func(_ struct{},
			a struct {
				Test uint8 `minVal:"-2"`
			}) {
		},
	}
	assert.ErrorIs(t, Validate(cmd, true), &ErrFailingParam{})
	cmd = gah.Cmd{
		Function: func(_ struct{},
			a struct {
				Test int8 `maxVal:"700"`
			}) {
		},
	}
	assert.ErrorIs(t, Validate(cmd, true), &ErrFailingParam{})
}

func TestValidateValueUnmarshallers(t *testing.T) {
	cmd := gah.Cmd{
		Function: func(f struct {
			Test gah.Cmd
		}, _ struct{}) {
		},
	}
	assert.ErrorIs(t, Validate(cmd, true), &ErrMissingValueUnmarshaller{})

	cmd = gah.Cmd{
		Function: func(_ struct{},
			a struct {
				Test gah.Cmd
			}) {
		},
	}
	assert.ErrorIs(t, Validate(cmd, true), &ErrMissingValueUnmarshaller{})
}

func TestValidateValuelessUnmarshallers(t *testing.T) {
	cmd := gah.Cmd{
		Function: func(f struct {
			Test string `takesVal:"false"`
		}, _ struct{}) {
		},
	}
	assert.ErrorIs(t, Validate(cmd, true), &ErrMissingValuelessUnmarshaller{})
}

func TestValidateSubcommandArgsOnCorrectType(t *testing.T) {
	cmd := gah.Cmd{
		Function: func(_ struct{}, a struct {
			SubcommandArgs []int `subcommandArgs:""`
		}) {
		},
	}
	assert.ErrorIs(t, Validate(cmd, true), &ErrSubcommandArgsOnIncorrectType{})
}

func TestValidateNoEmptyShortFlags(t *testing.T) {
	cmd := gah.Cmd{
		Function: func(f struct {
			Test1 bool `short:""`
		}, _ struct{}) {
		},
	}
	assert.ErrorIs(t, Validate(cmd, true), &ErrEmptyShortFlag{})
}

func TestValidateNoEmptyLongFlags(t *testing.T) {
	cmd := gah.Cmd{
		Function: func(f struct {
			Test1 bool `long:""`
		}, _ struct{}) {
		},
	}
	assert.ErrorIs(t, Validate(cmd, true), &ErrEmptyLongFlag{})
}

func TestValidateNoMultiRuneShortFlags(t *testing.T) {
	cmd := gah.Cmd{
		Function: func(f struct {
			Test1 bool `short:"tt"`
		}, _ struct{}) {
		},
	}
	assert.ErrorIs(t, Validate(cmd, true), &ErrMultiRuneShortFlag{})
}

func TestValidateNoConflictingShortFlags(t *testing.T) {
	cmd := gah.Cmd{
		Function: func(f struct {
			Test1 bool   `short:"t"`
			Test2 string `short:"t"`
		}, _ struct{}) {
		},
	}
	assert.ErrorIs(t, Validate(cmd, true), &ErrConflictingShortFlags{})
}

func TestValidateNoConflictingLongFlags(t *testing.T) {
	cmd := gah.Cmd{
		Function: func(f struct {
			Test1 bool   `long:"test"`
			Test2 string `long:"test"`
		}, _ struct{}) {
		},
	}
	assert.ErrorIs(t, Validate(cmd, true), &ErrConflictingLongFlags{})
	cmd = gah.Cmd{
		Function: func(f struct {
			Test1 bool `long:"test-2"`
			Test2 string
		}, _ struct{}) {
		},
	}
	assert.ErrorIs(t, Validate(cmd, true), &ErrConflictingLongFlags{})
}

func TestValidateNoConflictingSubcommands(t *testing.T) {
	cmd := gah.Cmd{
		Subcommands: []gah.Cmd{
			{
				Name: "test",
			},
			{
				Name: "test",
			},
		},
	}
	assert.ErrorIs(t, Validate(cmd, true), &ErrConflictingSubcommands{})
	cmd = gah.Cmd{
		Subcommands: []gah.Cmd{
			{
				Name:    "test-1",
				Aliases: []string{"test"},
			},
			{
				Name:    "test-2",
				Aliases: []string{"test"},
			},
		},
	}
	assert.ErrorIs(t, Validate(cmd, true), &ErrConflictingSubcommands{})
	cmd = gah.Cmd{
		Subcommands: []gah.Cmd{
			{
				Name: "test",
			},
			{
				Name:    "test-2",
				Aliases: []string{"test"},
			},
		},
	}
	assert.ErrorIs(t, Validate(cmd, true), &ErrConflictingSubcommands{})
	cmd = gah.Cmd{
		Subcommands: []gah.Cmd{
			{
				Name:    "test",
				Aliases: []string{"test"},
			},
		},
	}
	assert.ErrorIs(t, Validate(cmd, true), &ErrConflictingSubcommands{})
	cmd = gah.Cmd{
		Subcommands: []gah.Cmd{
			{
				Name:    "test1",
				Aliases: []string{"test2", "test2"},
			},
		},
	}
	assert.ErrorIs(t, Validate(cmd, true), &ErrConflictingSubcommands{})
}

func TestValidateNoFailingDefaults(t *testing.T) {
	cmd := gah.Cmd{
		Function: func(f struct {
			Test int `default:"not a number"`
		}, _ struct{}) {
		},
	}
	assert.ErrorIs(t, Validate(cmd, true), &ErrFailingDefault{})
}

func TestValidateOneOrFewerVariableArguments(t *testing.T) {
	cmd := gah.Cmd{
		Function: func(_ struct{}, a struct {
			Normal1   string
			Variable1 []string
			Normal2   int
			Variable2 []int
		}) {
		},
	}
	assert.ErrorIs(t, Validate(cmd, true), &ErrMultipleVariableArguments{})

	cmd = gah.Cmd{
		Function: func(_ struct{}, a struct {
			Normal1   string
			Variable1 []int `min:"3" max:"4"`
			Normal2   int
			Variable2 []string `max:"1"`
		}) {
		},
	}
	assert.ErrorIs(t, Validate(cmd, true), &ErrMultipleVariableArguments{})
}

func TestValidateNoArgsAndSubcommands(t *testing.T) {
	cmd := gah.Cmd{
		Function: func(_ struct{}, a struct {
			Arg string
		}) {
		},
		Subcommands: []gah.Cmd{},
	}
	assert.ErrorIs(t, Validate(cmd, true), &ErrArgsAndSubcommands{})
}
