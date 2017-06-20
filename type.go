
package glwapi

import (
	"encoding/xml"
)

type XmlD struct {
	Key		string
	Val		string
	Opt		bool
	IsCDATA		bool
	ExludeSign	bool
}



func (x XmlD) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if x.Opt == true && x.Val == "" {
		return nil
	}

	if x.IsCDATA == true {
		e.EncodeToken(xml.StartElement{Name: xml.Name{Local: x.Key}})
		e.EncodeToken(xml.Directive("[CDATA[" + x.Val+"]]"))
		e.EncodeToken(xml.EndElement{Name: xml.Name{Local: x.Key}})
	} else {
		e.EncodeElement(x.Val, xml.StartElement{Name: xml.Name{Local: x.Key}})
	}
	return nil
}


func (x * XmlD) UnmarshalXML(d * xml.Decoder, start xml.StartElement) error {
	for t, b := d.Token(); t != nil && b == nil; t, b = d.Token() {
		switch  t.(type) {
		case xml.EndElement:
			return nil
		case xml.CharData:
			x.Val = string(t.(xml.CharData))
			break
		}
	}
	return nil
}
