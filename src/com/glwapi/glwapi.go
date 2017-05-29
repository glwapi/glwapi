
package glwapi


var  G_API * GlwapiConfig = nil

type GlwapiConfig struct {
	WeChat       * WeChatConfig
	StateCheck     bool
	users        map[string]*WeChatUser
}

type WeChatConfig struct {
	AppId      string
	Secret     string
}


type  GlwapiHandler interface {
	Init()  	bool
	Destroy()  	bool
}



func D() * GlwapiConfig {
	if G_API == nil {
		G_API = &GlwapiConfig{WeChat : &WeChatConfig{}, users : map[string]*WeChatUser{}}
	}
	return G_API
}


