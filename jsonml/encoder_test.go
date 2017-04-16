package jsonml

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"testing"
)

func init() {
	fmt.Print()
}

func TestNewEncoder(t *testing.T) {
	xmlStr := `
<table class="MyTable" style="background-color:yellow">
<tr>
<td class="MyTD" style="border:1px solid black">
#5D28D1</td>
<td class="MyTD" style="background-color:red">
Example text here</td>
</tr>
<tr>
<td class="MyTD" style="border:1px solid black">
#AF44EF</td>
<td class="MyTD" style="background-color:green">
127310656</td>
</tr>
<tr>
<td class="MyTD" style="border:1px solid black">
#AAD034</td>
<td class="MyTD" style="background-color:blue">
&nbsp;
<span style="background-color:maroon">&copy;</span>
&nbsp;
</td>
</tr>
</table>`

	b := &bytes.Buffer{}
	xmlDec := xml.NewDecoder(bytes.NewBufferString(xmlStr))
	xmlDec.Entity = xml.HTMLEntity

	enc := NewEncoder(b)
	enc.ReadXml(xmlDec)

	expected := `["table",{"class":"MyTable","style":"background-color:yellow"},["tr",["td",{"class":"MyTD","style":"border:1px solid black"},"#5D28D1"],["td",{"class":"MyTD","style":"background-color:red"},"Example text here"]],["tr",["td",{"class":"MyTD","style":"border:1px solid black"},"#AF44EF"],["td",{"class":"MyTD","style":"background-color:green"},"127310656"]],["tr",["td",{"class":"MyTD","style":"border:1px solid black"},"#AAD034"],["td",{"class":"MyTD","style":"background-color:blue"}," ",["span",{"style":"background-color:maroon"},"©"]," "]]]`
	if string(b.Bytes()) != expected {
		t.Fail()
	}
}
