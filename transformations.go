package magiccfg

import (
	"reflect"
	"strings"
)

// adjustCase looks for struct tags on fields that are of type string or an
// aliased type, and adjusts the case of that field based on the tag.
//
// Tag values "lowercase", "lower", and "l" all work to produce a lower-case string value
// while "uppercase", "upper", and "u" all work to produce upper-case values.  These
// values are case-insensitive.
func adjustCase(field reflect.StructField, fld reflect.Value) {
	f := fld
	if fld.Kind() == reflect.Ptr {
		f = reflect.Indirect(fld)
	}
	// This only works on string values
	if f.Kind() != reflect.String {
		return
	}

	// Ensure that the tag value is lowercase for comparision with
	// the specific cases
	tag := strings.ToLower(field.Tag.Get("case"))
	if tag == "" {
		return
	}

	var lower bool

	switch tag {
	case "lower", "lowercase", "l":
		lower = true
	case "upper", "uppercase", "u":
		lower = false
	default:
		return
	}

	var val string
	if lower {
		val = strings.ToLower(f.String())
	} else {
		val = strings.ToUpper(f.String())
	}

	// Convert the altered string back to the original
	// field type before setting the field
	c := reflect.ValueOf(val).Convert(f.Type())
	f.Set(c)
}

// stripslash looks for the tag `stripslash:""` (note that the value is ignored).
// If present, a field that is of type string or an aliased type will have any following
// slashes removed.  This is useful for filepaths and urls because code that
// uses the configuration variable can be assured of the form and need not
// do its own check.
func stripslash(field reflect.StructField, fld reflect.Value) {
	f := fld
	if fld.Kind() == reflect.Ptr {
		f = reflect.Indirect(fld)
	}

	if _, ok := field.Tag.Lookup("stripslash"); !ok || f.Kind() != reflect.String {
		return
	}

	val := strings.TrimSuffix(f.String(), "/")
	c := reflect.ValueOf(val).Convert(f.Type())
	f.Set(c)
}

// addSlash looks for the tag `addslash:""` (note that the value is ignored).
// If present, a field that is of type string or an aliased type will be assured to end in
// a slash.  This is useful for filepaths and urls because code that
// uses the configuration variable can be assured of the form and need not
// do its own check.
func addslash(field reflect.StructField, fld reflect.Value) {
	f := fld
	if fld.Kind() == reflect.Ptr {
		f = reflect.Indirect(fld)
	}

	if _, ok := field.Tag.Lookup("addslash"); !ok || f.Kind() != reflect.String {
		return
	}

	val := f.String()
	if !strings.HasSuffix(val, "/") {
		val = val + "/"
	}
	c := reflect.ValueOf(val).Convert(f.Type())
	f.Set(c)
}

// stripscheme looks for the tag `stripscheme:""` (note that the value is ignored).
// If present, any common network protocol prefixes will be removed.  This is useful for
// server names/addresses so that code using the configuration variable can be assured
// of its form and need not make its own checks.
func stripscheme(field reflect.StructField, fld reflect.Value) {
	f := fld
	if fld.Kind() == reflect.Ptr {
		f = reflect.Indirect(fld)
	}

	if _, ok := field.Tag.Lookup("stripscheme"); !ok || f.Kind() != reflect.String {
		return
	}

	val := f.String()
	v := strings.TrimPrefix(val, "https://")
	v = strings.TrimPrefix(v, "HTTPS://")
	v = strings.TrimPrefix(v, "http://")
	v = strings.TrimPrefix(v, "HTTP://")
	v = strings.TrimPrefix(v, "tcp://")
	v = strings.TrimPrefix(v, "TCP://")
	v = strings.TrimPrefix(v, "udp://")
	v = strings.TrimPrefix(v, "UDP://")
	c := reflect.ValueOf(val).Convert(f.Type())
	f.Set(c)
}
