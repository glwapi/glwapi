

package glwapi

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "mime/multipart"
    "net/http"
    "os"
    "strings"
)


const (
    upload_temporary_url         = "https://api.weixin.qq.com/cgi-bin/media/upload"
    upload_permanent_art_url     = "https://api.weixin.qq.com/cgi-bin/material/add_news"
    upload_permanent_img_url     = "https://api.weixin.qq.com/cgi-bin/media/uploadimg"
    upload_permanent_oth_url     = "https://api.weixin.qq.com/cgi-bin/material/add_material"
    delete_permanent_url         = "https://api.weixin.qq.com/cgi-bin/material/del_material"
    get_temporary_url            = "https://api.weixin.qq.com/cgi-bin/media/get"
    get_permanent_url            = "https://api.weixin.qq.com/cgi-bin/material/get_material"
    get_res_statistic_url        = "https://api.weixin.qq.com/cgi-bin/material/get_materialcount"
)


const (
   R_TYPE_UNKOWN = -1
   R_TYPE_IMAGE  = iota
   R_TYPE_AUDIO
   R_TYPE_VIDEO
   R_TYPE_THUMB
   // only for permanent resource
   R_TYPE_ARTICLE
   R_TYPE_ARTICLE_IMAGE
)


const (
   ERR_R_CODE_PAR             = 50000
   ERR_R_CODE_NET_IO          = 50001
   ERR_R_CODE_READ_RESP       = 50002
   ERR_R_CODE_UNMARSHAL_RESP  = 50003
   ERR_R_CODE_PARSE_R_TYPE    = 50004
)

// Resource intrface
type IResource interface {
    /*
     * return resource type, must be one of 
     *   R_TYPE_IMAGE
     *   R_TYPE_AUDIO
     *   R_TYPE_VIDEO
     *   R_TYPE_THUMB
     *   R_TYPE_ARTICLE  only for permanent resource
     */
    T()     int

    // data of conent
    C()     []byte

    //return media id in wechat server
    MId()   string

    //return media url in wechat server
    MUrl()   string

    //set media url in wechat server
    SetMUrl(u string)

    //set media id of wechat server
    SetMId(id string)
}

//raw resource of wechat respond
type RawResource struct {

    Type        int

    MediaId     string

    MediaUrl    string

    Data        []byte

    ContentType string
}

func (r * RawResource) T() int {
    return r.Type
}

func (r * RawResource) C()  []byte {
    return r.Data
}

func (r * RawResource) MId() string {
    return r.MediaId
}

func (r * RawResource) MUrl() string {
    return r.MediaUrl
}

func (r * RawResource) SetMUrl(u string) {
    r.MediaUrl = u
}

func (r * RawResource) SetMId(id string) {
    r.MediaId = id
}


type FileResource struct {

    Type        int

    MediaId     string

    MediaUrl    string

    Data        []byte

}

func (f * FileResource) SetFilePath(path string) error {

    file, err := os.Open(path)
    if err != nil {
        return err
    }

    defer file.Close()
    fileContents, err := ioutil.ReadAll(file)
    if err != nil {
        return err
    }

    fi, err := file.Stat()
    if err != nil {
        return err
    }

    body := new(bytes.Buffer)
    writer := multipart.NewWriter(body)
    part, err := writer.CreateFormFile("", fi.Name())
    if err != nil {
        return err
    }
    part.Write(fileContents)

    err = writer.Close()
    if err != nil {
        return err
    }

    f.Data = body.Bytes()
    return nil
}

func (r * FileResource) T() int {
    return r.Type
}

func (r * FileResource) C()  []byte {
    return r.Data
}

func (r * FileResource) MId() string {
    return r.MediaId
}

func (r * FileResource) MUrl() string {
    return r.MediaUrl
}

func (r * FileResource) SetMUrl(u string) {
    r.MediaUrl = u
}

func (r * FileResource) SetMId(id string) {
    r.MediaId = id
}


