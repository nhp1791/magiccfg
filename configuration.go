package magiccfg

import (
	"fmt"
	"github.com/caarlos0/env/v11"
	"github.com/go-playground/validator/v10"
	"github.com/jessevdk/go-flags"
	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
	"os"
	"reflect"
)

type magicConfig[T any] struct {
	envPrefix           string
	emptySliceIndicator string
	listSeparator       string
	originalConfig      *T
	newConfig           any
	validator           *validator.Validate
	transformFuncs      []func(reflect.StructField, reflect.Value)
	timeFormats         []string
	constructionErrors  []error
}

type MagicConfigOptions struct {
	Prefix              string
	ListSeparator       string
	EmptySliceIndicator string
	ValidationFunctions map[string]ValidationFunction
	TransformFuncs      []func(reflect.StructField, reflect.Value)
	TimeFormats         []string
}

func NewMagicConfig[T any](config *T, options *MagicConfigOptions) (*magicConfig[T], error) {
	configVal := reflect.ValueOf(config)
	if configVal.Kind() != reflect.Ptr || configVal.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("config must be a pointer to a struct")
	}

	prefix := ""
	if options != nil {
		prefix = options.Prefix
	}

	listSeparator := ","
	if options != nil && options.ListSeparator != "" {
		listSeparator = options.ListSeparator
	}

	emptySliceIndicator := _emptyslice
	if options != nil && options.EmptySliceIndicator != "" {
		emptySliceIndicator = options.EmptySliceIndicator
	}

	timeFormats := []string{}
	if options != nil && options.TimeFormats != nil {
		timeFormats = options.TimeFormats
	}

	vd := validator.New(validator.WithRequiredStructEnabled())
	if options != nil {
		for n, f := range options.ValidationFunctions {
			if err := vd.RegisterValidation(n, (func(p validator.FieldLevel) bool)(f)); err != nil {
				return nil, err
			}
		}
	}

	var transformFuncs []func(reflect.StructField, reflect.Value)
	if options != nil {
		transformFuncs = options.TransformFuncs
	}

	fullConfig := makeFullConfig(configVal, prefix, nil, nil, nil).Interface()

	return &magicConfig[T]{
		originalConfig:      config,
		newConfig:           fullConfig,
		envPrefix:           prefix,
		listSeparator:       listSeparator,
		emptySliceIndicator: emptySliceIndicator,
		validator:           vd,
		transformFuncs:      transformFuncs,
		timeFormats:         timeFormats,
		constructionErrors:  []error{},
	}, nil
}

func (c *magicConfig[T]) ParseEnv() *magicConfig[T] {
	if err := env.Parse(c.newConfig); err != nil {
		c.constructionErrors = append(c.constructionErrors, err)
	}
	return c
}

func (c *magicConfig[T]) ParseFiles(files ...string) *magicConfig[T] {
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			c.constructionErrors = append(c.constructionErrors, err)
			continue
		}

		test := map[string]any{}

		isYAML := true
		if err := yaml.Unmarshal(data, &test); err != nil {
			if err := toml.Unmarshal(data, &test); err != nil {
				c.constructionErrors = append(c.constructionErrors, fmt.Errorf("failed to recognize file %s as YAML, JSON, or TOML", f))
				continue
			}
			isYAML = false
		}

		if isYAML {
			if err := yaml.Unmarshal(data, c.newConfig); err != nil {
				c.constructionErrors = append(c.constructionErrors, err)
			}
		} else {
			if err := toml.Unmarshal(data, c.newConfig); err != nil {
				c.constructionErrors = append(c.constructionErrors, err)
			}
		}
	}

	return c
}

func (c *magicConfig[T]) ParseFlags() *magicConfig[T] {
	parser := flags.NewParser(c.newConfig, flags.Default)
	args := splitArgs(reflect.ValueOf(c.newConfig), os.Args, c.listSeparator)
	args = removeTimes(reflect.ValueOf(c.newConfig), args, c.timeFormats)
	if _, err := parser.ParseArgs(args); err != nil {
		c.constructionErrors = append(c.constructionErrors, err)
	}

	return c
}

func (c *magicConfig[T]) Validate() *magicConfig[T] {
	transformValues(reflect.ValueOf(c.newConfig), c.transformFuncs)
	if errs := validateEnums(reflect.ValueOf(c.newConfig)); len(errs) > 0 {
		c.constructionErrors = append(c.constructionErrors, errs...)
	}
	if err := c.validator.Struct(c.newConfig); err != nil {
		c.constructionErrors = append(c.constructionErrors, err)
	}

	return c
}

func (c *magicConfig[T]) ApplyDefaults() *magicConfig[T] {
	setDefaults(
		reflect.ValueOf(c.newConfig),
		c.listSeparator,
		c.emptySliceIndicator,
		c.timeFormats,
	)
	return c
}

func (c *magicConfig[T]) Config() (*T, []error) {
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
