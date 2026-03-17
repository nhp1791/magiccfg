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

var (
	stringType                      = reflect.TypeOf("")
	boolType                        = reflect.TypeOf(true)
	wrappedBoolType                 = reflect.TypeOf(Bool(true))
	wrappedBoolKind                 = wrappedBoolType.Kind()
	boolSliceType                   = reflect.TypeOf([]bool{})
	boolPtrSliceType                = reflect.TypeOf([]*bool{})
	wrappedBoolSliceType            = reflect.TypeOf(BoolSlice{})
	wrappedBoolPtrSliceType         = reflect.TypeOf(BoolPtrSlice{})
	intType                         = reflect.TypeOf(0)
	int8Type                        = reflect.TypeOf(int8(0))
	int16Type                       = reflect.TypeOf(int16(0))
	int32Type                       = reflect.TypeOf(int32(0))
	int64Type                       = reflect.TypeOf(int64(0))
	uintType                        = reflect.TypeOf(uint(0))
	uint8Type                       = reflect.TypeOf(uint8(0))
	uint16Type                      = reflect.TypeOf(uint16(0))
	uint32Type                      = reflect.TypeOf(uint32(0))
	uint64Type                      = reflect.TypeOf(uint64(0))
	float32Type                     = reflect.TypeOf(float32(0))
	float64Type                     = reflect.TypeOf(float64(0))
	stringPtrMapType                = reflect.TypeOf(map[string]*string{})
	intPtrMapType                   = reflect.TypeOf(map[string]*int{})
	int8PtrMapType                  = reflect.TypeOf(map[string]*int8{})
	int16PtrMapType                 = reflect.TypeOf(map[string]*int16{})
	int32PtrMapType                 = reflect.TypeOf(map[string]*int32{})
	int64PtrMapType                 = reflect.TypeOf(map[string]*int64{})
	uintPtrMapType                  = reflect.TypeOf(map[string]*uint{})
	uint8PtrMapType                 = reflect.TypeOf(map[string]*uint8{})
	uint16PtrMapType                = reflect.TypeOf(map[string]*uint16{})
	uint32PtrMapType                = reflect.TypeOf(map[string]*uint32{})
	uint64PtrMapType                = reflect.TypeOf(map[string]*uint64{})
	float32PtrMapType               = reflect.TypeOf(map[string]*float32{})
	float64PtrMapType               = reflect.TypeOf(map[string]*float64{})
	durationType                    = reflect.TypeOf(time.Duration(0))
	durationSliceType               = reflect.TypeOf([]time.Duration{})
	durationMapType                 = reflect.TypeOf(map[string]time.Duration{})
	durationPtrType                 = reflect.TypeOf(new(time.Duration))
	durationSlicePtrType            = reflect.TypeOf([]*time.Duration{})
	durationMapPtrType              = reflect.TypeOf(map[string]*time.Duration{})
	parseableDurationType           = reflect.TypeOf(ParseableDuration(0))
	parseableDurationSliceType      = reflect.TypeOf([]ParseableDuration{})
	parseableDurationMapType        = reflect.TypeOf(map[string]ParseableDuration{})
	parseableDurationPtrType        = reflect.TypeOf(new(ParseableDuration))
	parseableDurationSlicePtrType   = reflect.TypeOf([]*ParseableDuration{})
	parseableDurationMapPtrType     = reflect.TypeOf(map[string]*ParseableDuration{})
	timeType                        = reflect.TypeOf(time.Now())
	timeSliceType                   = reflect.TypeOf([]time.Time{})
	timeMapType                     = reflect.TypeOf(map[string]time.Time{})
	timePtrType                     = reflect.TypeOf(new(time.Time))
	timeSlicePtrType                = reflect.TypeOf([]*time.Time{})
	timeMapPtrType                  = reflect.TypeOf(map[string]*time.Time{})
	parseableTimeType               = reflect.TypeOf(ParseableTime(time.Now()))
	parseableTimeSliceType          = reflect.TypeOf([]ParseableTime{})
	parseableTimeMapType            = reflect.TypeOf(map[string]ParseableTime{})
	parseableTimePtrType            = reflect.TypeOf(new(ParseableTime))
	parseableTimeSlicePtrType       = reflect.TypeOf([]*ParseableTime{})
	parseableTimeMapPtrType         = reflect.TypeOf(map[string]*ParseableTime{})
	urlType                         = reflect.TypeOf(url.URL{})
	urlSliceType                    = reflect.TypeOf([]url.URL{})
	urlMapType                      = reflect.TypeOf(map[string]url.URL{})
	urlPtrType                      = reflect.TypeOf(new(url.URL))
	urlSlicePtrType                 = reflect.TypeOf([]*url.URL{})
	urlMapPtrType                   = reflect.TypeOf(map[string]*url.URL{})
	parseableURLType                = reflect.TypeOf(ParseableURL{})
	parseableURLSliceType           = reflect.TypeOf([]ParseableURL{})
	parseableURLMapType             = reflect.TypeOf(map[string]ParseableURL{})
	parseableURLPtrType             = reflect.TypeOf(new(ParseableURL))
	parseableURLSlicePtrType        = reflect.TypeOf([]*ParseableURL{})
	parseableURLMapPtrType          = reflect.TypeOf(map[string]*ParseableURL{})
	packageTimeFormats              []string
	parseableTimeKind               = reflect.TypeOf(time.Time{}).Kind()
	complex128Type                  = reflect.TypeOf(complex128(0))
	complex128PtrType               = reflect.TypeOf(new(complex128))
	complex128SliceType             = reflect.TypeOf([]complex128{})
	complex128SlicePtrType          = reflect.TypeOf([]*complex128{})
	complex128MapType               = reflect.TypeOf(map[string]complex128{})
	complex128MapPtrType            = reflect.TypeOf(map[string]*complex128{})
	complex64Type                   = reflect.TypeOf(complex64(0))
	complex64PtrType                = reflect.TypeOf(new(complex64))
	complex64SliceType              = reflect.TypeOf([]complex64{})
	complex64SlicePtrType           = reflect.TypeOf([]*complex64{})
	complex64MapType                = reflect.TypeOf(map[string]complex64{})
	complex64MapPtrType             = reflect.TypeOf(map[string]*complex64{})
	boolPtrMapType                  = reflect.TypeOf(map[string]*bool{})
	wrappedBoolPtrMapType           = reflect.TypeOf(map[string]*Bool{})
	parseableComplex128Type         = reflect.TypeOf(ParseableComplex128(0))
	parseableComplex128PtrType      = reflect.TypeOf(new(ParseableComplex128))
	parseableComplex128SliceType    = reflect.TypeOf([]ParseableComplex128{})
	parseableComplex128SlicePtrType = reflect.TypeOf([]*ParseableComplex128{})
	parseableComplex128MapType      = reflect.TypeOf(map[string]ParseableComplex128{})
	parseableComplex128MapPtrType   = reflect.TypeOf(map[string]*ParseableComplex128{})
	parseableComplex64Type          = reflect.TypeOf(ParseableComplex64(0))
	parseableComplex64PtrType       = reflect.TypeOf(new(ParseableComplex64))
	parseableComplex64SliceType     = reflect.TypeOf([]ParseableComplex64{})
	parseableComplex64SlicePtrType  = reflect.TypeOf([]*ParseableComplex64{})
	parseableComplex64MapType       = reflect.TypeOf(map[string]ParseableComplex64{})
	parseableComplex64MapPtrType    = reflect.TypeOf(map[string]*ParseableComplex64{})
)

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

type Bool bool

func (b *Bool) UnmarshalFlag(s string) error {
	c, err := strconv.ParseBool(strings.TrimSpace(s))
	if err != nil {
		return err
	}

	*b = Bool(c)

	return nil
}

type BoolSlice []bool

func (b *BoolSlice) UnmarshalFlag(s string) error {
	c, err := strconv.ParseBool(strings.TrimSpace(s))
	if err != nil {
		return err
	}
	*b = append(*b, c)

	return nil
}

type BoolPtrSlice []*bool

func (b *BoolPtrSlice) UnmarshalFlag(s string) error {
	c, err := strconv.ParseBool(strings.TrimSpace(s))
	if err != nil {
		return err
	}
	*b = append(*b, &c)

	return nil
}