type ResourceStatistic struct {

   ErrCode        int       `json:"errcode"`
   ErrMsg         string    `json:"errmsg"`

   AudioCount     int       `json:"voice_count"`
   VideoCount     int       `json:"video_count"`
   ImageCount     int       `json:"image_count"`
   NewsCount      int       `json:"news_count"`
}


type ResourceError struct {
     ErrCode       int
     ErrMsg        string
}


type WechatResourceRet struct {
    ErrCode       int       `json:"errcode"`
    ErrMsg        int       `json:"errmsg"`
    Type          string    `json:"type"`
    MediaId       string    `json:"media_id"`
    TimeStamp     int       `json:"created_at"`
    MediaUrl      string    `json:"url"`
    VideoUrl      string    `json:"video_url"`
    DownUrl       string    `json:"down_url"`
}


func (r ResourceError) Error() string {
    return fmt.Sprintf("Err :%d,  Msg: %s", r.ErrCode, r.ErrMsg)
}


/*
 * Resource API interface:<br />
 *
 */
type ResourceManager  interface {

    //Upload new temporary resource
    // After upload successfully, will set media id and set media url if exist
    UploadTR(r IResource) error

    //Retrieve temporary resource according media id
    // return new resource which in wechat server
    GetTR(r IResource)  (IResource, error)

    //Upload new permanent resource
    // After upload successfully, will set media id and set media url if exist
    UploadPR(r IResource) error

    //Retrieve permanent resource according media id
    // return new resource which in wechat server
    GetPR(r IResource)  (IResource, error)

    //Delete permanent resource according to media id
    // Return nil for success, otherwise failed
    DelPR(r IResource) error

    //Get resources statistic information
    GetResourceStatistic() (*ResourceStatistic, error)
}



func (c * GlwapiConfig) UploadTR(r IResource) (error) {
    if c == nil || c.WeChat == nil ||  !c.WeChat.IsServerTokenValid() {
         return ResourceError{ErrCode :  ERR_R_CODE_PAR, ErrMsg : "Par error"}
    }

    var ty string
    switch (r.T ()) {
    case R_TYPE_IMAGE:
        ty = "image"
        break
    case R_TYPE_AUDIO:
        ty = "voice"
        break
    case R_TYPE_VIDEO:
        ty = "video"
        break
    case R_TYPE_THUMB:
        ty = "thumb"
        break
    default:
        return ResourceError{ErrCode :  ERR_R_CODE_NET_IO, ErrMsg : "Unsupport type for temporary"}
    }

    http_url := fmt.Sprintf("%s?access_token=%stype=%s", upload_temporary_url, c.WeChat.Token, ty)
    ret, er := c.send_request(http_url, r)
    if ret.ErrCode ==  0 {
         r.SetMId(ret.MediaId)
         r.SetMUrl(ret.MediaUrl)
         return nil
    } else {
         return er
    }
}


func (c * GlwapiConfig) GetTR(r IResource) (IResource, error) {
    if c == nil || c.WeChat == nil ||  !c.WeChat.IsServerTokenValid() {
         return nil, ResourceError{ErrCode :  ERR_R_CODE_PAR, ErrMsg : "Par error"}
    }
    http_url := fmt.Sprintf("%s?access_token=%s&media_id=%", get_temporary_url, c.WeChat.Token, r.MId())
    resp, err := http.Get(http_url)
    if err != nil {
         return nil, ResourceError{ErrCode :  ERR_R_CODE_NET_IO, ErrMsg : fmt.Sprintf("request fialed :%s", err)}
    }

    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
         return nil, ResourceError{ErrCode :  ERR_R_CODE_NET_IO, ErrMsg : fmt.Sprintf("request fialed :%d", resp.StatusCode)}
    }
    rr := new(RawResource)
    if data, err := ioutil.ReadAll(resp.Body); err == nil {
        rr.Data = data
        rr.ContentType = resp.Header.Get("Content-Type")
    } else {
        return nil, ResourceError{ErrCode :  ERR_R_CODE_READ_RESP, ErrMsg : fmt.Sprintf("retrieve data failed :%s", err)}
    }
    return rr, nil
}


