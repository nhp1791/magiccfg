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
	_envTag = "envTag"
)

func (c *MagicConfig[T]) empty() reflect.Value {
	newC := reflect.New(reflect.TypeOf(c.newConfig).Elem())
	populateEmptyStructs(newC)
	return newC
}

func recursiveStruct(c reflect.Value) bool {
	cType := c.Type()
	switch cType {
	case timePtrType, parseableTimePtrType, urlPtrType, parseableURLPtrType:
		return false
	}

	return c.Kind() == reflect.Ptr && cType.Elem().Kind() == reflect.Struct
}

func makeFullConfig(
	c reflect.Value,
	prefix string,
	commandLineNameConverter func(string) string,
	parentName *string,
	envParentName *string,
) (reflect.Value, map[string]map[string]int) {
	usages := map[string]map[string]int{
		_envTag: {},
		"long":  {},
		"short": {},
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

		var newParentName string
		if commandLineNameConverter != nil {
			newParentName = commandLineNameConverter(name)
		} else {
			newParentName = convertNameToCommandLine(name, parentName)
		}
		newEnvParentName := ""
		if envParentName != nil {
			newEnvParentName = fmt.Sprintf("%s_%s", *envParentName, formEnvName(name))
		} else {
			newEnvParentName = formEnvName(name)
		}

		wrappedBool := false
		if _, ok := field.Tag.Lookup("wrap"); ok {
			wrappedBool = true
		}

		tag := buildTag(field.Tag, name, prefix, structField, parentName, envParentName)

		fieldType := field.Type
		switch fld.Type() {
		case boolType:
			if wrappedBool {
				fieldType = wrappedBoolType
			}
		case boolSliceType:
			if wrappedBool {
				fieldType = wrappedBoolSliceType
			}
		case boolPtrSliceType:
			if wrappedBool {
				fieldType = wrappedBoolPtrSliceType
			}
		case durationType:
			fieldType = parseableDurationType
		case durationSliceType:
			fieldType = parseableDurationSliceType
		case durationMapType:
			fieldType = parseableDurationMapType
		case durationPtrType:
			fieldType = parseableDurationPtrType
		case durationSlicePtrType:
			fieldType = parseableDurationSlicePtrType
		case durationMapPtrType:
			fieldType = parseableDurationMapPtrType
		case timeType:
			fieldType = parseableTimeType
		case timeSliceType:
			fieldType = parseableTimeSliceType
		case timeMapType:
			fieldType = parseableTimeMapType
		case timePtrType:
			fieldType = parseableTimePtrType
		case timeSlicePtrType:
			fieldType = parseableTimeSlicePtrType
		case timeMapPtrType:
			fieldType = parseableTimeMapPtrType
		case urlType:
			fieldType = parseableURLType
		case urlSliceType:
			fieldType = parseableURLSliceType
		case urlMapType:
			fieldType = parseableURLMapType
		case urlPtrType:
			fieldType = parseableURLPtrType
		case urlSlicePtrType:
			fieldType = parseableURLSlicePtrType
		case urlMapPtrType:
			fieldType = parseableURLMapPtrType
		case complex128Type:
			fieldType = parseableComplex128Type
		case complex128PtrType:
			fieldType = parseableComplex128PtrType
		case complex64Type:
			fieldType = parseableComplex64Type
		case complex64PtrType:
			fieldType = parseableComplex64PtrType
		case complex64SliceType:
			fieldType = parseableComplex64SliceType
		case complex64SlicePtrType:
			fieldType = parseableComplex64SlicePtrType
		case complex64MapType:
			fieldType = parseableComplex64MapType
		case complex64MapPtrType:
			fieldType = parseableComplex64MapPtrType
		case complex128SliceType:
			fieldType = parseableComplex128SliceType
		case complex128SlicePtrType:
			fieldType = parseableComplex128SlicePtrType
		case complex128MapType:
			fieldType = parseableComplex128MapType
		case complex128MapPtrType:
			fieldType = parseableComplex128MapPtrType
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
			usages[_envTag][tag.Get(_envTag)]++
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

func makeValidateConfig(c reflect.Value) reflect.Value {
	if !recursiveStruct(c) {
		return reflect.Zero(c.Type())
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

		tag := field.Tag.Get("validate")

		fieldType := field.Type
		switch fld.Type() {
		case parseableDurationType:
			fieldType = durationType
		case parseableDurationPtrType:
			fieldType = durationPtrType
		case parseableTimeType:
			fieldType = timeType
		case parseableTimePtrType:
			fieldType = timePtrType
		case parseableURLType:
			fieldType = urlType
		case parseableURLPtrType:
			fieldType = urlPtrType
		}

		if structField {
			f := makeValidateConfig(fld)
			fieldType = f.Type()
		}

		newField = reflect.StructField{
			Name: name,
			Type: fieldType,
			Tag:  reflect.StructTag(fmt.Sprintf(`validate:"%s"`, tag)),
		}

		fields = append(fields, newField)
	}

	structType := reflect.StructOf(fields)
	newValue := reflect.New(structType)

	populateEmptyStructs(newValue)
	return newValue
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
	envName := convertNameToEnv(name, envParentName, prefix)
	dashName := convertNameToYAML(name)
	if !structField {
		if tag.Get(_envTag) == "" {
			sb.WriteString(fmt.Sprintf(` %s:"%s"`, _envTag, convertNameToEnv(name, envParentName, prefix)))
		} else {
			envName = tag.Get(_envTag)
		}
	}
	if tag.Get("yaml") == "" {
		sb.WriteString(fmt.Sprintf(` yaml:"%s"`, dashName))
	}
	if tag.Get("json") == "" {
		sb.WriteString(fmt.Sprintf(` json:"%s"`, dashName))
	}
	if tag.Get("toml") == "" {
		sb.WriteString(fmt.Sprintf(` toml:"%s"`, dashName))
	}
	if !structField && tag.Get("long") == "" {
		sb.WriteString(fmt.Sprintf(` long:"%s"`, convertNameToCommandLine(name, parentName)))
	}
	if tag.Get("description") == "" {
		if parentName != nil {
			sb.WriteString(fmt.Sprintf(` description:"%s: yaml/json/toml/xml: %s.%s, env: %s"`, name, *parentName, dashName, envName))
		} else {
			sb.WriteString(fmt.Sprintf(` description:"%s: yaml/json/toml/xml: %s, env: %s"`, name, dashName, envName))
		}
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

		if recursiveStruct(targetField) {
			merge(targetField, sourceField)
			continue
		}
		if !isZero(sourceField, "") && targetField.CanSet() {
			targetFieldType := targetField.Type()
			sourceFieldType := sourceField.Type()
			if sourceFieldType == targetFieldType {
				targetField.Set(sourceField)
				continue
			}
			switch sourceFieldType {
			case durationType, durationPtrType, parseableDurationType, parseableDurationPtrType,
				timeType, timePtrType, parseableTimeType, parseableTimePtrType, parseableURLType, urlType,
				urlPtrType, parseableURLPtrType, complex128Type, parseableComplex128Type,
				complex64Type, parseableComplex64Type, complex128PtrType, parseableComplex128PtrType,
				complex64PtrType, parseableComplex64PtrType, wrappedBoolType, boolType:
				targetField.Set(sourceField.Convert(targetFieldType))
			case complex64SliceType:
				sourceSlice := sourceField.Interface().([]complex64)
				targetSlice := make([]ParseableComplex64, len(sourceSlice))
				for i := range sourceSlice {
					targetSlice[i] = ParseableComplex64(sourceSlice[i])
				}
				targetField.Set(reflect.ValueOf(targetSlice))
			case complex64SlicePtrType:
				sourceSlice := sourceField.Interface().([]*complex64)
				targetSlice := make([]*ParseableComplex64, len(sourceSlice))
				for i := range sourceSlice {
					j := ParseableComplex64(*sourceSlice[i])
					targetSlice[i] = &j
				}
				targetField.Set(reflect.ValueOf(targetSlice))
			case complex64MapType:
				sourceMap := sourceField.Interface().(map[string]complex64)
				targetMap := make(map[string]ParseableComplex64, len(sourceMap))
				for k, v := range sourceMap {
					targetMap[k] = ParseableComplex64(v)
				}
				targetField.Set(reflect.ValueOf(targetMap))
			case complex64MapPtrType:
				sourceMap := sourceField.Interface().(map[string]*complex64)
				targetMap := make(map[string]*ParseableComplex64, len(sourceMap))
				for k, v := range sourceMap {
					j := ParseableComplex64(*v)
					targetMap[k] = &j
				}
				targetField.Set(reflect.ValueOf(targetMap))
			case complex128SliceType:
				sourceSlice := sourceField.Interface().([]complex128)
				targetSlice := make([]ParseableComplex128, len(sourceSlice))
				for i := range sourceSlice {
					targetSlice[i] = ParseableComplex128(sourceSlice[i])
				}
				targetField.Set(reflect.ValueOf(targetSlice))
			case complex128SlicePtrType:
				sourceSlice := sourceField.Interface().([]*complex128)
				targetSlice := make([]*ParseableComplex128, len(sourceSlice))
				for i := range sourceSlice {
					j := ParseableComplex128(*sourceSlice[i])
					targetSlice[i] = &j
				}
				targetField.Set(reflect.ValueOf(targetSlice))
			case complex128MapType:
				sourceMap := sourceField.Interface().(map[string]complex128)
				targetMap := make(map[string]ParseableComplex128, len(sourceMap))
				for k, v := range sourceMap {
					targetMap[k] = ParseableComplex128(v)
				}
				targetField.Set(reflect.ValueOf(targetMap))
			case complex128MapPtrType:
				sourceMap := sourceField.Interface().(map[string]*complex128)
				targetMap := make(map[string]*ParseableComplex128, len(sourceMap))
				for k, v := range sourceMap {
					j := ParseableComplex128(*v)
					targetMap[k] = &j
				}
				targetField.Set(reflect.ValueOf(targetMap))
			case parseableComplex64SliceType:
				sourceSlice := sourceField.Interface().([]ParseableComplex64)
				targetSlice := make([]complex64, len(sourceSlice))
				for i := range sourceSlice {
					targetSlice[i] = complex64(sourceSlice[i])
				}
				targetField.Set(reflect.ValueOf(targetSlice))
			case parseableComplex64SlicePtrType:
				sourceSlice := sourceField.Interface().([]*ParseableComplex64)
				targetSlice := make([]*complex64, len(sourceSlice))
				for i := range sourceSlice {
					j := complex64(*sourceSlice[i])
					targetSlice[i] = &j
				}
				targetField.Set(reflect.ValueOf(targetSlice))
			case parseableComplex64MapType:
				sourceMap := sourceField.Interface().(map[string]ParseableComplex64)
				targetMap := make(map[string]complex64, len(sourceMap))
				for k, v := range sourceMap {
					targetMap[k] = complex64(v)
				}
				targetField.Set(reflect.ValueOf(targetMap))
			case parseableComplex64MapPtrType:
				sourceMap := sourceField.Interface().(map[string]*ParseableComplex64)
				targetMap := make(map[string]*complex64, len(sourceMap))
				for k, v := range sourceMap {
					j := complex64(*v)
					targetMap[k] = &j
				}
				targetField.Set(reflect.ValueOf(targetMap))
			case parseableComplex128SliceType:
				sourceSlice := sourceField.Interface().([]ParseableComplex128)
				targetSlice := make([]complex128, len(sourceSlice))
				for i := range sourceSlice {
					targetSlice[i] = complex128(sourceSlice[i])
				}
				targetField.Set(reflect.ValueOf(targetSlice))
			case parseableComplex128SlicePtrType:
				sourceSlice := sourceField.Interface().([]*ParseableComplex128)
				targetSlice := make([]*complex128, len(sourceSlice))
				for i := range sourceSlice {
					j := complex128(*sourceSlice[i])
					targetSlice[i] = &j
				}
				targetField.Set(reflect.ValueOf(targetSlice))
			case parseableComplex128MapType:
				sourceMap := sourceField.Interface().(map[string]ParseableComplex128)
				targetMap := make(map[string]complex128, len(sourceMap))
				for k, v := range sourceMap {
					targetMap[k] = complex128(v)
				}
				targetField.Set(reflect.ValueOf(targetMap))
			case parseableComplex128MapPtrType:
				sourceMap := sourceField.Interface().(map[string]*ParseableComplex128)
				targetMap := make(map[string]*complex128, len(sourceMap))
				for k, v := range sourceMap {
					j := complex128(*v)
					targetMap[k] = &j
				}
				targetField.Set(reflect.ValueOf(targetMap))
			case durationSliceType:
				sourceSlice := sourceField.Interface().([]time.Duration)
				targetSlice := make([]ParseableDuration, len(sourceSlice))
				for i := range sourceSlice {
					targetSlice[i] = ParseableDuration(sourceSlice[i])
				}
				targetField.Set(reflect.ValueOf(targetSlice))
			case durationMapType:
				sourceMap := sourceField.Interface().(map[string]time.Duration)
				targetMap := make(map[string]ParseableDuration, len(sourceMap))
				for k, v := range sourceMap {
					targetMap[k] = ParseableDuration(v)
				}
				targetField.Set(reflect.ValueOf(targetMap))
			case parseableDurationSliceType:
				sourceSlice := sourceField.Interface().([]ParseableDuration)
				targetSlice := make([]time.Duration, len(sourceSlice))
				for i := range sourceSlice {
					targetSlice[i] = time.Duration(sourceSlice[i])
				}
				targetField.Set(reflect.ValueOf(targetSlice))
			case parseableDurationMapType:
				sourceMap := sourceField.Interface().(map[string]ParseableDuration)
				targetMap := make(map[string]time.Duration, len(sourceMap))
				for k, v := range sourceMap {
					targetMap[k] = time.Duration(v)
				}
				targetField.Set(reflect.ValueOf(targetMap))
			case durationSlicePtrType:
				sourceSlice := sourceField.Interface().([]*time.Duration)
				targetSlice := make([]*ParseableDuration, len(sourceSlice))
				for i := range sourceSlice {
					val := sourceSlice[i]
					if val == nil {
						targetSlice[i] = nil
					} else {
						v := ParseableDuration(*val)
						targetSlice[i] = &v
					}
					targetField.Set(reflect.ValueOf(targetSlice))
				}
			case durationMapPtrType:
				sourceMap := sourceField.Interface().(map[string]*time.Duration)
				targetMap := make(map[string]*ParseableDuration, len(sourceMap))
				for k, v := range sourceMap {
					if v == nil {
						targetMap[k] = nil
					} else {
						val := ParseableDuration(*v)
						targetMap[k] = &val
					}
				}
				targetField.Set(reflect.ValueOf(targetMap))
			case parseableDurationSlicePtrType:
				sourceSlice := sourceField.Interface().([]*ParseableDuration)
				targetSlice := make([]*time.Duration, len(sourceSlice))
				for i := range sourceSlice {
					val := sourceSlice[i]
					if val == nil {
						targetSlice[i] = nil
					} else {
						v := time.Duration(*val)
						targetSlice[i] = &v
					}
					targetField.Set(reflect.ValueOf(targetSlice))
				}
			case parseableDurationMapPtrType:
				sourceMap := sourceField.Interface().(map[string]*ParseableDuration)
				targetMap := make(map[string]*time.Duration, len(sourceMap))
				for k, v := range sourceMap {
					if v == nil {
						targetMap[k] = nil
					} else {
						val := time.Duration(*v)
						targetMap[k] = &val
					}
				}
				targetField.Set(reflect.ValueOf(targetMap))
			case timeSliceType:
				sourceSlice := sourceField.Interface().([]time.Time)
				targetSlice := make([]ParseableTime, len(sourceSlice))
				for i := range sourceSlice {
					targetSlice[i] = ParseableTime(sourceSlice[i])
				}
				targetField.Set(reflect.ValueOf(targetSlice))
			case timeMapType:
				sourceMap := sourceField.Interface().(map[string]time.Time)
				targetMap := make(map[string]ParseableTime, len(sourceMap))
				for k, v := range sourceMap {
					targetMap[k] = ParseableTime(v)
				}
				targetField.Set(reflect.ValueOf(targetMap))
			case parseableTimeSliceType:
				sourceSlice := sourceField.Interface().([]ParseableTime)
				targetSlice := make([]time.Time, len(sourceSlice))
				for i := range sourceSlice {
					targetSlice[i] = time.Time(sourceSlice[i])
				}
				targetField.Set(reflect.ValueOf(targetSlice))
			case parseableTimeMapType:
				sourceMap := sourceField.Interface().(map[string]ParseableTime)
				targetMap := make(map[string]time.Time, len(sourceMap))
				for k, v := range sourceMap {
					targetMap[k] = time.Time(v)
				}
				targetField.Set(reflect.ValueOf(targetMap))
			case timeSlicePtrType:
				sourceSlice := sourceField.Interface().([]*time.Time)
				targetSlice := make([]*ParseableTime, len(sourceSlice))
				for i := range sourceSlice {
					val := sourceSlice[i]
					if val == nil {
						targetSlice[i] = nil
					} else {
						v := ParseableTime(*val)
						targetSlice[i] = &v
					}
					targetField.Set(reflect.ValueOf(targetSlice))
				}
			case timeMapPtrType:
				sourceMap := sourceField.Interface().(map[string]*time.Time)
				targetMap := make(map[string]*ParseableTime, len(sourceMap))
				for k, v := range sourceMap {
					if v == nil {
						targetMap[k] = nil
					} else {
						val := ParseableTime(*v)
						targetMap[k] = &val
					}
				}
				targetField.Set(reflect.ValueOf(targetMap))
			case parseableTimeSlicePtrType:
				sourceSlice := sourceField.Interface().([]*ParseableTime)
				targetSlice := make([]*time.Time, len(sourceSlice))
				for i := range sourceSlice {
					val := sourceSlice[i]
					if val == nil {
						targetSlice[i] = nil
					} else {
						v := time.Time(*val)
						targetSlice[i] = &v
					}
				}
				targetField.Set(reflect.ValueOf(targetSlice))
			case parseableTimeMapPtrType:
				sourceMap := sourceField.Interface().(map[string]*ParseableTime)
				targetMap := make(map[string]*time.Time, len(sourceMap))
				for k, v := range sourceMap {
					if v == nil {
						targetMap[k] = nil
					} else {
						val := time.Time(*v)
						targetMap[k] = &val
					}
				}
				targetField.Set(reflect.ValueOf(targetMap))
			case wrappedBoolSliceType:
				sourceSlice := sourceField.Interface().(BoolSlice)
				targetSlice := make([]bool, len(sourceSlice))
				for k := range sourceSlice {
					targetSlice[k] = bool(sourceSlice[k])
				}
				targetField.Set(reflect.ValueOf(targetSlice))
			case wrappedBoolPtrSliceType:
				sourceSlice := sourceField.Interface().(BoolPtrSlice)
				targetSlice := make([]*bool, len(sourceSlice))
				for k := range sourceSlice {
					b := bool(*sourceSlice[k])
					targetSlice[k] = &b
				}
				targetField.Set(reflect.ValueOf(targetSlice))
			case parseableURLSliceType:
				sourceSlice := sourceField.Interface().([]ParseableURL)
				targetSlice := make([]url.URL, len(sourceSlice))
				for k := range sourceSlice {
					u := url.URL(sourceSlice[k])
					targetSlice[k] = u
				}
				targetField.Set(reflect.ValueOf(targetSlice))
			case urlSliceType:
				sourceSlice := sourceField.Interface().([]url.URL)
				targetSlice := make([]ParseableURL, len(sourceSlice))
				for k := range sourceSlice {
					p := ParseableURL(sourceSlice[k])
					targetSlice[k] = p
				}
				targetField.Set(reflect.ValueOf(targetSlice))
			case parseableURLSlicePtrType:
				sourceSlice := sourceField.Interface().([]*ParseableURL)
				targetSlice := make([]*url.URL, len(sourceSlice))
				for k := range sourceSlice {
					u := url.URL(*sourceSlice[k])
					targetSlice[k] = &u
				}
				targetField.Set(reflect.ValueOf(targetSlice))
			case urlSlicePtrType:
				sourceSlice := sourceField.Interface().([]*url.URL)
				targetSlice := make([]*ParseableURL, len(sourceSlice))
				for k := range sourceSlice {
					p := ParseableURL(*sourceSlice[k])
					targetSlice[k] = &p
				}
				targetField.Set(reflect.ValueOf(targetSlice))
			case parseableURLMapType:
				sourceMap := sourceField.Interface().(map[string]ParseableURL)
				targetMap := make(map[string]url.URL, len(sourceMap))
				for k, v := range sourceMap {
					u := url.URL(v)
					targetMap[k] = u
				}
				targetField.Set(reflect.ValueOf(targetMap))
			case urlMapType:
				sourceMap := sourceField.Interface().(map[string]url.URL)
				targetMap := make(map[string]ParseableURL, len(sourceMap))
				for k, v := range sourceMap {
					u := ParseableURL(v)
					targetMap[k] = u
				}
				targetField.Set(reflect.ValueOf(targetMap))
			case parseableURLMapPtrType:
				sourceMap := sourceField.Interface().(map[string]*ParseableURL)
				targetMap := make(map[string]*url.URL, len(sourceMap))
				for k, v := range sourceMap {
					u := url.URL(*v)
					targetMap[k] = &u
				}
				targetField.Set(reflect.ValueOf(targetMap))
			case urlMapPtrType:
				sourceMap := sourceField.Interface().(map[string]*url.URL)
				targetMap := make(map[string]*ParseableURL, len(sourceMap))
				for k, v := range sourceMap {
					u := ParseableURL(*v)
					targetMap[k] = &u
				}
				targetField.Set(reflect.ValueOf(targetMap))
			default:
				targetField.Set(sourceField)
			}
		}
	}
}

func isZero(v reflect.Value, tag string) bool {
	if tag == _zeroPointer {
		return true
	}
	switch v.Kind() {
	case reflect.Func, reflect.Struct, reflect.Ptr:
		if v.Kind() == reflect.Ptr && v.IsNil() {
			return true
		}
		switch v.Type() {
		case durationType, parseableDurationType:
			return v.Interface() == 0
		case timeType:
			return v.Interface().(time.Time).IsZero()
		case parseableTimeType:
			return time.Time(v.Interface().(ParseableTime)).IsZero()
		case urlType:
			u := v.Interface().(url.URL)
			return u.String() == ""
		case parseableURLType:
			u := url.URL(v.Interface().(ParseableURL))
			return u.String() == ""
		case durationPtrType, parseableDurationPtrType, timePtrType,
			parseableTimePtrType, urlPtrType, parseableURLPtrType:
			return v.IsZero()
		}
		return v.IsNil()
	case wrappedBoolKind:
		return false
	case reflect.Map, reflect.Slice, parseableTimeKind:
		return v.IsNil() || v.Len() == 0
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
	zeroPointer string,
	listSeparator string,
	keyValueSeparator string,
	timeFormats []string,
) {
	val := reflect.Indirect(v)
	for i := range val.NumField() {
		fld := val.Field(i)
		field := val.Type().Field(i)
		if recursiveStruct(fld) {
			setDefaults(fld, zeroPointer, listSeparator, keyValueSeparator, timeFormats)
			continue
		}
		tag := field.Tag.Get("def")
		if !isZero(fld, tag) && (fld.Type().Kind() != reflect.Map || fld.Len() > 0) {
			continue
		}

		if tag != "" {
			if fld.Type().Kind() == reflect.Map {
				fld.Set(reflect.MakeMap(fld.Type()))
			} else if fld.Type().Kind() == reflect.Slice {
				fld.Set(reflect.MakeSlice(fld.Type(), 0, 0))
			} else {
				switch fld.Type() {
				case parseableDurationType:
					if d, err := time.ParseDuration(tag); err == nil {
						x := ParseableDuration(d)
						fld.Set(reflect.ValueOf(x))
					}
					continue
				case parseableDurationPtrType:
					if tag == zeroPointer {
						tag = "0s"
					}
					if d, err := time.ParseDuration(tag); err == nil {
						x := ParseableDuration(d)
						fld.Set(reflect.ValueOf(&x))
					}
					continue
				case parseableTimeType:
					for _, timeFormat := range timeFormats {
						if t, err := time.Parse(timeFormat, tag); err == nil {
							x := ParseableTime(t)
							fld.Set(reflect.ValueOf(x))
							break
						}
					}
					continue
				case parseableTimePtrType:
					if tag == zeroPointer {
						t := new(time.Time)
						x := ParseableTime(*t)
						fld.Set(reflect.ValueOf(&x))
					} else {
						for _, timeFormat := range timeFormats {
							if t, err := time.Parse(timeFormat, tag); err == nil {
								x := ParseableTime(t)
								fld.Set(reflect.ValueOf(&x))
								break
							}
						}
					}
					continue
				case parseableURLType:
					if u, err := url.Parse(tag); err == nil {
						x := ParseableURL(*u)
						fld.Set(reflect.ValueOf(x))
					}
					continue
				case parseableURLPtrType:
					if tag == zeroPointer {
						tag = ""
					}
					if u, err := url.Parse(tag); err == nil {
						x := ParseableURL(*u)
						fld.Set(reflect.ValueOf(&x))
					}
					continue
				case complex128Type:
					if c, err := strconv.ParseComplex(tag, 128); err == nil {
						x := ParseableComplex128(c)
						fld.Set(reflect.ValueOf(x))
					}
				case complex128PtrType:
					if tag == zeroPointer {
						tag = "0"
					}
					if c, err := strconv.ParseComplex(tag, 128); err == nil {
						x := ParseableComplex128(c)
						fld.Set(reflect.ValueOf(&x))
					}
				case complex64Type:
					if c, err := strconv.ParseComplex(tag, 64); err == nil {
						x := ParseableComplex64(c)
						fld.Set(reflect.ValueOf(x))
					}
				case complex64PtrType:
					if tag == zeroPointer {
						tag = "0"
					}
					if c, err := strconv.ParseComplex(tag, 64); err == nil {
						x := ParseableComplex64(c)
						fld.Set(reflect.ValueOf(&x))
					}
				}

				if fld.Kind() == reflect.Ptr {
					fld.Set(reflect.New(field.Type.Elem()))
				} else {
					fld.Set(reflect.New(field.Type).Elem())
				}
			}
			setFieldValue(
				fld,
				tag,
				zeroPointer,
				listSeparator,
				keyValueSeparator,
				timeFormats,
			)
		}
	}
}

func setFieldValue(
	fld reflect.Value,
	val string,
	zeroPointer string,
	listSeparator string,
	keyValueSeparator string,
	timeFormats []string,
) {
	f := fld
	if fld.Kind() == reflect.Ptr {
		f = reflect.Indirect(fld)
	}

	switch f.Type().Kind() {
	case reflect.String:
		if val == zeroPointer {
			val = ""
		}
		c := reflect.ValueOf(val).Convert(f.Type())
		f.Set(c)
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
		if val == zeroPointer {
			val = "0"
		}
		i, err := strconv.ParseInt(val, 10, 0)
		if err != nil {
			return
		}
		c := reflect.ValueOf(int(i)).Convert(f.Type())
		f.Set(c)
	case reflect.Float32, reflect.Float64:
		if val == zeroPointer {
			val = "0.0"
		}
		v, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return
		}
		c := reflect.ValueOf(v).Convert(f.Type())
		f.Set(c)
	case reflect.Complex128, reflect.Complex64:
		if val == zeroPointer {
			val = "0+0i"
		}
		v, err := strconv.ParseComplex(val, 128)
		if err != nil {
			return
		}
		c := reflect.ValueOf(v).Convert(f.Type())
		f.Set(c)
	case reflect.Bool:
		if val == zeroPointer {
			val = "false"
		}
		b, err := strconv.ParseBool(val)
		if err != nil {
			return
		}
		c := reflect.ValueOf(b).Convert(f.Type())
		f.Set(c)
	case reflect.Map:
		var kvs []string
		vals := map[string]string{}
		if val == zeroPointer || val == "" {
			return
		}

		kvs = strings.Split(val, listSeparator)
		elemType := reflect.TypeOf(f.Interface()).Elem()

		switch elemType {
		case parseableTimeType, parseableTimePtrType, parseableURLType, parseableURLPtrType:
			for _, kv := range kvs {
				mapVals := strings.Split(kv, keyValueSeparator)
				value := strings.Join(mapVals[1:], keyValueSeparator)
				vals[mapVals[0]] = value
			}
		default:
			for _, kv := range kvs {
				mapVals := strings.Split(kv, keyValueSeparator)
				if len(mapVals) != 2 {
					continue
				}
				vals[mapVals[0]] = mapVals[1]
			}
		}
		c := reflect.MakeMap(f.Type())
		pointer := false
		elemKind := elemType.Kind()
		if elemKind == reflect.Ptr {
			pointer = true
			elemKind = elemType.Elem().Kind()
		}
		switch elemType {
		case parseableDurationType:
			for k, v := range vals {
				if d, err := time.ParseDuration(v); err == nil {
					pd := ParseableDuration(d)
					c.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(pd))
				}
			}
			f.Set(c)
			return
		case parseableDurationPtrType:
			for k, v := range vals {
				if d, err := time.ParseDuration(v); err == nil {
					pd := ParseableDuration(d)
					c.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(&pd))
				}
			}
			f.Set(c)
			return
		case parseableTimeType:
			for k, v := range vals {
				for _, timeFormat := range timeFormats {
					if t, err := time.Parse(timeFormat, v); err == nil {
						pt := ParseableTime(t)
						c.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(pt))
						break
					}
				}
			}
			f.Set(c)
			return
		case parseableTimePtrType:
			for k, v := range vals {
				for _, timeFormat := range timeFormats {
					if t, err := time.Parse(timeFormat, v); err == nil {
						pt := ParseableTime(t)
						c.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(&pt))
						break
					}
				}
			}
			f.Set(c)
			return
		case parseableURLType:
			for k, v := range vals {
				if u, err := url.Parse(v); err == nil {
					p := ParseableURL(*u)
					c.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(p))
				}
			}
			f.Set(c)
			return
		case parseableURLPtrType:
			for k, v := range vals {
				if u, err := url.Parse(v); err == nil {
					p := ParseableURL(*u)
					c.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(&p))
				}
			}
			f.Set(c)
			return
		}
		for k, v := range vals {
			switch elemKind {
			case reflect.String:
				if pointer {
					c.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(&v).Convert(f.Type().Elem()))
				} else {
					c.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v).Convert(f.Type().Elem()))
				}
			case reflect.Bool:
				b, err := strconv.ParseBool(v)
				if err != nil {
					continue
				}
				if pointer {
					c.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(&b).Convert(elemType))
				} else {
					c.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(b).Convert(elemType))
				}
			case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
				i, err := strconv.ParseInt(v, 10, 0)
				if err != nil {
					continue
				}
				if pointer {
					switch elemKind {
					case reflect.Int:
						value := int(i)
						c.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(&value))
					case reflect.Int8:
						value := int8(i)
						c.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(&value))
					case reflect.Int16:
						value := int16(i)
						c.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(&value))
					case reflect.Int32:
						value := int32(i)
						c.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(&value))
					case reflect.Uint:
						value := uint(i)
						c.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(&value))
					case reflect.Uint8:
						value := uint8(i)
						c.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(&value))
					case reflect.Uint16:
						value := uint16(i)
						c.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(&value))
					case reflect.Uint32:
						value := uint32(i)
						c.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(&value))
					case reflect.Uint64:
						value := uint64(i)
						c.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(&value))
					default:
						c.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(&i))
					}
				} else {
					c.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(i).Convert(elemType))
				}
			case reflect.Complex128, reflect.Complex64:
				i, err := strconv.ParseComplex(v, 128)
				if err != nil {
					continue
				}
				if pointer {
					switch elemKind {
					case reflect.Complex64:
						value := complex64(i)
						c.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(&value).Convert(elemType))
					default:
						c.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(&i).Convert(elemType))
					}
				} else {
					c.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(i).Convert(elemType))
				}
			case reflect.Float32, reflect.Float64:
				v, err := strconv.ParseFloat(v, 64)
				if err != nil {
					continue
				}
				if pointer {
					switch elemKind {
					case reflect.Float32:
						value := float32(v)
						c.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(&value))
					default:
						c.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(&v))
					}
				} else {
					c.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v).Convert(elemType))
				}
			default:
				continue
			}
		}
		f.Set(c)
	case reflect.Slice:
		var vals []string
		if val == "" || val == zeroPointer {
			return
		}

		vals = strings.Split(val, listSeparator)

		c := reflect.MakeSlice(f.Type(), 0, len(vals))
		elemType := reflect.TypeOf(f.Interface()).Elem()
		pointer := false
		elemKind := elemType.Kind()
		if elemKind == reflect.Ptr {
			pointer = true
			elemKind = elemType.Elem().Kind()
		}
		switch elemType {
		case parseableDurationType:
			for _, dur := range vals {
				if d, err := time.ParseDuration(dur); err == nil {
					pd := ParseableDuration(d)
					c = reflect.Append(c, reflect.ValueOf(pd))
				}
			}
			f.Set(c)
			return
		case parseableDurationPtrType:
			for _, dur := range vals {
				if d, err := time.ParseDuration(dur); err == nil {
					pd := ParseableDuration(d)
					c = reflect.Append(c, reflect.ValueOf(&pd))
				}
			}
			f.Set(c)
			return
		case parseableTimeType:
			for _, tm := range vals {
				for _, timeFormat := range timeFormats {
					if t, err := time.Parse(timeFormat, tm); err == nil {
						pt := ParseableTime(t)
						c = reflect.Append(c, reflect.ValueOf(pt))
						break
					}
				}
			}
			f.Set(c)
			return
		case parseableTimePtrType:
			for _, tm := range vals {
				for _, timeFormat := range timeFormats {
					if t, err := time.Parse(timeFormat, tm); err == nil {
						pt := ParseableTime(t)
						c = reflect.Append(c, reflect.ValueOf(&pt))
						break
					}
				}
			}
			f.Set(c)
			return
		case parseableURLType:
			for _, v := range vals {
				if u, err := url.Parse(v); err == nil {
					p := ParseableURL(*u)
					c = reflect.Append(c, reflect.ValueOf(p))
				}
			}
			f.Set(c)
			return
		case parseableURLPtrType:
			for _, v := range vals {
				if u, err := url.Parse(v); err == nil {
					p := ParseableURL(*u)
					c = reflect.Append(c, reflect.ValueOf(&p))
				}
			}
			f.Set(c)
			return
		}

		for _, v := range vals {
			switch elemKind {
			case reflect.String:
				if pointer {
					c = reflect.Append(c, reflect.ValueOf(&v).Convert(f.Type().Elem()))
				} else {
					c = reflect.Append(c, reflect.ValueOf(v).Convert(f.Type().Elem()))
				}
			case reflect.Bool:
				b, err := strconv.ParseBool(v)
				if err != nil {
					continue
				}
				if pointer {
					c = reflect.Append(c, reflect.ValueOf(&b).Convert(elemType))
				} else {
					c = reflect.Append(c, reflect.ValueOf(b).Convert(elemType))
				}
			case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
				i, err := strconv.ParseInt(v, 10, 0)
				if err != nil {
					continue
				}
				if pointer {
					switch elemKind {
					case reflect.Int:
						value := int(i)
						c = reflect.Append(c, reflect.ValueOf(&value))
					case reflect.Int8:
						value := int8(i)
						c = reflect.Append(c, reflect.ValueOf(&value))
					case reflect.Int16:
						value := int16(i)
						c = reflect.Append(c, reflect.ValueOf(&value))
					case reflect.Int32:
						value := int32(i)
						c = reflect.Append(c, reflect.ValueOf(&value))
					case reflect.Uint:
						value := uint(i)
						c = reflect.Append(c, reflect.ValueOf(&value))
					case reflect.Uint8:
						value := uint8(i)
						c = reflect.Append(c, reflect.ValueOf(&value))
					case reflect.Uint16:
						value := uint16(i)
						c = reflect.Append(c, reflect.ValueOf(&value))
					case reflect.Uint32:
						value := uint32(i)
						c = reflect.Append(c, reflect.ValueOf(&value))
					case reflect.Uint64:
						value := uint64(i)
						c = reflect.Append(c, reflect.ValueOf(&value))
					default:
						c = reflect.Append(c, reflect.ValueOf(&i))
					}
				} else {
					c = reflect.Append(c, reflect.ValueOf(i).Convert(elemType))
				}
			case reflect.Float32, reflect.Float64:
				v, err := strconv.ParseFloat(v, 64)
				if err != nil {
					continue
				}
				if pointer {
					switch elemKind {
					case reflect.Float32:
						value := float32(v)
						c = reflect.Append(c, reflect.ValueOf(&value))
					default:
						c = reflect.Append(c, reflect.ValueOf(&v))
					}
				} else {
					c = reflect.Append(c, reflect.ValueOf(v).Convert(elemType))
				}
			case reflect.Complex128, reflect.Complex64:
				i, err := strconv.ParseComplex(v, 128)
				if err != nil {
					continue
				}
				if pointer {
					switch elemKind {
					case reflect.Complex64:
						value := complex64(i)
						c = reflect.Append(c, reflect.ValueOf(&value).Convert(elemType))
					default:
						c = reflect.Append(c, reflect.ValueOf(&i).Convert(elemType))
					}
				} else {
					c = reflect.Append(c, reflect.ValueOf(i).Convert(elemType))
				}
			default:
				continue
			}
		}
		f.Set(c)
	default:
		return
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
		interpEnv(field, fld)

		for _, customFunc := range customFuncs {
			customFunc(field, fld)
		}
	}
}

