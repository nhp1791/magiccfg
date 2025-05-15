package magiccfg

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	_emptyslice = "emptyslice"
)

var (
	timePtrType     = reflect.TypeOf(new(time.Time))
	durationPtrType = reflect.TypeOf(new(time.Duration))
)

func (c *magicConfig[T]) empty() reflect.Value {
	newC := reflect.New(reflect.TypeOf(c.newConfig).Elem())
	populateEmptyStructs(newC)
	return newC
}

func recursiveStruct(c reflect.Value) bool {
	return c.Kind() == reflect.Ptr &&
		c.Type().Elem().Kind() == reflect.Struct &&
		c.Type() != timePtrType
}

func makeFullConfig(
	c reflect.Value,
	prefix string,
	parentName *string,
	envParentName *string,
	nameMap map[string]int,
) reflect.Value {
	if !recursiveStruct(c) {
		return reflect.Zero(c.Type())
	}

	populateEmptyStructs(c)

	if len(nameMap) == 0 {
		nameMap = makeNameMap(c, prefix)
	}

	fields := make([]reflect.StructField, 0)
	var newField reflect.StructField
	cVal := c.Elem()

	for i := range cVal.NumField() {
		field := cVal.Type().Field(i)
		fld := cVal.Field(i)
		name := field.Name
		structField := recursiveStruct(fld)

		newParentName := convertNameToCommandLine(name, parentName)
		newEnvParentName := ""
		if envParentName != nil {
			newEnvParentName = fmt.Sprintf("%s_%s", *envParentName, formEnvName(name))
		} else {
			newEnvParentName = formEnvName(name)
		}

		tag := buildTag(field.Tag, name, prefix, structField, parentName, envParentName, nameMap)

		fieldType := field.Type

		if structField {
			f := makeFullConfig(fld, prefix, &newParentName, &newEnvParentName, nameMap)
			fieldType = f.Type()
		}

		newField = reflect.StructField{
			Name: name,
			Type: fieldType,
			Tag:  tag,
		}

		fields = append(fields, newField)
	}

	structType := reflect.StructOf(fields)
	newValue := reflect.New(structType)

	populateEmptyStructs(newValue)
	return newValue
}

func makeNameMap(cVal reflect.Value, prefix string) map[string]int {
	c := reflect.Indirect(cVal)
	nameMap := make(map[string]int)
	for i := range c.NumField() {
		field := c.Type().Field(i)
		fld := c.Field(i)
		name := field.Name
		structField := recursiveStruct(fld)
		if !structField {
			nameMap[formEnvName(name)]++
			continue
		}
		subMap := makeNameMap(fld, prefix)
		for k := range subMap {
			nameMap[k]++
		}
	}
	return nameMap
}

func buildTag(
	tag reflect.StructTag,
	name string,
	prefix string,
	structField bool,
	parentName *string,
	envParentName *string,
	nameMap map[string]int,
) reflect.StructTag {
	var sb strings.Builder

	sb.WriteString(string(tag))
	if !structField && tag.Get("env") == "" {
		sb.WriteString(fmt.Sprintf(` env:"%s"`, convertNameToEnv(name, envParentName, prefix, nameMap)))
	}
	if tag.Get("yaml") == "" {
		sb.WriteString(fmt.Sprintf(` yaml:"%s"`, convertNameToYAML(name)))
	}
	if tag.Get("json") == "" {
		sb.WriteString(fmt.Sprintf(` json:"%s"`, convertNameToYAML(name)))
	}
	if tag.Get("toml") == "" {
		sb.WriteString(fmt.Sprintf(` toml:"%s"`, convertNameToYAML(name)))
	}
	if !structField && tag.Get("long") == "" {
		sb.WriteString(fmt.Sprintf(` long:"%s"`, convertNameToCommandLine(name, parentName)))
	}

	return reflect.StructTag(sb.String())
}

func populateEmptyStructs(c reflect.Value) {
	if !recursiveStruct(c) {
		return
	}

	cVal := reflect.Indirect(c)

	for i := range cVal.NumField() {
		field := cVal.Type().Field(i)
		fld := cVal.Field(i)

		if !recursiveStruct(fld) {
			continue
		}
		if fld.IsNil() {
			n := reflect.New(field.Type.Elem())
			populateEmptyStructs(n)
			fld.Set(n)
		}
	}
}

func merge(oldV, newV reflect.Value) {
	o := reflect.Indirect(oldV)
	n := reflect.Indirect(newV)
	for i := range o.NumField() {
		fld := o.Field(i)
		nFld := n.Field(i)
		if recursiveStruct(fld) {
			merge(fld, nFld)
			continue
		}

		if !nFld.IsNil() && fld.CanSet() {
			fld.Set(nFld)
		}
	}
}

