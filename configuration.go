// Package magiccfg provides
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

const (
	_zeroPointer       = "*0*"
	_listSeparator     = ","
	_keyValueSeparator = ":"
)

type MagicConfig[T any] struct {
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
}

type Options struct {
	ConfigFileEnvName    string
	ConfigFileShort      string
	ConfigFileLong       string
	Prefix               string
	ZeroPointer          string
	ListSeparator        string
	KeyValueSeparator    string
	IgnoreUnknownOptions bool
	ValidationFunctions  map[string]ValidationFunction
	TransformFuncs       []func(reflect.StructField, reflect.Value)
	TimeFormats          []string
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
	remainingArgs := os.Args
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

	fc, usages := makeFullConfig(configVal, prefix, nil, nil)
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
	}, nil
}

func (c *MagicConfig[T]) ParseEnv() *MagicConfig[T] {
	newConfig := c.empty().Interface()
	if err := env.ParseWithOptions(
		newConfig,
		env.Options{
			TagName: _envTag,
			FuncMap: map[reflect.Type]env.ParserFunc{
				parseableDurationType: UnmarshalDurationEnv,
				parseableTimeType:     UnmarshalTimeEnv,
			},
		},
	); err != nil {
		c.constructionErrors = append(c.constructionErrors, err)
	}
	merge(reflect.ValueOf(c.newConfig), reflect.ValueOf(newConfig))
	return c
}

func (c *MagicConfig[T]) ParseFiles() *MagicConfig[T] {
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
