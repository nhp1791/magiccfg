// Package magiccfg provides
package magiccfg

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/caarlos0/env/v11"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-yaml"
	"github.com/jessevdk/go-flags"
)

const (
	_zeroPointer       = "*0*"
	_listSeparator     = ","
	_keyValueSeparator = ":"
)

type MagicConfig[T any] struct {
	envPrefix                string
	zeroPointer              string
	listSeparator            string
	keyValueSeparator        string
	ignoreUnknownOptions     bool
	originalConfig           *T
	newConfig                any
	validateConfig           any
	validator                *validator.Validate
	transformFuncs           []func(reflect.StructField, reflect.Value)
	timeFormats              []string
	constructionErrors       []error
	configFiles              []string
	configFileVars           map[string]string
	remainingArgs            []string
	commandLineNameConverter func(string, *string) string
}

type Options struct {
	ConfigFileEnvName        string
	ConfigFileShort          string
	ConfigFileLong           string
	Prefix                   string
	ZeroPointer              string
	ListSeparator            string
	KeyValueSeparator        string
	IgnoreUnknownOptions     bool
	ValidationFunctions      map[string]ValidationFunction
	TransformFuncs           []func(reflect.StructField, reflect.Value)
	TimeFormats              []string
	CommandLineNameConverter func(string, *string) string
}

func NewMagicConfig[T any](config *T, options *Options) (*MagicConfig[T], error) {
	if options == nil {
		options = &Options{}
	}

	zeroPointer := _zeroPointer
	if options.ZeroPointer != "" {
		zeroPointer = options.ZeroPointer
	}

	listSeparator := _listSeparator
	if options.ListSeparator != "" {
		listSeparator = options.ListSeparator
	}

	var configFiles []string
	var remainingArgs []string
	var cliFiles []string
	if options.ConfigFileEnvName != "" {
		fileEnv, ok := os.LookupEnv(options.ConfigFileEnvName)
		if ok && fileEnv != "" {
			configFiles = strings.Split(fileEnv, listSeparator)
		}
	}
	if options.ConfigFileLong != "" || options.ConfigFileShort != "" {
		cliFiles, remainingArgs = locateCLIFiles(options, listSeparator)
		configFiles = append(configFiles, cliFiles...)
	}

	configVal := reflect.ValueOf(config)
	if configVal.Kind() != reflect.Ptr || configVal.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("config must be a pointer to a struct")
	}

	prefix := strings.ToUpper(options.Prefix)

	keyValueSeparator := _keyValueSeparator
	if options.KeyValueSeparator != "" {
		keyValueSeparator = options.KeyValueSeparator
	}

	ignoreUnknownOptions := options.IgnoreUnknownOptions

	var timeFormats []string
	if options.TimeFormats != nil {
		timeFormats = options.TimeFormats
	}
	packageTimeFormats = timeFormats

	vd := validator.New(validator.WithRequiredStructEnabled())
	for n, f := range options.ValidationFunctions {
		if err := vd.RegisterValidation(n, (func(p validator.FieldLevel) bool)(f)); err != nil {
			return nil, err
		}
	}

	transformFuncs := options.TransformFuncs

	fc, usages := makeFullConfig(configVal, prefix, options.CommandLineNameConverter, nil, nil)
	usages[_envTag][options.ConfigFileEnvName]++
	usages["long"][options.ConfigFileLong]++
	usages["short"][options.ConfigFileShort]++

	vc := makeValidateConfig(configVal)
	conflict := false
	var sb strings.Builder
	sb.WriteString("Cannot construct configuration due to the following conflicts:")
	for key, v := range usages {
		for k, n := range v {
			if k != "" && n > 1 {
				conflict = true
				sb.WriteString(fmt.Sprintf("\n%s variable %s defined %d times", key, k, n))
			}
			if key == "short" && k == "h" && n > 0 {
				conflict = true
				sb.WriteString(fmt.Sprintf("\nshort variable h conflicts with help"))
			}
			if key == "long" && k == "help" && n > 0 {
				conflict = true
				sb.WriteString(fmt.Sprintf("\nlong variable help conflicts with help"))
			}
		}
	}

	if conflict {
		return nil, errors.New(sb.String())
	}

	fullConfig := fc.Interface()

	return &MagicConfig[T]{
		originalConfig:       config,
		newConfig:            fullConfig,
		validateConfig:       vc.Interface(),
		envPrefix:            prefix,
		zeroPointer:          zeroPointer,
		listSeparator:        listSeparator,
		keyValueSeparator:    keyValueSeparator,
		ignoreUnknownOptions: ignoreUnknownOptions,
		validator:            vd,
		transformFuncs:       transformFuncs,
		timeFormats:          timeFormats,
		constructionErrors:   []error{},
		remainingArgs:        remainingArgs,
		configFiles:          configFiles,
		configFileVars: map[string]string{
			_envTag: options.ConfigFileEnvName,
			"short": options.ConfigFileShort,
			"long":  options.ConfigFileLong,
		},
		commandLineNameConverter: options.CommandLineNameConverter,
	}, nil
}

