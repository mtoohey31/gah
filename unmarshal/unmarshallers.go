package unmarshal

import (
	"errors"
	"fmt"
	"math/bits"
	"os"
	"reflect"
	"strconv"
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

var Unmarshallers = map[reflect.Type]reflect.Value{
	reflect.TypeOf(i): reflect.ValueOf(func(s string, t reflect.StructTag) (int, error) {
		i, err := strconv.ParseInt(s, 10, bits.UintSize)

		if err != nil {
			return int(i), err
		}

		if minStr, ok := t.Lookup("minVal"); ok {
			min, err := strconv.ParseInt(minStr, 10, bits.UintSize)
			if err != nil {
				panic(err)
			}
			if i < min {
				return int(i),
					errors.New(fmt.Sprintf("%d less than minimum: %d", i, min))
			}
		}

		if maxStr, ok := t.Lookup("maxVal"); ok {
			max, err := strconv.ParseInt(maxStr, 10, bits.UintSize)
			if err != nil {
				panic(err)
			}
			if i > max {
				return int(i),
					errors.New(fmt.Sprintf("%d greater than than maximum: %d", i, max))
			}
		}

		return int(i), nil
	}),
	reflect.TypeOf(i8): reflect.ValueOf(func(s string, t reflect.StructTag) (int8, error) {
		i, err := strconv.ParseInt(s, 10, 8)

		if err != nil {
			return int8(i), err
		}

		if minStr, ok := t.Lookup("minVal"); ok {
			min, err := strconv.ParseInt(minStr, 10, 8)
			if err != nil {
				panic(err)
			}
			if i < min {
				return int8(i),
					errors.New(fmt.Sprintf("%d less than minimum: %d", i, min))
			}
		}

		if maxStr, ok := t.Lookup("maxVal"); ok {
			max, err := strconv.ParseInt(maxStr, 10, 8)
			if err != nil {
				panic(err)
			}
			if i > max {
				return int8(i),
					errors.New(fmt.Sprintf("%d greater than than maximum: %d", i, max))
			}
		}

		return int8(i), nil
	}),
	reflect.TypeOf(i16): reflect.ValueOf(func(s string, t reflect.StructTag) (int16, error) {
		i, err := strconv.ParseInt(s, 10, 16)

		if err != nil {
			return int16(i), err
		}

		if minStr, ok := t.Lookup("minVal"); ok {
			min, err := strconv.ParseInt(minStr, 10, 16)
			if err != nil {
				panic(err)
			}
			if i < min {
				return int16(i),
					errors.New(fmt.Sprintf("%d less than minimum: %d", i, min))
			}
		}

		if maxStr, ok := t.Lookup("maxVal"); ok {
			max, err := strconv.ParseInt(maxStr, 10, 16)
			if err != nil {
				panic(err)
			}
			if i > max {
				return int16(i),
					errors.New(fmt.Sprintf("%d greater than than maximum: %d", i, max))
			}
		}

		return int16(i), nil
	}),
	reflect.TypeOf(i32): reflect.ValueOf(func(s string, t reflect.StructTag) (int32, error) {
		i, err := strconv.ParseInt(s, 10, 32)

		if err != nil {
			return int32(i), err
		}

		if minStr, ok := t.Lookup("minVal"); ok {
			min, err := strconv.ParseInt(minStr, 10, 32)
			if err != nil {
				panic(err)
			}
			if i < min {
				return int32(i),
					errors.New(fmt.Sprintf("%d less than minimum: %d", i, min))
			}
		}

		if maxStr, ok := t.Lookup("maxVal"); ok {
			max, err := strconv.ParseInt(maxStr, 10, 32)
			if err != nil {
				panic(err)
			}
			if i > max {
				return int32(i),
					errors.New(fmt.Sprintf("%d greater than than maximum: %d", i, max))
			}
		}

		return int32(i), nil
	}),
	reflect.TypeOf(i64): reflect.ValueOf(func(s string, t reflect.StructTag) (int64, error) {
		i, err := strconv.ParseInt(s, 10, 64)

		if err != nil {
			return i, err
		}

		if minStr, ok := t.Lookup("minVal"); ok {
			min, err := strconv.ParseInt(minStr, 10, 64)
			if err != nil {
				panic(err)
			}
			if i < min {
				return int64(i),
					errors.New(fmt.Sprintf("%d less than minimum: %d", i, min))
			}
		}

		if maxStr, ok := t.Lookup("maxVal"); ok {
			max, err := strconv.ParseInt(maxStr, 10, 64)
			if err != nil {
				panic(err)
			}
			if i > max {
				return int64(i),
					errors.New(fmt.Sprintf("%d greater than than maximum: %d", i, max))
			}
		}

		return i, err
	}),
	reflect.TypeOf(u): reflect.ValueOf(func(s string, t reflect.StructTag) (uint, error) {
		i, err := strconv.ParseUint(s, 10, bits.UintSize)

		if err != nil {
			return uint(i), err
		}

		if minStr, ok := t.Lookup("minVal"); ok {
			min, err := strconv.ParseUint(minStr, 10, bits.UintSize)
			if err != nil {
				panic(err)
			}
			if i < min {
				return uint(i),
					errors.New(fmt.Sprintf("%d less than minimum: %d", i, min))
			}
		}

		if maxStr, ok := t.Lookup("maxVal"); ok {
			max, err := strconv.ParseUint(maxStr, 10, bits.UintSize)
			if err != nil {
				panic(err)
			}
			if i > max {
				return uint(i),
					errors.New(fmt.Sprintf("%d greater than than maximum: %d", i, max))
			}
		}

		return uint(i), nil
	}),
	reflect.TypeOf(u8): reflect.ValueOf(func(s string, t reflect.StructTag) (uint8, error) {
		i, err := strconv.ParseUint(s, 10, 8)

		if err != nil {
			return uint8(i), err
		}

		if minStr, ok := t.Lookup("minVal"); ok {
			min, err := strconv.ParseUint(minStr, 10, 8)
			if err != nil {
				panic(err)
			}
			if i < min {
				return uint8(i),
					errors.New(fmt.Sprintf("%d less than minimum: %d", i, min))
			}
		}

		if maxStr, ok := t.Lookup("maxVal"); ok {
			max, err := strconv.ParseUint(maxStr, 10, 8)
			if err != nil {
				panic(err)
			}
			if i > max {
				return uint8(i),
					errors.New(fmt.Sprintf("%d greater than than maximum: %d", i, max))
			}
		}

		return uint8(i), nil
	}),
	reflect.TypeOf(u16): reflect.ValueOf(func(s string, t reflect.StructTag) (uint16, error) {
		i, err := strconv.ParseUint(s, 10, 16)

		if err != nil {
			return uint16(i), err
		}

		if minStr, ok := t.Lookup("minVal"); ok {
			min, err := strconv.ParseUint(minStr, 10, 16)
			if err != nil {
				panic(err)
			}
			if i < min {
				return uint16(i),
					errors.New(fmt.Sprintf("%d less than minimum: %d", i, min))
			}
		}

		if maxStr, ok := t.Lookup("maxVal"); ok {
			max, err := strconv.ParseUint(maxStr, 10, 16)
			if err != nil {
				panic(err)
			}
			if i > max {
				return uint16(i),
					errors.New(fmt.Sprintf("%d greater than than maximum: %d", i, max))
			}
		}

		return uint16(i), nil
	}),
	reflect.TypeOf(u32): reflect.ValueOf(func(s string, t reflect.StructTag) (uint32, error) {
		i, err := strconv.ParseUint(s, 10, 32)

		if err != nil {
			return uint32(i), err
		}

		if minStr, ok := t.Lookup("minVal"); ok {
			min, err := strconv.ParseUint(minStr, 10, 32)
			if err != nil {
				panic(err)
			}
			if i < min {
				return uint32(i),
					errors.New(fmt.Sprintf("%d less than minimum: %d", i, min))
			}
		}

		if maxStr, ok := t.Lookup("maxVal"); ok {
			max, err := strconv.ParseUint(maxStr, 10, 32)
			if err != nil {
				panic(err)
			}
			if i > max {
				return uint32(i),
					errors.New(fmt.Sprintf("%d greater than than maximum: %d", i, max))
			}
		}

		return uint32(i), nil
	}),
	reflect.TypeOf(u64): reflect.ValueOf(func(s string, t reflect.StructTag) (uint64, error) {
		i, err := strconv.ParseUint(s, 10, 64)

		if err != nil {
			return i, err
		}

		if minStr, ok := t.Lookup("minVal"); ok {
			min, err := strconv.ParseUint(minStr, 10, 64)
			if err != nil {
				panic(err)
			}
			if i < min {
				return uint64(i),
					errors.New(fmt.Sprintf("%d less than minimum: %d", i, min))
			}
		}

		if maxStr, ok := t.Lookup("maxVal"); ok {
			max, err := strconv.ParseUint(maxStr, 10, 64)
			if err != nil {
				panic(err)
			}
			if i > max {
				return uint64(i),
					errors.New(fmt.Sprintf("%d greater than than maximum: %d", i, max))
			}
		}

		return i, nil
	}),

	// TODO: add floats

	reflect.TypeOf(b): reflect.ValueOf(func(t reflect.StructTag) (bool, error) {
		_, invert := t.Lookup("invert")

		return !invert, nil
	}),

	reflect.TypeOf(s): reflect.ValueOf(func(s string, t reflect.StructTag) (string, error) {
		_, path := t.Lookup("path")
		if path {
			_, err := os.Stat(s)
			return s, err
		}

		// TODO: add file, dir, executable, socket, etc. available in unix test

		return s, nil
	}),
}
