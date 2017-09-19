
package glwapi

import (
	"fmt"
	"encoding/json"
	"net/http"
        "errors"
        "io/ioutil"
)

const (
	user_info_url = "https://api.weixin.qq.com/cgi-bin/user/info"
)

type UserManagementApi  interface {

    //Generate get user info url according token, openid, language
    // Return : https://api.weixin.qq.com/cgi-bin/user/info?access_token=ACCESS_TOKEN&openid=OPENID&lang=zh_CN  
    UserInfoURL(token string, openId  string, lang string) string


    //Handle Get User Info response and construct user data 
    // return error if something wrong
    HandleUserInfoResp(data []byte) (*WeChatUserP, error)

    //Get user info accoding to openid and language
    GetUserInfo(openId  string, lang string) (*WeChatUserP, error)
}

func (c * GlwapiConfig) UserInfoURL(token string, openId  string, lang string) string {
    return fmt.Sprintf("%s?access_token=%s&openid=%s&lang=%s", user_info_url, token, openId, lang)
}

func (c * GlwapiConfig) HandleUserInfoResp(data []byte) (*WeChatUserP, error) {
    if data == nil {
        return nil, errors.New("No data")
    }

    var user * WeChatUserP =&WeChatUserP{}
    if err := json.Unmarshal(data, user); err == nil {
        if user.Err > 0 {
            return nil, errors.New(fmt.Sprintf("Umaarshal failed: %d",user.Err))
        }
        return user, nil
    } else {
        return nil, err
    }

}


func (c * GlwapiConfig) GetUserInfo(openId  string, lang string) (*WeChatUserP, error) {
     if c.WeChat.IsServerTokenValid() {
          url := c.UserInfoURL(c.WeChat.Token, openId, lang)
          if resp, err := http.Get(url); err == nil {
               defer resp.Body.Close()
               if body, er := ioutil.ReadAll(resp.Body); er == nil {
                   return c.HandleUserInfoResp(body)
               } else {
                   return nil, er
               }
          } else {
              return nil, err
          }
     } else {
         return nil, errors.New("Server token is not valid")
     }
}

