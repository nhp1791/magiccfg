package magiccfg

import (
	"fmt"
	"reflect"
)

const (
	// ZeroPointer is used to
	_zeroPointer       = "*0*"
	_listSeparator     = ","
	_keyValueSeparator = ":"
)

type DefaultParseable interface {
	ParseDefault(string) any
}

type EnvironmentParseable interface {
	ParseEnvironment(string) (any, error)
}

type customType struct {
	basicType         reflect.Type
	kind              reflect.Kind
	defaultParser     func(string) any
	environmentParser func(string) (any, error)
}

type Options struct {
	ConfigFileEnvName    string
	ConfigFileShort      string
	ConfigFileLong       string
	EnvPrefix            string
	ZeroPointer          string
	ListSeparator        string
	KeyValueSeparator    string
	IgnoreUnknownOptions bool
	ValidationFunctions  map[string]ValidationFunction
	TransformFuncs       []func(reflect.StructField, reflect.Value)
	TimeFormats          []string
	CustomTypes          map[string]*customType
}

func (o *Options) AddCustomType(name string, t any) error {
	basic := reflect.TypeOf(t)
	pt := reflect.New(basic).Interface()
	pointer := reflect.ValueOf(pt)
	k := basic.Kind()

	if k != reflect.Struct {
		return fmt.Errorf("custom type must be struct or pointer to struct")
	}

	if o.CustomTypes == nil {
		o.CustomTypes = make(map[string]*customType)
	}

	ct := &customType{
		basicType: basic,
		kind:      k,
	}

	if _, ok := pt.(DefaultParseable); ok {
		ct.defaultParser = pointer.MethodByName("ParseDefault").Interface().(func(string) any)
	}

	if _, ok := pt.(EnvironmentParseable); ok {
		ct.environmentParser = pointer.MethodByName("ParseEnvironment").Interface().(func(string) (any, error))
	}

	o.CustomTypes[name] = ct

	return nil
}

func (o *Options) AddDefaultParser(name string, parser func(string) any) bool {
	ct, ok := o.CustomTypes[name]
	if !ok {
		return false
	}
	ct.defaultParser = parser
	return true
}
