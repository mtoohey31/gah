package validate

import (
	"reflect"
	"strconv"
	"testing"
	"unicode"
	"unicode/utf8"

	"mtoohey.com/gah"
	"mtoohey.com/gah/unmarshal"
)

func ValidateTest(c gah.Cmd, recursive bool, t *testing.T) {
	err := Validate(c, recursive)
	if err != nil {
		t.Fatal(err)
	}
}

func Validate(c gah.Cmd, recursive bool) error {
	if c.Function != nil {
		for _, v := range functionValidators {
			err := v(c)
			if err != nil {
				return err
			}
		}
	}

	if c.Subcommands != nil {
		for _, v := range subcommandValidators {
			err := v(c)
			if err != nil {
				return err
			}
		}

		if recursive {
			for _, subcommand := range c.Subcommands {
				err := Validate(subcommand, recursive)
				if err != nil {
					return err
				}
			}
		}
	}

	for _, v := range universalValidators {
		err := v(c)
		if err != nil {
			return err
		}
	}

	return nil
}

var subcommandValidators = []func(gah.Cmd) error{
	validateNoConflictingSubcommands,
}

var functionValidators = []func(gah.Cmd) error{
	validateFunctionIsFunction,
	validateFunctionTakesTwoArgs,
	validateFunctionTakesStructArgs,
	validateNoFailingParams,
	validateValueUmarshallers,
	validateValuelessUmarshallers,
	validateSubcommandArgsOnCorrectType,
	validateNoEmptyShortFlags,
	validateNoEmptyLongFlags,
	validateNoMultiRuneShortFlags,
	validateNoConflictingShortFlags,
	validateNoConflictingLongFlags,
	validateNoFailingDefaults,
	validateDynamicDefaultFlagsType,
	validateOneOrFewerVariableArguments,
}

var universalValidators = []func(gah.Cmd) error{
	validateNoArgsAndSubcommands,
}

func validateFunctionIsFunction(c gah.Cmd) error {
	functionKind := reflect.TypeOf(c.Function).Kind()

	if functionKind != reflect.Func {
		return &ErrFunctionIsNotFunction{functionKind: functionKind}
	}

	return nil
}

func validateFunctionTakesTwoArgs(c gah.Cmd) error {
	numIn := reflect.TypeOf(c.Function).NumIn()

	if numIn != 2 {
		return &ErrFunctionTakesNonTwoArgs{numFunctionArgs: numIn}
	}

	return nil
}

func validateFunctionTakesStructArgs(c gah.Cmd) error {
	inZeroKind := reflect.TypeOf(c.Function).In(0).Kind()

	if inZeroKind != reflect.Struct {
		return &ErrFunctionTakesNonStructArg{argumentIndex: 0, argumentKind: inZeroKind}
	}

	inOneKind := reflect.TypeOf(c.Function).In(1).Kind()

	if inOneKind != reflect.Struct {
		return &ErrFunctionTakesNonStructArg{argumentIndex: 1, argumentKind: inOneKind}
	}

	return nil
}

func validateNoFailingParams(c gah.Cmd) error {
	functionType := reflect.TypeOf(c.Function)

	for _, field := range reflect.VisibleFields(functionType.In(0)) {
		takesVal, found := field.Tag.Lookup("takesVal")
		if found {
			_, err := strconv.ParseBool(takesVal)
			if err != nil {
				return &ErrFailingParam{paramName: "takesVal", paramString: takesVal,
					flagName: field.Name, error: err}
			}
		}

		minVal, found := field.Tag.Lookup("minVal")
		if found {
			_, err := unmarshal.GetValueUnmarshaller(field.Type,
				c.CustomValueUnmarshallers)(minVal, "")
			if err != nil {
				return &ErrFailingParam{paramName: "minVal", paramString: minVal,
					flagName: field.Name, error: err}
			}
		}

		maxVal, found := field.Tag.Lookup("maxVal")
		if found {
			_, err := unmarshal.GetValueUnmarshaller(field.Type,
				c.CustomValueUnmarshallers)(maxVal, "")
			if err != nil {
				return &ErrFailingParam{paramName: "maxVal", paramString: maxVal,
					flagName: field.Name, error: err}
			}
		}
	}

	for _, field := range reflect.VisibleFields(functionType.In(1)) {
		min, found := field.Tag.Lookup("min")
		if found {
			_, err := strconv.Atoi(min)
			if err != nil {
				return &ErrFailingParam{paramName: "min", paramString: min,
					flagName: field.Name, error: err}
			}
		}

		max, found := field.Tag.Lookup("max")
		if found {
			_, err := strconv.Atoi(max)
			if err != nil {
				return &ErrFailingParam{paramName: "max", paramString: max,
					flagName: field.Name, error: err}
			}
		}

		minVal, found := field.Tag.Lookup("minVal")
		if found {
			_, err := unmarshal.GetValueUnmarshaller(field.Type,
				c.CustomValueUnmarshallers)(minVal, "")
			if err != nil {
				return &ErrFailingParam{paramName: "minVal", paramString: minVal,
					flagName: field.Name, error: err}
			}
		}

		maxVal, found := field.Tag.Lookup("maxVal")
		if found {
			_, err := unmarshal.GetValueUnmarshaller(field.Type,
				c.CustomValueUnmarshallers)(maxVal, "")
			if err != nil {
				return &ErrFailingParam{paramName: "maxVal", paramString: maxVal,
					flagName: field.Name, error: err}
			}
		}
	}

	return nil
}

