package magiccfg

import (
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/caarlos0/env/v11"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-yaml"
	"github.com/jessevdk/go-flags"
)

// MagicConfig is a generic that is temporarily used to take a user-defined struct, which may be nested
// with other structs to any degree, and allow a user to request that any combination of defaults, environment
// variables, command-line flags, and yaml/toml/json/xml files can be used to set values for the non-struct
// fields of the given configuration struct.  After construction, this object is eligible for garbage collection.
type MagicConfig[T any] struct {
	// types is a struct that holds all data types that are natively handled by magiccfg.
	types                *types
	envPrefix            string
	zeroPointer          string
	listSeparator        string
	keyValueSeparator    string
	ignoreUnknownOptions bool
	originalConfig       *T
	newConfig            any
	validateConfig       any
	validator            *validator.Validate
	transformFuncs       []func(reflect.StructField, reflect.Value)
	timeFormats          []string
	constructionErrors   []error
	configFiles          []string
	configFileVars       map[string]string
	remainingArgs        []string
	debugMenu            bool
	customTypes          map[string]*customType
}

func NewMagicConfig[T any](config *T, options *Options) (*MagicConfig[T], error) {
	if options == nil {
		options = &Options{}
	}

	typeStruct := newTypes()

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

	prefix := strings.ToUpper(options.EnvPrefix)

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

	c := &MagicConfig[T]{
		types:                typeStruct,
		originalConfig:       config,
		envPrefix:            prefix,
		zeroPointer:          zeroPointer,
		listSeparator:        listSeparator,
		keyValueSeparator:    keyValueSeparator,
		ignoreUnknownOptions: ignoreUnknownOptions,
		validator:            vd,
		transformFuncs:       transformFuncs,
		timeFormats:          timeFormats,
		constructionErrors:   []error{},
		debugMenu:            options.DebugMenu,
		remainingArgs:        remainingArgs,
		configFiles:          configFiles,
		configFileVars: map[string]string{
			_envTag: options.ConfigFileEnvName,
			"short": options.ConfigFileShort,
			"long":  options.ConfigFileLong,
		},
		customTypes: options.CustomTypes,
	}

	fc, usages := c.makeFullConfig(configVal, prefix, nil, nil)
	usages[_envTag][options.ConfigFileEnvName]++
	usages["long"][options.ConfigFileLong]++
	usages["short"][options.ConfigFileShort]++

	vc := c.makeValidateConfig(configVal)
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

	c.newConfig = fullConfig
	c.validateConfig = vc.Interface()

	return c, nil
}