func validateEnums(c reflect.Value) []error {
	var errs []error
	if !recursiveStruct(c) {
		return errs
	}
	cVal := c
	if cVal.Kind() == reflect.Ptr {
		cVal = reflect.Indirect(c)
	}
	for i := range cVal.NumField() {
		fld := cVal.Field(i)
		if recursiveStruct(fld) {
			if subErrs := validateEnums(fld); len(subErrs) > 0 {
				errs = append(errs, subErrs...)
			}
		}

		if isZero(fld, "") {
			continue
		}

		var t reflect.Type
		if fld.Kind() == reflect.Ptr {
			t = fld.Elem().Type()
		} else {
			t = fld.Type()
		}
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

	var newArgs []string

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
	var values []string
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
	var values []string

	if !strings.Contains(valArg, listSeparator) || !verifySliceOrMapType(c, arg) {
		values = append(values, valArg)
		return values
	}

	return strings.Split(valArg, listSeparator)
}

func verifySliceOrMapType(c reflect.Value, arg string) bool {
	if !recursiveStruct(c) {
		return false
	}

	bareArg := strings.TrimPrefix(strings.TrimPrefix(arg, "-"), "-")

	cVal := reflect.Indirect(c)
	for i := range cVal.NumField() {
		field := cVal.Type().Field(i)
		fld := cVal.Field(i)
		if recursiveStruct(fld) {
			if verifySliceOrMapType(fld, arg) {
				return true
			}
			continue
		}

		shortTag := field.Tag.Get("short")
		longTag := field.Tag.Get("long")

		if shortTag != bareArg && longTag != bareArg {
			continue
		}

		if fld.Type().Kind() == reflect.Slice || fld.Type().Kind() == reflect.Map {
			return true
		}

	}
	return false
}

func locateCLIFiles(options *Options, listSeparator string) ([]string, []string) {
	var args []string
	var files []string

	shortPrefix := options.ConfigFileShort
	if shortPrefix != "" && !strings.HasPrefix(shortPrefix, "-") {
		shortPrefix = "-" + options.ConfigFileShort
	}
	longPrefix := options.ConfigFileLong
	if longPrefix != "" && !strings.HasPrefix(longPrefix, "--") {
		longPrefix = "--" + options.ConfigFileLong
	}

	skip := false

	cliArgs := os.Args[1:]

	for i, arg := range cliArgs {
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

		if i < len(cliArgs)-1 {
			nextArg := cliArgs[i+1]
			files = append(files, strings.Split(nextArg, listSeparator)...)
			skip = true
		}
	}
	return files, args
}
