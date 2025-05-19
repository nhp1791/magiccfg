package magiccfg

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"reflect"
	"time"
)

var (
	durationPtrType          = reflect.TypeOf(new(time.Duration))
	sliceDurationType        = reflect.TypeOf([]*time.Duration{})
	xmlTimeDurationPtrType   = reflect.TypeOf(new(XMLTimeDuration))
	sliceXMLTimeDurationType = reflect.TypeOf([]*XMLTimeDuration{})
	timePtrType              = reflect.TypeOf(new(time.Time))
	xmlTimePtrType           = reflect.TypeOf(new(XMLTime))
	sliceTimeType            = reflect.TypeOf([]*time.Time{})
	sliceXMLTimeType         = reflect.TypeOf([]*XMLTime{})
	urlType                  = reflect.TypeOf(new(url.URL))
	xmlURLType               = reflect.TypeOf(new(XMLURL))
	sliceURLType             = reflect.TypeOf([]*url.URL{})
	sliceXMLURLType          = reflect.TypeOf([]*XMLURL{})
	packageTimeFormats       []string
)

// XMLTimeDuration enables time.Duration to be parsed by the encoding/xml package
type XMLTimeDuration time.Duration

func (x *XMLTimeDuration) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}
	dur, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*x = XMLTimeDuration(dur)
	return nil
}

func (x *XMLTimeDuration) UnmarshalFlag(value string) error {
	dur, err := time.ParseDuration(value)
	if err != nil {
		return err
	}
	*x = XMLTimeDuration(dur)
	return nil
}

// XMLTime enables time.Time to be parsed by the encoding/xml package
type XMLTime time.Time

func (x *XMLTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}
	if packageTimeFormats != nil {
		for _, v := range packageTimeFormats {
			t, err := time.Parse(v, s)
			if err == nil {
				*x = XMLTime(t)
				return nil
			}
		}
	}
	return fmt.Errorf("Could not parse XML time: %s", s)
}

func (x *XMLTime) UnmarshalFlag(s string) error {
	if packageTimeFormats != nil {
		for _, v := range packageTimeFormats {
			t, err := time.Parse(v, s)
			if err == nil {
				*x = XMLTime(t)
				return nil
			}
		}
	}
	return fmt.Errorf("Could not parse time: %s", s)

}

type XMLURL url.URL

func (x *XMLURL) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}
	u, err := url.Parse(s)
	if err != nil {
		return err
	}
	*x = XMLURL(*u)
	return nil
}

func (x *XMLURL) UnmarshalText(val []byte) error {
	u, err := url.Parse(string(val))
	if err != nil {
		return err
	}
	*x = XMLURL(*u)
	return nil
}

func (x *XMLURL) UnmarshalFlag(s string) error {
	u, err := url.Parse(s)
	if err != nil {
		return err
	}
	*x = XMLURL(*u)
	return nil
}