func isCapital(v rune) bool {
	return v >= 'A' && v <= 'Z'
}

func setDefaults(
	v reflect.Value,
	listSeparator string,
	emptySliceIndicator string,
	timeFormats []string,
) {
	val := reflect.Indirect(v)
	for i := range val.NumField() {
		fld := val.Field(i)
		field := val.Type().Field(i)

		if recursiveStruct(fld) {
			setDefaults(fld, listSeparator, emptySliceIndicator, timeFormats)
			continue
		}
		if !fld.IsNil() {
			continue
		}

		f := reflect.Indirect(fld)

		tag := field.Tag.Get("def")
		if tag != "" {
			if fld.Type().Kind() == reflect.Slice {
				fld.Set(reflect.MakeSlice(f.Type(), 0, 0))
			} else {
				if fld.Kind() == reflect.Ptr {
					switch fld.Type() {
					case durationPtrType:
						if v, err := time.ParseDuration(tag); err == nil {
							fld.Set(reflect.ValueOf(&v))
						}
						continue
					case timePtrType:
						for _, timeFormat := range timeFormats {
							if v, err := time.Parse(timeFormat, tag); err == nil {
								fld.Set(reflect.ValueOf(&v))
								break
							}
						}
						continue
					}
				}
				fld.Set(reflect.New(field.Type.Elem()))

			}

			f := reflect.Indirect(fld)

			SetValue(f, tag, listSeparator, emptySliceIndicator, timeFormats)
		}
	}
}

func SetValue(
	f reflect.Value,
	val string,
	listSeparator string,
	emptySliceIndicator string,
	timeFormats []string,
) {
	switch f.Type().Kind() {
	case reflect.String:
		c := reflect.ValueOf(val).Convert(f.Type())
		f.Set(c)
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
		i, err := strconv.ParseInt(val, 10, 0)
		if err != nil {
			return
		}
		c := reflect.ValueOf(int(i)).Convert(f.Type())
		f.Set(c)
	case reflect.Float32, reflect.Float64:
		v, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return
		}
		c := reflect.ValueOf(float64(v)).Convert(f.Type())
		f.Set(c)
	case reflect.Bool:
		b, err := strconv.ParseBool(val)
		if err != nil {
			return
		}
		c := reflect.ValueOf(b).Convert(f.Type())
		f.Set(c)
	case reflect.Slice:
		vals := []string{}
		if val != emptySliceIndicator {
			vals = strings.Split(val, listSeparator)
		}
		c := reflect.MakeSlice(f.Type(), 0, len(vals))
		elemType := reflect.TypeOf(f.Interface()).Elem()
		elemKind := elemType.Kind()
		switch elemType {
		case durationPtrType:
			for _, dur := range vals {
				if d, err := time.ParseDuration(dur); err == nil {
					c = reflect.Append(c, reflect.ValueOf(&d))
				}
			}
			f.Set(c)
			return
		case timePtrType:
			for _, tm := range vals {
				for _, timeFormat := range timeFormats {
					if t, err := time.Parse(timeFormat, tm); err == nil {
						c = reflect.Append(c, reflect.ValueOf(&t))
						break
					}
				}
			}
			f.Set(c)
			return
		}

		for _, v := range vals {
			switch elemKind {
			case reflect.String:
				c = reflect.Append(c, reflect.ValueOf(v).Convert(f.Type().Elem()))
			case reflect.Bool:
				b, err := strconv.ParseBool(v)
				if err != nil {
					continue
				}
				c = reflect.Append(c, reflect.ValueOf(b).Convert(elemType))
			case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
				i, err := strconv.ParseInt(v, 10, 0)
				if err != nil {
					continue
				}
				c = reflect.Append(c, reflect.ValueOf(i).Convert(elemType))
			case reflect.Float32, reflect.Float64:
				v, err := strconv.ParseFloat(val, 64)
				if err != nil {
					return
				}
				c = reflect.Append(c, reflect.ValueOf(v).Convert(elemType))
			}
		}
		f.Set(c)
	}
}

func transformValues(c reflect.Value, customFuncs []func(reflect.StructField, reflect.Value)) {
	if !recursiveStruct(c) {
		return
	}
	cVal := reflect.Indirect(c)
	for i := range cVal.NumField() {
		field := cVal.Type().Field(i)
		fld := cVal.Field(i)
		if recursiveStruct(fld) {
			transformValues(fld, customFuncs)
		}

		adjustCase(field, fld)
		stripslash(field, fld)
		stripprotocol(field, fld)

		for _, customFunc := range customFuncs {
			customFunc(field, fld)
		}
	}
}

