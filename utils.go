package magiccfg

import (
	"fmt"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	_emptyslice = "emptyslice"
	_emptymap   = "emptymap"
)

func (c *magicConfig[T]) empty() reflect.Value {
	newC := reflect.New(reflect.TypeOf(c.newConfig).Elem())
	populateEmptyStructs(newC)
	return newC
}

func recursiveStruct(c reflect.Value) bool {
	cType := c.Type()
	return c.Kind() == reflect.Ptr &&
		cType.Elem().Kind() == reflect.Struct &&
		cType != timePtrType &&
		cType != parseableTimePtrType &&
		cType != urlType &&
		cType != parseableURLType
}

func makeFullConfig(
	c reflect.Value,
	prefix string,
	parentName *string,
	envParentName *string,
) (reflect.Value, map[string]map[string]int) {
	usages := map[string]map[string]int{
		"env":   map[string]int{},
		"long":  map[string]int{},
		"short": map[string]int{},
	}

	if !recursiveStruct(c) {
		return reflect.Zero(c.Type()), usages
	}

	populateEmptyStructs(c)

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

		tag := buildTag(field.Tag, name, prefix, structField, parentName, envParentName)

		fieldType := field.Type
		switch fld.Type() {
		case durationPtrType:
			fieldType = parseableTimeDurationPtrType
		case sliceDurationType:
			fieldType = sliceParseableTimeDurationType
		case timePtrType:
			fieldType = parseableTimePtrType
		case sliceTimeType:
			fieldType = sliceParseableTimeType
		case urlType:
			fieldType = parseableURLType
		case sliceURLType:
			fieldType = sliceParseableURLType
		case mapDurationType:
			fieldType = mapParseableDurationType
		}

		if structField {
			f, subUsages := makeFullConfig(fld, prefix, &newParentName, &newEnvParentName)
			fieldType = f.Type()
			for key, v := range subUsages {
				for k, n := range v {
					usages[key][k] += n
				}
			}
		} else {
			usages["env"][tag.Get("env")]++
			usages["short"][tag.Get("short")]++
			usages["long"][tag.Get("long")]++
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
	return newValue, usages
}

func buildTag(
	tag reflect.StructTag,
	name string,
	prefix string,
	structField bool,
	parentName *string,
	envParentName *string,
) reflect.StructTag {
	var sb strings.Builder

	sb.WriteString(string(tag))
	if !structField && tag.Get("env") == "" {
		sb.WriteString(fmt.Sprintf(` env:"%s"`, convertNameToEnv(name, envParentName, prefix)))
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
	if tag.Get("description") == "" {
		sb.WriteString(fmt.Sprintf(` description:"%s: yaml/json/toml/xml: %s"`, name, convertNameToYAML(name)))
	}
	if tag.Get("xml") == "" {
		attr := false
		if tag.Get("xmlattr") != "" {
			if a, err := strconv.ParseBool(tag.Get("xmlattr")); err == nil {
				attr = a
			}
		}
		sb.WriteString(fmt.Sprintf(` xml:"%s`, convertNameToYAML(name)))
		if attr {
			sb.WriteString(`,attr`)
		}
		sb.WriteString(`"`)
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
	target := reflect.Indirect(oldV)
	source := reflect.Indirect(newV)
	for i := range target.NumField() {
		targetField := target.Field(i)
		sourceField := source.Field(i)
		name := target.Type().Field(i).Name
		_ = name
		if recursiveStruct(targetField) {
			merge(targetField, sourceField)
			continue
		}

		if !isZero(sourceField) && targetField.CanSet() {
			targetFieldType := targetField.Type()
			switch sourceField.Type() {
			case parseableTimeDurationPtrType:
				if targetFieldType == durationPtrType {
					targetField.Set(sourceField.Convert(durationPtrType))
				} else {
					targetField.Set(sourceField)
				}
			case sliceParseableTimeDurationType:
				if targetFieldType == sliceDurationType {
					newDurations := sourceField.Interface().([]*ParseableTimeDuration)
					durations := make([]*time.Duration, len(newDurations))
					for i, nd := range newDurations {
						if nd == nil {
							continue
						}
						d := time.Duration(*nd)
						durations[i] = &d
					}
					targetField.Set(reflect.ValueOf(durations))
				} else {
					targetField.Set(sourceField)
				}
			case mapParseableDurationType:
				if targetFieldType == mapDurationType {
					newDurationMap := sourceField.Interface().(map[string]*ParseableTimeDuration)
					durationMap := make(map[string]*time.Duration, len(newDurationMap))
					for k, nd := range newDurationMap {
						if nd == nil {
							continue
						}
						d := time.Duration(*nd)
						durationMap[k] = &d
					}
					targetField.Set(reflect.ValueOf(durationMap))
				} else {
					targetField.Set(sourceField)
				}
			case parseableTimePtrType:
				if targetFieldType == timePtrType {
					targetField.Set(sourceField.Convert(timePtrType))
				} else {
					targetField.Set(sourceField)
				}
			case sliceParseableTimeType:
				if targetFieldType == sliceTimeType {
					newTimes := sourceField.Interface().([]*ParseableTime)
					times := make([]*time.Time, len(newTimes))
					for i, nt := range newTimes {
						if nt == nil {
							continue
						}
						t := time.Time(*nt)
						times[i] = &t
					}
					targetField.Set(reflect.ValueOf(times))
				} else {
					targetField.Set(sourceField)
				}
			case parseableURLType:
				if targetFieldType == urlType {
					targetField.Set(sourceField.Convert(urlType))
				} else {
					targetField.Set(sourceField)
				}
			case sliceParseableURLType:
				if targetFieldType == sliceURLType {
					newURLs := sourceField.Interface().([]*ParseableURL)
					urls := make([]*url.URL, len(newURLs))
					for i, nu := range newURLs {
						if nu == nil {
							continue
						}
						u := url.URL(*nu)
						urls[i] = &u
					}
					targetField.Set(reflect.ValueOf(urls))
				} else {
					targetField.Set(sourceField)
				}
			default:
				targetField.Set(sourceField)
			}
		}
	}
}

func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Func, reflect.Map, reflect.Slice, reflect.Struct:
		return v.IsNil()
	default:
		zero := reflect.Zero(v.Type())
		return v.Interface() == zero.Interface()
	}
}

func isCapital(v rune) bool {
	return v >= 'A' && v <= 'Z'
}

func setDefaults(
	v reflect.Value,
	listSeparator string,
	keyValueSeparator string,
	emptySliceIndicator string,
	emptyMapIndicator string,
	timeFormats []string,
) {
	val := reflect.Indirect(v)
	for i := range val.NumField() {
		fld := val.Field(i)
		field := val.Type().Field(i)
		name := field.Name
		_ = name
		if recursiveStruct(fld) {
			setDefaults(fld, listSeparator, keyValueSeparator, emptySliceIndicator, emptyMapIndicator, timeFormats)
			continue
		}
		if !isZero(fld) && (fld.Type().Kind() != reflect.Map || fld.Len() > 0) {
			continue
		}

		f := reflect.Indirect(fld)

		tag := field.Tag.Get("def")
		if tag != "" {
			if fld.Type().Kind() == reflect.Map {
				fld.Set(reflect.MakeMap(fld.Type()))
			} else if fld.Type().Kind() == reflect.Slice {
				fld.Set(reflect.MakeSlice(f.Type(), 0, 0))
			} else {
				if fld.Kind() == reflect.Ptr {
					switch fld.Type() {
					case durationPtrType:
						if d, err := time.ParseDuration(tag); err == nil {
							fld.Set(reflect.ValueOf(&d))
						}
						continue
					case parseableTimeDurationPtrType:
						if d, err := time.ParseDuration(tag); err == nil {
							x := ParseableTimeDuration(d)
							fld.Set(reflect.ValueOf(&x))
						}
						continue
					case timePtrType:
						for _, timeFormat := range timeFormats {
							if t, err := time.Parse(timeFormat, tag); err == nil {
								fld.Set(reflect.ValueOf(&t))
								break
							}
						}
						continue
					case parseableTimePtrType:
						for _, timeFormat := range timeFormats {
							if t, err := time.Parse(timeFormat, tag); err == nil {
								x := ParseableTime(t)
								fld.Set(reflect.ValueOf(&x))
								break
							}
						}
						continue
					case urlType:
						if u, err := url.Parse(tag); err == nil {
							fld.Set(reflect.ValueOf(u))
						}
						continue
					}
				}
				fld.Set(reflect.New(field.Type.Elem()))
			}

			f := reflect.Indirect(fld)

			setFieldValue(f, tag, listSeparator, keyValueSeparator, emptySliceIndicator, emptyMapIndicator, timeFormats)
		}
	}
}

func setFieldValue(
	f reflect.Value,
	val string,
	listSeparator string,
	keyValueSeparator string,
	emptySliceIndicator string,
	emptyMapIndicator string,
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
	case reflect.Map:
		kvs := []string{}
		vals := map[string]string{}
		if val != emptyMapIndicator {
			kvs = strings.Split(val, listSeparator)
		}
		for _, kv := range kvs {
			mapVals := strings.Split(kv, keyValueSeparator)
			if len(mapVals) != 2 {
				continue
			}
			vals[mapVals[0]] = mapVals[1]
		}
		c := reflect.MakeMap(f.Type())
		elemType := reflect.TypeOf(f.Interface()).Elem()
		elemKind := elemType.Kind()
		switch elemType {
		case durationPtrType:
			for k, v := range vals {
				if d, err := time.ParseDuration(v); err == nil {
					c.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(&d))
				}
			}
			f.Set(c)
			return
		case timePtrType:
			for k, v := range vals {
				for _, timeFormat := range timeFormats {
					if t, err := time.Parse(timeFormat, v); err == nil {
						c.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(&t))
						break
					}
				}
			}
			f.Set(c)
			return
		case urlType:
			for k, v := range vals {
				if u, err := url.Parse(v); err == nil {
					c.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(u))
				}
			}
			f.Set(c)
			return
		}
		for k, v := range vals {
			switch elemKind {
			case reflect.String:
				c.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v).Convert(f.Type().Elem()))
			case reflect.Bool:
				b, err := strconv.ParseBool(v)
				if err != nil {
					continue
				}
				c.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(b).Convert(elemType))
			case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
				i, err := strconv.ParseInt(v, 10, 0)
				if err != nil {
					continue
				}
				c.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(i).Convert(elemType))
			case reflect.Float32, reflect.Float64:
				v, err := strconv.ParseFloat(val, 64)
				if err != nil {
					return
				}
				c.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v).Convert(elemType))
			}
		}
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
		case urlType:
			for _, v := range vals {
				if u, err := url.Parse(v); err == nil {
					c = reflect.Append(c, reflect.ValueOf(u))
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
		addslash(field, fld)
		stripscheme(field, fld)

		for _, customFunc := range customFuncs {
			customFunc(field, fld)
		}
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

		mt := m.Type
		if mt.NumIn() != 1 || mt.NumOut() != 1 {
			continue
		}
		if mt.Out(0).Kind() != reflect.Bool {
			continue
		}

		f := reflect.Indirect(fld)
		if f.Type().Kind() == reflect.Slice {
			for i := range f.Len() {
				element := f.Index(i)
				e := m.Func.Call([]reflect.Value{element})
				if len(e) > 0 {
					v, ok := e[0].Interface().(bool)
					if ok && !v {
						errs = append(errs, fmt.Errorf("value %s for enum field '%s' is not valid", element.String(), t.Name()))
					}
				}
			}
		} else if f.Type().Kind() == reflect.Map {
			for _, k := range f.MapKeys() {
				element := f.MapIndex(k)
				e := m.Func.Call([]reflect.Value{element})
				if len(e) > 0 {
					v, ok := e[0].Interface().(bool)
					if ok && !v {
						errs = append(errs, fmt.Errorf("value %s for enum field '%s' is not valid", element.String(), t.Name()))
					}
				}
			}
		} else {
			e := m.Func.Call([]reflect.Value{f})
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

	if nextArg == "" || strings.HasPrefix(nextArg, "-") || strings.HasPrefix(nextArg, "--") {
		return fixedArg, []string{}, false
	}

	values = extractValues(c, fixedArg, nextArg, listSeparator)
	return fixedArg, values, true
}

func extractValues(c reflect.Value, arg string, valArg string, listSeparator string) []string {
	values := []string{}

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

func locateCLIFiles(options *Options, listSeparator string) ([]string, []string) {
	args := []string{}
	files := []string{}

	shortPrefix := options.ConfigFileShort
	if shortPrefix != "" && !strings.HasPrefix(shortPrefix, "-") {
		shortPrefix = "-" + options.ConfigFileShort
	}
	longPrefix := options.ConfigFileLong
	if longPrefix != "" && !strings.HasPrefix(longPrefix, "--") {
		longPrefix = "--" + options.ConfigFileLong
	}

	skip := false

	for i, arg := range os.Args {
		if skip {
			skip = false
			continue
		}
		if (shortPrefix == "" || !strings.HasPrefix(arg, shortPrefix)) &&
			(longPrefix == "" || !strings.HasPrefix(arg, longPrefix)) {
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
