package glwapi_test


import (
	"testing"
	"com/glwapi"
	"encoding/json"
)

func TestAuthToken(t * testing.T) {

	g := glwapi.D()
	g.WeChat.AppId = "1234"
	g.WeChat.Secret = "1234"
	str := g.TokenURL()
	t.Logf("%s", str)
	t.Logf("%s", "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=1234&secret=1234")
	if str =="https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=1234&secret=1234" {
		t.Logf("TEST AuthTokenURL ok")
	} else {
		t.Errorf(" url mismatch")
	}
}


type TokenJson struct {
	A   string   `json:"access_token"`
	E   int      `json:"expires_in"`
	EM  string   `json:"errmsg"`
	ER  int      `json:"errcode"`
}

func TestAuthTokenResp(t * testing.T) {
	g := glwapi.D()
	var tj = TokenJson{A: "aaa", E:7200, ER: 40029}
	b, e1 := json.Marshal(&tj)
	t.Logf("%s", string(b[:len(b)]))
	if e1 != nil {
		t.Errorf(" prepare test data failed" + e1.Error())
		return
	}
	w, e := g.HandleTokenAuthResp(b)
	if e == nil {
		t.Errorf(" should not success")
		return
	}
	
	tj = TokenJson{A: "aaa", E:7200}
	b, e1 = json.Marshal(&tj)
	if e1 != nil {
		t.Errorf(" prepare test data failed" + e1.Error())
		return
	}
	w, e = g.HandleTokenAuthResp(b)
	if e == nil {
		if w.Token == "aaa" {
			t.Logf("TEST  ok")
		} else {
			t.Errorf("token extract failed")
			return
		}
	} else {
		t.Errorf("failed:" + e.Error())
	}
}