func validateValueUmarshallers(c gah.Cmd) (err error) {
	functionType := reflect.TypeOf(c.Function)

	var currentValueType reflect.Type

	defer func() {
		if r := recover(); r != nil {
			err = &ErrMissingValueUnmarshaller{valueType: currentValueType}
		}
	}()

	for _, field := range reflect.VisibleFields(functionType.In(0)) {
		if unmarshal.TakesValue(field) {
			currentValueType = field.Type
			unmarshal.GetValueUnmarshaller(field.Type, nil)
		}
	}

	for _, field := range reflect.VisibleFields(functionType.In(1)) {
		if unmarshal.TakesValue(field) {
			switch field.Type.Kind() {
			case reflect.Slice:
				currentValueType = field.Type.Elem()
			case reflect.Array:
				currentValueType = field.Type.Elem()
			default:
				currentValueType = field.Type
			}
			unmarshal.GetValueUnmarshaller(currentValueType, nil)
		}
	}

	return nil
}

func validateValuelessUmarshallers(c gah.Cmd) (err error) {
	var currentValueType reflect.Type

	defer func() {
		if r := recover(); r != nil {
			err = &ErrMissingValuelessUnmarshaller{valueType: currentValueType}
		}
	}()

	for _, field := range reflect.VisibleFields(reflect.TypeOf(c.Function).In(0)) {
		if !unmarshal.TakesValue(field) {
			currentValueType = field.Type
			unmarshal.GetValuelessUnmarshaller(field.Type, nil)
		}
	}

	return nil
}

func validateSubcommandArgsOnCorrectType(c gah.Cmd) error {
	for _, field := range reflect.VisibleFields(reflect.TypeOf(c.Function).In(1)) {
		_, found := field.Tag.Lookup("subcommandArgs")
		if found && field.Type != reflect.TypeOf([]string{}) {
			return &ErrSubcommandArgsOnIncorrectType{}
		}
	}

	return nil
}

func validateNoEmptyShortFlags(c gah.Cmd) error {
	for _, field := range reflect.VisibleFields(reflect.TypeOf(c.Function).In(0)) {
		short, found := field.Tag.Lookup("short")
		if found {
			if utf8.RuneCountInString(short) == 0 {
				return &ErrEmptyShortFlag{flagName: field.Name}
			}
		}
	}

	return nil
}

func validateNoEmptyLongFlags(c gah.Cmd) error {
	for _, field := range reflect.VisibleFields(reflect.TypeOf(c.Function).In(0)) {
		long, found := field.Tag.Lookup("long")
		if found {
			if utf8.RuneCountInString(long) == 0 {
				return &ErrEmptyLongFlag{flagName: field.Name}
			}
		}
	}

	return nil
}

func validateNoMultiRuneShortFlags(c gah.Cmd) error {
	for _, field := range reflect.VisibleFields(reflect.TypeOf(c.Function).In(0)) {
		short, found := field.Tag.Lookup("short")
		if found {
			if utf8.RuneCountInString(short) > 1 {
				return &ErrMultiRuneShortFlag{flagName: field.Name, shortFlag: short}
			}
		}
	}

	return nil
}

func validateNoConflictingShortFlags(c gah.Cmd) error {
	var shortSoFar [][2]string

	for _, field := range reflect.VisibleFields(reflect.TypeOf(c.Function).In(0)) {
		short, found := field.Tag.Lookup("short")
		if found {
			for _, otherShort := range shortSoFar {
				if short == otherShort[0] {
					return &ErrConflictingShortFlags{flagNames: []string{
						otherShort[1], field.Name}}
				}
			}
		}

		shortSoFar = append(shortSoFar, [2]string{short, field.Name})
	}

	return nil
}

