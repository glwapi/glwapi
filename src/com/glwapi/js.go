
package glwapi

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	js_ticket_url = "https://api.weixin.qq.com/cgi-bin/ticket/getticket?"
)

// JS Auth Error Message
type JSAuthError struct {
	Err      string
}

// JS Auth Error Message implements error interface
func (js * JSAuthError) Error() string {
	return js.Err
}

// JS struct 
// For Web AUTH
type WeChatJS struct {

	//WeChat configuration
	WeChat   *WeChatConfig

	//nonce
	Nonce     string

	//Signature 
	Sign      string

	//Timestamp of string
	Ts        string

	//Timestamp of time 
	Stamp     time.Time

	//token of js ticket
	Ticket    string

	//Auth url from client
	U         string
}

//String interface
func (js * WeChatJS) String() string {
	var s string
	s = fmt.Sprintf("{WeChat: %p, Nonce: %s, Sign : %s, Ts : %s, Ticket : %s, U : %s}", js.WeChat, js.Nonce, js.Sign, js.Ts, js.Ticket, js.U)
	return s
}

// JS Ticket response wrapper
type jsticket struct {
	Errcode    int     `json:"errcode"`
	Errmsg     string  `json:"errmsg"`
	Ticket     string  `json:"ticket"`
}

// Get JS Auth nonce data
func (js * WeChatJS) N() string {
	if len(js.Nonce) <= 0 {
		js.Nonce = R(32)
	}
	return js.Nonce
}


//Get JS Auth sign data
func (js *WeChatJS)  S() string {
	if len(js.Sign) <= 0 {
		ts := js.T()
		no := js.N()
		u  := js.U
		str := "jsapi_ticket="+js.Ticket+"&noncestr="+no+"&timestamp="+ts+"&url=" + u
		js.Sign = Sha1(str)
	}

	return js.Sign
}


//Get JS timestamp data
func (js *WeChatJS)  T() string {
	if len(js.Ts) <= 0 {
		js.Ts = fmt.Sprintf("%d", js.Stamp.Unix())
	}
	return js.Ts
}

/*
   JS Auth API
   Before JS auth, you have get server token by Auth interface first
   Follow step to get JS Auth:
      1.  GetJSTicketURL   to generate JS Ticket URL
      2.  HandleTicketResp  to extract ticket from  response
      3.  GetJSAuth to get JS Auth struct for web
 */
type JsApi interface {

	// Generate URL for JS AUTH ticket
	// return ticket url if anythings goes well
	// otherwise error is not nil
	GetJSTicketURL()       (string, error)

	// generat JS Auth struct for web
	//  ticket:  js auth ticket
	//  url   :  auth domain
	GetJSAuth(ticket string, url string)  (*WeChatJS, error)

	// Extract ticket from response	
	// error is not nil if something wrong
	HandleTicketResp(d []byte)   (string, error)
}

// JsApi interface of WeChatConfig implements
// To get ja auth object for js
func (w * WeChatConfig) GetJSAuth(ticket string, url string)  (* WeChatJS, error) {
	var js *WeChatJS = &WeChatJS{WeChat: w, Stamp : time.Now(), Ticket: ticket, U : url}
	return js, nil
}


// JsApi interface of WeChatConfig implements
// Generate url that to get js api token 
// return error for failed to generation
func (w * WeChatConfig) GetJSTicketURL() (string, error) {
	if w.IsExpired() == false && w.IsTokened() == true {
		str := fmt.Sprintf("%saccess_token=%s&type=jsapi", js_ticket_url, w.Token)
		return str, nil
	} else {
		return "", &JSAuthError{" Expired or not get token yet"}
	}
}


// JsApi interface of GlwapiConfig implements
// To get ja auth object for js
func (c * GlwapiConfig) GetJSAuth(ticket string, url string)  (* WeChatJS, error) {
	return	c.WeChat.GetJSAuth(ticket, url)
}


// JsApi interface of GlwapiConfig implements
// Generate url that to get js api token 
// return error for failed to generation
func (c * GlwapiConfig) GetJSTicketURL() (string, error) {
	return	c.WeChat.GetJSTicketURL()
}

// JsApi interface of GlwapiConfig implements
// return error for failed extract
func (c * GlwapiConfig) HandleTicketResp(d []byte)   (string, error) {
	if d == nil {
		return "", &JSAuthError{" No data available"}
	}

	var tic jsticket = jsticket{}
	if err := json.Unmarshal(d, &tic); err == nil {
		if tic.Errcode > 0 {
			return "", &JSAuthError{" get ticket failed:"+tic.Errmsg}
		}
		ts := time.Now()
		c.WeChat.JSToken = tic.Ticket
		c.WeChat.JSRtime = &ts
		return tic.Ticket, nil
	} else {
		return "", &JSAuthError{" parse response failed " + err.Error()}
	}
}



