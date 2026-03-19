package magiccfg

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func (c *MagicConfig[T]) parseStringPointerMap(v string) (any, error) {
	result := make(map[string]*string)

	kvs := strings.Split(v, c.listSeparator)
	for _, kv := range kvs {
		vals := strings.Split(kv, c.keyValueSeparator)
		if len(vals) != 2 {
			continue
		}
		result[vals[0]] = &vals[1]
	}

	return result, nil
}

func (c *MagicConfig[T]) parseBoolPointerMap(v string) (any, error) {
	pt := reflect.PointerTo(c.types.boolType)

	resultType := reflect.MapOf(c.types.stringType, pt)
	result := reflect.MakeMap(resultType)

	kvs := strings.Split(v, c.listSeparator)
	for _, kv := range kvs {
		vals := strings.Split(kv, c.keyValueSeparator)
		if len(vals) != 2 {
			continue
		}

		b, err := strconv.ParseBool(vals[1])
		if err != nil {
			continue
		}
		result.SetMapIndex(reflect.ValueOf(vals[0]), reflect.ValueOf(&b))
	}

	r, ok := result.Interface().(map[string]*bool)
	if !ok {
		return nil, fmt.Errorf("failed to parse bool pointer map")
	}
	return r, nil
}

func (c *MagicConfig[T]) parseIntPointerMap(v string) (any, error) {
	result := c.parsePointerMap(v, c.types.intType)
	r, ok := result.Interface().(map[string]*int)
	if !ok {
		return nil, fmt.Errorf("failed to parse int pointer map")
	}
	return r, nil
}

func (c *MagicConfig[T]) parseInt8PointerMap(v string) (any, error) {
	result := c.parsePointerMap(v, c.types.int8Type)
	r, ok := result.Interface().(map[string]*int8)
	if !ok {
		return nil, fmt.Errorf("failed to parse int pointer map")
	}
	return r, nil
}

func (c *MagicConfig[T]) parseInt16PointerMap(v string) (any, error) {
	result := c.parsePointerMap(v, c.types.int16Type)
	r, ok := result.Interface().(map[string]*int16)
	if !ok {
		return nil, fmt.Errorf("failed to parse int pointer map")
	}
	return r, nil
}

func (c *MagicConfig[T]) parseInt32PointerMap(v string) (any, error) {
	result := c.parsePointerMap(v, c.types.int32Type)
	r, ok := result.Interface().(map[string]*int32)
	if !ok {
		return nil, fmt.Errorf("failed to parse int pointer map")
	}
	return r, nil
}

func (c *MagicConfig[T]) parseInt64PointerMap(v string) (any, error) {
	result := c.parsePointerMap(v, c.types.int64Type)
	r, ok := result.Interface().(map[string]*int64)
	if !ok {
		return nil, fmt.Errorf("failed to parse int pointer map")
	}
	return r, nil
}

func (c *MagicConfig[T]) parseUintPointerMap(v string) (any, error) {
	result := c.parsePointerMap(v, c.types.uintType)
	r, ok := result.Interface().(map[string]*uint)
	if !ok {
		return nil, fmt.Errorf("failed to parse int pointer map")
	}
	return r, nil
}

func (c *MagicConfig[T]) parseUint8PointerMap(v string) (any, error) {
	result := c.parsePointerMap(v, c.types.uint8Type)
	r, ok := result.Interface().(map[string]*uint8)
	if !ok {
		return nil, fmt.Errorf("failed to parse int pointer map")
	}
	return r, nil
}

func (c *MagicConfig[T]) parseUint16PointerMap(v string) (any, error) {
	result := c.parsePointerMap(v, c.types.uint16Type)
	r, ok := result.Interface().(map[string]*uint16)
	if !ok {
		return nil, fmt.Errorf("failed to parse int pointer map")
	}
	return r, nil
}

func (c *MagicConfig[T]) parseUint32PointerMap(v string) (any, error) {
	result := c.parsePointerMap(v, c.types.uint32Type)
	r, ok := result.Interface().(map[string]*uint32)
	if !ok {
		return nil, fmt.Errorf("failed to parse int pointer map")
	}
	return r, nil
}

func (c *MagicConfig[T]) parseUint64PointerMap(v string) (any, error) {
	result := c.parsePointerMap(v, c.types.uint64Type)
	r, ok := result.Interface().(map[string]*uint64)
	if !ok {
		return nil, fmt.Errorf("failed to parse int pointer map")
	}
	return r, nil
}

