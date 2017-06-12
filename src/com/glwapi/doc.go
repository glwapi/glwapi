/*
A tool wechat API of GO implement
This Tool Doesn't contain http server part
Example:
Initalize glwapi configuration:
    g := glwapi.D()
    ret := g.Init(data.Appid, data.Secret, data.MchId, data.MchApiKey)
    if ret != true {
         log.Fatal(" initalize glwapi failed")
    }

Get Token:

 g := glwapi.D()
 if g.IsInit() == false {
       LE("glwapi doesn't inital yet")
       return
 }
 token_url := g.TokenURL()
 r, err :=  http.Get(token_url)
 if err == nil {
       defer r.Body.Close()
       rdata, _ := ioutil.ReadAll(r.Body)
       _, e := g.HandleTokenAuthResp([]byte(rdata))
       if e != nil {
            LE("Handle auth response error:%s", e)
       } else {
            LI("Get auth token successfully: %s\n", g)
       }
 }

   
*/
package glwapi
