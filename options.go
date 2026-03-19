package magiccfg

import "reflect"

const (
	// ZeroPointer is used to
	_zeroPointer       = "*0*"
	_listSeparator     = ","
	_keyValueSeparator = ":"
)

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
}
