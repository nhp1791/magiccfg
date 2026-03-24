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
	c.populateEmptyStructs(newC)
	return newC
}

func (c *MagicConfig[T]) recursiveStruct(v reflect.Value) bool {
	vType := v.Type()
	k := v.Kind()
	switch vType {
	case c.types.timePtrType, c.types.parseableTimePtrType, c.types.urlPtrType, c.types.parseableURLPtrType:
		return false
	default:
		for _, t := range c.customTypes {
			if vType == t.basicType {
				return false
			}
			if k == reflect.Ptr {
				ut := vType.Elem()
				if ut == t.basicType {
					return false
				}
			}
		}
	}

	return v.Kind() == reflect.Ptr && vType.Elem().Kind() == reflect.Struct
}

func (c *MagicConfig[T]) makeFullConfig(
	v reflect.Value,
	prefix string,
	parentName *string,
	envParentName *string,
) (reflect.Value, map[string]map[string]int) {
	usages := map[string]map[string]int{
		_envTag: {},
		"long":  {},
		"short": {},
	}

	if !c.recursiveStruct(v) {
		return reflect.Zero(v.Type()), usages
	}

	c.populateEmptyStructs(v)

	fields := make([]reflect.StructField, 0)
	var newField reflect.StructField
	vVal := v.Elem()

	for i := range vVal.NumField() {
		field := vVal.Type().Field(i)
		fld := vVal.Field(i)
		name := field.Name
		structField := c.recursiveStruct(fld)

		newParentName := convertNameToCommandLine(name, parentName)

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
		case c.types.boolType:
			if wrappedBool {
				fieldType = c.types.wrappedBoolType
			}
		case c.types.boolPtrType:
			if wrappedBool {
				fieldType = c.types.wrappedBoolType
			}
		case c.types.boolSliceType:
			if wrappedBool {
				fieldType = c.types.wrappedBoolSliceType
			}
		case c.types.boolPtrSliceType:
			if wrappedBool {
				fieldType = c.types.wrappedBoolSliceType
			}
		case c.types.durationType:
			fieldType = c.types.parseableDurationType
		case c.types.durationSliceType:
			fieldType = c.types.parseableDurationSliceType
		case c.types.durationMapType:
			fieldType = c.types.parseableDurationMapType
		case c.types.durationPtrType:
			fieldType = c.types.parseableDurationPtrType
		case c.types.durationSlicePtrType:
			fieldType = c.types.parseableDurationSlicePtrType
		case c.types.durationMapPtrType:
			fieldType = c.types.parseableDurationMapPtrType
		case c.types.timeType:
			fieldType = c.types.parseableTimeType
		case c.types.timeSliceType:
			fieldType = c.types.parseableTimeSliceType
		case c.types.timeMapType:
			fieldType = c.types.parseableTimeMapType
		case c.types.timePtrType:
			fieldType = c.types.parseableTimePtrType
		case c.types.timeSlicePtrType:
			fieldType = c.types.parseableTimeSlicePtrType
		case c.types.timeMapPtrType:
			fieldType = c.types.parseableTimeMapPtrType
		case c.types.urlType:
			fieldType = c.types.parseableURLType
		case c.types.urlSliceType:
			fieldType = c.types.parseableURLSliceType
		case c.types.urlMapType:
			fieldType = c.types.parseableURLMapType
		case c.types.urlPtrType:
			fieldType = c.types.parseableURLPtrType
		case c.types.urlSlicePtrType:
			fieldType = c.types.parseableURLSlicePtrType
		case c.types.urlMapPtrType:
			fieldType = c.types.parseableURLMapPtrType
		case c.types.complex128Type:
			fieldType = c.types.parseableComplex128Type
		case c.types.complex128PtrType:
			fieldType = c.types.parseableComplex128PtrType
		case c.types.complex64Type:
			fieldType = c.types.parseableComplex64Type
		case c.types.complex64PtrType:
			fieldType = c.types.parseableComplex64PtrType
		case c.types.complex64SliceType:
			fieldType = c.types.parseableComplex64SliceType
		case c.types.complex64SlicePtrType:
			fieldType = c.types.parseableComplex64SlicePtrType
		case c.types.complex64MapType:
			fieldType = c.types.parseableComplex64MapType
		case c.types.complex64MapPtrType:
			fieldType = c.types.parseableComplex64MapPtrType
		case c.types.complex128SliceType:
			fieldType = c.types.parseableComplex128SliceType
		case c.types.complex128SlicePtrType:
			fieldType = c.types.parseableComplex128SlicePtrType
		case c.types.complex128MapType:
			fieldType = c.types.parseableComplex128MapType
		case c.types.complex128MapPtrType:
			fieldType = c.types.parseableComplex128MapPtrType
		}

		for _, ct := range c.customTypes {
			if fieldType == ct.listType && ct.parseableListType != nil {
				fieldType = ct.parseableListType
			}
			if fieldType == ct.listPointerType && ct.parseableListPointerType != nil {
				fieldType = ct.parseableListPointerType
			}
			if fieldType == ct.mapType && ct.parseableMapType != nil {
				fieldType = ct.parseableMapType
			}
			if fieldType == ct.mapPointerType && ct.parseableMapPointerType != nil {
				fieldType = ct.parseableMapPointerType
			}
		}

		if structField {
			f, subUsages := c.makeFullConfig(fld, prefix, &newParentName, &newEnvParentName)
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

	c.populateEmptyStructs(newValue)
	return newValue, usages
}

func (c *MagicConfig[T]) makeValidateConfig(v reflect.Value) reflect.Value {
	if !c.recursiveStruct(v) {
		return reflect.Zero(v.Type())
	}

	c.populateEmptyStructs(v)

	fields := make([]reflect.StructField, 0)
	var newField reflect.StructField
	vVal := v.Elem()

	for i := range vVal.NumField() {
		field := vVal.Type().Field(i)
		fld := vVal.Field(i)
		name := field.Name
		structField := c.recursiveStruct(fld)

		tag := field.Tag.Get("validate")

		fieldType := field.Type
		switch fld.Type() {
		case c.types.parseableDurationType:
			fieldType = c.types.durationType
		case c.types.parseableDurationPtrType:
			fieldType = c.types.durationPtrType
		case c.types.parseableTimeType:
			fieldType = c.types.timeType
		case c.types.parseableTimePtrType:
			fieldType = c.types.timePtrType
		case c.types.parseableURLType:
			fieldType = c.types.urlType
		case c.types.parseableURLPtrType:
			fieldType = c.types.urlPtrType
		}

		if structField {
			f := c.makeValidateConfig(fld)
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

	c.populateEmptyStructs(newValue)
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

func (c *MagicConfig[T]) populateEmptyStructs(v reflect.Value) {
	if !c.recursiveStruct(v) {
		return
	}

	vVal := reflect.Indirect(v)

	for i := range vVal.NumField() {
		field := vVal.Type().Field(i)
		fld := vVal.Field(i)
		name := field.Name
		_ = name
		if !c.recursiveStruct(fld) {
			continue
		}
		if fld.IsNil() {
			n := reflect.New(field.Type.Elem())
			c.populateEmptyStructs(n)
			fld.Set(n)
		}
	}
}

func (c *MagicConfig[T]) merge(oldV, newV reflect.Value) {
	target := reflect.Indirect(oldV)
	source := reflect.Indirect(newV)
	for i := range target.NumField() {
		targetField := target.Field(i)
		sourceField := source.Field(i)

		if c.recursiveStruct(targetField) {
			c.merge(targetField, sourceField)
			continue
		}
		if !c.isZero(sourceField, "") && targetField.CanSet() {
			targetFieldType := targetField.Type()
			sourceFieldType := sourceField.Type()
			if sourceFieldType == targetFieldType {
				targetField.Set(sourceField)
				continue
			}
			custom := false
			for _, ct := range c.customTypes {
				if sourceFieldType == ct.parseableMapType && ct.mergeMap != nil {
					ct.mergeMap(sourceField, targetField)
					custom = true
					break
				}
			}
			if custom {
				continue
			}
			switch sourceFieldType {
			case c.types.durationType, c.types.durationPtrType, c.types.parseableDurationType,
				c.types.parseableDurationPtrType, c.types.timeType, c.types.timePtrType,
				c.types.parseableTimeType, c.types.parseableTimePtrType, c.types.parseableURLType,
				c.types.urlType, c.types.urlPtrType, c.types.parseableURLPtrType, c.types.complex128Type,
				c.types.parseableComplex128Type, c.types.complex64Type, c.types.parseableComplex64Type,
				c.types.complex128PtrType, c.types.parseableComplex128PtrType, c.types.complex64PtrType,
				c.types.parseableComplex64PtrType:
				targetField.Set(sourceField.Convert(targetFieldType))
			case c.types.boolType:
				if targetFieldType == c.types.wrappedBoolType {
					sourceValue := sourceField.Interface().(bool)
					targetValue := Bool{Value: &sourceValue}
					targetField.Set(reflect.ValueOf(targetValue))
				} else {
					targetField.Set(sourceField.Convert(targetFieldType))
				}
			case c.types.boolPtrType:
				if targetFieldType == c.types.wrappedBoolType {
					sourceValue := sourceField.Interface().(*bool)
					targetValue := Bool{Value: sourceValue}
					targetField.Set(reflect.ValueOf(targetValue))
				} else {
					targetField.Set(sourceField.Convert(targetFieldType))
				}
			case c.types.boolSliceType:
				if targetFieldType == c.types.wrappedBoolSliceType {
					sourceSlice := sourceField.Interface().([]bool)
					targetSlice := make([]Bool, len(sourceSlice))
					for k := range sourceSlice {
						targetSlice[k] = Bool{Value: &sourceSlice[k]}
					}
					targetField.Set(reflect.ValueOf(targetSlice))
				} else {
					targetField.Set(sourceField.Convert(targetFieldType))
				}
			case c.types.boolPtrSliceType:
				if targetFieldType == c.types.wrappedBoolSliceType {
					sourceSlice := sourceField.Interface().([]*bool)
					targetSlice := make([]Bool, len(sourceSlice))
					for k := range sourceSlice {
						targetSlice[k] = Bool{Value: sourceSlice[k]}
					}
					targetField.Set(reflect.ValueOf(targetSlice))
				} else {
					targetField.Set(sourceField.Convert(targetFieldType))
				}
			case c.types.complex64SliceType:
				sourceSlice := sourceField.Interface().([]complex64)
				targetSlice := make([]ParseableComplex64, len(sourceSlice))
				for i := range sourceSlice {
					targetSlice[i] = ParseableComplex64(sourceSlice[i])
				}
				targetField.Set(reflect.ValueOf(targetSlice))
			case c.types.complex64SlicePtrType:
				sourceSlice := sourceField.Interface().([]*complex64)
				targetSlice := make([]*ParseableComplex64, len(sourceSlice))
				for i := range sourceSlice {
					j := ParseableComplex64(*sourceSlice[i])
					targetSlice[i] = &j
				}
				targetField.Set(reflect.ValueOf(targetSlice))
			case c.types.complex64MapType:
				sourceMap := sourceField.Interface().(map[string]complex64)
				targetMap := make(map[string]ParseableComplex64, len(sourceMap))
				for k, v := range sourceMap {
					targetMap[k] = ParseableComplex64(v)
				}
				targetField.Set(reflect.ValueOf(targetMap))
			case c.types.complex64MapPtrType:
				sourceMap := sourceField.Interface().(map[string]*complex64)
				targetMap := make(map[string]*ParseableComplex64, len(sourceMap))
				for k, v := range sourceMap {
					j := ParseableComplex64(*v)
					targetMap[k] = &j
				}
				targetField.Set(reflect.ValueOf(targetMap))
			case c.types.complex128SliceType:
				sourceSlice := sourceField.Interface().([]complex128)
				targetSlice := make([]ParseableComplex128, len(sourceSlice))
				for i := range sourceSlice {
					targetSlice[i] = ParseableComplex128(sourceSlice[i])
				}
				targetField.Set(reflect.ValueOf(targetSlice))
			case c.types.complex128SlicePtrType:
				sourceSlice := sourceField.Interface().([]*complex128)
				targetSlice := make([]*ParseableComplex128, len(sourceSlice))
				for i := range sourceSlice {
					j := ParseableComplex128(*sourceSlice[i])
					targetSlice[i] = &j
				}
				targetField.Set(reflect.ValueOf(targetSlice))
			case c.types.complex128MapType:
				sourceMap := sourceField.Interface().(map[string]complex128)
				targetMap := make(map[string]ParseableComplex128, len(sourceMap))
				for k, v := range sourceMap {
					targetMap[k] = ParseableComplex128(v)
				}
				targetField.Set(reflect.ValueOf(targetMap))
			case c.types.complex128MapPtrType:
				sourceMap := sourceField.Interface().(map[string]*complex128)
				targetMap := make(map[string]*ParseableComplex128, len(sourceMap))
				for k, v := range sourceMap {
					j := ParseableComplex128(*v)
					targetMap[k] = &j
				}
				targetField.Set(reflect.ValueOf(targetMap))
			case c.types.parseableComplex64SliceType:
				sourceSlice := sourceField.Interface().([]ParseableComplex64)
				targetSlice := make([]complex64, len(sourceSlice))
				for i := range sourceSlice {
					targetSlice[i] = complex64(sourceSlice[i])
				}
				targetField.Set(reflect.ValueOf(targetSlice))
			case c.types.parseableComplex64SlicePtrType:
				sourceSlice := sourceField.Interface().([]*ParseableComplex64)
				targetSlice := make([]*complex64, len(sourceSlice))
				for i := range sourceSlice {
					j := complex64(*sourceSlice[i])
					targetSlice[i] = &j
				}
				targetField.Set(reflect.ValueOf(targetSlice))
			case c.types.parseableComplex64MapType:
				sourceMap := sourceField.Interface().(map[string]ParseableComplex64)
				targetMap := make(map[string]complex64, len(sourceMap))
				for k, v := range sourceMap {
					targetMap[k] = complex64(v)
				}
				targetField.Set(reflect.ValueOf(targetMap))
			case c.types.parseableComplex64MapPtrType:
				sourceMap := sourceField.Interface().(map[string]*ParseableComplex64)
				targetMap := make(map[string]*complex64, len(sourceMap))
				for k, v := range sourceMap {
					j := complex64(*v)
					targetMap[k] = &j
				}
				targetField.Set(reflect.ValueOf(targetMap))
			case c.types.parseableComplex128SliceType:
				sourceSlice := sourceField.Interface().([]ParseableComplex128)
				targetSlice := make([]complex128, len(sourceSlice))
				for i := range sourceSlice {
					targetSlice[i] = complex128(sourceSlice[i])
				}
				targetField.Set(reflect.ValueOf(targetSlice))
			case c.types.parseableComplex128SlicePtrType:
				sourceSlice := sourceField.Interface().([]*ParseableComplex128)
				targetSlice := make([]*complex128, len(sourceSlice))
				for i := range sourceSlice {
					j := complex128(*sourceSlice[i])
					targetSlice[i] = &j
				}
				targetField.Set(reflect.ValueOf(targetSlice))
			case c.types.parseableComplex128MapType:
				sourceMap := sourceField.Interface().(map[string]ParseableComplex128)
				targetMap := make(map[string]complex128, len(sourceMap))
				for k, v := range sourceMap {
					targetMap[k] = complex128(v)
				}
				targetField.Set(reflect.ValueOf(targetMap))
			case c.types.parseableComplex128MapPtrType:
				sourceMap := sourceField.Interface().(map[string]*ParseableComplex128)
				targetMap := make(map[string]*complex128, len(sourceMap))
				for k, v := range sourceMap {
					j := complex128(*v)
					targetMap[k] = &j
				}
				targetField.Set(reflect.ValueOf(targetMap))
			case c.types.durationSliceType:
				sourceSlice := sourceField.Interface().([]time.Duration)
				targetSlice := make([]ParseableDuration, len(sourceSlice))
				for i := range sourceSlice {
					targetSlice[i] = ParseableDuration(sourceSlice[i])
				}
				targetField.Set(reflect.ValueOf(targetSlice))
			case c.types.durationMapType:
				sourceMap := sourceField.Interface().(map[string]time.Duration)
				targetMap := make(map[string]ParseableDuration, len(sourceMap))
				for k, v := range sourceMap {
					targetMap[k] = ParseableDuration(v)
				}
				targetField.Set(reflect.ValueOf(targetMap))
			case c.types.parseableDurationSliceType:
				sourceSlice := sourceField.Interface().([]ParseableDuration)
				targetSlice := make([]time.Duration, len(sourceSlice))
				for i := range sourceSlice {
					targetSlice[i] = time.Duration(sourceSlice[i])
				}
				targetField.Set(reflect.ValueOf(targetSlice))
			case c.types.parseableDurationMapType:
				sourceMap := sourceField.Interface().(map[string]ParseableDuration)
				targetMap := make(map[string]time.Duration, len(sourceMap))
				for k, v := range sourceMap {
					targetMap[k] = time.Duration(v)
				}
				targetField.Set(reflect.ValueOf(targetMap))
			case c.types.durationSlicePtrType:
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
			case c.types.durationMapPtrType:
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
			case c.types.parseableDurationSlicePtrType:
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
			case c.types.parseableDurationMapPtrType:
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
			case c.types.timeSliceType:
				sourceSlice := sourceField.Interface().([]time.Time)
				targetSlice := make([]ParseableTime, len(sourceSlice))
				for i := range sourceSlice {
					targetSlice[i] = ParseableTime(sourceSlice[i])
				}
				targetField.Set(reflect.ValueOf(targetSlice))
			case c.types.timeMapType:
				sourceMap := sourceField.Interface().(map[string]time.Time)
				targetMap := make(map[string]ParseableTime, len(sourceMap))
				for k, v := range sourceMap {
					targetMap[k] = ParseableTime(v)
				}
				targetField.Set(reflect.ValueOf(targetMap))
			case c.types.parseableTimeSliceType:
				sourceSlice := sourceField.Interface().([]ParseableTime)
				targetSlice := make([]time.Time, len(sourceSlice))
				for i := range sourceSlice {
					targetSlice[i] = time.Time(sourceSlice[i])
				}
				targetField.Set(reflect.ValueOf(targetSlice))
			case c.types.parseableTimeMapType:
				sourceMap := sourceField.Interface().(map[string]ParseableTime)
				targetMap := make(map[string]time.Time, len(sourceMap))
				for k, v := range sourceMap {
					targetMap[k] = time.Time(v)
				}
				targetField.Set(reflect.ValueOf(targetMap))
			case c.types.timeSlicePtrType:
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
			case c.types.timeMapPtrType:
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
			case c.types.parseableTimeSlicePtrType:
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
			case c.types.parseableTimeMapPtrType:
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
			case c.types.parseableURLSliceType:
				sourceSlice := sourceField.Interface().([]ParseableURL)
				targetSlice := make([]url.URL, len(sourceSlice))
				for k := range sourceSlice {
					u := url.URL(sourceSlice[k])
					targetSlice[k] = u
				}
				targetField.Set(reflect.ValueOf(targetSlice))
			case c.types.urlSliceType:
				sourceSlice := sourceField.Interface().([]url.URL)
				targetSlice := make([]ParseableURL, len(sourceSlice))
				for k := range sourceSlice {
					p := ParseableURL(sourceSlice[k])
					targetSlice[k] = p
				}
				targetField.Set(reflect.ValueOf(targetSlice))
			case c.types.parseableURLSlicePtrType:
				sourceSlice := sourceField.Interface().([]*ParseableURL)
				targetSlice := make([]*url.URL, len(sourceSlice))
				for k := range sourceSlice {
					u := url.URL(*sourceSlice[k])
					targetSlice[k] = &u
				}
				targetField.Set(reflect.ValueOf(targetSlice))
			case c.types.urlSlicePtrType:
				sourceSlice := sourceField.Interface().([]*url.URL)
				targetSlice := make([]*ParseableURL, len(sourceSlice))
				for k := range sourceSlice {
					p := ParseableURL(*sourceSlice[k])
					targetSlice[k] = &p
				}
				targetField.Set(reflect.ValueOf(targetSlice))
			case c.types.parseableURLMapType:
				sourceMap := sourceField.Interface().(map[string]ParseableURL)
				targetMap := make(map[string]url.URL, len(sourceMap))
				for k, v := range sourceMap {
					u := url.URL(v)
					targetMap[k] = u
				}
				targetField.Set(reflect.ValueOf(targetMap))
			case c.types.urlMapType:
				sourceMap := sourceField.Interface().(map[string]url.URL)
				targetMap := make(map[string]ParseableURL, len(sourceMap))
				for k, v := range sourceMap {
					u := ParseableURL(v)
					targetMap[k] = u
				}
				targetField.Set(reflect.ValueOf(targetMap))
			case c.types.parseableURLMapPtrType:
				sourceMap := sourceField.Interface().(map[string]*ParseableURL)
				targetMap := make(map[string]*url.URL, len(sourceMap))
				for k, v := range sourceMap {
					u := url.URL(*v)
					targetMap[k] = &u
				}
				targetField.Set(reflect.ValueOf(targetMap))
			case c.types.urlMapPtrType:
				sourceMap := sourceField.Interface().(map[string]*url.URL)
				targetMap := make(map[string]*ParseableURL, len(sourceMap))
				for k, v := range sourceMap {
					u := ParseableURL(*v)
					targetMap[k] = &u
				}
				targetField.Set(reflect.ValueOf(targetMap))
			case c.types.wrappedBoolType:
				sourceValue := sourceField.Interface().(Bool)
				if sourceValue.Value != nil {
					switch targetFieldType {
					case c.types.boolType:
						targetValue := *sourceValue.Value
						targetField.Set(reflect.ValueOf(targetValue))
					case c.types.boolPtrType:
						targetValue := sourceValue.Value
						targetField.Set(reflect.ValueOf(targetValue))
					}
				}
			case c.types.wrappedBoolSliceType:
				sourceSlice := sourceField.Interface().([]Bool)
				switch targetFieldType {
				case c.types.boolSliceType:
					targetSlice := make([]bool, len(sourceSlice))
					for k := range sourceSlice {
						v := sourceSlice[k]
						if v.Value != nil {
							targetSlice[k] = *v.Value
						}
					}
					targetField.Set(reflect.ValueOf(targetSlice))
				case c.types.boolPtrSliceType:
					targetSlice := make([]*bool, len(sourceSlice))
					for k := range sourceSlice {
						v := sourceSlice[k]
						if v.Value != nil {
							targetSlice[k] = v.Value
						}
					}
					targetField.Set(reflect.ValueOf(targetSlice))
				}
			default:
				targetField.Set(sourceField)
			}
		}
	}
}

func (c *MagicConfig[T]) isZero(v reflect.Value, tag string) bool {
	if tag == _zeroPointer {
		return true
	}
	kind := v.Kind()
	switch kind {
	case reflect.Func, reflect.Struct, reflect.Ptr:
		if kind == reflect.Ptr && v.IsNil() {
			return true
		}
		vType := v.Type()
		switch vType {
		case c.types.durationType, c.types.parseableDurationType:
			return v.Interface() == 0
		case c.types.timeType:
			return v.Interface().(time.Time).IsZero()
		case c.types.parseableTimeType:
			return time.Time(v.Interface().(ParseableTime)).IsZero()
		case c.types.urlType:
			u := v.Interface().(url.URL)
			return u.String() == ""
		case c.types.parseableURLType:
			u := url.URL(v.Interface().(ParseableURL))
			return u.String() == ""
		case c.types.durationPtrType, c.types.parseableDurationPtrType, c.types.timePtrType,
			c.types.parseableTimePtrType, c.types.urlPtrType, c.types.parseableURLPtrType:
			return v.IsZero()
		case c.types.wrappedBoolType:
			b := v.Interface().(Bool)
			return b.Value == nil
		default:
			for _, t := range c.customTypes {
				//if kind == reflect.Ptr {
				//	vType = vType.Elem()
				//}
				if vType == t.basicType {
					return v.IsZero()
				}
			}
		}
		return v.IsNil()
	//case c.types.wrappedBoolKind:
	//	return false
	case reflect.Map, reflect.Slice, c.types.parseableTimeKind:
		return v.IsNil() || v.Len() == 0
	default:
		zero := reflect.Zero(v.Type())
		return v.Interface() == zero.Interface()
	}
}

func isCapital(v rune) bool {
	return v >= 'A' && v <= 'Z'
}

func (c *MagicConfig[T]) setDefaults(
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
		if c.recursiveStruct(fld) {
			c.setDefaults(fld, zeroPointer, listSeparator, keyValueSeparator, timeFormats)
			continue
		}
		tag := field.Tag.Get("def")

		fType := fld.Type()
		kind := fType.Kind()

		if !c.isZero(fld, tag) &&
			(kind != reflect.Map || fld.Len() > 0) {
			continue
		}
		if tag != "" {
			if kind == reflect.Map {
				fld.Set(reflect.MakeMap(fld.Type()))
			} else if kind == reflect.Slice {
				fld.Set(reflect.MakeSlice(fld.Type(), 0, 0))
			} else {
				switch fType {
				case c.types.parseableDurationType:
					if d, err := time.ParseDuration(tag); err == nil {
						x := ParseableDuration(d)
						fld.Set(reflect.ValueOf(x))
					}
					continue
				case c.types.parseableDurationPtrType:
					if tag == zeroPointer {
						tag = "0s"
					}
					if d, err := time.ParseDuration(tag); err == nil {
						x := ParseableDuration(d)
						fld.Set(reflect.ValueOf(&x))
					}
					continue
				case c.types.parseableTimeType:
					for _, timeFormat := range timeFormats {
						if t, err := time.Parse(timeFormat, tag); err == nil {
							x := ParseableTime(t)
							fld.Set(reflect.ValueOf(x))
							break
						}
					}
					continue
				case c.types.parseableTimePtrType:
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
				case c.types.parseableURLType:
					if u, err := url.Parse(tag); err == nil {
						x := ParseableURL(*u)
						fld.Set(reflect.ValueOf(x))
					}
					continue
				case c.types.parseableURLPtrType:
					if tag == zeroPointer {
						tag = ""
					}
					if u, err := url.Parse(tag); err == nil {
						x := ParseableURL(*u)
						fld.Set(reflect.ValueOf(&x))
					}
					continue
				case c.types.complex128Type:
					if c, err := strconv.ParseComplex(tag, 128); err == nil {
						x := ParseableComplex128(c)
						fld.Set(reflect.ValueOf(x))
					}
					continue
				case c.types.complex128PtrType:
					if tag == zeroPointer {
						tag = "0"
					}
					if c, err := strconv.ParseComplex(tag, 128); err == nil {
						x := ParseableComplex128(c)
						fld.Set(reflect.ValueOf(&x))
					}
					continue
				case c.types.complex64Type:
					if c, err := strconv.ParseComplex(tag, 64); err == nil {
						x := ParseableComplex64(c)
						fld.Set(reflect.ValueOf(x))
					}
					continue
				case c.types.complex64PtrType:
					if tag == zeroPointer {
						tag = "0"
					}
					if c, err := strconv.ParseComplex(tag, 64); err == nil {
						x := ParseableComplex64(c)
						fld.Set(reflect.ValueOf(&x))
					}
					continue
				default:
					attempt := false
					for _, t := range c.customTypes {
						uType := fType
						if kind == reflect.Ptr {
							uType = uType.Elem()
						}
						if uType == t.basicType {
							attempt = true
							if t.defaultParser != nil {
								result := t.defaultParser(tag)
								if r, ok := result.(error); ok {
									println(r.Error())
									break
								}
								if kind == reflect.Ptr {
									p := reflect.New(reflect.TypeOf(result))
									p.Elem().Set(reflect.ValueOf(result))
									fld.Set(p)
								} else {
									fld.Set(reflect.ValueOf(result))
								}
							}
						}
					}
					if attempt {
						continue
					}
				}

				if fld.Kind() == reflect.Ptr {
					fld.Set(reflect.New(field.Type.Elem()))
				} else {
					fld.Set(reflect.New(field.Type).Elem())
				}
			}
			c.setFieldValue(
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

func (c *MagicConfig[T]) setFieldValue(
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

	fType := f.Type()
	kind := fType.Kind()

	if fType == c.types.wrappedBoolType {
		if val == zeroPointer || val == "" {
			fls := false
			f.Set(reflect.ValueOf(Bool{Value: &fls}))
		} else if v, err := strconv.ParseBool(val); err == nil {
			f.Set(reflect.ValueOf(Bool{Value: &v}))
		}
		return
	}

	switch kind {
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

		attempt := false
		for _, t := range c.customTypes {
			uType := elemType
			if elemType.Kind() == reflect.Ptr {
				uType = uType.Elem()
			}
			if uType == t.basicType {
				attempt = true
				if elemType.Kind() == reflect.Ptr && t.defaultMapPointerParser != nil {
					result := t.defaultMapPointerParser(val)
					if r, ok := result.(error); ok {
						println(r.Error())
						break
					}
					f.Set(reflect.ValueOf(result))
				} else if t.defaultMapParser != nil {
					result := t.defaultMapParser(val)
					if r, ok := result.(error); ok {
						println(r.Error())
						break
					}
					f.Set(reflect.ValueOf(result))
				}
			}
		}
		if attempt {
			return
		}

		switch elemType {
		case c.types.parseableTimeType, c.types.parseableTimePtrType,
			c.types.parseableURLType, c.types.parseableURLPtrType:
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
		m := reflect.MakeMap(f.Type())
		pointer := false
		elemKind := elemType.Kind()
		if elemKind == reflect.Ptr {
			pointer = true
			elemKind = elemType.Elem().Kind()
		}
		switch elemType {
		case c.types.parseableDurationType:
			for k, v := range vals {
				if d, err := time.ParseDuration(v); err == nil {
					pd := ParseableDuration(d)
					m.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(pd))
				}
			}
			f.Set(m)
			return
		case c.types.parseableDurationPtrType:
			for k, v := range vals {
				if d, err := time.ParseDuration(v); err == nil {
					pd := ParseableDuration(d)
					m.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(&pd))
				}
			}
			f.Set(m)
			return
		case c.types.parseableTimeType:
			for k, v := range vals {
				for _, timeFormat := range timeFormats {
					if t, err := time.Parse(timeFormat, v); err == nil {
						pt := ParseableTime(t)
						m.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(pt))
						break
					}
				}
			}
			f.Set(m)
			return
		case c.types.parseableTimePtrType:
			for k, v := range vals {
				for _, timeFormat := range timeFormats {
					if t, err := time.Parse(timeFormat, v); err == nil {
						pt := ParseableTime(t)
						m.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(&pt))
						break
					}
				}
			}
			f.Set(m)
			return
		case c.types.parseableURLType:
			for k, v := range vals {
				if u, err := url.Parse(v); err == nil {
					p := ParseableURL(*u)
					m.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(p))
				}
			}
			f.Set(m)
			return
		case c.types.parseableURLPtrType:
			for k, v := range vals {
				if u, err := url.Parse(v); err == nil {
					p := ParseableURL(*u)
					m.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(&p))
				}
			}
			f.Set(m)
			return
		}
		for k, v := range vals {
			switch elemKind {
			case reflect.String:
				if pointer {
					m.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(&v).Convert(f.Type().Elem()))
				} else {
					m.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v).Convert(f.Type().Elem()))
				}
			case reflect.Bool:
				b, err := strconv.ParseBool(v)
				if err != nil {
					continue
				}
				if pointer {
					m.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(&b).Convert(elemType))
				} else {
					m.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(b).Convert(elemType))
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
						m.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(&value))
					case reflect.Int8:
						value := int8(i)
						m.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(&value))
					case reflect.Int16:
						value := int16(i)
						m.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(&value))
					case reflect.Int32:
						value := int32(i)
						m.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(&value))
					case reflect.Uint:
						value := uint(i)
						m.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(&value))
					case reflect.Uint8:
						value := uint8(i)
						m.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(&value))
					case reflect.Uint16:
						value := uint16(i)
						m.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(&value))
					case reflect.Uint32:
						value := uint32(i)
						m.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(&value))
					case reflect.Uint64:
						value := uint64(i)
						m.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(&value))
					default:
						m.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(&i))
					}
				} else {
					m.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(i).Convert(elemType))
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
						m.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(&value).Convert(elemType))
					default:
						m.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(&i).Convert(elemType))
					}
				} else {
					m.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(i).Convert(elemType))
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
						m.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(&value))
					default:
						m.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(&v))
					}
				} else {
					m.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v).Convert(elemType))
				}
			default:
				continue
			}
		}
		f.Set(m)
	case reflect.Slice:
		var vals []string
		if val == "" || val == zeroPointer {
			return
		}

		vals = strings.Split(val, listSeparator)

		s := reflect.MakeSlice(f.Type(), 0, len(vals))
		elemType := reflect.TypeOf(f.Interface()).Elem()
		pointer := false
		elemKind := elemType.Kind()
		if elemKind == reflect.Ptr {
			pointer = true
			elemKind = elemType.Elem().Kind()
		}

		attempt := false
		for _, t := range c.customTypes {
			uType := elemType
			if pointer {
				uType = uType.Elem()
			}
			if uType == t.basicType {
				attempt = true
				if elemType.Kind() == reflect.Ptr && t.defaultListPointerParser != nil {
					result := t.defaultListPointerParser(val)
					if r, ok := result.(error); ok {
						println(r.Error())
						break
					}
					f.Set(reflect.ValueOf(result))
				} else if t.defaultListParser != nil {
					result := t.defaultListParser(val)
					if r, ok := result.(error); ok {
						println(r.Error())
						break
					}
					f.Set(reflect.ValueOf(result))
				}
			}
		}
		if attempt {
			return
		}
		switch elemType {
		case c.types.wrappedBoolType:
			for _, v := range vals {
				if value, err := strconv.ParseBool(v); err == nil {
					b := Bool{Value: &value}
					s = reflect.Append(s, reflect.ValueOf(b))
				}
			}
			f.Set(s)
			return
		case c.types.parseableDurationType:
			for _, dur := range vals {
				if d, err := time.ParseDuration(dur); err == nil {
					pd := ParseableDuration(d)
					s = reflect.Append(s, reflect.ValueOf(pd))
				}
			}
			f.Set(s)
			return
		case c.types.parseableDurationPtrType:
			for _, dur := range vals {
				if d, err := time.ParseDuration(dur); err == nil {
					pd := ParseableDuration(d)
					s = reflect.Append(s, reflect.ValueOf(&pd))
				}
			}
			f.Set(s)
			return
		case c.types.parseableTimeType:
			for _, tm := range vals {
				for _, timeFormat := range timeFormats {
					if t, err := time.Parse(timeFormat, tm); err == nil {
						pt := ParseableTime(t)
						s = reflect.Append(s, reflect.ValueOf(pt))
						break
					}
				}
			}
			f.Set(s)
			return
		case c.types.parseableTimePtrType:
			for _, tm := range vals {
				for _, timeFormat := range timeFormats {
					if t, err := time.Parse(timeFormat, tm); err == nil {
						pt := ParseableTime(t)
						s = reflect.Append(s, reflect.ValueOf(&pt))
						break
					}
				}
			}
			f.Set(s)
			return
		case c.types.parseableURLType:
			for _, v := range vals {
				if u, err := url.Parse(v); err == nil {
					p := ParseableURL(*u)
					s = reflect.Append(s, reflect.ValueOf(p))
				}
			}
			f.Set(s)
			return
		case c.types.parseableURLPtrType:
			for _, v := range vals {
				if u, err := url.Parse(v); err == nil {
					p := ParseableURL(*u)
					s = reflect.Append(s, reflect.ValueOf(&p))
				}
			}
			f.Set(s)
			return
		}

		for _, v := range vals {
			switch elemKind {
			case reflect.String:
				if pointer {
					s = reflect.Append(s, reflect.ValueOf(&v).Convert(f.Type().Elem()))
				} else {
					s = reflect.Append(s, reflect.ValueOf(v).Convert(f.Type().Elem()))
				}
			case reflect.Bool:
				b, err := strconv.ParseBool(v)
				if err != nil {
					continue
				}
				if pointer {
					s = reflect.Append(s, reflect.ValueOf(&b).Convert(elemType))
				} else {
					s = reflect.Append(s, reflect.ValueOf(b).Convert(elemType))
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
						s = reflect.Append(s, reflect.ValueOf(&value))
					case reflect.Int8:
						value := int8(i)
						s = reflect.Append(s, reflect.ValueOf(&value))
					case reflect.Int16:
						value := int16(i)
						s = reflect.Append(s, reflect.ValueOf(&value))
					case reflect.Int32:
						value := int32(i)
						s = reflect.Append(s, reflect.ValueOf(&value))
					case reflect.Uint:
						value := uint(i)
						s = reflect.Append(s, reflect.ValueOf(&value))
					case reflect.Uint8:
						value := uint8(i)
						s = reflect.Append(s, reflect.ValueOf(&value))
					case reflect.Uint16:
						value := uint16(i)
						s = reflect.Append(s, reflect.ValueOf(&value))
					case reflect.Uint32:
						value := uint32(i)
						s = reflect.Append(s, reflect.ValueOf(&value))
					case reflect.Uint64:
						value := uint64(i)
						s = reflect.Append(s, reflect.ValueOf(&value))
					default:
						s = reflect.Append(s, reflect.ValueOf(&i))
					}
				} else {
					s = reflect.Append(s, reflect.ValueOf(i).Convert(elemType))
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
						s = reflect.Append(s, reflect.ValueOf(&value))
					default:
						s = reflect.Append(s, reflect.ValueOf(&v))
					}
				} else {
					s = reflect.Append(s, reflect.ValueOf(v).Convert(elemType))
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
						s = reflect.Append(s, reflect.ValueOf(&value).Convert(elemType))
					default:
						s = reflect.Append(s, reflect.ValueOf(&i).Convert(elemType))
					}
				} else {
					s = reflect.Append(s, reflect.ValueOf(i).Convert(elemType))
				}
			default:
				continue
			}
		}
		f.Set(s)
	default:
		return
	}
}

