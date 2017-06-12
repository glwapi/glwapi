
package glwapi

import (
	"fmt"
	"time"
)

//Wechat security 
type WeChatSec struct {
	Nonce		string
	Sign		string
	SignType	string
}

//WeChat config
type WeChatConfig struct {
	AppId      string
	Secret     string
	Token      string
	Expires    int
	MchId      string
	MchApiKey  string
	Rtime      *time.Time
	JSToken	   string
	JSRtime    *time.Time
}

//Check Wechat appid or secret available or not
func (w * WeChatConfig) IsValid() bool {
	return w.AppId != "" && w.Secret != ""
}

//Check wechat config payment configuration available or not
func (w * WeChatConfig) IsPaymentValid() bool {
	return w.IsValid() && w.MchId != "" && w.MchApiKey != ""
}

//check JS token available or not
func (w * WeChatConfig) IsGotJSToken() bool {
	//FIXME check expired time
	return w.JSToken == ""
}

func (w * WeChatConfig) String() string {
	return fmt.Sprintf("{AppId : %s, Secret : %s , MchId : %s, Token : %s, Expires :%d}", w.AppId, w.Secret, w.MchId, w.Token, w.Expires)
}



type WeChatUser struct {
	Token	    string      `json:"access_token"`
	Expires     int	        `json:"expires_in"`
	FToken      string      `json:"refresh_token"`
	OpenId      string      `json:"openid"`
	Scope       string      `json:"scope"`
	Err         int         `json:"errcode"`
	P           * WeChatUserP
}


func (u * WeChatUser) String() string {
	d := fmt.Sprintf("Token: %s E:%d  Ftoken:%s  OpenId: %s", u.Token, u.Expires, u.FToken, u.OpenId)
	return d
}

type WeChatUserP struct {
	OpenId      string      `json:"openid"`
	NickName    string      `json:"nickname"`
	Sex         string      `json:"sex"`
	Province    string      `json:"province"`
	City        string      `json:"city"`
	Country     string      `json:"country"`
	ImgUrl      string      `json:"headimgurl"`
	UnionId     string      `json:"unionid"`
	Err         int         `json:"errcode"`
}