func pascalToKebab(s string) string {
	if len(s) == 0 {
		return ""
	}

	runes := []rune(s)
	res := []rune{unicode.ToLower(runes[0])}
	runes = runes[1:]

	for _, r := range runes {
		if unicode.IsUpper(r) {
			res = append(res, '-', unicode.ToLower(r))
		} else if unicode.IsDigit(r) {
			res = append(res, '-', r)
		} else {
			res = append(res, r)
		}
	}

	return string(res)
}

func validateNoConflictingLongFlags(c gah.Cmd) error {
	var longSoFar [][2]string

	for _, field := range reflect.VisibleFields(reflect.TypeOf(c.Function).In(0)) {
		long, found := field.Tag.Lookup("long")
		if !found {
			long = pascalToKebab(field.Name)
		}

		for _, otherLong := range longSoFar {
			if long == otherLong[0] {
				return &ErrConflictingLongFlags{flagNames: []string{
					otherLong[1], field.Name}}
			}
		}

		longSoFar = append(longSoFar, [2]string{long, field.Name})
	}

	return nil
}

func validateNoConflictingSubcommands(c gah.Cmd) error {
	var namesSoFar [][2]string

	for _, subcommand := range c.Subcommands {
		for _, otherName := range namesSoFar {
			if subcommand.Name == otherName[0] {
				return &ErrConflictingSubcommands{subcommandNames: []string{
					subcommand.Name, otherName[1]}, aliasOrName: subcommand.Name}
			}
		}

		namesSoFar = append(namesSoFar, [2]string{subcommand.Name, subcommand.Name})

		for _, alias := range subcommand.Aliases {
			for _, otherName := range namesSoFar {
				if alias == otherName[0] {
					return &ErrConflictingSubcommands{subcommandNames: []string{
						alias, otherName[1]}, aliasOrName: alias}
				}
			}

			namesSoFar = append(namesSoFar, [2]string{alias, subcommand.Name})
		}
	}

	return nil
}

func validateNoFailingDefaults(c gah.Cmd) error {
	for _, field := range reflect.VisibleFields(reflect.TypeOf(c.Function).In(0)) {
		defaultStr, found := field.Tag.Lookup("default")
		if found {
			_, err := unmarshal.GetValueUnmarshaller(field.Type,
				c.CustomValueUnmarshallers)(defaultStr, field.Tag)
			if err != nil {
				return &ErrFailingDefault{defaultString: defaultStr,
					flagName: field.Name, error: err}
			}
		}
	}

	return nil
}
func validateDynamicDefaultFlagsType(c gah.Cmd) error {
	if c.DefaultFlags == nil {
		return nil
	}

	flagType := reflect.TypeOf(c.Function).In(0)
	dynamicDefaultFlagsType := reflect.TypeOf(c.DefaultFlags)

	if flagType != dynamicDefaultFlagsType {
		return &ErrMismatchedDynamicDefaultFlags{
			actual:   dynamicDefaultFlagsType,
			expected: flagType,
		}
	}

	return nil
}

func validateOneOrFewerVariableArguments(c gah.Cmd) error {
	var variableSoFar []string

	for _, field := range reflect.VisibleFields(reflect.TypeOf(c.Function).In(1)) {
		if field.Type.Kind() != reflect.Slice {
			continue
		}

		min := 0
		minStr, found := field.Tag.Lookup("minVal")
		if found {
			var err error
			min, err = strconv.Atoi(minStr)
			if err != nil {
				panic(err)
			}
		}

		max := ^int(0)
		maxStr, found := field.Tag.Lookup("maxVal")
		if found {
			var err error
			min, err = strconv.Atoi(maxStr)
			if err != nil {
				panic(err)
			}
		}

		if min != max {
			variableSoFar = append(variableSoFar, field.Name)
		}

		if len(variableSoFar) > 1 {
			return &ErrMultipleVariableArguments{argumentNames: variableSoFar}
		}
	}

	return nil
}

func validateNoArgsAndSubcommands(c gah.Cmd) error {
	if len(reflect.VisibleFields(reflect.TypeOf(c.Function).In(1))) != 0 && c.Subcommands != nil {
		return &ErrArgsAndSubcommands{}
	}

	return nil
}
