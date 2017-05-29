
package glwapi

import (
     "fmt"
)


type WeChatUser struct {
	Token 	    string      `json:"access_token"`
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
