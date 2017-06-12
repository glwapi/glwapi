
package glwapi

import (
	"encoding/json"
	"fmt"
)


const (
	token_wechat_url = "https://api.weixin.qq.com/cgi-bin/token"
)

// Auth Error Message
type AuthError struct {
	Err    string
}

// Auth Error Message implments error interface
func (ae  AuthError) Error() string {
	return ae.Err
}

/*
  WeChat auth response struct wrapper
*/
type AuthRespJson struct {
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	AccToken    string `json:"access_token"`
	Expires     int    `json:"expires_in"`
}

/*
   Auth API interface
 */
type AuthApi interface {
	// Generate TOKEN URL
	TokenURL()    string

	// Handle Token URL response and return struct of WeChatConfig if succeeed
	// Otherwise error is not nil
	HandleTokenAuthResp(d []byte)  (*WeChatConfig, error)
}


func (g * GlwapiConfig) TokenURL() string {
	return g.WeChat.TokenURL()
}

func (w * WeChatConfig) TokenURL() string {
	var u string
	u = fmt.Sprintf("%s?grant_type=client_credential&appid=%s&secret=%s", token_wechat_url, w.AppId, w.Secret)
	return u
}

func (w * WeChatConfig) HandleTokenAuthResp(d []byte)  (*WeChatConfig, error) {
	if d == nil {
		return nil, AuthError{"No data avaialble"}
	}
	var arj AuthRespJson = AuthRespJson{}
	if err := json.Unmarshal(d, &arj); err == nil {
		if arj.ErrCode > 0 {
			return nil, AuthError{"Auth failed:" + arj.ErrMsg}
		} else {
			w.Token = arj.AccToken
			w.Expires = arj.Expires
		}
		return w, nil
	} else {
		return nil, AuthError{"Can not parse response:" + err.Error()}
	}
}

func (g * GlwapiConfig) HandleTokenAuthResp(d []byte)  (*WeChatConfig, error) {
	return g.WeChat.HandleTokenAuthResp(d)
}



func (w *WeChatConfig) IsExpired() bool {
	//TODO add check 
	return false
}

func (w *WeChatConfig) IsTokened() bool {
	if len(w.Token) > 0 {
		return true
	} else {
		return false
	}
}


