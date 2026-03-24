package magiccfg

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/goccy/go-yaml"
)

var packageTimeFormats []string

type types struct {
	stringType                      reflect.Type
	boolType                        reflect.Type
	boolPtrType                     reflect.Type
	boolSliceType                   reflect.Type
	boolPtrSliceType                reflect.Type
	wrappedBoolType                 reflect.Type
	wrappedBoolSliceType            reflect.Type
	intType                         reflect.Type
	int8Type                        reflect.Type
	int16Type                       reflect.Type
	int32Type                       reflect.Type
	int64Type                       reflect.Type
	uintType                        reflect.Type
	uint8Type                       reflect.Type
	uint16Type                      reflect.Type
	uint32Type                      reflect.Type
	uint64Type                      reflect.Type
	float32Type                     reflect.Type
	float64Type                     reflect.Type
	stringPtrMapType                reflect.Type
	intPtrMapType                   reflect.Type
	int8PtrMapType                  reflect.Type
	int16PtrMapType                 reflect.Type
	int32PtrMapType                 reflect.Type
	int64PtrMapType                 reflect.Type
	uintPtrMapType                  reflect.Type
	uint8PtrMapType                 reflect.Type
	uint16PtrMapType                reflect.Type
	uint32PtrMapType                reflect.Type
	uint64PtrMapType                reflect.Type
	float32PtrMapType               reflect.Type
	float64PtrMapType               reflect.Type
	durationType                    reflect.Type
	durationSliceType               reflect.Type
	durationMapType                 reflect.Type
	durationPtrType                 reflect.Type
	durationSlicePtrType            reflect.Type
	durationMapPtrType              reflect.Type
	parseableDurationType           reflect.Type
	parseableDurationSliceType      reflect.Type
	parseableDurationMapType        reflect.Type
	parseableDurationPtrType        reflect.Type
	parseableDurationSlicePtrType   reflect.Type
	parseableDurationMapPtrType     reflect.Type
	timeType                        reflect.Type
	timeSliceType                   reflect.Type
	timeMapType                     reflect.Type
	timePtrType                     reflect.Type
	timeSlicePtrType                reflect.Type
	timeMapPtrType                  reflect.Type
	parseableTimeType               reflect.Type
	parseableTimeSliceType          reflect.Type
	parseableTimeMapType            reflect.Type
	parseableTimePtrType            reflect.Type
	parseableTimeSlicePtrType       reflect.Type
	parseableTimeMapPtrType         reflect.Type
	urlType                         reflect.Type
	urlSliceType                    reflect.Type
	urlMapType                      reflect.Type
	urlPtrType                      reflect.Type
	urlSlicePtrType                 reflect.Type
	urlMapPtrType                   reflect.Type
	parseableURLType                reflect.Type
	parseableURLSliceType           reflect.Type
	parseableURLMapType             reflect.Type
	parseableURLPtrType             reflect.Type
	parseableURLSlicePtrType        reflect.Type
	parseableURLMapPtrType          reflect.Type
	parseableTimeKind               reflect.Kind
	complex128Type                  reflect.Type
	complex128PtrType               reflect.Type
	complex128SliceType             reflect.Type
	complex128SlicePtrType          reflect.Type
	complex128MapType               reflect.Type
	complex128MapPtrType            reflect.Type
	complex64Type                   reflect.Type
	complex64PtrType                reflect.Type
	complex64SliceType              reflect.Type
	complex64SlicePtrType           reflect.Type
	complex64MapType                reflect.Type
	complex64MapPtrType             reflect.Type
	boolPtrMapType                  reflect.Type
	parseableComplex128Type         reflect.Type
	parseableComplex128PtrType      reflect.Type
	parseableComplex128SliceType    reflect.Type
	parseableComplex128SlicePtrType reflect.Type
	parseableComplex128MapType      reflect.Type
	parseableComplex128MapPtrType   reflect.Type
	parseableComplex64Type          reflect.Type
	parseableComplex64PtrType       reflect.Type
	parseableComplex64SliceType     reflect.Type
	parseableComplex64SlicePtrType  reflect.Type
	parseableComplex64MapType       reflect.Type
	parseableComplex64MapPtrType    reflect.Type
}

