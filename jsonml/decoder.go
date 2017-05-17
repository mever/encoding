package jsonml

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/pkg/errors"
	"io"
)

// A Decoder represents a JSONML parser reading a particular input stream.
// The parser assumes that its input is encoded in UTF-8.
type Decoder struct {
	jsonDec *json.Decoder
	xmlEnc  *xml.Encoder
	level   []string

	Prefix, Indent string
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{jsonDec: json.NewDecoder(r)}
}

func (d *Decoder) WriteXml(e *xml.Encoder) error {
	t, err := d.jsonDec.Token()
	if err != nil {
		return newTokenError(err)
	}

	if del, ok := t.(json.Delim); !ok || string(del) != "[" {
		return newParseError("[", t)
	}

	d.xmlEnc = e
	e.Indent(d.Prefix, d.Indent)
	err = d.convertElement(t)
	if err == io.EOF || err == nil {
		return e.Flush()
	}

	return err
}

func (d *Decoder) convertElement(t json.Token) error {

	if debug {
		fmt.Printf("convertElement: %s\n", t)
	}

	switch v := t.(type) {
	case string:
		if err := d.convertTextNode(v); err != nil {
			return err
		} else {
			return d.convertElementList(nil)
		}

	case json.Delim:
		if string(v) == "[" {
			if tagNameToken, err := d.jsonDec.Token(); err == nil {
				if tagName, ok := tagNameToken.(string); ok {
					return d.convertStartTag(tagName)
				} else {
					return newParseError("tag-name", tagNameToken)
				}
			} else {
				return errors.WithStack(err)
			}
		}
	}
	return newParseError("[ or textNode", t)
}

func (d *Decoder) convertElementList(t json.Token) error {
	if t == nil {
		var err error
		t, err = d.jsonDec.Token()
		if err != nil {
			return newTokenError(err)
		}
	}

	if debug {
		fmt.Printf("convertElementList: %s\n", t)
	}

	if del, ok := t.(json.Delim); ok && string(del) == "]" {

		lvl := d.level[len(d.level)-1]
		d.level = d.level[:len(d.level)-1]

		if debug {
			fmt.Printf("end element: %s\n", lvl)
		}

		ee := xml.EndElement{}
		ee.Name.Local = lvl
		if err := d.xmlEnc.EncodeToken(ee); err != nil {
			return errors.WithStack(err)
		}

		if len(d.level) == 0 {
			return nil
		}

		return d.convertElementList(nil)
	}

	return d.convertElement(t)
}

func (d *Decoder) convertStartTag(tagName string) error {
	t, err := d.jsonDec.Token()
	if err != nil {
		return newTokenError(err)
	}

	if debug {
		fmt.Printf("convertStartTag %s: %s\n", tagName, t)
	}

	se := xml.StartElement{}
	se.Name.Local = tagName
	d.level = append(d.level, tagName)

	if del, ok := t.(json.Delim); ok && string(del) == "{" {
		t = nil
		if se.Attr, err = d.getAttributes(); err != nil {
			return err
		}
	}

	if err := d.xmlEnc.EncodeToken(se); err != nil {
		return errors.WithStack(err)
	}

	return d.convertElementList(t)
}

func (d *Decoder) convertTextNode(t string) error {
	if err := d.xmlEnc.EncodeToken(xml.CharData(t)); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (d *Decoder) getAttributes() ([]xml.Attr, error) {
	t, err := d.jsonDec.Token()
	if err != nil {
		return []xml.Attr{}, newTokenError(err)
	}

	if del, ok := t.(json.Delim); ok && string(del) == "}" {
		return []xml.Attr{}, nil
	}

	if attrName, ok := t.(string); ok {
		a := xml.Attr{}
		a.Name.Local = attrName
		if av, err := d.getAttrValue(); err != nil {
			return []xml.Attr{}, err
		} else {
			a.Value = av
		}

		attributes, err := d.getAttributes()
		return append(attributes, a), err
	}

	return []xml.Attr{}, newParseError("} or attribute name", t)
}

func (d *Decoder) getAttrValue() (string, error) {
	t, err := d.jsonDec.Token()
	if err != nil {
		return "", newTokenError(err)
	}

	if s, ok := t.(string); ok {
		return s, nil
	} else {
		return "", newParseError("string (others are not yet implemented!)", t)
	}
}