func (c * GlwapiConfig) UploadPR(r IResource) (error) {
    if c == nil || c.WeChat == nil ||  !c.WeChat.IsServerTokenValid() {
         return ResourceError{ErrCode :  ERR_R_CODE_PAR, ErrMsg : "Par error"}
    }

    var api_url string
    var ty       string
    switch (r.T ()) {
    case R_TYPE_IMAGE:
        api_url = upload_permanent_oth_url
        ty      = "image"
        break
    case R_TYPE_AUDIO:
        api_url = upload_permanent_oth_url
        ty      = "voice"
        break
    case R_TYPE_VIDEO:
        api_url = upload_permanent_oth_url
        ty      = "video"
        break
    case R_TYPE_THUMB:
        api_url = upload_permanent_oth_url
        ty      = "thumb"
        break
    case R_TYPE_ARTICLE:
        api_url = upload_permanent_art_url
        ty      = "art"
    case R_TYPE_ARTICLE_IMAGE:
        api_url = upload_permanent_img_url
        ty      = "art_image"
    default:
        return ResourceError{ErrCode :  ERR_R_CODE_NET_IO, ErrMsg : "Unsupport type for temporary"}
    }

    http_url := fmt.Sprintf("%s?access_token=%stype=%s", api_url, c.WeChat.Token, ty)

    ret, er := c.send_request(http_url, r)
    if ret.ErrCode ==  0 {
         r.SetMId(ret.MediaId)
         r.SetMUrl(ret.MediaUrl)
         return nil
    } else {
         return er
    }
}

func (c * GlwapiConfig) GetPR(r IResource) (IResource, error) {
    if c == nil || c.WeChat == nil ||  !c.WeChat.IsServerTokenValid() {
         return nil, ResourceError{ErrCode :  ERR_R_CODE_PAR, ErrMsg : "Par error"}
    }
    http_url := fmt.Sprintf("%s?access_token=%s&media_id=%", get_permanent_url, c.WeChat.Token, r.MId())
    resp, err := http.Get(http_url)
    if err != nil {
         return nil, ResourceError{ErrCode :  ERR_R_CODE_NET_IO, ErrMsg : fmt.Sprintf("request fialed :%s", err)}
    }

    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
         return nil, ResourceError{ErrCode :  ERR_R_CODE_NET_IO, ErrMsg : fmt.Sprintf("request fialed :%d", resp.StatusCode)}
    }

    rr := new(RawResource)
    if data, err := ioutil.ReadAll(resp.Body); err == nil {
        rr.Data = data
        rr.ContentType = resp.Header.Get("Content-Type")
    } else {
        return nil, ResourceError{ErrCode :  ERR_R_CODE_READ_RESP, ErrMsg : fmt.Sprintf("retrieve data failed :%s", err)}
    }
    return rr, nil
}



func (c * GlwapiConfig) DeletePR(r IResource) (error) {
    if c == nil || c.WeChat == nil ||  !c.WeChat.IsServerTokenValid() {
         return ResourceError{ErrCode :  ERR_R_CODE_PAR, ErrMsg : "Par error"}
    }
    ret, er := c.send_request(delete_permanent_url, r)
    if ret.ErrCode ==  0 {
         return nil
    } else {
         return er
    }
}