func newTypes() *types {
	t := true

	return &types{
		stringType:                      reflect.TypeOf(""),
		boolType:                        reflect.TypeOf(true),
		boolPtrType:                     reflect.TypeOf(&t),
		boolSliceType:                   reflect.TypeOf([]bool{}),
		boolPtrSliceType:                reflect.TypeOf([]*bool{}),
		boolPtrMapType:                  reflect.TypeOf(map[string]*bool{}),
		wrappedBoolType:                 reflect.TypeOf(Bool{}),
		wrappedBoolSliceType:            reflect.TypeOf([]Bool{}),
		intType:                         reflect.TypeOf(0),
		int8Type:                        reflect.TypeOf(int8(0)),
		int16Type:                       reflect.TypeOf(int16(0)),
		int32Type:                       reflect.TypeOf(int32(0)),
		int64Type:                       reflect.TypeOf(int64(0)),
		uintType:                        reflect.TypeOf(uint(0)),
		uint8Type:                       reflect.TypeOf(uint8(0)),
		uint16Type:                      reflect.TypeOf(uint16(0)),
		uint32Type:                      reflect.TypeOf(uint32(0)),
		uint64Type:                      reflect.TypeOf(uint64(0)),
		float32Type:                     reflect.TypeOf(float32(0)),
		float64Type:                     reflect.TypeOf(float64(0)),
		stringPtrMapType:                reflect.TypeOf(map[string]*string{}),
		intPtrMapType:                   reflect.TypeOf(map[string]*int{}),
		int8PtrMapType:                  reflect.TypeOf(map[string]*int8{}),
		int16PtrMapType:                 reflect.TypeOf(map[string]*int16{}),
		int32PtrMapType:                 reflect.TypeOf(map[string]*int32{}),
		int64PtrMapType:                 reflect.TypeOf(map[string]*int64{}),
		uintPtrMapType:                  reflect.TypeOf(map[string]*uint{}),
		uint8PtrMapType:                 reflect.TypeOf(map[string]*uint8{}),
		uint16PtrMapType:                reflect.TypeOf(map[string]*uint16{}),
		uint32PtrMapType:                reflect.TypeOf(map[string]*uint32{}),
		uint64PtrMapType:                reflect.TypeOf(map[string]*uint64{}),
		float32PtrMapType:               reflect.TypeOf(map[string]*float32{}),
		float64PtrMapType:               reflect.TypeOf(map[string]*float64{}),
		durationType:                    reflect.TypeOf(time.Duration(0)),
		durationSliceType:               reflect.TypeOf([]time.Duration{}),
		durationMapType:                 reflect.TypeOf(map[string]time.Duration{}),
		durationPtrType:                 reflect.TypeOf(new(time.Duration)),
		durationSlicePtrType:            reflect.TypeOf([]*time.Duration{}),
		durationMapPtrType:              reflect.TypeOf(map[string]*time.Duration{}),
		parseableDurationType:           reflect.TypeOf(ParseableDuration(0)),
		parseableDurationSliceType:      reflect.TypeOf([]ParseableDuration{}),
		parseableDurationMapType:        reflect.TypeOf(map[string]ParseableDuration{}),
		parseableDurationPtrType:        reflect.TypeOf(new(ParseableDuration)),
		parseableDurationSlicePtrType:   reflect.TypeOf([]*ParseableDuration{}),
		parseableDurationMapPtrType:     reflect.TypeOf(map[string]*ParseableDuration{}),
		timeType:                        reflect.TypeOf(time.Now()),
		timeSliceType:                   reflect.TypeOf([]time.Time{}),
		timeMapType:                     reflect.TypeOf(map[string]time.Time{}),
		timePtrType:                     reflect.TypeOf(new(time.Time)),
		timeSlicePtrType:                reflect.TypeOf([]*time.Time{}),
		timeMapPtrType:                  reflect.TypeOf(map[string]*time.Time{}),
		parseableTimeType:               reflect.TypeOf(ParseableTime(time.Now())),
		parseableTimeSliceType:          reflect.TypeOf([]ParseableTime{}),
		parseableTimeMapType:            reflect.TypeOf(map[string]ParseableTime{}),
		parseableTimePtrType:            reflect.TypeOf(new(ParseableTime)),
		parseableTimeSlicePtrType:       reflect.TypeOf([]*ParseableTime{}),
		parseableTimeMapPtrType:         reflect.TypeOf(map[string]*ParseableTime{}),
		urlType:                         reflect.TypeOf(url.URL{}),
		urlSliceType:                    reflect.TypeOf([]url.URL{}),
		urlMapType:                      reflect.TypeOf(map[string]url.URL{}),
		urlPtrType:                      reflect.TypeOf(new(url.URL)),
		urlSlicePtrType:                 reflect.TypeOf([]*url.URL{}),
		urlMapPtrType:                   reflect.TypeOf(map[string]*url.URL{}),
		parseableURLType:                reflect.TypeOf(ParseableURL{}),
		parseableURLSliceType:           reflect.TypeOf([]ParseableURL{}),
		parseableURLMapType:             reflect.TypeOf(map[string]ParseableURL{}),
		parseableURLPtrType:             reflect.TypeOf(new(ParseableURL)),
		parseableURLSlicePtrType:        reflect.TypeOf([]*ParseableURL{}),
		parseableURLMapPtrType:          reflect.TypeOf(map[string]*ParseableURL{}),
		parseableTimeKind:               reflect.TypeOf(time.Time{}).Kind(),
		complex128Type:                  reflect.TypeOf(complex128(0)),
		complex128PtrType:               reflect.TypeOf(new(complex128)),
		complex128SliceType:             reflect.TypeOf([]complex128{}),
		complex128SlicePtrType:          reflect.TypeOf([]*complex128{}),
		complex128MapType:               reflect.TypeOf(map[string]complex128{}),
		complex128MapPtrType:            reflect.TypeOf(map[string]*complex128{}),
		complex64Type:                   reflect.TypeOf(complex64(0)),
		complex64PtrType:                reflect.TypeOf(new(complex64)),
		complex64SliceType:              reflect.TypeOf([]complex64{}),
		complex64SlicePtrType:           reflect.TypeOf([]*complex64{}),
		complex64MapType:                reflect.TypeOf(map[string]complex64{}),
		complex64MapPtrType:             reflect.TypeOf(map[string]*complex64{}),
		parseableComplex128Type:         reflect.TypeOf(ParseableComplex128(0)),
		parseableComplex128PtrType:      reflect.TypeOf(new(ParseableComplex128)),
		parseableComplex128SliceType:    reflect.TypeOf([]ParseableComplex128{}),
		parseableComplex128SlicePtrType: reflect.TypeOf([]*ParseableComplex128{}),
		parseableComplex128MapType:      reflect.TypeOf(map[string]ParseableComplex128{}),
		parseableComplex128MapPtrType:   reflect.TypeOf(map[string]*ParseableComplex128{}),
		parseableComplex64Type:          reflect.TypeOf(ParseableComplex64(0)),
		parseableComplex64PtrType:       reflect.TypeOf(new(ParseableComplex64)),
		parseableComplex64SliceType:     reflect.TypeOf([]ParseableComplex64{}),
		parseableComplex64SlicePtrType:  reflect.TypeOf([]*ParseableComplex64{}),
		parseableComplex64MapType:       reflect.TypeOf(map[string]ParseableComplex64{}),
		parseableComplex64MapPtrType:    reflect.TypeOf(map[string]*ParseableComplex64{}),
	}
}