func (c *MagicConfig[T]) parseFloat32PointerMap(v string) (any, error) {
	result := c.parsePointerMap(v, c.types.float32Type)
	r, ok := result.Interface().(map[string]*float32)
	if !ok {
		return nil, fmt.Errorf("failed to parse float pointer map")
	}
	return r, nil
}

func (c *MagicConfig[T]) parseFloat64PointerMap(v string) (any, error) {
	result := c.parsePointerMap(v, c.types.float64Type)
	r, ok := result.Interface().(map[string]*float64)
	if !ok {
		return nil, fmt.Errorf("failed to parse float pointer map")
	}
	return r, nil
}

func (c *MagicConfig[T]) parseComplex64MapPtr(v string) (any, error) {
	result := c.parsePointerMap(v, c.types.parseableComplex64Type)
	r, ok := result.Interface().(map[string]*ParseableComplex64)
	if !ok {
		return nil, fmt.Errorf("failed to parse complex pointer map")
	}
	return r, nil
}

func (c *MagicConfig[T]) parseComplex128MapPtr(v string) (any, error) {
	result := c.parsePointerMap(v, c.types.parseableComplex128Type)
	r, ok := result.Interface().(map[string]*ParseableComplex128)
	if !ok {
		return nil, fmt.Errorf("failed to parse complex pointer map")
	}
	return r, nil
}

func (c *MagicConfig[T]) parsePointerMap(v string, t reflect.Type) reflect.Value {
	pt := reflect.PointerTo(t)

	resultType := reflect.MapOf(c.types.stringType, pt)
	result := reflect.MakeMap(resultType)

	k := t.Kind()

	kvs := strings.Split(v, c.listSeparator)
	for _, kv := range kvs {
		vals := strings.Split(kv, c.keyValueSeparator)
		if len(vals) != 2 {
			continue
		}
		switch k {
		case reflect.Int:
			i, err := strconv.Atoi(vals[1])
			if err != nil {
				continue
			}
			result.SetMapIndex(reflect.ValueOf(vals[0]), reflect.ValueOf(&i))
		case reflect.Int8:
			i, err := strconv.ParseInt(vals[1], 10, 8)
			if err != nil {
				continue
			}
			j := int8(i)
			result.SetMapIndex(reflect.ValueOf(vals[0]), reflect.ValueOf(&j))
		case reflect.Int16:
			i, err := strconv.ParseInt(vals[1], 10, 16)
			if err != nil {
				continue
			}
			j := int16(i)
			result.SetMapIndex(reflect.ValueOf(vals[0]), reflect.ValueOf(&j))
		case reflect.Int32:
			i, err := strconv.ParseInt(vals[1], 10, 32)
			if err != nil {
				continue
			}
			j := int32(i)
			result.SetMapIndex(reflect.ValueOf(vals[0]), reflect.ValueOf(&j))
		case reflect.Int64:
			i, err := strconv.ParseInt(vals[1], 10, 64)
			if err != nil {
				continue
			}
			j := int64(i)
			result.SetMapIndex(reflect.ValueOf(vals[0]), reflect.ValueOf(&j))
		case reflect.Uint:
			i, err := strconv.ParseUint(vals[1], 10, 64)
			if err != nil {
				continue
			}
			j := uint(i)
			result.SetMapIndex(reflect.ValueOf(vals[0]), reflect.ValueOf(&j))
		case reflect.Uint8:
			i, err := strconv.ParseUint(vals[1], 10, 8)
			if err != nil {
				continue
			}
			j := uint8(i)
			result.SetMapIndex(reflect.ValueOf(vals[0]), reflect.ValueOf(&j))
		case reflect.Uint16:
			i, err := strconv.ParseUint(vals[1], 10, 16)
			if err != nil {
				continue
			}
			j := uint16(i)
			result.SetMapIndex(reflect.ValueOf(vals[0]), reflect.ValueOf(&j))
		case reflect.Uint32:
			i, err := strconv.ParseUint(vals[1], 10, 32)
			if err != nil {
				continue
			}
			j := uint32(i)
			result.SetMapIndex(reflect.ValueOf(vals[0]), reflect.ValueOf(&j))
		case reflect.Uint64:
			i, err := strconv.ParseUint(vals[1], 10, 64)
			if err != nil {
				continue
			}
			result.SetMapIndex(reflect.ValueOf(vals[0]), reflect.ValueOf(&i))
		case reflect.Float32:
			i, err := strconv.ParseFloat(vals[1], 32)
			if err != nil {
				continue
			}
			j := float32(i)
			result.SetMapIndex(reflect.ValueOf(vals[0]), reflect.ValueOf(&j))
		case reflect.Float64:
			i, err := strconv.ParseFloat(vals[1], 64)
			if err != nil {
				continue
			}
			result.SetMapIndex(reflect.ValueOf(vals[0]), reflect.ValueOf(&i))
		case reflect.Complex64:
			i, err := strconv.ParseComplex(vals[1], 64)
			if err != nil {
				continue
			}
			j := ParseableComplex64(complex64(i))
			result.SetMapIndex(reflect.ValueOf(vals[0]), reflect.ValueOf(&j))
		case reflect.Complex128:
			i, err := strconv.ParseComplex(vals[1], 128)
			if err != nil {
				continue
			}
			j := ParseableComplex128(i)
			result.SetMapIndex(reflect.ValueOf(vals[0]), reflect.ValueOf(&j))
		default:
			return result
		}
	}

	return result
}