func (c * GlwapiConfig)  GetResourceStatistic() (*ResourceStatistic, error) {
    if c == nil || c.WeChat == nil ||  !c.WeChat.IsServerTokenValid() {
         return nil, ResourceError{ErrCode :  ERR_R_CODE_PAR, ErrMsg : "Par error"}
    }


    http_url := fmt.Sprintf("%s?access_token=%s", get_res_statistic_url, c.WeChat.Token)
    resp, err := http.Get(http_url)
    if err != nil {
         return nil, ResourceError{ErrCode :  ERR_R_CODE_NET_IO, ErrMsg : fmt.Sprintf("request fialed :%s", err)}
    }

    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
         return nil, ResourceError{ErrCode :  ERR_R_CODE_NET_IO, ErrMsg : fmt.Sprintf("request fialed :%d", resp.StatusCode)}
    }
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
         return nil, ResourceError{ErrCode :  ERR_R_CODE_NET_IO, ErrMsg : fmt.Sprintf("read io  fialed :%s", err)}
    }

    ret := ResourceStatistic{}
    if er := json.Unmarshal(body, &ret); er == nil {
        if ret.ErrCode != 0 {
            return nil, ResourceError{ErrCode :  ret.ErrCode, ErrMsg : ret.ErrMsg}
        } else {
            return &ret, nil
        }
    } else {
        return nil, ResourceError{ErrCode :  ERR_R_CODE_UNMARSHAL_RESP, ErrMsg : fmt.Sprintf("unmarshal response failed  %s:\n%s", err, body)}
    }

}

func (c * GlwapiConfig) send_request(uri string, r IResource) (*WechatResourceRet, error) {

    client := &http.Client{}
    req, err  := http.NewRequest("POST", uri, bytes.NewReader(r.C()))
    if err != nil {
        return nil, ResourceError{ErrCode :  ERR_R_CODE_NET_IO, ErrMsg : fmt.Sprintf("upload failed can not construct request %s", err)}
    }


    resp, err := client.Do(req)
    if err != nil {
        return nil, ResourceError{ErrCode :  ERR_R_CODE_NET_IO, ErrMsg : fmt.Sprintf("upload failed status code :%d", resp.StatusCode)}
    }

    defer resp.Body.Close()
    body, er := ioutil.ReadAll(resp.Body)
    if er != nil {
        return nil, ResourceError{ErrCode :  ERR_R_CODE_READ_RESP, ErrMsg : fmt.Sprintf("upload failed status code :%d", resp.StatusCode)}
    }


    ret := WechatResourceRet{}
    if er = json.Unmarshal(body, &ret); er != nil {
        return nil, ResourceError{ErrCode :  ERR_R_CODE_UNMARSHAL_RESP, ErrMsg : fmt.Sprintf("unmarshal response failed  %s:\n%s", err, body)}
    }

    if ret.ErrCode ==  0 {
        return &ret, nil
    } else {
        return nil, ResourceError{ErrCode :  ret.ErrCode, ErrMsg : fmt.Sprintf("wechat return error %s", ret.ErrMsg)}
    }
}


func TestRawResource(rr * RawResource) (int, error) {
    ret := WechatResourceRet{}
    ct := strings.ToLower(rr.ContentType)

    if strings.Contains(ct, "image/") {
        return R_TYPE_IMAGE, nil
    } else if strings.Contains(ct, "video/") {
        return R_TYPE_VIDEO, nil
    } else if strings.Contains(ct, "audio/") {
        return R_TYPE_AUDIO, nil
    } else if strings.Contains(ct, "application/json") {
        if er := json.Unmarshal(rr.Data, &ret); er == nil {
            if ret.ErrCode != 0 {
                return R_TYPE_UNKOWN, ResourceError{ErrCode :  ret.ErrCode, ErrMsg : fmt.Sprintf("wechat return error %s", ret.ErrMsg)}
            } else if ret.DownUrl != "" || ret.VideoUrl != "" {
                return R_TYPE_VIDEO, nil
	    }
        }
    } else {
         return R_TYPE_UNKOWN,  ResourceError{ErrCode :  ERR_R_CODE_PARSE_R_TYPE, ErrMsg : fmt.Sprintf("unknow content type %s", ct)}
    }


    return R_TYPE_UNKOWN, nil
}

