package magiccfg

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"time"
)

var (
	durationType             = reflect.TypeOf(time.Duration(0))
	durationPtrType          = reflect.TypeOf(new(time.Duration))
	parseableDurationType    = reflect.TypeOf(ParseableDuration(0))
	parseableDurationPtrType = reflect.TypeOf(new(ParseableDuration))
	timeType                 = reflect.TypeOf(time.Now())
	timePtrType              = reflect.TypeOf(new(time.Time))
	parseableTimeType        = reflect.TypeOf(ParseableTime(time.Now()))
	parseableTimePtrType     = reflect.TypeOf(new(ParseableTime))
	urlType                  = reflect.TypeOf(url.URL{})
	urlPtrType               = reflect.TypeOf(new(url.URL))
	parseableURLType         = reflect.TypeOf(ParseableURL{})
	parseableURLPtrType      = reflect.TypeOf(new(ParseableURL))
	packageTimeFormats       []string
	parseableTimeKind        = reflect.TypeOf(time.Time{}).Kind()
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