func adjustCase(field reflect.StructField, fld reflect.Value) {
	tag := field.Tag.Get("case")
	if tag == "" {
		return
	}
	var lower bool

	if tag == "lower" {
		lower = true
	} else if tag == "upper" {
		lower = false
	} else {
		return
	}

	f := reflect.Indirect(fld)
	if f.Kind() == reflect.String {
		var val string
		if lower {
			val = strings.ToLower(f.String())
		} else {
			val = strings.ToUpper(f.String())
		}
		c := reflect.ValueOf(val).Convert(f.Type())
		f.Set(c)
	}
}

func stripslash(field reflect.StructField, fld reflect.Value) {
	if _, ok := field.Tag.Lookup("stripslash"); !ok {
		return
	}

	f := reflect.Indirect(fld)
	if f.Kind() == reflect.String {
		val := strings.TrimSuffix(f.String(), "/")
		c := reflect.ValueOf(val).Convert(f.Type())
		f.Set(c)
	}
}

func stripprotocol(field reflect.StructField, fld reflect.Value) {
	if _, ok := field.Tag.Lookup("stripprotocol"); !ok {
		return
	}
	f := reflect.Indirect(fld)
	if f.Kind() == reflect.String {
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
}

func validateEnums(c reflect.Value) []error {
	errs := []error{}
	if !recursiveStruct(c) {
		return errs
	}
	cVal := reflect.Indirect(c)
	for i := range cVal.NumField() {
		fld := cVal.Field(i)
		if recursiveStruct(fld) {
			if subErrs := validateEnums(fld); len(subErrs) > 0 {
				errs = append(errs, subErrs...)
			}
		}

		t := fld.Type().Elem()
		if t.Kind() != reflect.String || t.Name() == "string" {
			continue
		}

		m, ok := t.MethodByName("ValidateEnum")
		if !ok {
			continue
		}

		f := reflect.Indirect(fld)
		if f.Type().Kind() == reflect.Slice {
			for i := range f.Len() {
				element := f.Index(i)
				e := m.Func.Call([]reflect.Value{element, reflect.ValueOf(element.String())})
				if len(e) > 0 {
					v, ok := e[0].Interface().(bool)
					if ok && !v {
						errs = append(errs, fmt.Errorf("value %s for enum field '%s' is not valid", element.String(), t.Name()))
					}
				}
			}
		} else {
			e := m.Func.Call([]reflect.Value{f, reflect.ValueOf(f.String())})
			if len(e) > 0 {
				v, ok := e[0].Interface().(bool)
				if ok && !v {
					errs = append(errs, fmt.Errorf("value %s for enum field '%s' is not valid", f.String(), t.Name()))
				}
			}
		}
	}
	return errs
}

func splitArgs(c reflect.Value, args []string, listSeparator string) []string {
	if !recursiveStruct(c) {
		return args
	}

	newArgs := []string{}

	skip := false
	var rawValues []string
	var fixedArg string

	for i, arg := range args {
		if skip {
			skip = false
			continue
		}
		if !strings.HasPrefix(arg, "--") && !strings.HasPrefix(arg, "-") {
			newArgs = append(newArgs, arg)
			continue
		}

		fixedArg, rawValues, skip = getValues(c, arg, args, i, listSeparator)

		if len(rawValues) == 0 {
			newArgs = append(newArgs, fixedArg)
			continue
		}
		for _, rawValue := range rawValues {
			newArgs = append(newArgs, fixedArg)
			newArgs = append(newArgs, rawValue)
		}
	}

	return newArgs
}

func getValues(c reflect.Value, arg string, args []string, i int, listSeparator string) (string, []string, bool) {
	values := []string{}
	fixedArg := arg

	if strings.Contains(arg, "=") {
		vals := strings.Split(arg, "=")
		if len(vals) == 2 {
			fixedArg = vals[0]
			values = extractValues(c, fixedArg, vals[1], listSeparator)
		}

		return fixedArg, values, false
	}

	nextArg := ""
	if i < len(args)-2 {
		nextArg = args[i+1]
	}

	values = extractValues(c, fixedArg, nextArg, listSeparator)
	return fixedArg, values, true
}

func extractValues(c reflect.Value, arg string, valArg string, listSeparator string) []string {
	values := []string{}
	if valArg == "" || strings.HasPrefix(valArg, "-") || strings.HasPrefix(valArg, "--") {
		return values
	}
	if !strings.Contains(valArg, listSeparator) || !verifySliceType(c, arg) {
		values = append(values, valArg)
		return values
	}

	return strings.Split(valArg, listSeparator)
}

func verifySliceType(c reflect.Value, arg string) bool {
	if !recursiveStruct(c) {
		return false
	}

	bareArg := strings.TrimPrefix(strings.TrimPrefix(arg, "-"), "-")

	cVal := reflect.Indirect(c)
	for i := range cVal.NumField() {
		field := cVal.Type().Field(i)
		fld := cVal.Field(i)
		if recursiveStruct(fld) {
			if verifySliceType(fld, arg) {
				return true
			}
			continue
		}

		shortTag := field.Tag.Get("short")
		longTag := field.Tag.Get("long")

		if shortTag != bareArg && longTag != bareArg {
			continue
		}

		if fld.Type().Kind() == reflect.Slice {
			return true
		}

	}
	return false
}

func removeTimes(
	c reflect.Value,
	args []string,
	timeFormats []string,
) []string {
	newArgs := []string{}

	skip := false
	timeSetter := createTimeSetter()

	for i, arg := range args {
		if skip {
			skip = false
			continue
		}
		if !strings.HasPrefix(arg, "--") && !strings.HasPrefix(arg, "-") {
			newArgs = append(newArgs, arg)
			continue
		}

		var write bool
		write, skip = lookForTimes(c, arg, args, i, timeFormats, timeSetter)
		if write {
			newArgs = append(newArgs, arg)
		}
	}

	return newArgs
}

func lookForTimes(
	c reflect.Value,
	arg string,
	args []string,
	i int,
	timeFormats []string,
	timeSetter func(reflect.Value, string, time.Time) bool,
) (write bool, skip bool) {
	var t time.Time
	var err error

	nextArg := ""
	if i < len(args)-1 {
		nextArg = args[i+1]
	}

	if nextArg == "" || strings.HasPrefix(nextArg, "-") || strings.HasPrefix(nextArg, "--") {
		return true, false
	}

	for _, timeFormat := range timeFormats {
		t, err = time.Parse(timeFormat, nextArg)
		if err != nil {
			continue
		}
		if !t.IsZero() && timeSetter(c, arg, t) {
			return false, true
		}
	}

	return true, false
}

func createTimeSetter() func(reflect.Value, string, time.Time) bool {
	initialized := false

	var timeSetter func(c reflect.Value, arg string, t time.Time) bool

	timeSetter = func(c reflect.Value, arg string, t time.Time) bool {
		if !recursiveStruct(c) {
			return false
		}
		bareArg := strings.TrimPrefix(strings.TrimPrefix(arg, "-"), "-")
		cVal := reflect.Indirect(c)
		for i := range cVal.NumField() {
			field := cVal.Type().Field(i)
			fld := cVal.Field(i)
			if recursiveStruct(fld) {
				if timeSetter(fld, arg, t) {
					return true
				}
				continue
			}

			shortTag := field.Tag.Get("short")
			longTag := field.Tag.Get("long")

			if shortTag != bareArg && longTag != bareArg {
				continue
			}

			if fld.Type() == timePtrType {
				fld.Set(reflect.ValueOf(&t))
				return true
			}
			f := reflect.Indirect(fld)

			elemType := reflect.TypeOf(f.Interface()).Elem()
			if fld.Type().Kind() == reflect.Slice && elemType == timePtrType {
				if fld.IsNil() || !initialized {
					fld.Set(reflect.MakeSlice(f.Type(), 0, 0))
					initialized = true
				}
				c := reflect.Append(fld, reflect.ValueOf(&t))
				fld.Set(c)
				return true
			}
		}
		return false
	}

	return timeSetter
}

func locateCLIFiles(options *Options, listSeparator string) ([]string, []string) {
	args := []string{}
	files := []string{}

	shortPrefix := options.ConfigFileShort
	if !strings.HasPrefix(shortPrefix, "-") {
		shortPrefix = "-" + options.ConfigFileShort
	}
	longPrefix := options.ConfigFileLong
	if !strings.HasPrefix(longPrefix, "--") {
		longPrefix = "--" + options.ConfigFileLong
	}

	skip := false

	for i, arg := range os.Args {
		if skip {
			skip = false
			continue
		}
		if !strings.HasPrefix(arg, shortPrefix) && !strings.HasPrefix(arg, longPrefix) {
			args = append(args, arg)
			continue
		}

		if strings.Contains(arg, "=") {
			vals := strings.Split(arg, "=")
			if len(vals) != 2 {
				args = append(args, arg)
				continue
			}
			files = append(files, strings.Split(vals[1], listSeparator)...)
			continue
		}

		if i < len(os.Args)-1 {
			nextArg := os.Args[i+1]
			files = append(files, strings.Split(nextArg, listSeparator)...)
			skip = true
		}
	}
	return files, args
}
