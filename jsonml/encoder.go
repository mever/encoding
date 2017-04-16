package jsonml

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"strings"
)

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

type Encoder struct {
	w io.Writer
}

func (e *Encoder) ReadXml(d *xml.Decoder) error {
	var err error
	var t xml.Token
	var newElement bytes.Buffer
	writeSeparator := false
	firstElement := true
	for {

		if newElement.Len() > 0 {
			if writeSeparator && !firstElement {
				e.w.Write([]byte(","))
			}

			firstElement = false
			e.w.Write(newElement.Bytes())
		}

		t, err = d.RawToken()
		if nil != err {
			if io.EOF != err {
				return err
			}
			break
		}

		newElement.Reset()
		writeSeparator = false
		switch v := t.(type) {
		case xml.StartElement:
			if debug {
				fmt.Printf("StartElement: %s\n", v)
			}

			if data, err := json.Marshal(v.Name.Local); err == nil {
				writeSeparator = true
				newElement.WriteString("[")
				newElement.Write(data)

				if len(v.Attr) > 0 {
					newElement.WriteString(",")
					attributes := make(map[string]string, len(v.Attr))
					for _, attr := range v.Attr {
						attributes[attr.Name.Local] = attr.Value
					}

					if data, err = json.Marshal(attributes); nil == err {
						newElement.Write(data)
					} else {
						return err
					}
				}
			} else {
				return err
			}

		case xml.EndElement:
			if debug {
				fmt.Printf("EndElement: %s\n", v)
			}
			newElement.WriteString("]")

		case xml.CharData:
			if data, _ := json.Marshal(strings.Trim(string(v), " \n\r")); len(data) > 2 {
				newElement.Write(data)
				writeSeparator = true
			}

		case xml.Comment, xml.ProcInst, xml.Directive:
		default:
			return errors.New("must never be executed")
		}
	}

	return nil
}