func (c *MagicConfig[T]) ParseEnv() *MagicConfig[T] {
	newConfig := c.empty().Interface()
	if err := env.ParseWithOptions(
		newConfig,
		env.Options{
			TagName: _envTag,
			FuncMap: map[reflect.Type]env.ParserFunc{
				stringPtrMapType:              c.ParseStringPointerMap,
				intPtrMapType:                 c.ParseIntPointerMap,
				int8PtrMapType:                c.ParseInt8PointerMap,
				int16PtrMapType:               c.ParseInt16PointerMap,
				int32PtrMapType:               c.ParseInt32PointerMap,
				int64PtrMapType:               c.ParseInt64PointerMap,
				uintPtrMapType:                c.ParseUintPointerMap,
				uint8PtrMapType:               c.ParseUint8PointerMap,
				uint16PtrMapType:              c.ParseUint16PointerMap,
				uint32PtrMapType:              c.ParseUint32PointerMap,
				uint64PtrMapType:              c.ParseUint64PointerMap,
				float32PtrMapType:             c.ParseFloat32PointerMap,
				float64PtrMapType:             c.ParseFloat64PointerMap,
				boolPtrMapType:                c.ParseBoolPointerMap,
				wrappedBoolPtrMapType:         c.ParseWrappedBoolPointerMap,
				parseableComplex64Type:        c.ParseParseableComplex64,
				parseableComplex128Type:       c.ParseParseableComplex128,
				parseableComplex64MapPtrType:  c.ParseComplex64MapPtr,
				parseableComplex128MapPtrType: c.ParseComplex128MapPtr,
				parseableDurationMapPtrType:   c.ParseDurationPointerMap,
				parseableTimeMapPtrType:       c.ParseTimePointerMap,
				parseableURLMapType:           c.ParseURLMap,
				parseableURLMapPtrType:        c.ParseURLPointerMap,
				parseableDurationType:         UnmarshalDurationEnv,
				parseableTimeType:             UnmarshalTimeEnv,
			},
		},
	); err != nil {
		c.constructionErrors = append(c.constructionErrors, err)
	}
	merge(reflect.ValueOf(c.newConfig), reflect.ValueOf(newConfig))
	return c
}

func (c *MagicConfig[T]) ParseFiles() *MagicConfig[T] {
	if len(c.configFiles) == 0 {
		return c
	}
	newConfig := c.empty().Interface()
	for _, f := range c.configFiles {
		data, err := os.ReadFile(f)
		if err != nil {
			c.constructionErrors = append(c.constructionErrors, err)
			continue
		}

		test := map[string]any{}
		nc := c.empty().Interface()
		isYAML := true
		isXML := false
		if err := yaml.Unmarshal(data, &test); err != nil {
			if err := toml.Unmarshal(data, &test); err != nil {
				if err := xml.Unmarshal(data, nc); err != nil {
					println(err.Error())
					c.constructionErrors = append(c.constructionErrors, fmt.Errorf("failed to recognize file %s as YAML, JSON, TOML, or XML", f))
					continue
				} else {
					isXML = true
				}
			}
			isYAML = false
		}

		if isYAML {
			if err := yaml.Unmarshal(data, nc); err != nil {
				c.constructionErrors = append(c.constructionErrors, err)
			}
		} else if !isXML {
			if err := toml.Unmarshal(data, nc); err != nil {
				c.constructionErrors = append(c.constructionErrors, err)
			}
		}
		merge(reflect.ValueOf(newConfig), reflect.ValueOf(nc))
	}
	merge(reflect.ValueOf(c.newConfig), reflect.ValueOf(newConfig))
	return c
}

