
package glwapi

import (
	"fmt"
	"encoding/json"
	"net/url"
	"net/http"
	"time"
)

const (
	code_url = "https://open.weixin.qq.com/connect/oauth2/authorize"
	token_url = "https://api.weixin.qq.com/sns/oauth2/access_token"
)


//  Error message for auth
type UserAuthError struct {
	Err   string
}

// error interface implement
func (err *UserAuthError) Error() string { return err.Err }



//  Related User Auth API interface:
//    Whole auth process follow below sequence:
//    1.  Code URL
//    2.  HandleCodeURLResponse
//    3.  TokenURL
//    4.  HandleTokenURLResponse
//    5.  TokenURL
//    6.  HandleFTokenURLResponse
type UserAuthApi  interface {

	//Generate Code get URL for wechat auth first step 
	//return  url state error
        //If redirect url, will throw error
        //If no state, will autmatically generate state info
	CodeURL(r string, s string)      (string, string, error)

	// Generate Token get URL for wechat auth 3th step 
	//return  url  error
        //    If No code, will throw error
	UserTokenURL(code string)     (string, error)

	//Generate refresh taoken URL
	UserFTokenURL(ftoken string)  (string, error)

	//Extract code from wechat response and return code, state_code,error
	HandleCodeURLResponse(r * http.Request)   (string, string, error)

	//Check code and state code legal or not
	HandleCode(code string, state string)   (string, string, error)

	//Extract token from response data
	//If extract succeed, Return user object
	HandleUserTokenURLResponse(data []byte)   (*WeChatUser, error)

	//Extract refresh token from response data
	//If extract succeed, Return user object
	HandleUserFTokenURLResponse(data []byte)   (*WeChatUser, error)
}




// Generate Code get URL for wechat auth first setp 
// return  url state error
func (c * GlwapiConfig) CodeURL(r string, s string)  (string, string, error) {
	var st string
	if  len(r) <=0 {
		return "", "", &UserAuthError{"No redirect URL"}
	}
	ur := url.QueryEscape(r)
	if len(s) <= 0 {
		st = time.Now().String()
	} else {
		st = s;
	}

	u := fmt.Sprintf("%s?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_userinfo&state=%s#wechat_redirect", code_url, c.WeChat.AppId, ur, st)
	return u, st, nil
}


// Generate Token get URL for wechat auth 3th step 
//return  url  error
//    If No code, will throw error
func (c * GlwapiConfig) UserTokenURL(code string) (string, error) {
	if len(code) <= 0 {
		return "", &UserAuthError{"No available code"}
	}
	u := fmt.Sprintf("%s?appid=%s&secret=%s&code=%s&grant_type=authorization_code", token_url, c.WeChat.AppId, c.WeChat.Secret, code)
	return u, nil
}

func (c * GlwapiConfig) UserFTokenURL(ftoken string) (string, error) {
	if len(ftoken) <= 0 {
		return "", &UserAuthError{"No available refresh toekn"}
	}
	u := fmt.Sprintf("%s?appid=%s&grant_type=refresh_token&refresh_token=%s", token_url, c.WeChat.AppId, ftoken)
	return u, nil
}

//Extract code from wechat response and return code, state_code,error
func (c * GlwapiConfig) HandleCodeURLResponse(r * http.Request)   (string, string, error) {
	if r == nil {
		return "", "", &UserAuthError{"No request "}
	}
	return c.HandleCode(r.FormValue("code"), r.FormValue("state"))
}

func (c * GlwapiConfig)	HandleCode(code string, state string)   (string, string, error) {
	//TODO check state value
	if c.StateCheck == false {
		return code, state, nil
	} else {
		//TODO Check state
		return code, state, nil
	}
}


//Extract token from response data
//If extract succeed, Return user object
func (c * GlwapiConfig)	HandleUserTokenURLResponse(data []byte)   (*WeChatUser, error) {
	var user * WeChatUser =&WeChatUser{}
	if data == nil {
		return nil, &UserAuthError{"No data"}
	}
	if err := json.Unmarshal(data, user); err == nil {
		if user.Err > 0 {
			er := fmt.Sprintf("Umaarshal failed: %d",user.Err)
			return nil, &UserAuthError{er}
		}
		var u * WeChatUser = c.users[user.OpenId]
		if u == nil {
			c.users[user.OpenId] = user
			u = user
		} else {
			u.Token = user.Token
			u.FToken = user.FToken
			u.Expires = user.Expires
			u.OpenId = user.OpenId
			u.Err = user.Err
		}
		return u, nil
	} else {
		return nil, &UserAuthError{"Umaarshal failed:"+err.Error()}
	}
}


//Extract refresh token from response data
//If extract succeed, Return user object
func (c * GlwapiConfig)	HandleUserFTokenURLResponse(data []byte)   (*WeChatUser, error) {
	var user * WeChatUser =&WeChatUser{}
	if data == nil {
		return nil, &UserAuthError{"No data"}
	}
	if err := json.Unmarshal(data, user); err == nil {
		if user.Err > 0 {
			er := fmt.Sprintf("Umaarshal failed: %d",user.Err)
			return nil, &UserAuthError{er}
		}
		var u * WeChatUser = c.users[user.OpenId]
		if u == nil {
			c.users[user.OpenId] = user
			u = user
		} else {
			u.Token = user.Token
			u.FToken = user.FToken
			u.Expires = user.Expires
			u.OpenId = user.OpenId
			u.Err = user.Err
		}
		return u, nil
	} else {
		return nil, &UserAuthError{"Umaarshal failed:"+err.Error()}
	}
}