// ParseableDuration enables time.Duration to be parsed by the encoding/xml package
type ParseableDuration time.Duration

func (p *ParseableDuration) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}
	dur, err := time.ParseDuration(strings.TrimSpace(s))
	if err != nil {
		return err
	}
	*p = ParseableDuration(dur)
	return nil
}

func (p *ParseableDuration) UnmarshalFlag(value string) error {
	dur, err := time.ParseDuration(strings.TrimSpace(value))
	if err != nil {
		return err
	}
	*p = ParseableDuration(dur)
	return nil
}

func (p *ParseableDuration) UnmarshalTOML(value any) error {
	v, ok := value.(string)
	if !ok {
		return fmt.Errorf("%s is not a string", value)
	}
	dur, err := time.ParseDuration(strings.TrimSpace(v))
	if err != nil {
		return err
	}
	*p = ParseableDuration(dur)
	return nil
}

func (p *ParseableDuration) UnmarshalYAML(value []byte) error {
	dur, err := time.ParseDuration(strings.TrimSpace(strings.ReplaceAll(string(value), "\"", "")))
	if err != nil {
		return err
	}
	*p = ParseableDuration(dur)
	return nil
}

func UnmarshalDurationEnv(v string) (any, error) {
	dur, err := time.ParseDuration(v)
	if err != nil {
		return "", err
	}
	return ParseableDuration(dur), nil
}

func UnmarshalTimeEnv(v string) (any, error) {
	for _, format := range packageTimeFormats {
		if t, err := time.Parse(format, v); err == nil {
			pt := ParseableTime(t)
			return pt, nil
		}
	}
	return nil, fmt.Errorf("no time format found that could parse %s", v)
}

// ParseableTime enables time.Time to be parsed by the encoding/xml package
type ParseableTime time.Time

func (p *ParseableTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}
	v := strings.ReplaceAll(strings.TrimSpace(s), "\"", "")
	if packageTimeFormats != nil {
		for _, f := range packageTimeFormats {
			t, err := time.Parse(f, v)
			if err == nil {
				*p = ParseableTime(t)
				return nil
			}
		}
	}
	return fmt.Errorf("could not parse time: %s", s)
}

func (p *ParseableTime) UnmarshalFlag(s string) error {
	if packageTimeFormats != nil {
		for _, v := range packageTimeFormats {
			t, err := time.Parse(v, s)
			if err == nil {
				*p = ParseableTime(t)
				return nil
			}
		}
	}
	return fmt.Errorf("could not parse time: %s", s)
}

