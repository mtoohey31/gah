package gah

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"mtoohey.com/gah/unmarshal"
)

// TODO: allow custom default values provided in tags
// TODO: allow passing of custom unmarshallers
// TODO: refactor help and version into flags/subcommands so they're treated
// like normal values

func (c Cmd) SimpleEval() {
	err := c.Eval(os.Args, []string{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[31%v\033[0m\n", err)
		os.Exit(1)
	}
}

func (c Cmd) Eval(args []string, parentNames []string) error {
	if hasSubcommads(c) {
		return evalSubcommand(c, args, parentNames)
	} else {
		return evalAndRun(c, args, parentNames)
	}
}

func (c Cmd) EvalMulticall(args []string) {
	wanted := path.Base(args[0])
	for _, subcommand := range c.Content.([]Cmd) {
		if wanted == subcommand.Name {
			subcommand.Eval(args, []string{})
			return
		}
	}
}

func hasSubcommads(c Cmd) bool {
	return reflect.TypeOf(c.Content).AssignableTo(reflect.TypeOf([]Cmd{}))
}

func evalSubcommand(c Cmd, args []string, parentNames []string) error {
	if len(args) < 2 {
		return ErrExpectedSubcommand
	}
	arg := args[1]
	for _, subcommand := range c.Content.([]Cmd) {
		if arg == subcommand.Name {
			return subcommand.Eval(args[1:], append(parentNames, c.Name))
		}

		for _, alias := range subcommand.Aliases {
			if arg == alias {
				return subcommand.Eval(args[1:], append(parentNames, c.Name))
			}
		}
	}

	if arg == "-h" || arg == "--help" || arg == "help" {
		if len(args) > 2 {
			for _, subcommand := range c.Content.([]Cmd) {
				if args[2] == subcommand.Name {
					subcommand.PrintHelp(append(parentNames, c.Name))
					return nil
				}

				for _, alias := range subcommand.Aliases {
					if args[2] == alias {
						subcommand.PrintHelp(append(parentNames, c.Name))
						return nil
					}
				}
			}

			if args[2] == "help" {
				Cmd{Name: "help",
					Description: "Print this help message or the help message of the given subcommand",
					Content:     func(_ struct{}, _ struct{}) {},
				}.PrintHelp(append(parentNames, c.Name))
				return nil
			}
		}

		c.PrintHelp(parentNames)
		return nil
	}

	if arg == "-v" || arg == "--version" {
		if c.Version == "" {
			return &ErrUnexpectedFlag{flag: arg}
		} else {
			println(c.Version)
			return nil
		}
	}

	return &ErrInvalidSubcommand{subcommand: arg}
}

func evalAndRun(c Cmd, inputArgs []string, parentNames []string) error {
	flagsType := reflect.TypeOf(c.Content).In(0)
	flags := reflect.New(flagsType)
	argsType := reflect.TypeOf(c.Content).In(1)
	args := reflect.New(argsType)

	validShort, validLong := getValidFlags(flagsType)
	remainingArgs := getArgs(argsType)

	doubleDashEncountered := false

	for i := 1; i < len(inputArgs); i++ {
		arg := inputArgs[i]

		if strings.HasPrefix(arg, "--") && len(arg) > 2 && !doubleDashEncountered {
			eqIndex := strings.IndexRune(arg, '=')
			var flagName string
			if eqIndex == -1 {
				flagName = arg[2:]
			} else {
				flagName = arg[2:eqIndex]
			}

			field, ok := validLong[flagName]
			if !ok {
				return trySalvageBuiltinLong(c, flagName, parentNames)
			}

			unmarshaller, ok := unmarshal.Unmarshallers[field.Type]
			if !ok {
				panic(fmt.Sprintf("no unmarshaller for flag --%s", flagName))
			}

			if unmarshaller.Type().NumIn() == 1 {
				if eqIndex != -1 {
					return unexpectedFlagValueLong(flagName, arg[eqIndex+1:])
				}

				res := unmarshaller.Call([]reflect.Value{reflect.ValueOf(field.Tag)})
				if !res[1].IsNil() {
					return unmarshallingFlagLong(flagName, res[1].Interface().(error))
				}
				flags.Elem().FieldByIndex(field.Index).Set(res[0])
			} else {
				var flagValue string
				if eqIndex == -1 {
					if i == len(inputArgs)-1 {
						return expectedFlagValueLong(flagName)
					}

					i++
					flagValue = inputArgs[i]
				} else {
					flagValue = arg[eqIndex+1:]
				}

				res := unmarshaller.Call([]reflect.Value{reflect.ValueOf(flagValue),
					reflect.ValueOf(field.Tag)})
				if !res[1].IsNil() {
					return unmarshallingFlagLong(flagName, res[1].Interface().(error))
				}
				flags.Elem().FieldByIndex(field.Index).Set(res[0])
			}
		} else if strings.HasPrefix(arg, "-") && len(arg) > 1 && !doubleDashEncountered {
			if arg == "--" {
				doubleDashEncountered = true
				continue
			}

			eqIndex := strings.IndexRune(arg, '=')
			var flagRunes []rune
			if eqIndex == -1 {
				flagRunes = []rune(arg[1:])
			} else {
				flagRunes = []rune(arg[1:eqIndex])
			}

			for j := 0; j < len(flagRunes); j++ {
				flagRune := flagRunes[j]

				field, ok := validShort[flagRune]
				if !ok {
					return trySalvageBuiltinShort(c, flagRune, parentNames)
				}

				unmarshaller, ok := unmarshal.Unmarshallers[field.Type]
				if !ok {
					panic(fmt.Sprintf("no unmarshaller for flag -%c", flagRune))
				}

				if unmarshaller.Type().NumIn() == 1 {
					if j == len(flagRunes)-1 && eqIndex != -1 {
						return unexpectedFlagValueShort(flagRune, arg[eqIndex+1:])
					}

					res := unmarshaller.Call([]reflect.Value{reflect.ValueOf(field.Tag)})
					if !res[1].IsNil() {
						return unmarshallingFlagShort(flagRune, res[1].Interface().(error))
					}
					flags.Elem().FieldByIndex(field.Index).Set(res[0])
				} else {
					var flagValue string
					if j == len(flagRunes)-1 {
						if eqIndex == -1 {
							if i == len(inputArgs)-1 {
								return expectedFlagValueShort(flagRune)
							}

							i++
							flagValue = inputArgs[i]
						} else {
							flagValue = arg[eqIndex+1:]
						}
					} else {
						flagValue = string(flagRunes[j+1:])
						j = len(flagRunes) - 1
					}

					res := unmarshaller.Call([]reflect.Value{reflect.ValueOf(flagValue),
						reflect.ValueOf(field.Tag)})
					if !res[1].IsNil() {
						return unmarshallingFlagShort(flagRune, res[1].Interface().(error))
					}
					flags.Elem().FieldByIndex(field.Index).Set(res[0])
				}
			}
		} else {
			if len(remainingArgs) == 0 {
				if arg == "help" {
					// NOTE: we don't need to search for subcommand names here because
					// there are no subcommands in this evaulation case
					c.PrintHelp(parentNames)
					return nil
				} else {
					return &ErrUnexpectedArgument{argument: arg}
				}
			}

			if remainingArgs[0].MaxReached(args) {
				remainingArgs = remainingArgs[1:]
				i--
				continue
			}

			res := remainingArgs[0].Unmarshaller().Call(
				[]reflect.Value{reflect.ValueOf(arg),
					reflect.ValueOf(remainingArgs[0].Field().Tag)})
			if !res[1].IsNil() {
				if remainingArgs[0].MinReached(args) {
					remainingArgs = remainingArgs[1:]
					i--
					continue
				}

				return &ErrUnmarshallingArgument{
					name:  strings.ToUpper(remainingArgs[0].Field().Name),
					value: arg, error: res[1].Interface().(error)}
			}

			remainingArgs[0].Update(args, res[0])
		}
	}

	for _, arg := range remainingArgs {
		if !arg.MinReached(args) {
			return &ErrExpectedArgumentValue{name: strings.ToUpper(arg.Field().Name)}
		}
	}

	reflect.ValueOf(c.Content).Call([]reflect.Value{
		reflect.Indirect(flags), reflect.Indirect(args)})

	return nil
}

func trySalvageBuiltinLong(c Cmd, flagName string, parentNames []string) error {
	if flagName == "help" {
		c.PrintHelp(parentNames)
		return nil
	} else if flagName == "version" && c.Version != "" {
		println(c.Version)
		return nil
	} else {
		return unexpectedLong(flagName)
	}
}

func trySalvageBuiltinShort(c Cmd, flagRune rune, parentNames []string) error {
	if flagRune == 'h' {
		c.PrintHelp(parentNames)
		return nil
	} else if flagRune == 'v' && c.Version != "" {
		println(c.Version)
		return nil
	} else {
		return unexpectedShort(flagRune)
	}
}

func getValidFlags(flagsType reflect.Type) (map[rune]reflect.StructField, map[string]reflect.StructField) {
	validShort := make(map[rune]reflect.StructField)
	validLong := make(map[string]reflect.StructField)

	for _, field := range reflect.VisibleFields(flagsType) {
		short, found := field.Tag.Lookup("short")
		if found {
			runes := []rune(short)
			if len(runes) != 1 {
				// TODO: provide error to developer about invalid short flag
			}
			validShort[runes[0]] = field
		}

		long, found := field.Tag.Lookup("long")
		if found {
			validLong[long] = field
		} else {
			validLong[pascalToKebab(field.Name)] = field
		}
	}

	return validShort, validLong
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

type argInfo interface {
	MinReached(reflect.Value) bool
	MaxReached(reflect.Value) bool
	Field() reflect.StructField
	Unmarshaller() reflect.Value
	Update(reflect.Value, reflect.Value)
	Optional() bool
	Multiple() bool
}

type sliceArgInfo struct {
	min   int
	max   int
	field reflect.StructField
}

func (i *sliceArgInfo) MinReached(f reflect.Value) bool {
	return f.Elem().FieldByIndex(i.field.Index).Len() >= i.min
}

func (i *sliceArgInfo) MaxReached(f reflect.Value) bool {
	return f.Elem().FieldByIndex(i.field.Index).Len() >= i.max
}

func (i *sliceArgInfo) Field() reflect.StructField { return i.field }

func (i *sliceArgInfo) Unmarshaller() reflect.Value {
	u, ok := unmarshal.Unmarshallers[i.field.Type.Elem()]
	if ok {
		return u
	} else {
		panic(fmt.Sprintf("no unmarshaller for type %s", i.field.Type.Name()))
	}
}

func (i *sliceArgInfo) Update(f reflect.Value, v reflect.Value) {
	// TODO: performance could be improved here by creating a slice with capacity
	// of max if there is one, or capacity of min if there is no max then
	// appending after it's reached instead of setting the current index which
	// could be tracked inside this struct
	f.Elem().FieldByIndex(i.field.Index).Set(
		reflect.Append(f.Elem().FieldByIndex(i.field.Index), v))
}

func (i *sliceArgInfo) Optional() bool {
	return i.min == 0
}

func (i *sliceArgInfo) Multiple() bool {
	return i.max > 1
}

type arrayArgInfo struct {
	curr  int
	field reflect.StructField
}

func (i *arrayArgInfo) MinReached(v reflect.Value) bool {
	return v.Elem().FieldByIndex(i.field.Index).Len() == i.curr
}

func (i *arrayArgInfo) MaxReached(v reflect.Value) bool {
	return v.Elem().FieldByIndex(i.field.Index).Len() == i.curr
}

func (i *arrayArgInfo) Field() reflect.StructField { return i.field }

func (i *arrayArgInfo) Unmarshaller() reflect.Value {
	u, ok := unmarshal.Unmarshallers[i.field.Type.Elem()]
	if ok {
		return u
	} else {
		panic(fmt.Sprintf("no unmarshaller for type %s", i.field.Type.Name()))
	}
}

func (i *arrayArgInfo) Update(f reflect.Value, v reflect.Value) {
	f.Elem().FieldByIndex(i.field.Index).Index(i.curr).Set(v)
	i.curr++
}

// NOTE: this assumes nobody's passed an array of length 0, we should deal with
// that somewhere
func (i *arrayArgInfo) Optional() bool { return false }

func (i *arrayArgInfo) Multiple() bool { return i.field.Type.Len() > 1 }

type defaultArgInfo struct {
	set   bool
	field reflect.StructField
}

func (i *defaultArgInfo) MinReached(_ reflect.Value) bool {
	return i.set
}

func (i *defaultArgInfo) MaxReached(_ reflect.Value) bool {
	return i.set
}

func (i *defaultArgInfo) Field() reflect.StructField { return i.field }

func (i *defaultArgInfo) Unmarshaller() reflect.Value {
	u, ok := unmarshal.Unmarshallers[i.field.Type]
	if ok {
		return u
	} else {
		panic(fmt.Sprintf("no unmarshaller for type %s", i.field.Type.Name()))
	}
}

func (i *defaultArgInfo) Optional() bool { return false }

func (i *defaultArgInfo) Multiple() bool { return false }

func (i *defaultArgInfo) Update(f reflect.Value, v reflect.Value) {
	f.Elem().FieldByIndex(i.field.Index).Set(v)
	i.set = true
}

func getArgs(argsType reflect.Type) []argInfo {
	var argInfoItems = make([]argInfo, len(reflect.VisibleFields(argsType)))

	for i, field := range reflect.VisibleFields(argsType) {
		switch field.Type.Kind() {
		case reflect.Slice:
			minStr, found := field.Tag.Lookup("min")
			var min int
			if found {
				var err error
				min, err = strconv.Atoi(minStr)
				if err != nil {
					panic(err)
				}
			} else {
				min = 0
			}

			maxStr, found := field.Tag.Lookup("max")
			var max int
			if found {
				var err error
				max, err = strconv.Atoi(maxStr)
				if err != nil {
					panic(err)
				}
			} else {
				max = ^int(0)
			}

			argInfoItems[i] = &sliceArgInfo{min: min, max: max, field: field}
		case reflect.Array:
			argInfoItems[i] = &arrayArgInfo{field: field}
		default:
			argInfoItems[i] = &defaultArgInfo{field: field}
		}
	}

	return argInfoItems
}

func (c Cmd) PrintHelp(parentNames []string) {
	println(strings.Join(append(parentNames, c.Name), "-") + " " + c.Version)
	if c.Author != "" {
		println(c.Author)
	}
	if c.Description != "" {
		println(c.Description)
	}
	if hasSubcommads(c) {
		println("\nUSAGE:\n\t" + strings.Join(append(parentNames, c.Name), " ") + " [SUBCOMMAND]")
		if c.Version == "" {
			println("\nFLAGS:\n\t-h, --help Prints help information")
		} else {
			println("\nFLAGS:\n\t-h, --help    Prints help information\n\t-v, --version Prints version information")
		}
		println("\nSUBCOMMANDS:")
		maxSubcommandNameLength := 0
		for _, subcommand := range c.Content.([]Cmd) {
			l := len(strings.Join(
				append([]string{subcommand.Name}, subcommand.Aliases...), ", "))
			if l > maxSubcommandNameLength {
				maxSubcommandNameLength = l
			}
		}
		for _, subcommand := range c.Content.([]Cmd) {
			s := strings.Join(append([]string{subcommand.Name},
				subcommand.Aliases...), ", ")
			l := len(s)
			println("\t" + s + strings.Repeat(" ",
				1+maxSubcommandNameLength-l) + subcommand.Description)
		}
	} else {
		print("\nUSAGE:\n\t" + c.Name)
		args := getArgs(reflect.TypeOf(c.Content).In(1))
		for _, arg := range args {
			if arg.Optional() {
				if arg.Multiple() {
					print(" [..." + strings.ToUpper(arg.Field().Name) + "]")
				} else {
					print(" [" + strings.ToUpper(arg.Field().Name) + "]")
				}
			} else {
				if arg.Multiple() {
					print(" ..." + strings.ToUpper(arg.Field().Name))
				} else {
					print(" " + strings.ToUpper(arg.Field().Name))
				}
			}
		}
		// TODO: print flags
	}
}
