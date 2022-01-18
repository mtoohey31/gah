package unmarshal

import (
	"errors"
	"fmt"
	"math/bits"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var i int
var i8 int8
var i16 int16
var i32 int32
var i64 int64
var u uint
var u8 uint8
var u16 uint16
var u32 uint32
var u64 uint64
var b bool
var s string

// TODO: add support fur custom unmarshallers, which can be set as the default
// for a type, or can be registered and called by name
// TODO: add support for enums somehow

var defaultsToNoValue = []reflect.Type{reflect.TypeOf(b)}

func TakesValue(f reflect.StructField) bool {
	s, found := f.Tag.Lookup("takesVal")

	if found {
		takesVal, err := unmarshalBool(s)

		if err == nil {
			return takesVal
		} else {
			panic(err)
		}
	}

	for _, v := range defaultsToNoValue {
		if f.Type == v {
			return false
		}
	}

	return true
}

func unmarshalBool(s string) (bool, error) {
	trimmed := strings.ToLower(strings.TrimSpace(s))

	for _, v := range []string{"0", "n", "no", "f", "false"} {
		if trimmed == v {
			return false, nil
		}
	}

	for _, v := range []string{"1", "y", "yes", "t", "true"} {
		if trimmed == v {
			return true, nil
		}
	}

	return false, errors.New(fmt.Sprintf("invalid syntax \"%s\"", s))
}

var ValueUnmarshallers = map[reflect.Type]func(string, reflect.StructTag) (reflect.Value, error){
	reflect.TypeOf(i): func(s string, t reflect.StructTag) (reflect.Value, error) {
		i, err := strconv.ParseInt(s, 10, bits.UintSize)

		if err != nil {
			return reflect.ValueOf(int(i)), err
		}

		if minStr, ok := t.Lookup("minVal"); ok {
			min, err := strconv.ParseInt(minStr, 10, bits.UintSize)
			if err != nil {
				panic(err)
			}
			if i < min {
				return reflect.ValueOf(int(i)),
					errors.New(fmt.Sprintf("%d less than minimum: %d", i, min))
			}
		}

		if maxStr, ok := t.Lookup("maxVal"); ok {
			max, err := strconv.ParseInt(maxStr, 10, bits.UintSize)
			if err != nil {
				panic(err)
			}
			if i > max {
				return reflect.ValueOf(int(i)),
					errors.New(fmt.Sprintf("%d greater than than maximum: %d", i, max))
			}
		}

		return reflect.ValueOf(int(i)), nil
	},
	reflect.TypeOf(i8): func(s string, t reflect.StructTag) (reflect.Value, error) {
		i, err := strconv.ParseInt(s, 10, 8)

		if err != nil {
			return reflect.ValueOf(int8(i)), err
		}

		if minStr, ok := t.Lookup("minVal"); ok {
			min, err := strconv.ParseInt(minStr, 10, 8)
			if err != nil {
				panic(err)
			}
			if i < min {
				return reflect.ValueOf(int8(i)),
					errors.New(fmt.Sprintf("%d less than minimum: %d", i, min))
			}
		}

		if maxStr, ok := t.Lookup("maxVal"); ok {
			max, err := strconv.ParseInt(maxStr, 10, 8)
			if err != nil {
				panic(err)
			}
			if i > max {
				return reflect.ValueOf(int8(i)),
					errors.New(fmt.Sprintf("%d greater than than maximum: %d", i, max))
			}
		}

		return reflect.ValueOf(int8(i)), nil
	},
	reflect.TypeOf(i16): func(s string, t reflect.StructTag) (reflect.Value, error) {
		i, err := strconv.ParseInt(s, 10, 16)

		if err != nil {
			return reflect.ValueOf(int16(i)), err
		}

		if minStr, ok := t.Lookup("minVal"); ok {
			min, err := strconv.ParseInt(minStr, 10, 16)
			if err != nil {
				panic(err)
			}
			if i < min {
				return reflect.ValueOf(int16(i)),
					errors.New(fmt.Sprintf("%d less than minimum: %d", i, min))
			}
		}

		if maxStr, ok := t.Lookup("maxVal"); ok {
			max, err := strconv.ParseInt(maxStr, 10, 16)
			if err != nil {
				panic(err)
			}
			if i > max {
				return reflect.ValueOf(int16(i)),
					errors.New(fmt.Sprintf("%d greater than than maximum: %d", i, max))
			}
		}

		return reflect.ValueOf(int16(i)), nil
	},
	reflect.TypeOf(i32): func(s string, t reflect.StructTag) (reflect.Value, error) {
		i, err := strconv.ParseInt(s, 10, 32)

		if err != nil {
			return reflect.ValueOf(int32(i)), err
		}

		if minStr, ok := t.Lookup("minVal"); ok {
			min, err := strconv.ParseInt(minStr, 10, 32)
			if err != nil {
				panic(err)
			}
			if i < min {
				return reflect.ValueOf(int32(i)),
					errors.New(fmt.Sprintf("%d less than minimum: %d", i, min))
			}
		}

		if maxStr, ok := t.Lookup("maxVal"); ok {
			max, err := strconv.ParseInt(maxStr, 10, 32)
			if err != nil {
				panic(err)
			}
			if i > max {
				return reflect.ValueOf(int32(i)),
					errors.New(fmt.Sprintf("%d greater than than maximum: %d", i, max))
			}
		}

		return reflect.ValueOf(int32(i)), nil
	},
	reflect.TypeOf(i64): func(s string, t reflect.StructTag) (reflect.Value, error) {
		i, err := strconv.ParseInt(s, 10, 64)

		if err != nil {
			return reflect.ValueOf(i), err
		}

		if minStr, ok := t.Lookup("minVal"); ok {
			min, err := strconv.ParseInt(minStr, 10, 64)
			if err != nil {
				panic(err)
			}
			if i < min {
				return reflect.ValueOf(i),
					errors.New(fmt.Sprintf("%d less than minimum: %d", i, min))
			}
		}

		if maxStr, ok := t.Lookup("maxVal"); ok {
			max, err := strconv.ParseInt(maxStr, 10, 64)
			if err != nil {
				panic(err)
			}
			if i > max {
				return reflect.ValueOf(i),
					errors.New(fmt.Sprintf("%d greater than than maximum: %d", i, max))
			}
		}

		return reflect.ValueOf(i), err
	},
	reflect.TypeOf(u): func(s string, t reflect.StructTag) (reflect.Value, error) {
		i, err := strconv.ParseUint(s, 10, bits.UintSize)

		if err != nil {
			return reflect.ValueOf(uint(i)), err
		}

		if minStr, ok := t.Lookup("minVal"); ok {
			min, err := strconv.ParseUint(minStr, 10, bits.UintSize)
			if err != nil {
				panic(err)
			}
			if i < min {
				return reflect.ValueOf(uint(i)),
					errors.New(fmt.Sprintf("%d less than minimum: %d", i, min))
			}
		}

		if maxStr, ok := t.Lookup("maxVal"); ok {
			max, err := strconv.ParseUint(maxStr, 10, bits.UintSize)
			if err != nil {
				panic(err)
			}
			if i > max {
				return reflect.ValueOf(uint(i)),
					errors.New(fmt.Sprintf("%d greater than than maximum: %d", i, max))
			}
		}

		return reflect.ValueOf(uint(i)), nil
	},
	reflect.TypeOf(u8): func(s string, t reflect.StructTag) (reflect.Value, error) {
		i, err := strconv.ParseUint(s, 10, 8)

		if err != nil {
			return reflect.ValueOf(uint8(i)), err
		}

		if minStr, ok := t.Lookup("minVal"); ok {
			min, err := strconv.ParseUint(minStr, 10, 8)
			if err != nil {
				panic(err)
			}
			if i < min {
				return reflect.ValueOf(uint8(i)),
					errors.New(fmt.Sprintf("%d less than minimum: %d", i, min))
			}
		}

		if maxStr, ok := t.Lookup("maxVal"); ok {
			max, err := strconv.ParseUint(maxStr, 10, 8)
			if err != nil {
				panic(err)
			}
			if i > max {
				return reflect.ValueOf(uint8(i)),
					errors.New(fmt.Sprintf("%d greater than than maximum: %d", i, max))
			}
		}

		return reflect.ValueOf(uint8(i)), nil
	},
	reflect.TypeOf(u16): func(s string, t reflect.StructTag) (reflect.Value, error) {
		i, err := strconv.ParseUint(s, 10, 16)

		if err != nil {
			return reflect.ValueOf(uint16(i)), err
		}

		if minStr, ok := t.Lookup("minVal"); ok {
			min, err := strconv.ParseUint(minStr, 10, 16)
			if err != nil {
				panic(err)
			}
			if i < min {
				return reflect.ValueOf(uint16(i)),
					errors.New(fmt.Sprintf("%d less than minimum: %d", i, min))
			}
		}

		if maxStr, ok := t.Lookup("maxVal"); ok {
			max, err := strconv.ParseUint(maxStr, 10, 16)
			if err != nil {
				panic(err)
			}
			if i > max {
				return reflect.ValueOf(uint16(i)),
					errors.New(fmt.Sprintf("%d greater than than maximum: %d", i, max))
			}
		}

		return reflect.ValueOf(uint16(i)), nil
	},
	reflect.TypeOf(u32): func(s string, t reflect.StructTag) (reflect.Value, error) {
		i, err := strconv.ParseUint(s, 10, 32)

		if err != nil {
			return reflect.ValueOf(uint32(i)), err
		}

		if minStr, ok := t.Lookup("minVal"); ok {
			min, err := strconv.ParseUint(minStr, 10, 32)
			if err != nil {
				panic(err)
			}
			if i < min {
				return reflect.ValueOf(uint32(i)),
					errors.New(fmt.Sprintf("%d less than minimum: %d", i, min))
			}
		}

		if maxStr, ok := t.Lookup("maxVal"); ok {
			max, err := strconv.ParseUint(maxStr, 10, 32)
			if err != nil {
				panic(err)
			}
			if i > max {
				return reflect.ValueOf(uint32(i)),
					errors.New(fmt.Sprintf("%d greater than than maximum: %d", i, max))
			}
		}

		return reflect.ValueOf(uint32(i)), nil
	},
	reflect.TypeOf(u64): func(s string, t reflect.StructTag) (reflect.Value, error) {
		i, err := strconv.ParseUint(s, 10, 64)

		if err != nil {
			return reflect.ValueOf(i), err
		}

		if minStr, ok := t.Lookup("minVal"); ok {
			min, err := strconv.ParseUint(minStr, 10, 64)
			if err != nil {
				panic(err)
			}
			if i < min {
				return reflect.ValueOf(i),
					errors.New(fmt.Sprintf("%d less than minimum: %d", i, min))
			}
		}

		if maxStr, ok := t.Lookup("maxVal"); ok {
			max, err := strconv.ParseUint(maxStr, 10, 64)
			if err != nil {
				panic(err)
			}
			if i > max {
				return reflect.ValueOf(i),
					errors.New(fmt.Sprintf("%d greater than than maximum: %d", i, max))
			}
		}

		return reflect.ValueOf(i), nil
	},

	// TODO: add floats

	reflect.TypeOf(b): func(s string, t reflect.StructTag) (reflect.Value, error) {
		_, invert := t.Lookup("invert")
		b, err := unmarshalBool(s)

		if err == nil {
			return reflect.ValueOf(invert != b), nil
		} else {
			return reflect.ValueOf(false), err
		}
	},

	reflect.TypeOf(s): func(s string, t reflect.StructTag) (reflect.Value, error) {
		_, path := t.Lookup("path")
		if path {
			_, err := os.Stat(s)
			return reflect.ValueOf(s), err
		}

		// TODO: add file, dir, executable, socket, etc. available in unix test

		return reflect.ValueOf(s), nil
	},
}

var ValuelessUnmarshallers = map[reflect.Type]func(reflect.Value, reflect.StructTag) (reflect.Value, error){
	reflect.TypeOf(b): func(_ reflect.Value, t reflect.StructTag) (reflect.Value, error) {
		_, invert := t.Lookup("invert")

		return reflect.ValueOf(!invert), nil
	},
}
