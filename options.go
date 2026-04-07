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

type DefaultListParseable interface {
	ParseListDefault(string) any
}

type DefaultListPointerParseable interface {
	ParseListPointerDefault(string) any
}

type DefaultMapParseable interface {
	ParseMapDefault(val string) any
}

type DefaultMapPointerParseable interface {
	ParseMapPointerDefault(val string) any
}

type EnvironmentParseable interface {
	ParseEnvironment(string) (any, error)
}

type EnvironmentListParseable interface {
	ParseListEnvironment(string) (any, error)
}

type EnvironmentListPointerParseable interface {
	ParseListPointerEnvironment(string) (any, error)
}

type EnvironmentMapParseable interface {
	ParseMapEnvironment(string) (any, error)
}

type EnvironmentMapPointerParseable interface {
	ParseMapPointerEnvironment(string) (any, error)
}

type ParseableListType interface {
	GetParseableListType() reflect.Type
}

type ParseableListPointerType interface {
	GetParseableListPointerType() reflect.Type
}

type ParseableMapType interface {
	GetParseableMapType() reflect.Type
}

type ParseableMapPointerType interface {
	GetParseableMapPointerType() reflect.Type
}

type MapMergeable interface {
	MergeMap(reflect.Value, reflect.Value)
}

type customType struct {
	basicType                    reflect.Type
	kind                         reflect.Kind
	listType                     reflect.Type
	listPointerType              reflect.Type
	mapType                      reflect.Type
	mapPointerType               reflect.Type
	defaultParser                func(string) any
	defaultListParser            func(string) any
	defaultListPointerParser     func(string) any
	defaultMapParser             func(string) any
	defaultMapPointerParser      func(string) any
	environmentParser            func(string) (any, error)
	environmentListParser        func(string) (any, error)
	environmentListPointerParser func(string) (any, error)
	environmentMapParser         func(string) (any, error)
	environmentMapPointerParser  func(string) (any, error)
	parseableListType            reflect.Type
	parseableListPointerType     reflect.Type
	parseableMapType             reflect.Type
	parseableMapPointerType      reflect.Type
	mergeMap                     func(reflect.Value, reflect.Value)
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
	DebugMenu            bool
	CustomTypes          map[string]*customType
}

func (o *Options) AddCustomType(name string, t any) error {
	basic := reflect.TypeOf(t)
	listType := reflect.MakeSlice(reflect.SliceOf(basic), 0, 0).Type()
	listPointerType := reflect.MakeSlice(reflect.SliceOf(reflect.PointerTo(basic)), 0, 0).Type()
	mapType := reflect.MakeMap(reflect.MapOf(reflect.TypeOf(""), basic)).Type()
	mapPointerType := reflect.MakeMap(reflect.MapOf(reflect.TypeOf(""), reflect.PointerTo(basic))).Type()
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
		basicType:       basic,
		listType:        listType,
		listPointerType: listPointerType,
		mapType:         mapType,
		mapPointerType:  mapPointerType,
		kind:            k,
	}

	if _, ok := pt.(DefaultParseable); ok {
		ct.defaultParser = pointer.MethodByName("ParseDefault").Interface().(func(string) any)
	}

	if _, ok := pt.(DefaultListParseable); ok {
		ct.defaultListParser = pointer.MethodByName("ParseListDefault").Interface().(func(string) any)
	}

	if _, ok := pt.(DefaultListPointerParseable); ok {
		ct.defaultListPointerParser = pointer.MethodByName("ParseListPointerDefault").Interface().(func(string) any)
	}

	if _, ok := pt.(DefaultMapParseable); ok {
		ct.defaultMapParser = pointer.MethodByName("ParseMapDefault").Interface().(func(string) any)
	}

	if _, ok := pt.(DefaultMapPointerParseable); ok {
		ct.defaultMapPointerParser = pointer.MethodByName("ParseMapPointerDefault").Interface().(func(string) any)
	}

	if _, ok := pt.(EnvironmentParseable); ok {
		ct.environmentParser = pointer.MethodByName("ParseEnvironment").Interface().(func(string) (any, error))
	}

	if _, ok := pt.(EnvironmentListParseable); ok {
		ct.environmentListParser = pointer.MethodByName("ParseListEnvironment").Interface().(func(string) (any, error))
	}

	if _, ok := pt.(EnvironmentListPointerParseable); ok {
		ct.environmentListPointerParser = pointer.MethodByName("ParseListPointerEnvironment").Interface().(func(string) (any, error))
	}

	if _, ok := pt.(EnvironmentMapParseable); ok {
		ct.environmentMapParser = pointer.MethodByName("ParseMapEnvironment").Interface().(func(string) (any, error))
	}

	if _, ok := pt.(EnvironmentMapPointerParseable); ok {
		ct.environmentMapPointerParser = pointer.MethodByName("ParseMapPointerEnvironment").Interface().(func(string) (any, error))
	}

	if _, ok := pt.(ParseableListType); ok {
		f := pointer.MethodByName("GetParseableListType").Interface().(func() reflect.Type)
		ct.parseableListType = f()
	}

	if _, ok := pt.(ParseableListPointerType); ok {
		f := pointer.MethodByName("GetParseableListPointerType").Interface().(func() reflect.Type)
		ct.parseableListPointerType = f()
	}

	if _, ok := pt.(ParseableMapType); ok {
		f := pointer.MethodByName("GetParseableMapType").Interface().(func() reflect.Type)
		ct.parseableMapType = f()
	}

	if _, ok := pt.(ParseableMapPointerType); ok {
		f := pointer.MethodByName("GetParseableMapPointerType").Interface().(func() reflect.Type)
		ct.parseableMapPointerType = f()
	}

	if _, ok := pt.(MapMergeable); ok {
		ct.mergeMap = pointer.MethodByName("MergeMap").Interface().(func(reflect.Value, reflect.Value))
	}

	o.CustomTypes[name] = ct

	return nil
}
