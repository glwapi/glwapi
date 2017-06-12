
package glwapi_test

import (
	"com/glwapi"
	"testing"
	"encoding/json"
)


func TestJSTicketURL(t * testing.T) {
	g := glwapi.D()
	g.WeChat.Token = "1234"
	str, e := g.GetJSTicketURL()
	if e == nil {
		if str =="https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=1234&type=jsapi" {
			t.Logf("Test ok")
		} else {
			t.Errorf("Failed: js ticket url mismatch"+ str)
		}
	} else {
		t.Errorf("Failed: "+ e.Error())
	}
}


// JS Ticket response wrapper
type jsticket struct {
         Errcode    int     `json:"errcode"`
         Errmsg     string  `json:"errmsg"`
         Ticket     string  `json:"ticket"`
}


func TestHandleTicketResp(t * testing.T) {
	g := glwapi.D()
	jt := jsticket{Errcode : 0, Errmsg :"aaa", Ticket: "1234"}
	b, e := json.Marshal(&jt)
	t.Logf("%s", string(b[:len(b)]))
	if e == nil {
		str, e1 := g.HandleTicketResp(b)
		if e1 != nil {
			t.Errorf(" handle response failed:"+ e1.Error())
		} else {
			if str == "1234" {
				t.Logf("Test ok")
			} else {
				t.Errorf(" ticket mismatch:"+ str)
			}
			
		}
	} else {
		t.Errorf("Prepared test data failed: " + e.Error())
	}

	jt = jsticket{Errcode : 1234, Errmsg :"aaa", Ticket: "1234"}
	b, e = json.Marshal(&jt)
	t.Logf("%s", string(b[:len(b)]))
	if e == nil {
		_, e1 := g.HandleTicketResp(b)
		if e1 != nil {	
			t.Logf("Test ok")
		} else {
			t.Errorf("should not succeed")
		}
	} else {
		t.Errorf("prepared test data failed: " + e.Error())
	}



}