func (p *ParseableTime) UnmarshalYAML(value []byte) error {
	v := strings.ReplaceAll(strings.TrimSpace(string(value)), "\"", "")
	for _, format := range packageTimeFormats {
		if t, err := time.Parse(format, v); err == nil {
			*p = ParseableTime(t)
			return nil
		}
	}
	return fmt.Errorf("could not parse time: %s", value)
}

func (p *ParseableTime) UnmarshalTOML(value any) error {
	v, ok := value.(string)
	if !ok {
		return fmt.Errorf("%s is not a string", value)
	}
	v = strings.ReplaceAll(strings.TrimSpace(v), "\"", "")
	for _, format := range packageTimeFormats {
		if t, err := time.Parse(format, v); err == nil {
			*p = ParseableTime(t)
			return nil
		}
	}
	return fmt.Errorf("could not parse time: %s", value)
}

type ParseableURL url.URL

func (p *ParseableURL) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}
	u, err := url.Parse(strings.TrimSpace(s))
	if err != nil {
		return err
	}
	*p = ParseableURL(*u)
	return nil
}

func (p *ParseableURL) UnmarshalText(val []byte) error {
	u, err := url.Parse(strings.TrimSpace(string(val)))
	if err != nil {
		return err
	}
	*p = ParseableURL(*u)
	return nil
}

func (p *ParseableURL) UnmarshalFlag(s string) error {
	u, err := url.Parse(strings.TrimSpace(s))
	if err != nil {
		return err
	}
	*p = ParseableURL(*u)
	return nil
}

type ParseableComplex128 complex128

func (n *ParseableComplex128) UnmarshalYAML(data []byte) error {
	var s string
	if err := yaml.Unmarshal(data, &s); err != nil {
		return err
	}

	c, err := strconv.ParseComplex(strings.TrimSpace(s), 128)
	if err != nil {
		return err
	}

	*n = ParseableComplex128(c)

	return nil
}

func (n *ParseableComplex128) UnmarshalFlag(s string) error {
	s = strings.TrimPrefix(s, "+")
	c, err := strconv.ParseComplex(strings.TrimSpace(s), 128)
	if err != nil {
		return err
	}

	*n = ParseableComplex128(c)

	return nil
}

func (n *ParseableComplex128) UnmarshalTOML(value any) error {
	v, ok := value.(string)
	if !ok {
		return fmt.Errorf("%s is not a string", value)
	}
	c, err := strconv.ParseComplex(strings.TrimSpace(v), 128)
	if err != nil {
		return err
	}

	*n = ParseableComplex128(c)

	return nil
}

func (n *ParseableComplex128) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}
	c, err := strconv.ParseComplex(strings.TrimSpace(s), 128)
	if err != nil {
		return err
	}
	*n = ParseableComplex128(c)
	return nil
}

type ParseableComplex64 complex64

func (n *ParseableComplex64) UnmarshalYAML(data []byte) error {
	var s string
	if err := yaml.Unmarshal(data, &s); err != nil {
		return err
	}

	c, err := strconv.ParseComplex(strings.TrimSpace(s), 64)
	if err != nil {
		return err
	}

	*n = ParseableComplex64(c)

	return nil
}

func (n *ParseableComplex64) UnmarshalFlag(s string) error {
	s = strings.TrimPrefix(s, "+")
	c, err := strconv.ParseComplex(strings.TrimSpace(s), 128)
	if err != nil {
		return err
	}

	*n = ParseableComplex64(c)

	return nil
}

func (n *ParseableComplex64) UnmarshalTOML(value any) error {
	v, ok := value.(string)
	if !ok {
		return fmt.Errorf("%s is not a string", value)
	}
	c, err := strconv.ParseComplex(strings.TrimSpace(v), 64)
	if err != nil {
		return err
	}

	r := complex64(c)

	*n = ParseableComplex64(r)

	return nil
}

func (n *ParseableComplex64) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}
	c, err := strconv.ParseComplex(strings.TrimSpace(s), 64)
	if err != nil {
		return err
	}
	r := complex64(c)
	*n = ParseableComplex64(r)
	return nil
}

type Bool struct {
	Value *bool
}

func (b *Bool) UnmarshalFlag(s string) error {
	v, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}
	*b = Bool{
		Value: &v,
	}
	return nil
}

func (b *Bool) UnmarshalYAML(data []byte) error {
	v, err := strconv.ParseBool(strings.TrimSpace(string(data)))
	if err != nil {
		return err
	}
	*b = Bool{
		Value: &v,
	}
	return nil
}

func (b *Bool) UnmarshalTOML(value any) error {
	v, ok := value.(bool)
	if !ok {
		return fmt.Errorf("%s is not a bool", value)
	}

	*b = Bool{Value: &v}

	return nil
}

func (b *Bool) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}
	c, err := strconv.ParseBool(strings.TrimSpace(s))
	if err != nil {
		return err
	}

	*b = Bool{Value: &c}
	return nil
}