func (c *MagicConfig[T]) parseParseableComplex64(v string) (any, error) {
	n, err := strconv.ParseComplex(v, 64)
	if err != nil {
		return nil, err
	}
	r := ParseableComplex64(n)
	return r, nil
}

func (c *MagicConfig[T]) parseParseableComplex128(v string) (any, error) {
	n, err := strconv.ParseComplex(v, 128)
	if err != nil {
		return nil, err
	}
	r := ParseableComplex128(n)
	return r, nil
}

func (c *MagicConfig[T]) parseDurationPointerMap(v string) (any, error) {
	result := make(map[string]*ParseableDuration)

	kvs := strings.Split(v, c.listSeparator)
	for _, kv := range kvs {
		vals := strings.Split(kv, c.keyValueSeparator)
		if len(vals) != 2 {
			continue
		}
		d, err := time.ParseDuration(vals[1])
		if err != nil {
			return nil, err
		}
		p := ParseableDuration(d)
		result[vals[0]] = &p
	}

	return result, nil
}

func (c *MagicConfig[T]) parseTimePointerMap(v string) (any, error) {
	result := make(map[string]*ParseableTime)

	kvs := strings.Split(v, c.listSeparator)
	for _, kv := range kvs {
		vals := strings.Split(kv, c.keyValueSeparator)
		if (c.keyValueSeparator != ":" && len(vals) != 2) ||
			(c.keyValueSeparator == ":" && len(vals) != 4) {
			continue
		}
		cand := vals[1]
		if c.keyValueSeparator == ":" {
			cand = strings.Join(vals[1:], c.keyValueSeparator)
		}
		var parseErr error
		for _, format := range packageTimeFormats {
			if t, err := time.Parse(format, cand); err == nil {
				p := ParseableTime(t)
				result[vals[0]] = &p
				parseErr = nil
				break
			} else {
				parseErr = err
			}
		}
		if parseErr != nil {
			return nil, parseErr
		}
	}

	return result, nil
}

func (c *MagicConfig[T]) parseURLMap(v string) (any, error) {
	result := make(map[string]ParseableURL)

	kvs := strings.Split(v, c.listSeparator)
	for _, kv := range kvs {
		vals := strings.Split(kv, c.keyValueSeparator)
		size := len(vals)
		if (c.keyValueSeparator != ":" && size != 2) ||
			(c.keyValueSeparator == ":" && size != 3 && size != 4) {
			continue
		}
		cand := vals[1]
		if c.keyValueSeparator == ":" {
			cand = strings.Join(vals[1:], c.keyValueSeparator)
		}
		u, err := url.Parse(cand)
		if err != nil {
			return nil, err
		}

		p := ParseableURL(*u)
		result[vals[0]] = p
	}

	return result, nil
}

func (c *MagicConfig[T]) parseURLPointerMap(v string) (any, error) {
	result := make(map[string]*ParseableURL)

	kvs := strings.Split(v, c.listSeparator)
	for _, kv := range kvs {
		vals := strings.Split(kv, c.keyValueSeparator)
		size := len(vals)
		if (c.keyValueSeparator != ":" && size != 2) ||
			(c.keyValueSeparator == ":" && size != 3 && size != 4) {
			continue
		}
		cand := vals[1]
		if c.keyValueSeparator == ":" {
			cand = strings.Join(vals[1:], c.keyValueSeparator)
		}
		u, err := url.Parse(cand)
		if err != nil {
			return nil, err
		}

		p := ParseableURL(*u)
		result[vals[0]] = &p
	}

	return result, nil
}