func (c *MagicConfig[T]) ParseEnv() *MagicConfig[T] {
	newConfig := c.empty().Interface()
	funcMap := map[reflect.Type]env.ParserFunc{
		c.types.stringPtrMapType:              c.parseStringPointerMap,
		c.types.intPtrMapType:                 c.parseIntPointerMap,
		c.types.int8PtrMapType:                c.parseInt8PointerMap,
		c.types.int16PtrMapType:               c.parseInt16PointerMap,
		c.types.int32PtrMapType:               c.parseInt32PointerMap,
		c.types.int64PtrMapType:               c.parseInt64PointerMap,
		c.types.uintPtrMapType:                c.parseUintPointerMap,
		c.types.uint8PtrMapType:               c.parseUint8PointerMap,
		c.types.uint16PtrMapType:              c.parseUint16PointerMap,
		c.types.uint32PtrMapType:              c.parseUint32PointerMap,
		c.types.uint64PtrMapType:              c.parseUint64PointerMap,
		c.types.float32PtrMapType:             c.parseFloat32PointerMap,
		c.types.float64PtrMapType:             c.parseFloat64PointerMap,
		c.types.boolPtrMapType:                c.parseBoolPointerMap,
		c.types.parseableComplex64Type:        c.parseParseableComplex64,
		c.types.parseableComplex128Type:       c.parseParseableComplex128,
		c.types.parseableComplex64MapPtrType:  c.parseComplex64MapPtr,
		c.types.parseableComplex128MapPtrType: c.parseComplex128MapPtr,
		c.types.parseableDurationMapPtrType:   c.parseDurationPointerMap,
		c.types.parseableTimeMapPtrType:       c.parseTimePointerMap,
		c.types.parseableURLMapType:           c.parseURLMap,
		c.types.parseableURLMapPtrType:        c.parseURLPointerMap,
		c.types.wrappedBoolType:               c.parseWrappedBoolType,
		c.types.parseableDurationType:         UnmarshalDurationEnv,
		c.types.parseableTimeType:             UnmarshalTimeEnv,
	}

	for _, ct := range c.customTypes {
		if ct.environmentParser != nil {
			funcMap[ct.basicType] = ct.environmentParser
		}
		if ct.environmentListParser != nil {
			typeSlice := reflect.MakeSlice(reflect.SliceOf(ct.basicType), 0, 0)
			funcMap[typeSlice.Type()] = ct.environmentListParser
		}
		if ct.environmentListPointerParser != nil {
			typeSlice := reflect.MakeSlice(reflect.SliceOf(reflect.PointerTo(ct.basicType)), 0, 0)
			funcMap[typeSlice.Type()] = ct.environmentListPointerParser
		}
		if ct.environmentMapParser != nil {
			typeMap := reflect.MakeMap(reflect.MapOf(c.types.stringType, ct.basicType))
			funcMap[typeMap.Type()] = ct.environmentMapParser
		}
		if ct.environmentMapPointerParser != nil {
			typeMap := reflect.MakeMap(reflect.MapOf(c.types.stringType, reflect.PointerTo(ct.basicType)))
			funcMap[typeMap.Type()] = ct.environmentMapPointerParser
		}
	}

	if err := env.ParseWithOptions(
		newConfig,
		env.Options{
			TagName: _envTag,
			FuncMap: funcMap,
		},
	); err != nil {
		c.constructionErrors = append(c.constructionErrors, err)
	}
	c.merge(reflect.ValueOf(c.newConfig), reflect.ValueOf(newConfig))
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
		c.merge(reflect.ValueOf(newConfig), reflect.ValueOf(nc))
	}
	c.merge(reflect.ValueOf(c.newConfig), reflect.ValueOf(newConfig))
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
	args := c.splitArgs(reflect.ValueOf(newConfig), c.remainingArgs, c.listSeparator)
	if _, err := parser.ParseArgs(args); err != nil {
		var e *flags.Error
		if errors.As(err, &e) && errors.Is(e.Type, flags.ErrHelp) {
			e.Message = "Help Menu Displayed"
		}
		c.constructionErrors = append(c.constructionErrors, err)
	}
	c.merge(reflect.ValueOf(c.newConfig), reflect.ValueOf(newConfig))
	return c
}

func (c *MagicConfig[T]) Validate() *MagicConfig[T] {
	c.transformValues(reflect.ValueOf(c.newConfig), c.transformFuncs)
	if errs := c.validateEnums(reflect.ValueOf(c.newConfig)); len(errs) > 0 {
		c.constructionErrors = append(c.constructionErrors, errs...)
	}
	c.merge(reflect.ValueOf(c.validateConfig), reflect.ValueOf(c.newConfig))

	if err := c.validator.Struct(c.validateConfig); err != nil {
		c.constructionErrors = append(c.constructionErrors, err)
	}
	c.merge(reflect.ValueOf(c.newConfig), reflect.ValueOf(c.validateConfig))
	return c
}

func (c *MagicConfig[T]) ApplyDefaults() *MagicConfig[T] {
	c.setDefaults(
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
	c.merge(o, f)
	n, ok := o.Interface().(*T)
	if !ok {
		c.constructionErrors = append(c.constructionErrors, fmt.Errorf("could not construct final configuration"))
		return nil, c.constructionErrors
	}
	return n, c.constructionErrors
}
