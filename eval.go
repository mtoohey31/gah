package gah

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"

	"github.com/gobeam/stringy"

	"mtoohey.com/gah/unmarshal"
)

// TODO: allow custom default values provided in tags
// TODO: allow passing of custom unmarshallers
// TODO: take aliases into account

func (c Cmd) SimpleEval() {
	err := c.Eval(os.Args)

	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[31%v\033[0m\n", err)
		os.Exit(1)
	}
}

func (c Cmd) Eval(args []string) error {
	if hasSubcommads(c) {
		return evalSubcommand(c, args, "")
	} else {
		return evalAndRun(c, args, "")
	}
}

func hasSubcommads(c Cmd) bool {
	return reflect.TypeOf(c.Content).AssignableTo(reflect.TypeOf([]Cmd{}))
}

func evalSubcommand(c Cmd, args []string, parentName string) error {
	if len(args) < 2 {
		return ErrExpectedSubcommand
	}
	arg := args[1]
	for _, subcommand := range c.Content.([]Cmd) {
		if arg == subcommand.Name {
			return subcommand.Eval(args[1:])
		}

		for _, alias := range subcommand.Aliases {
			if arg == alias {
				return subcommand.Eval(args[1:])
			}
		}
	}

	if arg == "-h" || arg == "--help" || arg == "help" {
		if len(args) == 2 {
			c.PrintHelp(parentName)
			return nil
		} else {
			for _, subcommand := range c.Content.([]Cmd) {
				if args[2] == subcommand.Name {
					subcommand.PrintHelp(parentName)
					return nil
				}

				for _, alias := range subcommand.Aliases {
					if arg == alias {
						subcommand.PrintHelp(parentName)
						return nil
					}
				}
			}

			if args[2] == "help" {
				// TODO: print help help
			}
		}
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

func evalAndRun(c Cmd, inputArgs []string, parentName string) error {
	flagsType := reflect.TypeOf(c.Content).In(0)
	flags := reflect.New(flagsType)
	argsType := reflect.TypeOf(c.Content).In(1)
	args := reflect.New(argsType)

	validShort, validLong := getValidFlags(flagsType)
	remainingArgs := getArgs(argsType)

	for i := 1; i < len(inputArgs); i++ {
		arg := inputArgs[i]
		if strings.HasPrefix(arg, "--") && len(arg) > 2 {
			eqIndex := strings.IndexRune(arg, '=')

			var flagName string

			if eqIndex == -1 {
				flagName = arg[2:]

				field, ok := validLong[flagName]

				if !ok {
					return trySalvageBuiltinLong(c, flagName, parentName)
				}

				unmarshaller, ok := unmarshal.Unmarshallers[field.Type]

				if !ok {
					panic(fmt.Sprintf("no unmarshaller for flag --%s", flagName))
				}

				if unmarshaller.Type().NumIn() == 2 {
					i++
					var flagValue string
					if i < len(inputArgs) {
						flagValue = inputArgs[i]
					} else {
						return expectedFlagValueLong(flagName)
					}

					res := unmarshaller.Call([]reflect.Value{reflect.ValueOf(flagValue),
						reflect.ValueOf(field.Tag)})

					if !res[1].IsNil() {
						return unmarshallingFlagLong(flagName, res[1].Interface().(error))
					}

					flags.Elem().FieldByIndex(field.Index).Set(res[0])
				} else {
					res := unmarshaller.Call([]reflect.Value{reflect.ValueOf(field.Tag)})

					if !res[1].IsNil() {
						return unmarshallingFlagLong(flagName, res[1].Interface().(error))
					}

					flags.Elem().FieldByIndex(field.Index).Set(res[0])
				}
			} else {
				flagName = arg[2:eqIndex]

				field, ok := validLong[flagName]

				if !ok {
					return trySalvageBuiltinLong(c, flagName, parentName)
				}

				flagValue := arg[eqIndex+1:]

				unmarshaller, ok := unmarshal.Unmarshallers[field.Type]

				if !ok {
					panic(fmt.Sprintf("no unmarshaller for --%s", flagName))
				}

				res := unmarshaller.Call([]reflect.Value{reflect.ValueOf(flagValue),
					reflect.ValueOf(field.Tag)})

				if !res[1].IsNil() {
					return unmarshallingFlagLong(flagName, res[1].Interface().(error))
				}

				flags.Elem().FieldByIndex(field.Index).Set(res[0])
			}
		} else if strings.HasPrefix(arg, "-") && len(arg) > 1 {
			// TODO: handle --

			eqIndex := strings.IndexRune(arg, '=')

			if eqIndex == -1 {
				shortRunes := []rune(arg[1:])

				for j := 0; j < len(shortRunes); j++ {
					flagRune := shortRunes[j]

					field, ok := validShort[flagRune]

					if !ok {
						return trySalvageBuiltinShort(c, flagRune, parentName)
					}

					unmarshaller, ok := unmarshal.Unmarshallers[field.Type]

					if !ok {
						panic(fmt.Sprintf("no unmarshaller for flag -%c", flagRune))
					}

					if unmarshaller.Type().NumIn() == 2 {
						if j == len(shortRunes)-1 {
							i++
							var flagValue string
							if i < len(inputArgs) {
								flagValue = inputArgs[i]
							} else {
								return expectedFlagValueShort(flagRune)
							}

							res := unmarshaller.Call([]reflect.Value{reflect.ValueOf(flagValue),
								reflect.ValueOf(field.Tag)})

							if !res[1].IsNil() {
								return unmarshallingFlagShort(flagRune, res[1].Interface().(error))
							}

							flags.Elem().FieldByIndex(field.Index).Set(res[0])
						} else {
							// TODO: try parsing the rest of the runes as a value here
							return expectedFlagValueShort(flagRune)
						}
					} else {
						res := unmarshaller.Call([]reflect.Value{reflect.ValueOf(field.Tag)})

						if !res[1].IsNil() {
							return unmarshallingFlagShort(flagRune, res[1].Interface().(error))
						}

						flags.Elem().FieldByIndex(field.Index).Set(res[0])
					}
				}
			} else {
				shortRunes := []rune(arg[1:eqIndex])
				flagValue := arg[eqIndex+1:]

				for j := 0; j < len(shortRunes); j++ {
					flagRune := shortRunes[j]

					field, ok := validShort[flagRune]

					if !ok {
						return trySalvageBuiltinShort(c, flagRune, parentName)
					}

					unmarshaller, ok := unmarshal.Unmarshallers[field.Type]

					if !ok {
						panic(fmt.Sprintf("no unmarshaller for flag -%c", flagRune))
					}

					if unmarshaller.Type().NumIn() == 2 {
						if j == len(shortRunes)-1 {
							if i < len(inputArgs) {
								flagValue = inputArgs[i]
							} else {
								return expectedFlagValueShort(flagRune)
							}

							res := unmarshaller.Call([]reflect.Value{reflect.ValueOf(flagValue),
								reflect.ValueOf(field.Tag)})

							if !res[1].IsNil() {
								return unmarshallingFlagShort(flagRune, res[1].Interface().(error))
							}

							flags.Elem().FieldByIndex(field.Index).Set(res[0])
						} else {
							return expectedFlagValueShort(flagRune)
						}
					} else {
						return unexpectedFlagValueShort(flagRune, flagValue)
					}
				}
			}
		} else {
			if len(remainingArgs) == 0 {
				if arg == "help" {
					c.PrintHelp(parentName)
					return nil
				} else {
					return &ErrUnexpectedArgument{argument: arg}
				}
			}

			field := remainingArgs[0].field
			t := remainingArgs[0].t
			// BUG: this will panic for non-collection types, we should check for
			// that first
			if args.Elem().FieldByIndex(field.Index).Len() == remainingArgs[0].max {
				remainingArgs = remainingArgs[1:]
				i--
				continue
			}

			unmarshaller, ok := unmarshal.Unmarshallers[t]

			if !ok {
				panic(fmt.Sprintf("no unmarshaller for argument type: %v", t))
			}

			res := unmarshaller.Call([]reflect.Value{reflect.ValueOf(arg),
				reflect.ValueOf(field.Tag)})

			if !res[1].IsNil() {
				// BUG: this will panic for non-collection types, we should check for
				// that first
				if args.Elem().FieldByIndex(field.Index).Len() >= remainingArgs[0].min {
					remainingArgs = remainingArgs[1:]
					i--
					continue
				} else {
					return &ErrUnmarshallingArgument{name: strings.ToUpper(t.Name()),
						value: arg, error: res[1].Interface().(error)}
				}
			}

			// TODO: refactor this to handle different kinds of collections
			if field.Type.Kind() == reflect.Slice {
				args.Elem().FieldByIndex(field.Index).Set(
					reflect.Append(args.Elem().FieldByIndex(field.Index), res[0]))
			} else {
				args.Elem().FieldByIndex(field.Index).Set(res[0])
				remainingArgs = remainingArgs[1:]
			}
		}
	}

	reflect.ValueOf(c.Content).Call([]reflect.Value{
		reflect.Indirect(flags), reflect.Indirect(args)})

	return nil
}

func trySalvageBuiltinLong(c Cmd, flagName string, parentName string) error {
	if flagName == "help" {
		c.PrintHelp(parentName)
		return nil
	} else if flagName == "version" && c.Version != "" {
		println(c.Version)
		return nil
	} else {
		return unexpectedLong(flagName)
	}
}

func trySalvageBuiltinShort(c Cmd, flagRune rune, parentName string) error {
	if flagRune == 'h' {
		c.PrintHelp(parentName)
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
			validLong[stringy.New(field.Name).KebabCase().ToLower()] = field
		}
	}

	return validShort, validLong
}

type argInfo struct {
	min   int
	max   int
	field reflect.StructField
	t     reflect.Type
}

func getArgs(argsType reflect.Type) []argInfo {
	var argInfoItems = make([]argInfo, len(reflect.VisibleFields(argsType)))

	for i, field := range reflect.VisibleFields(argsType) {
		if field.Type.Kind() == reflect.Slice {
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

			argInfoItems[i] = argInfo{min: min, max: max, field: field, t: field.Type.Elem()}
		} else {
			argInfoItems[i] = argInfo{min: 1, max: 1, field: field, t: field.Type}
		}
	}

	return argInfoItems
}

func (c Cmd) EvalMulticall(args []string) {
	wanted := path.Base(args[0])

	for _, subcommand := range c.Content.([]Cmd) {
		if wanted == subcommand.Name {
			subcommand.Eval(args)
			return
		}
	}
}

func (c Cmd) PrintHelp(parentName string) {
	println(parentName + "-" + c.Name + " " + c.Version)
	if c.Author != "" {
		println(c.Author)
	}
	if c.Description != "" {
		println(c.Description)
	}
	if hasSubcommads(c) {
		println("\nUSAGE:\n\t" + c.Name + " [SUBCOMMAND]")
		if c.Version == "" {
			println("\nFLAGS:\n\t-h, --help Prints help information")
		} else {
			println("\nFLAGS:\n\t-h, --help    Prints help information\n\t-v, --version Prints version information")
		}
		println("\nSUBCOMMANDS:")
		maxSubcommandNameLength := 0
		for _, subcommand := range c.Content.([]Cmd) {
			l := len(subcommand.Name)
			if l > maxSubcommandNameLength {
				maxSubcommandNameLength = l
			}
		}
		for _, subcommand := range c.Content.([]Cmd) {
			l := len(subcommand.Name)
			println("\t" + subcommand.Name + strings.Repeat(" ",
				1+maxSubcommandNameLength-l) + subcommand.Description)
		}
	} else {
		print("\nUSAGE:\n\t" + c.Name)
		args := getArgs(reflect.TypeOf(c.Content).In(1))
		for _, arg := range args {
			if arg.min <= 0 && arg.max == 1 {
				print(" [" + strings.ToUpper(arg.field.Name) + "]")
			} else if arg.min <= 0 {
				print(" [..." + strings.ToUpper(arg.field.Name) + "]")
			} else if arg.max == 1 {
				print(" " + strings.ToUpper(arg.field.Name))
			} else {
				print(" ..." + strings.ToUpper(arg.field.Name))
			}
		}
		// TODO: print flags
	}
}