func (c *MagicConfig[T]) ParseFlags() *MagicConfig[T] {
	if len(c.remainingArgs) == 0 {
		return c
	}
	newConfig := c.empty().Interface()
	var parser *flags.Parser
	if c.ignoreUnknownOptions {
		parser = flags.NewParser(newConfig, flags.Default|flags.IgnoreUnknown)
	} else {
		parser = flags.NewParser(newConfig, flags.Default)
	}
	args := splitArgs(reflect.ValueOf(newConfig), c.remainingArgs, c.listSeparator)
	if _, err := parser.ParseArgs(args); err != nil {
		var e *flags.Error
		if errors.As(err, &e) && errors.Is(e.Type, flags.ErrHelp) {
			e.Message = "Help Menu Displayed"
		}
		c.constructionErrors = append(c.constructionErrors, err)
	}
	merge(reflect.ValueOf(c.newConfig), reflect.ValueOf(newConfig))
	return c
}

func (c *MagicConfig[T]) Validate() *MagicConfig[T] {
	transformValues(reflect.ValueOf(c.newConfig), c.transformFuncs)
	if errs := validateEnums(reflect.ValueOf(c.newConfig)); len(errs) > 0 {
		c.constructionErrors = append(c.constructionErrors, errs...)
	}
	merge(reflect.ValueOf(c.validateConfig), reflect.ValueOf(c.newConfig))

	if err := c.validator.Struct(c.validateConfig); err != nil {
		c.constructionErrors = append(c.constructionErrors, err)
	}
	merge(reflect.ValueOf(c.newConfig), reflect.ValueOf(c.validateConfig))
	return c
}

func (c *MagicConfig[T]) ApplyDefaults() *MagicConfig[T] {
	setDefaults(
		reflect.ValueOf(c.newConfig),
		c.zeroPointer,
		c.listSeparator,
		c.keyValueSeparator,
		c.timeFormats,
	)
	return c
}

func (c *MagicConfig[T]) Config() (*T, []error) {
	f := reflect.ValueOf(c.newConfig)
	o := reflect.ValueOf(c.originalConfig)
	merge(o, f)
	n, ok := o.Interface().(*T)
	if !ok {
		c.constructionErrors = append(c.constructionErrors, fmt.Errorf("could not construct final configuration"))
		return nil, c.constructionErrors
	}
	return n, c.constructionErrors
}

