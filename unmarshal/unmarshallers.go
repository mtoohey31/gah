package unmarshal

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/bits"
	"net"
	"os"
	"reflect"
	"strconv"
	"time"
)

// TODO: add support for enums somehow
// TODO: add tests for all unmarshallers
// TODO: support all the types that pflag does

var defaultsToNoValue = []reflect.Type{reflect.TypeOf(false)}

func TakesValue(f reflect.StructField) bool {
	s, found := f.Tag.Lookup("takesVal")

	if found {
		takesVal, err := strconv.ParseBool(s)

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

func GetValueUnmarshaller(t reflect.Type, g reflect.StructTag,
	c CustomValueUnmarshallers) ValueUnmarshaller {
	u, found := c[t]
	if found {
		return u
	}
	u, found = valueUnmarshallers[t]
	if found {
		return u
	}
	panic(fmt.Sprintf("no value unmarshaller for type %s", t.Name()))
}

func GetValuelessUnmarshaller(t reflect.Type, g reflect.StructTag,
	c CustomValuelessUnmarshallers) ValuelessUnmarshaller {
	u, found := c[t]
	if found {
		return u
	}
	u, found = valuelessUnmarshallers[t]
	if found {
		return u
	}
	panic(fmt.Sprintf("no valueless unmarshaller for type %s", t.Name()))
}

type ValueUnmarshaller = func(string, reflect.StructTag) (reflect.Value, error)
type ValuelessUnmarshaller = func(reflect.Value, reflect.StructTag) (reflect.Value, error)

type CustomValueUnmarshallers = map[reflect.Type]ValueUnmarshaller
type CustomValuelessUnmarshallers = map[reflect.Type]ValuelessUnmarshaller

var valueUnmarshallers = map[reflect.Type]ValueUnmarshaller{
	reflect.TypeOf(int(0)): func(s string, t reflect.StructTag) (reflect.Value, error) {
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
	reflect.TypeOf(int8(0)): func(s string, t reflect.StructTag) (reflect.Value, error) {
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
	reflect.TypeOf(int16(0)): func(s string, t reflect.StructTag) (reflect.Value, error) {
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
	reflect.TypeOf(int32(0)): func(s string, t reflect.StructTag) (reflect.Value, error) {
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
	reflect.TypeOf(int64(0)): func(s string, t reflect.StructTag) (reflect.Value, error) {
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
	reflect.TypeOf(uint(0)): func(s string, t reflect.StructTag) (reflect.Value, error) {
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
	reflect.TypeOf(uint8(0)): func(s string, t reflect.StructTag) (reflect.Value, error) {
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
	reflect.TypeOf(uint16(0)): func(s string, t reflect.StructTag) (reflect.Value, error) {
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
	reflect.TypeOf(uint32(0)): func(s string, t reflect.StructTag) (reflect.Value, error) {
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
	reflect.TypeOf(uint64(0)): func(s string, t reflect.StructTag) (reflect.Value, error) {
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
	reflect.TypeOf(float32(0.0)): func(s string, t reflect.StructTag) (reflect.Value, error) {
		i, err := strconv.ParseFloat(s, 32)

		if err != nil {
			return reflect.ValueOf(float32(i)), err
		}

		if minStr, ok := t.Lookup("minVal"); ok {
			min, err := strconv.ParseFloat(minStr, 32)
			if err != nil {
				panic(err)
			}
			if i < min {
				return reflect.ValueOf(float32(i)),
					errors.New(fmt.Sprintf("%f less than minimum: %f", i, min))
			}
		}

		if maxStr, ok := t.Lookup("maxVal"); ok {
			max, err := strconv.ParseFloat(maxStr, 32)
			if err != nil {
				panic(err)
			}
			if i > max {
				return reflect.ValueOf(float32(i)),
					errors.New(fmt.Sprintf("%f greater than than maximum: %f", i, max))
			}
		}

		return reflect.ValueOf(float32(i)), nil
	},
	reflect.TypeOf(float64(0.0)): func(s string, t reflect.StructTag) (reflect.Value, error) {
		i, err := strconv.ParseFloat(s, 64)

		if err != nil {
			return reflect.ValueOf(float64(i)), err
		}

		if minStr, ok := t.Lookup("minVal"); ok {
			min, err := strconv.ParseFloat(minStr, 64)
			if err != nil {
				panic(err)
			}
			if i < min {
				return reflect.ValueOf(float64(i)),
					errors.New(fmt.Sprintf("%f less than minimum: %f", i, min))
			}
		}

		if maxStr, ok := t.Lookup("maxVal"); ok {
			max, err := strconv.ParseFloat(maxStr, 64)
			if err != nil {
				panic(err)
			}
			if i > max {
				return reflect.ValueOf(float64(i)),
					errors.New(fmt.Sprintf("%f greater than than maximum: %f", i, max))
			}
		}

		return reflect.ValueOf(float64(i)), nil
	},

	reflect.TypeOf(false): func(s string, t reflect.StructTag) (reflect.Value, error) {
		_, invert := t.Lookup("invert")
		b, err := strconv.ParseBool(s)

		if err == nil {
			return reflect.ValueOf(invert != b), nil
		} else {
			return reflect.ValueOf(false), err
		}
	},

	reflect.TypeOf(""): func(s string, t reflect.StructTag) (reflect.Value, error) {
		_, path := t.Lookup("path")
		if path {
			_, err := os.Stat(s)
			return reflect.ValueOf(s), err
		}

		// TODO: add file, dir, executable, socket, etc. available in unix test

		return reflect.ValueOf(s), nil
	},

	reflect.TypeOf([]byte{}): func(s string, t reflect.StructTag) (reflect.Value, error) {
		bytes, err := hex.DecodeString(s)
		return reflect.ValueOf(bytes), err
	},

	reflect.TypeOf(time.Duration(0)): func(s string, t reflect.StructTag) (reflect.Value, error) {
		d, err := time.ParseDuration(s)
		return reflect.ValueOf(d), err
	},

	reflect.TypeOf(net.IP([]byte{})): func(s string, t reflect.StructTag) (reflect.Value, error) {
		ip := net.ParseIP(s)
		if ip == nil {
			return reflect.ValueOf(nil), errors.New(fmt.Sprintf(
				"string \"%s\" is not a valid IP address", s))
		} else {
			return reflect.ValueOf(ip), nil
		}
	},
	reflect.TypeOf(net.IPNet{}): func(s string, t reflect.StructTag) (reflect.Value, error) {
		_, ipNet, err := net.ParseCIDR(s)
		return reflect.ValueOf(ipNet), err
	},
}

var valuelessUnmarshallers = map[reflect.Type]ValuelessUnmarshaller{
	reflect.TypeOf(false): func(_ reflect.Value, t reflect.StructTag) (reflect.Value, error) {
		_, invert := t.Lookup("invert")

		return reflect.ValueOf(!invert), nil
	},

	reflect.TypeOf(int(0)): func(v reflect.Value, t reflect.StructTag) (reflect.Value, error) {
		return reflect.ValueOf(int(v.Int() + 1)), nil
	},
	reflect.TypeOf(int8(0)): func(v reflect.Value, t reflect.StructTag) (reflect.Value, error) {
		return reflect.ValueOf(int8(v.Int() + 1)), nil
	},
	reflect.TypeOf(int16(0)): func(v reflect.Value, t reflect.StructTag) (reflect.Value, error) {
		return reflect.ValueOf(int16(v.Int() + 1)), nil
	},
	reflect.TypeOf(int32(0)): func(v reflect.Value, t reflect.StructTag) (reflect.Value, error) {
		return reflect.ValueOf(int32(v.Int() + 1)), nil
	},
	reflect.TypeOf(int64(0)): func(v reflect.Value, t reflect.StructTag) (reflect.Value, error) {
		return reflect.ValueOf(v.Int() + 1), nil
	},
	reflect.TypeOf(uint(0)): func(v reflect.Value, t reflect.StructTag) (reflect.Value, error) {
		return reflect.ValueOf(uint(v.Uint() + 1)), nil
	},
	reflect.TypeOf(uint8(0)): func(v reflect.Value, t reflect.StructTag) (reflect.Value, error) {
		return reflect.ValueOf(uint8(v.Uint() + 1)), nil
	},
	reflect.TypeOf(uint16(0)): func(v reflect.Value, t reflect.StructTag) (reflect.Value, error) {
		return reflect.ValueOf(uint16(v.Uint() + 1)), nil
	},
	reflect.TypeOf(uint32(0)): func(v reflect.Value, t reflect.StructTag) (reflect.Value, error) {
		return reflect.ValueOf(uint32(v.Uint() + 1)), nil
	},
	reflect.TypeOf(uint64(0)): func(v reflect.Value, t reflect.StructTag) (reflect.Value, error) {
		return reflect.ValueOf(v.Uint() + 1), nil
	},
}
