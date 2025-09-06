package types

import "encoding/xml"

type XMLTVBool bool

func (b *XMLTVBool) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	if b == nil {
		return xml.Attr{}, nil
	}

	attributeValue := "no"
	if *b {
		attributeValue = "yes"
	}

	return xml.Attr{Name: name, Value: attributeValue}, nil
}

func (b *XMLTVBool) UnmarshalXMLAttr(attr xml.Attr) error {
	*b = XMLTVBool(false)
	if attr.Value == "yes" {
		*b = XMLTVBool(true)
	}

	return nil
}

func (b *XMLTVBool) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if b == nil {
		return e.EncodeElement(nil, start)
	}

	var elementText string
	switch {
	case start.Name.Local == "new":
		elementText = ""
	case bool(*b):
		elementText = "yes"
	default:
		elementText = "no"
	}

	return e.EncodeElement(elementText, start)
}

func (b *XMLTVBool) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var elementText string
	if err := d.DecodeElement(&elementText, &start); err != nil {
		return err
	}

	var v bool
	switch {
	case start.Name.Local == "new":
		v = true
	case elementText == "yes":
		v = true
	}

	*b = XMLTVBool(v)

	return nil
}
