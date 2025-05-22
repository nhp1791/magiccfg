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
	timePtrType              = reflect.TypeOf(new(time.Time))
	parseableTimePtrType     = reflect.TypeOf(new(ParseableTime))
	urlType                  = reflect.TypeOf(new(url.URL))
	parseableURLType         = reflect.TypeOf(new(ParseableURL))
	packageTimeFormats       []string
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

func UnmarshalDurationPtrEnv(v string) (any, error) {
	dur, err := time.ParseDuration(v)
	if err != nil {
		return "", err
	}
	pd := ParseableDuration(dur)
	return &pd, nil
}

// XMLTime enables time.Time to be parsed by the encoding/xml package
type ParseableTime time.Time

func (p *ParseableTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}
	if packageTimeFormats != nil {
		for _, v := range packageTimeFormats {
			t, err := time.Parse(v, strings.TrimSpace(s))
			if err == nil {
				*p = ParseableTime(t)
				return nil
			}
		}
	}
	return fmt.Errorf("Could not parse time: %s", s)
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
	return fmt.Errorf("Could not parse time: %s", s)

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
