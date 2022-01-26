package gah

import (
	"fmt"
	"math"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"mtoohey.com/gah/unmarshal"
)

func (c Cmd) SimpleEval() {
	err := c.Eval(os.Args, []string{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[31m%v\033[0m\n", err)
		os.Exit(1)
	}
}

func (c Cmd) EvalMulticall(args []string) {
	wanted := path.Base(args[0])
	for _, subcommand := range c.Subcommands {
		if wanted == subcommand.Name {
			subcommand.Eval(args, []string{})
			return
		}
	}
}

func (c Cmd) Eval(inputArgs []string, parentNames []string) error {
	var flagsType reflect.Type
	var argsType reflect.Type
	if c.Function == nil {
		flagsType = reflect.TypeOf(struct{}{})
		argsType = flagsType
	} else {
		flagsType = reflect.TypeOf(c.Function).In(0)
		argsType = reflect.TypeOf(c.Function).In(1)
	}
	flags := reflect.New(flagsType)
	var positionalArgs []string

	var enrichedSubcommands []Cmd
	if c.Subcommands != nil {
		helpFound := false
		for _, subcommand := range c.Subcommands {
			if subcommand.Name == "help" {
				helpFound = true
				break
			}

			for _, alias := range subcommand.Aliases {
				if alias == "help" {
					helpFound = true
					break
				}
			}
			if helpFound {
				break
			}
		}

		if !helpFound {
			// TODO: can I make this an append instead? implement tests that ensure
			// the subcommands aren't mutated in the return value
			enrichedSubcommands = make([]Cmd, len(c.Subcommands)+1)
			copy(enrichedSubcommands, c.Subcommands)
			enrichedSubcommands[len(enrichedSubcommands)-1] = Cmd{
				Name: "help",
				Function: func(_ struct{}, a struct {
					SubcommandName []string `min:"0" max:"1"`
				}) {
					if len(a.SubcommandName) > 0 {
						for _, subcommand := range c.Subcommands {
							if subcommand.Name == a.SubcommandName[0] {
								subcommand.PrintHelp(append(parentNames, c.Name))
								return
							}

							for _, alias := range subcommand.Aliases {
								if alias == a.SubcommandName[0] {
									subcommand.PrintHelp(append(parentNames, c.Name))
									return
								}
							}
						}
					}

					c.PrintHelp(parentNames)
				},
			}
		}
	}

	allFlags := getFlags(flagsType)
	validShort, validLong := getFlagMaps(allFlags)

	for i := 1; i < len(inputArgs); i++ {
		arg := inputArgs[i]

		if strings.HasPrefix(arg, "--") && len(arg) > 2 {
			eqIndex := strings.IndexRune(arg, '=')
			var flagName string
			if eqIndex == -1 {
				flagName = arg[2:]
			} else {
				flagName = arg[2:eqIndex]
			}

			flag, ok := validLong[flagName]
			if !ok {
				return trySalvageBuiltinLong(c, flagName, parentNames)
			}

			if unmarshal.TakesValue(flag.field) {
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

				unmarshaller := unmarshal.GetValueUnmarshaller(flag.field.Type,
					c.CustomValueUnmarshallers)

				res, err := unmarshaller(flagValue, flag.field.Tag)
				if err != nil {
					return unmarshallingFlagLong(flagName, err)
				}
				flags.Elem().FieldByIndex((*flag).field.Index).Set(res)
				flag.set = true
			} else {
				if eqIndex != -1 {
					return unexpectedFlagValueLong(flagName, arg[eqIndex+1:])
				}

				unmarshaller := unmarshal.GetValuelessUnmarshaller(flag.field.Type,
					c.CustomValuelessUnmarshallers)

				res, err := unmarshaller(flags.Elem().FieldByIndex(flag.field.Index),
					flag.field.Tag)
				if err != nil {
					return unmarshallingFlagLong(flagName, err)
				}
				flags.Elem().FieldByIndex((*flag).field.Index).Set(res)
				flag.set = true
			}
		} else if strings.HasPrefix(arg, "-") && len(arg) > 1 {
			if arg == "--" {
				positionalArgs = append(positionalArgs, inputArgs[i+1:]...)
				break
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

				flag, ok := validShort[flagRune]
				if !ok {
					return trySalvageBuiltinShort(c, flagRune, parentNames)
				}

				if unmarshal.TakesValue(flag.field) {
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

					unmarshaller := unmarshal.GetValueUnmarshaller(flag.field.Type,
						c.CustomValueUnmarshallers)

					res, err := unmarshaller(flagValue, flag.field.Tag)
					if err != nil {
						return unmarshallingFlagShort(flagRune, err)
					}
					flags.Elem().FieldByIndex(flag.field.Index).Set(res)
					flag.set = true
				} else {
					if j == len(flagRunes)-1 && eqIndex != -1 {
						return unexpectedFlagValueShort(flagRune, arg[eqIndex+1:])
					}

					unmarshaller := unmarshal.GetValuelessUnmarshaller(flag.field.Type,
						c.CustomValuelessUnmarshallers)
					if !ok {
						panic(fmt.Sprintf("no valueless unmarshaller for flag -%c", flagRune))
					}

					res, err := unmarshaller(flags.Elem().FieldByIndex(flag.field.Index),
						flag.field.Tag)
					if err != nil {
						return unmarshallingFlagShort(flagRune, err)
					}
					flags.Elem().FieldByIndex(flag.field.Index).Set(res)
					flag.set = true
				}
			}
		} else {
			if c.Subcommands == nil {
				positionalArgs = append(positionalArgs, arg)
			} else {
				for _, subcommand := range enrichedSubcommands {
					if arg == subcommand.Name {
						if c.Function != nil {
							reflect.ValueOf(c.Function).Call([]reflect.Value{reflect.Indirect(flags),
								reflect.Indirect(reflect.New(reflect.TypeOf(c.Function).In(1)))})
						}
						return subcommand.Eval(inputArgs[i:], append(parentNames, c.Name))
					}

					for _, alias := range subcommand.Aliases {
						if arg == alias {
							if c.Function != nil {
								reflect.ValueOf(c.Function).Call([]reflect.Value{reflect.Indirect(flags),
									reflect.Indirect(reflect.New(reflect.TypeOf(c.Function).In(1)))})
							}
							return subcommand.Eval(inputArgs[i:], append(parentNames, c.Name))
						}
					}
				}

				return &ErrInvalidSubcommand{subcommand: arg}
			}
		}
	}

	if c.Subcommands != nil {
		return &ErrExpectedSubcommand{}
	}

	args := reflect.New(argsType)
	argInfo := getArgs(argsType)

	minArgs := 0
	maxArgs := 0
	for _, arg := range argInfo {
		minArgs += arg.Min()
		maxArgs += arg.Max()
	}

	if len(positionalArgs) < minArgs {
		for _, arg := range argInfo {
			minArgs -= arg.Min()
			if minArgs < 0 {
				return &ErrExpectedArgumentValue{name: strings.ToUpper(arg.Field().Name)}
			}
		}
	} else if len(positionalArgs) > maxArgs {
		return &ErrUnexpectedArgument{argument: positionalArgs[maxArgs]}
	}

	additionalVariableArgs := len(positionalArgs) - minArgs

	i := 0
	for _, info := range argInfo {
		var numToTake int
		if info.Min() != info.Max() {
			numToTake = info.Min() + additionalVariableArgs
		} else {
			numToTake = info.Min()
		}

		for j := 0; j < numToTake; j++ {
			res, err := info.Unmarshaller(c.CustomValueUnmarshallers)(positionalArgs[i], info.Field().Tag)
			if err != nil {
				return &ErrUnmarshallingArgument{
					name:  strings.ToUpper(info.Field().Name),
					value: positionalArgs[i], error: err}
			}

			i++
			info.Update(args, res)
		}
	}

	if c.DefaultFlags != nil {
		for _, flag := range allFlags {
			flag.SetDefaultIfUnset(flags, c.DefaultFlags, c.CustomValueUnmarshallers)
		}
	}

	reflect.ValueOf(c.Function).Call([]reflect.Value{
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

type flagInfo struct {
	field reflect.StructField
	set   bool
}

func (i *flagInfo) SetDefaultIfUnset(f reflect.Value, d interface{}, c unmarshal.CustomValueUnmarshallers) {
	if i.set {
		return
	}

	f.Elem().FieldByIndex(i.field.Index).Set(reflect.ValueOf(d).FieldByIndex(i.field.Index))
	i.set = true
}

func getFlags(flagsType reflect.Type) []flagInfo {
	visibleFields := reflect.VisibleFields(flagsType)
	flagInfoItems := make([]flagInfo, len(visibleFields))

	for i, field := range visibleFields {
		flagInfoItems[i] = flagInfo{field: field}
	}

	return flagInfoItems
}

func getFlagMaps(flags []flagInfo) (map[rune]*flagInfo, map[string]*flagInfo) {
	validShort := make(map[rune]*flagInfo)
	validLong := make(map[string]*flagInfo)

	for i := range flags {
		short, found := flags[i].field.Tag.Lookup("short")
		if found {
			runes := []rune(short)
			validShort[runes[0]] = &flags[i]
		}

		long, found := flags[i].field.Tag.Lookup("long")
		if found {
			validLong[long] = &flags[i]
		} else {
			validLong[pascalToKebab(flags[i].field.Name)] = &flags[i]
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
	Min() int
	Max() int
	Field() reflect.StructField
	Unmarshaller(c unmarshal.CustomValueUnmarshallers,
	) func(string, reflect.StructTag) (reflect.Value, error)
	Update(reflect.Value, reflect.Value)
	Optional() bool
	Multiple() bool
}

type sliceArgInfo struct {
	min   int
	max   int
	field reflect.StructField
}

func (i *sliceArgInfo) Min() int { return i.min }

func (i *sliceArgInfo) Max() int { return i.max }

func (i *sliceArgInfo) Field() reflect.StructField { return i.field }

func (i *sliceArgInfo) Unmarshaller(c unmarshal.CustomValueUnmarshallers,
) func(string, reflect.StructTag) (reflect.Value, error) {
	return unmarshal.GetValueUnmarshaller(i.field.Type.Elem(), c)
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

func (i *arrayArgInfo) Min() int { return i.field.Type.Len() }

func (i *arrayArgInfo) Max() int { return i.field.Type.Len() }

func (i *arrayArgInfo) Field() reflect.StructField { return i.field }

func (i *arrayArgInfo) Unmarshaller(c unmarshal.CustomValueUnmarshallers,
) func(string, reflect.StructTag) (reflect.Value, error) {
	return unmarshal.GetValueUnmarshaller(i.field.Type.Elem(), c)
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

func (i *defaultArgInfo) Min() int { return 1 }

func (i *defaultArgInfo) Max() int { return 1 }

func (i *defaultArgInfo) Field() reflect.StructField { return i.field }

func (i *defaultArgInfo) Unmarshaller(c unmarshal.CustomValueUnmarshallers,
) func(string, reflect.StructTag) (reflect.Value, error) {
	return unmarshal.GetValueUnmarshaller(i.field.Type, c)
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
			if !unmarshal.ElementWise(field) {
				argInfoItems[i] = &defaultArgInfo{field: field}
				continue
			}

			min := 0
			minStr, found := field.Tag.Lookup("min")
			if found {
				var err error
				min, err = strconv.Atoi(minStr)
				if err != nil {
					panic(err)
				}
			}

			max := math.MaxInt
			maxStr, found := field.Tag.Lookup("max")
			if found {
				var err error
				max, err = strconv.Atoi(maxStr)
				if err != nil {
					panic(err)
				}
			}

			argInfoItems[i] = &sliceArgInfo{min: min, max: max, field: field}
		case reflect.Array:
			if !unmarshal.ElementWise(field) {
				argInfoItems[i] = &defaultArgInfo{field: field}
				continue
			}

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
	if c.Subcommands != nil {
		println("\nUSAGE:\n\t" + strings.Join(append(parentNames, c.Name), " ") + " [SUBCOMMAND]")
		if c.Version == "" {
			println("\nFLAGS:\n\t-h, --help Prints help information")
		} else {
			println("\nFLAGS:\n\t-h, --help    Prints help information\n\t-v, --version Prints version information")
		}
		println("\nSUBCOMMANDS:")
		maxSubcommandNameLength := 0
		for _, subcommand := range c.Subcommands {
			l := len(strings.Join(
				append([]string{subcommand.Name}, subcommand.Aliases...), ", "))
			if l > maxSubcommandNameLength {
				maxSubcommandNameLength = l
			}
		}
		for _, subcommand := range c.Subcommands {
			s := strings.Join(append([]string{subcommand.Name},
				subcommand.Aliases...), ", ")
			l := len(s)
			println("\t" + s + strings.Repeat(" ",
				1+maxSubcommandNameLength-l) + subcommand.Description)
		}
	} else {
		print("\nUSAGE:\n\t" + c.Name)
		args := getArgs(reflect.TypeOf(c.Function).In(1))
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
		// to ensure there's a new line at the end of the usage line
		println()
		// TODO: print flags
	}
}
