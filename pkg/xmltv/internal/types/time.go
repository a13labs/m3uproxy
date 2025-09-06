package types

import (
	"encoding/xml"
	"time"
)

type XMLTVTime struct{ time.Time }

func (t *XMLTVTime) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	if t == nil {
		return xml.Attr{}, nil
	}

	if t.IsZero() {
		return xml.Attr{}, nil
	}

	return xml.Attr{Name: name, Value: t.Format("20060102150405 -0700")}, nil
}

func (t *XMLTVTime) UnmarshalXMLAttr(attr xml.Attr) error {
	if attr.Value == "" {
		*t = XMLTVTime{}

		return nil
	}

	tt, err := time.Parse("20060102150405 -0700", attr.Value)
	if err != nil {
		tt, err = time.Parse("20060102", attr.Value)
		if err != nil {
			return err
		}
	}

	*t = XMLTVTime{tt}

	return nil
}

func (t *XMLTVTime) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if t == nil {
		return e.EncodeElement(nil, start)
	}

	if t.IsZero() {
		return e.EncodeElement("", start)
	}

	return e.EncodeElement(t.Format("20060102"), start)
}

func (t *XMLTVTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	if err := d.DecodeElement(&v, &start); err != nil {
		return err
	}

	if v == "" {
		*t = XMLTVTime{}

		return nil
	}

	tt, err := time.Parse("20060102", v)
	if err != nil {
		return err
	}

	*t = XMLTVTime{tt}

	return nil
}
