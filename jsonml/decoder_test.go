package jsonml

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"
	"testing"
)

func init() {
	fmt.Print()
}

func TestDecoder_getAttributes(t *testing.T) {
	jsonml := `["table",{"a": "b", "1": "2"},["foo"]]`
	d := NewDecoder(strings.NewReader(jsonml))
	var b bytes.Buffer
	xmlEncoder := xml.NewEncoder(&b)
	err := d.WriteXml(xmlEncoder)

	if err != nil {
		PrintErrorTrace(err)
		t.Error("Should not return an error")
	}

	xmlEncoder.Flush()
	expected := `<table 1="2" a="b"><foo></foo></table>`
	if string(b.Bytes()) != expected {
		t.Error("Did not match expected output")
	}
}

func TestDecoder_WriteXml(t *testing.T) {
	jsonml := `
["table",{"class":"MyTable","style":"background-color:yellow"},["tr",["td",{"class":"MyTD","style":"border:1px solid black"},"#550758"],["td",{"class":"MyTD","style":"background-color:red"},"Example text here"]],["tr",["td",{"class":"MyTD","style":"border:1px solid black"},"#993101"],["td",{"class":"MyTD","style":"background-color:green"},"127624015"]],["tr",["td",{"class":"MyTD","style":"border:1px solid black"},"#E33D87"],["td",{"class":"MyTD","style":"background-color:blue"},"\u00a0",["span",{"style":"background-color:maroon"},"\u00a9"],"\u00a0"]]]
`
	var b bytes.Buffer
	xmlEncoder := xml.NewEncoder(&b)
	d := NewDecoder(strings.NewReader(jsonml))
	err := d.WriteXml(xmlEncoder)
	if err != nil {
		PrintErrorTrace(err)
		t.Error("Should not return an error")
	}

	xmlEncoder.Flush()
	expected := `<table style="background-color:yellow" class="MyTable"><tr><td style="border:1px solid black" class="MyTD">#550758</td><td style="background-color:red" class="MyTD">Example text here</td></tr><tr><td style="border:1px solid black" class="MyTD">#993101</td><td style="background-color:green" class="MyTD">127624015</td></tr><tr><td style="border:1px solid black" class="MyTD">#E33D87</td><td style="background-color:blue" class="MyTD"> <span style="background-color:maroon">©</span> </td></tr></table>`
	if string(b.Bytes()) != expected {
		t.Error("Did not match expected output")
	}
}

func TestDecoder_WriteXml_failStart(t *testing.T) {
	jsonml := `{"foo": "bar"}`
	var b bytes.Buffer
	xmlEncoder := xml.NewEncoder(&b)
	d := NewDecoder(strings.NewReader(jsonml))
	err := d.WriteXml(xmlEncoder)
	if err == nil {
		t.Error("Should return an error")
	}

	if !IsParseError(err) {
		t.Error("Should be a parse error")
	}
}