func (c *MagicConfig[T]) ParseStringPointerMap(v string) (any, error) {
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

func (c *MagicConfig[T]) ParseBoolPointerMap(v string) (any, error) {
	pt := reflect.PointerTo(boolType)

	resultType := reflect.MapOf(stringType, pt)
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

func (c *MagicConfig[T]) ParseWrappedBoolPointerMap(v string) (any, error) {
	pt := reflect.PointerTo(wrappedBoolType)

	resultType := reflect.MapOf(stringType, pt)
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
		w := Bool(b)
		result.SetMapIndex(reflect.ValueOf(vals[0]), reflect.ValueOf(&w))
	}

	r, ok := result.Interface().(map[string]*Bool)
	if !ok {
		return nil, fmt.Errorf("failed to parse Bool pointer map")
	}
	return r, nil
}

func (c *MagicConfig[T]) ParseIntPointerMap(v string) (any, error) {
	result := c.parsePointerMap(v, intType)
	r, ok := result.Interface().(map[string]*int)
	if !ok {
		return nil, fmt.Errorf("failed to parse int pointer map")
	}
	return r, nil
}

func (c *MagicConfig[T]) ParseInt8PointerMap(v string) (any, error) {
	result := c.parsePointerMap(v, int8Type)
	r, ok := result.Interface().(map[string]*int8)
	if !ok {
		return nil, fmt.Errorf("failed to parse int pointer map")
	}
	return r, nil
}

func (c *MagicConfig[T]) ParseInt16PointerMap(v string) (any, error) {
	result := c.parsePointerMap(v, int16Type)
	r, ok := result.Interface().(map[string]*int16)
	if !ok {
		return nil, fmt.Errorf("failed to parse int pointer map")
	}
	return r, nil
}

func (c *MagicConfig[T]) ParseInt32PointerMap(v string) (any, error) {
	result := c.parsePointerMap(v, int32Type)
	r, ok := result.Interface().(map[string]*int32)
	if !ok {
		return nil, fmt.Errorf("failed to parse int pointer map")
	}
	return r, nil
}

func (c *MagicConfig[T]) ParseInt64PointerMap(v string) (any, error) {
	result := c.parsePointerMap(v, int64Type)
	r, ok := result.Interface().(map[string]*int64)
	if !ok {
		return nil, fmt.Errorf("failed to parse int pointer map")
	}
	return r, nil
}

func (c *MagicConfig[T]) ParseUintPointerMap(v string) (any, error) {
	result := c.parsePointerMap(v, uintType)
	r, ok := result.Interface().(map[string]*uint)
	if !ok {
		return nil, fmt.Errorf("failed to parse int pointer map")
	}
	return r, nil
}

func (c *MagicConfig[T]) ParseUint8PointerMap(v string) (any, error) {
	result := c.parsePointerMap(v, uint8Type)
	r, ok := result.Interface().(map[string]*uint8)
	if !ok {
		return nil, fmt.Errorf("failed to parse int pointer map")
	}
	return r, nil
}

func (c *MagicConfig[T]) ParseUint16PointerMap(v string) (any, error) {
	result := c.parsePointerMap(v, uint16Type)
	r, ok := result.Interface().(map[string]*uint16)
	if !ok {
		return nil, fmt.Errorf("failed to parse int pointer map")
	}
	return r, nil
}

func (c *MagicConfig[T]) ParseUint32PointerMap(v string) (any, error) {
	result := c.parsePointerMap(v, uint32Type)
	r, ok := result.Interface().(map[string]*uint32)
	if !ok {
		return nil, fmt.Errorf("failed to parse int pointer map")
	}
	return r, nil
}

func (c *MagicConfig[T]) ParseUint64PointerMap(v string) (any, error) {
	result := c.parsePointerMap(v, uint64Type)
	r, ok := result.Interface().(map[string]*uint64)
	if !ok {
		return nil, fmt.Errorf("failed to parse int pointer map")
	}
	return r, nil
}

func (c *MagicConfig[T]) ParseFloat32PointerMap(v string) (any, error) {
	result := c.parsePointerMap(v, float32Type)
	r, ok := result.Interface().(map[string]*float32)
	if !ok {
		return nil, fmt.Errorf("failed to parse float pointer map")
	}
	return r, nil
}

func (c *MagicConfig[T]) ParseFloat64PointerMap(v string) (any, error) {
	result := c.parsePointerMap(v, float64Type)
	r, ok := result.Interface().(map[string]*float64)
	if !ok {
		return nil, fmt.Errorf("failed to parse float pointer map")
	}
	return r, nil
}

func (c *MagicConfig[T]) ParseComplex64MapPtr(v string) (any, error) {
	result := c.parsePointerMap(v, parseableComplex64Type)
	r, ok := result.Interface().(map[string]*ParseableComplex64)
	if !ok {
		return nil, fmt.Errorf("failed to parse complex pointer map")
	}
	return r, nil
}

func (c *MagicConfig[T]) ParseComplex128MapPtr(v string) (any, error) {
	result := c.parsePointerMap(v, parseableComplex128Type)
	r, ok := result.Interface().(map[string]*ParseableComplex128)
	if !ok {
		return nil, fmt.Errorf("failed to parse complex pointer map")
	}
	return r, nil
}

func (c *MagicConfig[T]) parsePointerMap(v string, t reflect.Type) reflect.Value {
	pt := reflect.PointerTo(t)

	resultType := reflect.MapOf(stringType, pt)
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

func (c *MagicConfig[T]) ParseParseableComplex64(v string) (any, error) {
	n, err := strconv.ParseComplex(v, 64)
	if err != nil {
		return nil, err
	}
	r := ParseableComplex64(n)
	return r, nil
}

func (c *MagicConfig[T]) ParseParseableComplex128(v string) (any, error) {
	n, err := strconv.ParseComplex(v, 128)
	if err != nil {
		return nil, err
	}
	r := ParseableComplex128(n)
	return r, nil
}

func (c *MagicConfig[T]) ParseDurationPointerMap(v string) (any, error) {
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

func (c *MagicConfig[T]) ParseTimePointerMap(v string) (any, error) {
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

func (c *MagicConfig[T]) ParseURLMap(v string) (any, error) {
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

func (c *MagicConfig[T]) ParseURLPointerMap(v string) (any, error) {
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
