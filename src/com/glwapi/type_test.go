
package glwapi_test

import (
	"testing"
	"com/glwapi"
	"encoding/xml"
)


type TestD struct {
	T1  glwapi.XmlD `xml: "ta"`
}

func TestMarshalType(t * testing.T) {
	var td TestD =  TestD{T1 : glwapi.XmlD{Key: "a",Val: "b", Opt: true, IsCDATA:true}}
	b, e := xml.MarshalIndent(td, " ", "  ")
	if e == nil {
		t.Logf(" Test ok data ==>\n %s\n", string(b[:len(b)]))
	} else {
		t.Errorf("Test Failed %s", e.Error())
	}

	// test without value
	td =  TestD{T1 : glwapi.XmlD{Key: "a",Val: "", Opt: true, IsCDATA:false}}
	b, e = xml.MarshalIndent(td, " ", "  ")
	if e == nil {
		t.Logf(" Test ok data ==>\n %s\n", string(b[:len(b)]))
	} else {
		t.Errorf("Test Failed %s", e.Error())
	}

	td =  TestD{T1 : glwapi.XmlD{Key: "a",Val: "eee", Opt: true, IsCDATA:false}}
	b, e = xml.MarshalIndent(td, " ", "  ")
	if e == nil {
		t.Logf(" Test ok data ==>\n %s\n", string(b[:len(b)]))
	} else {
		t.Errorf("Test Failed %s", e.Error())
	}

}
