package magiccfg

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"reflect"
	"time"
)

var (
	durationPtrType                = reflect.TypeOf(new(time.Duration))
	sliceDurationType              = reflect.TypeOf([]*time.Duration{})
	parseableTimeDurationPtrType   = reflect.TypeOf(new(ParseableTimeDuration))
	sliceParseableTimeDurationType = reflect.TypeOf([]*ParseableTimeDuration{})
	timePtrType                    = reflect.TypeOf(new(time.Time))
	parseableTimePtrType           = reflect.TypeOf(new(ParseableTime))
	sliceTimeType                  = reflect.TypeOf([]*time.Time{})
	sliceParseableTimeType         = reflect.TypeOf([]*ParseableTime{})
	urlType                        = reflect.TypeOf(new(url.URL))
	parseableURLType               = reflect.TypeOf(new(ParseableURL))
	sliceURLType                   = reflect.TypeOf([]*url.URL{})
	sliceParseableURLType          = reflect.TypeOf([]*ParseableURL{})
	mapDurationType                = reflect.TypeOf(map[string]*time.Duration{})
	mapParseableDurationType       = reflect.TypeOf(map[string]*ParseableTimeDuration{})
	packageTimeFormats             []string
)

// ParseableTimeDuration enables time.Duration to be parsed by the encoding/xml package
type ParseableTimeDuration time.Duration

func (p *ParseableTimeDuration) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}
	dur, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*p = ParseableTimeDuration(dur)
	return nil
}

func (p *ParseableTimeDuration) UnmarshalFlag(value string) error {
	dur, err := time.ParseDuration(value)
	if err != nil {
		return err
	}
	*p = ParseableTimeDuration(dur)
	return nil
}

// ParseableTime enables time.Time to be parsed by the encoding/xml package
type ParseableTime time.Time

func (p *ParseableTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}
	if packageTimeFormats != nil {
		for _, v := range packageTimeFormats {
			t, err := time.Parse(v, s)
			if err == nil {
				*p = ParseableTime(t)
				return nil
			}
		}
	}
	return fmt.Errorf("Could not parse XML time: %s", s)
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
	u, err := url.Parse(s)
	if err != nil {
		return err
	}
	*p = ParseableURL(*u)
	return nil
}

func (p *ParseableURL) UnmarshalText(val []byte) error {
	u, err := url.Parse(string(val))
	if err != nil {
		return err
	}
	*p = ParseableURL(*u)
	return nil
}

func (p *ParseableURL) UnmarshalFlag(s string) error {
	u, err := url.Parse(s)
	if err != nil {
		return err
	}
	*p = ParseableURL(*u)
	return nil
}