func (c *MagicConfig[T]) transformValues(v reflect.Value, customFuncs []func(reflect.StructField, reflect.Value)) {
	if !c.recursiveStruct(v) {
		return
	}
	vVal := reflect.Indirect(v)
	for i := range vVal.NumField() {
		field := vVal.Type().Field(i)
		fld := vVal.Field(i)
		if c.recursiveStruct(fld) {
			c.transformValues(fld, customFuncs)
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

func (c *MagicConfig[T]) validateEnums(v reflect.Value) []error {
	var errs []error
	if !c.recursiveStruct(v) {
		return errs
	}
	vVal := v
	if vVal.Kind() == reflect.Ptr {
		vVal = reflect.Indirect(v)
	}
	for i := range vVal.NumField() {
		fld := vVal.Field(i)
		if c.recursiveStruct(fld) {
			if subErrs := c.validateEnums(fld); len(subErrs) > 0 {
				errs = append(errs, subErrs...)
			}
		}

		if c.isZero(fld, "") {
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

func (c *MagicConfig[T]) splitArgs(v reflect.Value, args []string, listSeparator string) []string {
	if !c.recursiveStruct(v) {
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

		fixedArg, rawValues, skip = c.getValues(v, arg, args, i, listSeparator)

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

func (c *MagicConfig[T]) getValues(v reflect.Value, arg string, args []string, i int, listSeparator string) (string, []string, bool) {
	var values []string
	fixedArg := arg

	if strings.Contains(arg, "=") {
		vals := strings.Split(arg, "=")
		if len(vals) == 2 {
			fixedArg = vals[0]
			values = c.extractValues(v, fixedArg, vals[1], listSeparator)
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

	values = c.extractValues(v, fixedArg, nextArg, listSeparator)
	return fixedArg, values, true
}

func (c *MagicConfig[T]) extractValues(v reflect.Value, arg string, valArg string, listSeparator string) []string {
	var values []string

	if !strings.Contains(valArg, listSeparator) || !c.verifySliceOrMapType(v, arg) {
		values = append(values, valArg)
		return values
	}

	return strings.Split(valArg, listSeparator)
}

func (c *MagicConfig[T]) verifySliceOrMapType(v reflect.Value, arg string) bool {
	if !c.recursiveStruct(v) {
		return false
	}

	bareArg := strings.TrimPrefix(strings.TrimPrefix(arg, "-"), "-")

	vVal := reflect.Indirect(v)
	for i := range vVal.NumField() {
		field := vVal.Type().Field(i)
		fld := vVal.Field(i)
		if c.recursiveStruct(fld) {
			if c.verifySliceOrMapType(fld, arg) {
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
