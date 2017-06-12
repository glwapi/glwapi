
package glwapi


import (
	"fmt"
)

var  G_API * GlwapiConfig = nil

func D() * GlwapiConfig {
	if G_API == nil {
		G_API = new(GlwapiConfig)
		G_API.WeChat = new(WeChatConfig)
		G_API.users = make(map[string]*WeChatUser)
	}
	return G_API
}

type  GlwapiHandler interface {
	Init(appid string, secret string, mchid string, apikey string)  	bool
	Destroy()  	bool

	RegEventListener(h * GlwapiEventsNotifier)
}

type GlwapiEvent interface {
	Type()	int	
	Data()	interface{}
}

type GlwapiEventsNotifier interface {
	OnEvent(e GlwapiEvent)
}


type GlwapiConfig struct {
	WeChat		* WeChatConfig
	StateCheck	bool
	SignCheck	bool
	users		map[string]*WeChatUser
	init		bool
	n		* GlwapiEventsNotifier
}

func (c * GlwapiConfig) IsInit() bool {
	return c.init
}

func (c * GlwapiConfig) String() string {
	return fmt.Sprintf("{WeChat : %s, n : %x}", c.WeChat, c.n)
}


func (c * GlwapiConfig) Init(appid string, secret string, mchid string, apikey string) bool {
	c.WeChat.AppId  = appid
	c.WeChat.Secret = secret
	c.WeChat.MchId  = mchid
	c.WeChat.MchApiKey = apikey
	c.init = true
	return true
}


func (c * GlwapiConfig) Destroy() bool {
	return false
}


func (c * GlwapiConfig) RegEventListener(h * GlwapiEventsNotifier) {
	c.n = h
}

