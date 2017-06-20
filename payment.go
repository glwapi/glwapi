
package glwapi

import (
	"encoding/xml"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	ORDER_CREATE_URL string = "https://api.mch.weixin.qq.com/pay/unifiedorder"
)

// Generate Create order xml according to create order document
// URL: PAYMENT_URL_CO
func EncodeCreateOrderXml(o * WeChatOrder) (string, error) {
	if o == nil {
		return "", pe("Order is nil")
	}
	if o.WeChat == nil {
		return "", pe("WeChat Config is nil")
	}

	if o.Body == "" {
		return "", pe("No Body description")
	}
	if o.OrderNo == "" {
		return "", pe("No Order No description")
	}

	if len(o.OrderNo) > 32 {
		return "", pe("OrderNo lenght overed max limition(32)")
	}
	if o.NotifiedUrl == "" {
		return "", pe("No return URL")
	}
	return o.xorder().Encode()
}


// Generate query order xml according to query order document
// URL: PAYMENT_URL_QO
func EncodeQueryOrderXml(o * WeChatOrder) (string, error) {
	if o == nil {
		return "", pe("Order is nil")
	}
	if o.WeChat == nil {
		return "", pe("WeChat Config is nil")
	}
	return o.qorder().Encode()
}

// Generate close order xml according to close order document
// URL: PAYMENT_URL_DO
func EncodeCloseOrderXml(o * WeChatOrder) (string, error) {
	if o == nil {
		return "", pe("Order is nil")
	}
	if o.WeChat == nil {
		return "", pe("WeChat Config is nil")
	}
	return o.dorder().Encode()
}

// Generate request refund xml according to request refund document
// URL: PAYMENT_URL_RR
func EncodeRequestRefundXml(o * WeChatOrder) (string, error) {
	if o == nil {
		return "", pe("Order is nil")
	}
	if o.WeChat == nil {
		return "", pe("WeChat Config is nil")
	}
	return o.rrefund().Encode()
}

// Generate refund query xml according to refund query document
// URL: PAYMENT_URL_QRR
func EncodeQueryRefundXml(o * WeChatOrder) (string, error) {
	if o == nil {
		return "", pe("Order is nil")
	}
	if o.WeChat == nil {
		return "", pe("WeChat Config is nil")
	}
	return o.qrefund().Encode()
}


// Generate download bill xml according to bill download xml document
// URL: PAYMENT_URL_DB
// Notice not implement yet
func EncodeDownloadBillXml(o * WeChatOrder) (string, error) {
	return "", pe("Not implement yet")
}


// Generate respond payment result  xml
func EncodePaymentResultRespXml(code string, msg string) (string, error) {
	type Mar struct {
		XMLName    xml.Name  `xml:"xml"`
		Rcode   string `xml:"return_code"`
		RMsg    string `xml:"return_msg"`
		}
	var resp = Mar{Rcode : code, RMsg : msg}

	output, err := xml.MarshalIndent(resp , " ", "   ")
	if err != nil {
		return "", pe("Output xml failed " + err.Error())
	} else {
		return string(output[:len(output)]), nil
	}
}


// Decode order create response
// Populate WeChatOrder data automatically, if goes well
// otherwise return error
// Populate:
//      PrepayId
func DecodeCreateOrderResp(o * WeChatOrder, data []byte) error {
	var resp * xmlCreateOrderResp = &xmlCreateOrderResp{}
	fmt.Printf(" Create Order Resp :\n%s\n", data)
	if _, e := resp.Decode(data); e == nil {
		if resp.Code != CO_SUCCESS || resp.ECode != ""{
			return pe("Decode Order Failed:" + resp.Code + "  msg:" + resp.Msg+"  err_code:"+resp.ECode)
		}
		o.PrepayId = resp.PrepayId
		//FIXME check sign
		return nil
	} else {
		return pe("Failed to Decode Order:" + e.Error())
	}
}


// Decode order query response
// Populate WeChatOrder data automatically, if goes well
// otherwise return error
// Populate:
//      TradeState,OrderNo,TransId,TradeType,TotalFee,Cash,Attach
func DecodeQueryOrderResp(o * WeChatOrder, data []byte) error {
	if data == nil {
		return pe("No data for parse")
	}
	var oo queryOrderResult = queryOrderResult{}
	if _, e := oo.Decode(data); e != nil {
		return e
	} else {
		if oo.RCode.Val != QO_RES_SUCCESS || oo.ErrCode.Val !="" {
			return pe("Query Order filed : " + oo.RCode.Val+"  errCode:"+ oo.ErrCode.Val)
		}
		o.TradeState = oo.TradeState.Val
		if o.TradeState == TRADE_STATE_SUCCESS {
			populate(o, xmlOrderData(oo))
		}
		return nil
	}
}




// Decode payment result response
// Populate WeChatOrder data automatically, if goes well
// otherwise return error
// Populate:
//      TradeState,OrderNo,TransId,TradeType,TotalFee,Cash,Attach
func DecodePaymentResult(o * WeChatOrder, data []byte) error {
	if data == nil {
		return pe("No data for parse")
	}
	var oo paymentResult = paymentResult{}
	if _, e := oo.Decode(data); e != nil {
		return e
	} else {
		if oo.RCode.Val != QO_RES_SUCCESS || oo.ErrCode.Val !=""  {
			return pe("Query Order filed : " + oo.RCode.Val+"  errCode:"+ oo.ErrCode.Val)
		}
		o.TransId = oo.TransId.Val
		if o.TransId != "" {
			populate(o, xmlOrderData(oo))
		} else {
			return pe("Query Order filed : " + oo.RCode.Val+"  errCode:"+ oo.ErrCode.Val +" No TransId")
		}
		return nil
	}
}



//Decode close order response, do not populate any data to WeChatOrder
// return nil for success, otherwise failed
func DecodeCloseOrderResp(o * WeChatOrder, data []byte) error {
	var s = struct {
			Rcode    string `xml:"return_code"`
			RMsg     string `xml:"return_msg"`
			AppId    string `xml:"appid"`
			MchId    string `xml:"mch_id"`
			N        string `xml:"nonce_str"`
			Sign     string `xml:"sign"`
			ResCode  string `xml:"result_code"`
			ResMsg   string `xml:"result_msg"`
			}{}
	if err := xml.Unmarshal(data, &s); err == nil {
		if s.Rcode != DO_RES_SUCCESS || s.ResCode !=""  {
			return pe("Close Order failed:" + s.Rcode+"  errCode:"+ s.ResCode)
		}
		return nil
	} else {
		return pe("Decode failed: "+ err.Error())
	}
}


//Decode refund query response 
// return nil for success, otherwise failed
//Populate:
//    TransId OrderNo Cash CashType Coupon Refund
func DecodeQueryRefundResp(o * WeChatOrder, data []byte) error {
	var rr queryRefundRespXmlData = queryRefundRespXmlData{}
	if _, e := rr.Decode(data); e != nil {
		return e
	} else {
		o.TransId = rr.TransId.Val
		o.OrderNo = rr.MchOrderId.Val
		o.Cash,_ = strconv.Atoi(rr.Cash.Val)
		o.CashType = rr.CashType.Val
		o.Coupon   = rr.Coupon.Coupon
		o.Refund   = rr.Refund.Refund
		return nil
	}

}

//Decode refund application response 
// return nil for success, otherwise failed
//Populate:
//    TansId OrderNo Cash CashType Coupon Refund, TotalFee
func DecodeRequestRefundResp(o * WeChatOrder, data []byte) error {
	var rr requestRefundRespXmlData = requestRefundRespXmlData{}
	if _, e := rr.Decode(data); e != nil {
		return e
	} else {
		if rr.RCode.Val == RR_RES_SUCCESS  {
			populate(o, xmlOrderData(rr))
			if o.Refund == nil {
				o.Refund = &OrderRefund{}
			}
			o.Refund.RefundId = rr.WeChatRefundId.Val
			o.Refund.Fee,_ = strconv.Atoi(rr.RefundFee.Val)
		}
		return nil
	}
}




func populate(o * WeChatOrder, x xmlOrderData) {
	o.OrderNo    = x.MchOrderId.Val
	o.TransId    = x.TransId.Val
	o.TradeType  = x.TradeType.Val
	o.TotalFee,_ = strconv.Atoi(x.TotalFee.Val)
	o.Cash, _    = strconv.Atoi(x.TotalFee.Val)
	o.Attach     = x.Attach.Val
	o.Coupon     = x.Coupon.Coupon
	//TODO update trade time
	if o.User == nil {
		o.User = &WeChatUser{}
		o.User.OpenId = x.OpenId.Val
	}
}


func sign(v interface{}, key string) string {

	var signStr string
	rval := reflect.ValueOf(v)
	n := rval.NumField()
	names := make([]string, 0, n)
	xmlds := map[string]XmlD{}
	for idx:= 0;  idx < n; idx++ {
		v := rval.Field(idx)
		if v.Type().Name() == "XmlD" {
			xmd := v.Interface().(XmlD)
			xmlds[xmd.Key] = xmd
			names = append(names,xmd.Key)
		}
	}
	sort.Strings(names)
	for k, _ := range names {
		if (xmlds[names[k]].Opt == true && xmlds[names[k]].Val == "") || xmlds[names[k]].ExludeSign == true {
			continue
		}
		signStr +=  xmlds[names[k]].Key +"=" + xmlds[names[k]].Val +"&"
	}
	signStr = signStr +"key=" + key
	signStr = Md5(string(signStr))
	return strings.ToUpper(signStr)
}


func check_sign(v interface{}, key string) (bool, error) {
	return true, nil
}


func nxml(v interface{}) (string, error) {
	if output, err := xml.MarshalIndent(v, " ", "   "); err != nil {
		return "", pe("Output xml failed " + err.Error())
	} else {
		return string(output[:len(output)]), nil
	}
}

func pe(msg string) error {
	return &WeChatPaymentError{msg}
}

type WeChatPaymentError struct {
	Err	string
}

func (p * WeChatPaymentError) Error() string {
	return p.Err
}


type XmlEncode interface {

	// Generate XML according wechat spec
	Encode()  (string, error)
}

type XmlDecode interface {

	// Extract data from response
	Decode(xml []byte)  (interface{}, error)
}


type JSOrderSign struct {
	AppId	string
	TS	string
	N	string
	P	string
	S	string
	ST	string
}


type OrderRefund  struct {
	RefundId        string
	Fee             int
	CouponRefund	[]*PaymentRefund
}

type PaymentCoupon struct {
	CouponType	string
	CouponId	string
	CouponFee	int
}

type PaymentRefund struct {
	Fee		int
	Id		string
	Account		string
	RAccount	string
	RTime		*time.Time
	Status		string
	Channel		string
	Count		int
}

type WeChatOrder struct {
	WeChat		* WeChatConfig
	User		* WeChatUser
	Device		string
	Body		string
	Detail		string
	Attach		string
	OrderNo		string
	FeeType		string
	TotalFee	int
	Cash		int
	CashType	string
	Ip		string
	Ts		*time.Time
	Expires		*time.Time
	NotifiedUrl	string
	TradeType	string
	Pid		string
	Limit		string
	PrepayId	string
	TransId		string
	TradeState	string
	RefundFee	int
	EndTime		*time.Time
	Coupon		[]*PaymentCoupon
	Refund		*OrderRefund
}


func (o * WeChatOrder) JsSign() * JSOrderSign {
	var  jss *JSOrderSign = new(JSOrderSign)
	jss.AppId = o.WeChat.AppId
	jss.N     = R(32)
	jss.TS    = fmt.Sprintf("%d", time.Now().Unix())
	jss.ST     = "MD5"
	jss.P     = "prepay_id=" + o.PrepayId
	sinStr := ("appId=" +jss.AppId +"&nonceStr="+jss.N +"&package="+jss.P+"&signType="+jss.ST+"&timeStamp="+jss.TS+"&key="+o.WeChat.MchApiKey)
	jss.S     = strings.ToUpper(Md5(sinStr))
	fmt.Printf("SignStr: %s\n", sinStr)
	fmt.Printf("Sign: %s\n", jss.S)
	return jss
}






func (o * WeChatOrder) xorder() * xmlCreateOrder {
	var x xmlCreateOrder = xmlCreateOrder {
		o        :  o,
		AppId    :  XmlD{Key : "appid",            Val : o.WeChat.AppId},
		MchId    :  XmlD{Key : "mch_id",           Val : o.WeChat.MchId},
		Device   :  XmlD{Key : "device_info",      Val : o.Device,               Opt : true},
		Nonce    :  XmlD{Key : "nonce_str",        Val : R(32)},
		Sign     :  XmlD{Key : "sign",             Val : "",		 ExludeSign : true},
		SignType :  XmlD{Key : "sign_type",        Val : "MD5",             Opt : true},
		Body     :  XmlD{Key : "body",             Val : o.Body},
		Detail   :  XmlD{Key : "detail",           Val : o.Detail,               Opt : true},
		Attach   :  XmlD{Key : "attach",           Val : o.Attach,               Opt : true,  IsCDATA: true},
		OrderId  :  XmlD{Key : "out_trade_no",     Val : o.OrderNo},
		FeeType  :  XmlD{Key : "fee_type",         Val : o.FeeType,              Opt : true},
		Fee      :  XmlD{Key : "total_fee",        Val : ""},
		Ip       :  XmlD{Key : "spbill_create_ip", Val : o.Ip},
		Ts       :  XmlD{Key : "time_start",       Val : "",                     Opt : true},
		Expire   :  XmlD{Key : "time_expire",      Val : "",                     Opt : true},
		Tag      :  XmlD{Key : "goods_tag",        Val : "",                     Opt : true},
		Url      :  XmlD{Key : "notify_url",       Val : o.NotifiedUrl},
		Type     :  XmlD{Key : "trade_type",       Val : o.TradeType},
		Limit    :  XmlD{Key : "limit_pay",        Val : o.Limit,                Opt : true},
		OpenId   :  XmlD{Key : "openid",           Val : "",                     Opt : true},
	}

	if o.User != nil {
		x.OpenId.Val = o.User.OpenId
	}
	x.Fee.Val = fmt.Sprintf("%d", o.TotalFee)
	x.Ts.Val  = fmt.Sprintf("%d", o.Ts.Unix())
	if o.Expires != nil {
		x.Expire.Val  = fmt.Sprintf("%d", o.Expires.Unix())
	}
	return &x
}

func (o * WeChatOrder) qorder() * xmlQueryOrder {
	var x xmlQueryOrder = xmlQueryOrder {
		o        :  o,
		AppId    :  XmlD{Key : "appid",            Val : o.WeChat.AppId                   },
		MchId    :  XmlD{Key : "mch_id",           Val : o.WeChat.MchId                   },
		Nonce    :  XmlD{Key : "nonce_str",        Val : R(32)                            },
		Sign     :  XmlD{Key : "sign",             Val : "",		 ExludeSign : true},
		SignType :  XmlD{Key : "sign_type",        Val : "MD5",          Opt : true       },
		TransId  :  XmlD{Key : "out_trade_no",     Val : o.OrderNo,                       },
	}

	return &x;
}

func (o * WeChatOrder) dorder() * xmlCloseOrder {
	var x xmlCloseOrder = xmlCloseOrder {
		o        :  o,
		AppId    :  XmlD{Key : "appid",            Val : o.WeChat.AppId                   },
		MchId    :  XmlD{Key : "mch_id",           Val : o.WeChat.MchId                   },
		Nonce    :  XmlD{Key : "nonce_str",        Val : R(32)                            },
		Sign     :  XmlD{Key : "sign",             Val : "",		 ExludeSign : true},
		SignType :  XmlD{Key : "sign_type",        Val : "MD5",          Opt : true       },
		OrderId  :  XmlD{Key : "out_trade_no",     Val : o.OrderNo,                       },
	}

	return &x;
}

func (o * WeChatOrder) rrefund() * xmlRequestRefund {
	var x xmlRequestRefund =xmlRequestRefund {
		o        :  o,
		AppId    :  XmlD{Key : "appid",            Val : o.WeChat.AppId                   },
		MchId    :  XmlD{Key : "mch_id",           Val : o.WeChat.MchId                   },
		Nonce    :  XmlD{Key : "nonce_str",        Val : R(32)                            },
		Sign     :  XmlD{Key : "sign",             Val : "",		 ExludeSign : true},
		SignType :  XmlD{Key : "sign_type",        Val : "MD5",          Opt : true       },
		OrderId  :  XmlD{Key : "out_trade_no",     Val : o.OrderNo,                       },
		RefundId :  XmlD{Key : "out_refund_no",    Val : o.OrderNo,                       },
		Fee      :  XmlD{Key : "total_fee",        Val : "",                              },
		RFee     :  XmlD{Key : "refund_fee",       Val : string(o.RefundFee),             },
	}

	x.Fee.Val  = fmt.Sprintf("%d", o.TotalFee )
	x.RFee.Val = fmt.Sprintf("%d", o.RefundFee )

	return &x;
}


func (o * WeChatOrder) qrefund() * xmlQueryRefund {
	var x xmlQueryRefund =xmlQueryRefund {
		o        :  o,
		AppId    :  XmlD{Key : "appid",            Val : o.WeChat.AppId                   },
		MchId    :  XmlD{Key : "mch_id",           Val : o.WeChat.MchId                   },
		Nonce    :  XmlD{Key : "nonce_str",        Val : R(32)                            },
		Sign     :  XmlD{Key : "sign",             Val : "",		 ExludeSign : true},
		SignType :  XmlD{Key : "sign_type",        Val : "MD5",          Opt : true       },
		OrderId  :  XmlD{Key : "out_trade_no",     Val : o.OrderNo,                       },
	}

	return &x;
}



type xmlCreateOrder struct {
	o          * WeChatOrder    `xml:"-"`
	XMLName    xml.Name  `xml:"xml"`
	AppId      XmlD
	MchId      XmlD
	Device     XmlD
	Nonce      XmlD
	Sign       XmlD
	SignType   XmlD
	Body       XmlD
	Detail     XmlD
	Attach     XmlD
	OrderId    XmlD
	FeeType    XmlD
	Fee        XmlD
	Ip         XmlD
	Ts         XmlD
	Expire     XmlD
	Tag        XmlD
	Url        XmlD
	Type       XmlD
	Limit      XmlD
	OpenId     XmlD
}



func (x * xmlCreateOrder) Encode() (string, error) {
	if e := x.sign(); e != nil {
		return "", e
	}
	return nxml(x)
}

func (x * xmlCreateOrder) sign()  error {
	if  x.o == nil || x.o.WeChat == nil {
		return &WeChatPaymentError{"WeChat or User is nil "}
	}
	x.Sign.Val = sign(*x, x.o.WeChat.MchApiKey)
	return nil
}


type  xmlQueryOrder struct {
	o		* WeChatOrder	`xml:"-"`
	XMLName		xml.Name	`xml:"xml"`
	AppId		XmlD
	MchId		XmlD
	Nonce		XmlD
	Sign		XmlD
	SignType	XmlD
	TransId		XmlD
}

func (x * xmlQueryOrder) Encode() (string, error) {
	if e := x.sign(); e != nil {
		return "", e
	}
	if x.TransId.Val == "" {
		return "", pe("TransId is nil can not generate")
	}
	return nxml(x)
}

func (x * xmlQueryOrder) sign()  error {
	if  x.o == nil || x.o.WeChat == nil {
		return &WeChatPaymentError{"WeChat or User is nil "}
	}

	x.Sign.Val = sign(*x, x.o.WeChat.MchApiKey)
	return nil
}

type  xmlCloseOrder struct {
	o		* WeChatOrder	`xml:"-"`
	XMLName		xml.Name	`xml:"xml"`
	AppId		XmlD
	MchId		XmlD
	Nonce		XmlD
	Sign		XmlD
	SignType	XmlD
	OrderId		XmlD
}

func (x * xmlCloseOrder) Encode() (string, error) {
	if e := x.sign(); e != nil {
		return "", e
	}
	if x.OrderId.Val == "" {
		return "", pe("OrderId is nil can not generate")
	}
	return nxml(x)
}

func (x * xmlCloseOrder) sign()  error {
	if  x.o == nil || x.o.WeChat == nil {
		return &WeChatPaymentError{"WeChat or User is nil "}
	}

	x.Sign.Val = sign(*x, x.o.WeChat.Secret)
	return nil
}

type xmlRequestRefund struct {
	o		*WeChatOrder
	XMLName		xml.Name	`xml:"xml"`
	AppId		XmlD
	MchId		XmlD
	Nonce		XmlD
	Sign		XmlD
	SignType	XmlD
	OrderId		XmlD
	RefundId	XmlD
	Fee		XmlD
	RFee		XmlD
}

func (x *xmlRequestRefund) Encode() (string, error) {
	if e := x.sign(); e != nil {
		return "", e
	}
	if x.OrderId.Val == "" {
		return "", pe("OrderId is nil can not generate")
	}
	if x.RefundId.Val == "" {
		return "", pe("RefundId is nil can not generate")
	}
	return nxml(x)
}

func (x * xmlRequestRefund) sign()  error {
	if  x.o == nil || x.o.WeChat == nil {
		return &WeChatPaymentError{"WeChat or User is nil "}
	}

	x.Sign.Val = sign(*x, x.o.WeChat.Secret)
	return nil
}



type xmlQueryRefund struct {
	o		*WeChatOrder
	XMLName		xml.Name	`xml:"xml"`
	AppId		XmlD
	MchId		XmlD
	Nonce		XmlD
	Sign		XmlD
	SignType	XmlD
	OrderId		XmlD
}

func (x * xmlQueryRefund) Encode() (string, error) {
	if e := x.sign(); e != nil {
		return "", e
	}
	if x.OrderId.Val == "" {
		return "", pe("OrderId is nil can not generate")
	}
	return nxml(x)
}



func (x * xmlQueryRefund) sign()  error {
	if  x.o == nil || x.o.WeChat == nil {
		return &WeChatPaymentError{"WeChat or User is nil "}
	}

	x.Sign.Val = sign(*x, x.o.WeChat.Secret)
	return nil
}


/////////Response definition



type xmlCreateOrderResp struct {
	PrepayId       string `xml:"prepay_id"`
	Code           string `xml:"return_code"`
	Msg            string `xml:"return_msg"`
	Nonstr         string `xml:"nonce_str"`
	Sign           string `xml:"sign"`
	ECode          string `xml:"err_code"`
}


func (resp * xmlCreateOrderResp) Decode(data []byte) (interface{}, error) {
	if err := xml.Unmarshal(data, resp); err == nil {
		return resp, nil
	} else {
		return nil, err
	}
}

type xmlCloseOrderResp struct {
	OrderId        string `xml:"out_trade_no"`
	Code           string `xml:"return_code"`
	Msg            string `xml:"return_msg"`
	Nonstr         string `xml:"nonce_str"`
	Sign           string `xml:"sign"`
}


func (resp * xmlCloseOrderResp) Decode(data []byte) (interface{}, error) {
	if err := xml.Unmarshal(data, resp); err == nil {
		return resp, nil
	} else {
		return nil, err
	}
}



func getArrIdx(str string, prefix string) int {
	var n int
	var sf string = str[len(prefix) :]
	if idx := strings.LastIndex(sf, "_"); idx != -1 {
		n,_= strconv.Atoi(sf[:idx])
	} else {
		n,_= strconv.Atoi(sf)
	}
	return n
}

func createCouponRefunc(r * OrderRefund, n int) {
	for ;len(r.CouponRefund) <= n; {
		newc := new(PaymentRefund)
		r.CouponRefund = append(r.CouponRefund, newc)
	}
}


type xmlRefundData struct {
	Refund		*OrderRefund
}


func (f * xmlRefundData) UnmarshalXML(d * xml.Decoder, start xml.StartElement) error {
	name := start.Name.Local
	t, b := d.Token()
	if b != nil {
		return b
	}
	if f.Refund == nil {
		f.Refund = new(OrderRefund)
	}
	switch t.(type) {
	case xml.CharData:
		//FIXME optime performance
		if strings.Contains(name, "refund_status_") {
			n := getArrIdx(name, "refund_status_")
			createCouponRefunc(f.Refund, n)
			f.Refund.CouponRefund[n].Status = string(t.(xml.CharData))
		} else if strings.Contains(name, "refund_account_") {
			n := getArrIdx(name, "refund_account_")
			createCouponRefunc(f.Refund, n)
			f.Refund.CouponRefund[n].Account = string(t.(xml.CharData))
		} else if strings.Contains(name, "refund_recv_accout_") {
			n := getArrIdx(name, "refund_recv_accout_")
			createCouponRefunc(f.Refund, n)
			f.Refund.CouponRefund[n].RAccount = string(t.(xml.CharData))
		} else if strings.Contains(name, "refund_success_time_") {
			n := getArrIdx(name, "refund_success_time_")
			createCouponRefunc(f.Refund, n)
		} else if strings.Contains(name, "coupon_refund_fee_") {
			n := getArrIdx(name, "coupon_refund_fee_")
			createCouponRefunc(f.Refund, n)
			f.Refund.CouponRefund[n].Fee, _ = strconv.Atoi(string(t.(xml.CharData)))
		} else if strings.Contains(name, "coupon_refund_count_") {
			n := getArrIdx(name, "coupon_refund_count_")
			createCouponRefunc(f.Refund, n)
			f.Refund.CouponRefund[n].Count, _ = strconv.Atoi(string(t.(xml.CharData)))
		} else if strings.Contains(name, "coupon_refund_id_") {
			n := getArrIdx(name, "coupon_refund_Id_")
			createCouponRefunc(f.Refund, n)
			f.Refund.CouponRefund[n].Id = string(t.(xml.CharData))
		}

	}
	return nil
}


type xmlCouponData struct {
	Coupon		[]*PaymentCoupon
}

func (c * xmlCouponData) UnmarshalXML(d * xml.Decoder, start xml.StartElement) error {
	name := start.Name.Local
	t, b := d.Token()
	if b != nil {
		return b
	}
	switch t.(type) {
	case xml.CharData:
		//FIXME should only find one time and return flag
		if strings.Contains(name, "coupon_fee_") {
			idx := strings.LastIndex(name, "_")
			n,_ := strconv.Atoi(name[idx+1:])
			for ; len(c.Coupon) <= n; {
				newc := new(PaymentCoupon)
				c.Coupon = append(c.Coupon, newc)
			}
			strfee,_ := t.(xml.CharData)
			c.Coupon[n].CouponFee,_ = strconv.Atoi(string(strfee))
		} else if strings.Contains(name, "coupon_id_") {
			idx := strings.LastIndex(name, "_")
			n,_ := strconv.Atoi(name[idx+1:])
			for ; len(c.Coupon) <= n; {
				newc := new(PaymentCoupon)
				c.Coupon = append(c.Coupon, newc)
			}
			c.Coupon[n].CouponId = string(t.(xml.CharData))
		} else if strings.Contains(name, "coupon_type_") {
			idx := strings.LastIndex(name, "_")
			n,_ := strconv.Atoi(name[idx+1:])
			for ; len(c.Coupon) <= n; {
				newc := new(PaymentCoupon)
				c.Coupon = append(c.Coupon, newc)
			}
			c.Coupon[n].CouponType = string(t.(xml.CharData))
		}
	default:
		return pe("Unknown type")
	}
	return nil
}

type  xmlOrderData struct {
	RCode		XmlD		`xml:"result_code"`
	ErrCode		XmlD		`xml:"err_code"`
	ErrCodeMsg	XmlD		`xml:"err_code_des"`
	AppId		XmlD		`xml:"appid"`
	MchId		XmlD		`xml:"mch_id"`
	Device		XmlD		`xml:"device_info"`
	Nonce		XmlD		`xml:"nonce_str"`
	Sign		XmlD		`xml:"sign"`
	SignType	XmlD		`xml:"sign_type"`
	OpenId		XmlD		`xml:"openid"`
	TradeState      XmlD            `xml:"trade_state"`
	TradeType       XmlD            `xml:"trade_type"`
	BankType        XmlD            `xml:"bank_type"`
	Attach          XmlD            `xml:"attach"`
	TotalFee	XmlD		`xml:"total_fee"`
	SettlementFee	XmlD		`xml:"settlement_total_fee"`
	FeeType		XmlD		`xml:"fee_type"`
	Cash		XmlD		`xml:"cash_fee"`
	CashType	XmlD		`xml:"cash_fee_type"`
	TransId		XmlD		`xml:"transaction_id"`
	MchOrderId	XmlD		`xml:"out_trade_no"`
	WeChatRefundId	XmlD		`xml:"refund_id"`
	RefundId	XmlD		`xml:"out_refund_no"`
	RefundFee	XmlD		`xml:"refund_fee"`
	SettlementRFee	XmlD		`xml:"settlement_refund_fee"`
	Coupon          xmlCouponData
	//TODO refund 
	Refund          xmlRefundData
	l		int
}

func (c * xmlOrderData) UnmarshalXML(d * xml.Decoder, start xml.StartElement) error {
	var level int = 1
	var t xml.Token
	var b error
	var fm map[string]string = make(map[string]string)
	val := reflect.ValueOf(c).Elem()
	typ := val.Type()
	for i, n := 0, typ.NumField(); i < n; i++ {
		f := typ.Field(i)
		if f.Type.String() == "glwapi.XmlD" {
			fm[f.Tag.Get("xml")] = f.Name
		}
	}
	for t, b = d.Token(); t != nil && b == nil; t,b = d.Token() {
		switch t.(type) {
		case xml.StartElement:
			level +=1
			xs,_ := t.(xml.StartElement)
			if strings.HasPrefix(xs.Name.Local, "coupon_refund_") == true {
				c.Refund.UnmarshalXML(d, xs)
			} else if strings.HasPrefix(xs.Name.Local, "refund_status_") == true {
				c.Refund.UnmarshalXML(d, xs)
			} else if strings.HasPrefix(xs.Name.Local, "refund_recv_accout_") == true {
				c.Refund.UnmarshalXML(d, xs)
			} else if strings.HasPrefix(xs.Name.Local, "refund_account_") == true {
				c.Refund.UnmarshalXML(d, xs)
			} else if strings.HasPrefix(xs.Name.Local, "coupon_") == true {
				c.Coupon.UnmarshalXML(d, xs)
			} else {
				fname, b := fm[xs.Name.Local]
				if b == true {
					xmdptr := &XmlD{}
					xmdptr.UnmarshalXML(d, xs)
					val.FieldByName(fname).Set(reflect.ValueOf(*xmdptr))
				} else {
					//TODO log no such name in object
				}
			}
			break
		case xml.EndElement:
			level -= 1
			if level == 0 {
				break
			}
		}
	}
	return nil

}

type paymentResult xmlOrderData


func (pr * paymentResult) Decode(data []byte) (interface{}, error) {
	var x = xmlOrderData(*pr)
	if err := xml.Unmarshal(data, &x); err == nil {
		*pr = paymentResult(x);
		//TODO check sign
		return pr, nil
	} else {
		return nil, err
	}

}


type queryOrderResult  xmlOrderData

func (pr * queryOrderResult) Decode(data []byte) (interface{}, error) {
	var x = xmlOrderData(*pr)
	if err := xml.Unmarshal(data, &x); err == nil {
		*pr = queryOrderResult(x);
		return pr, nil
	} else {
		return nil, err
	}
}


type  requestRefundRespXmlData xmlOrderData

func (rr * requestRefundRespXmlData) Decode(data []byte) (interface{}, error) {
	var x = xmlOrderData(*rr)
	if err := xml.Unmarshal(data, &x); err == nil {
		*rr = requestRefundRespXmlData(x);
		return rr, nil
	} else {
		return nil, err
	}
}

type  queryRefundRespXmlData xmlOrderData

func (qr * queryRefundRespXmlData) Decode(data []byte) (interface{}, error) {
	var x = xmlOrderData(*qr)
	if err := xml.Unmarshal(data, &x); err == nil {
		*qr = queryRefundRespXmlData(x);
		return qr, nil
	} else {
		return nil, err
	}
}





type OrderApi interface {
	CreateOrder1(user * WeChatUser, orderId string, desc string, detail string, additional string, ip string, url string, fee int) (*WeChatOrder, error)
}

func (w * WeChatConfig) CreateOrder1(user * WeChatUser, orderId string, desc string, detail string, additional string, ip string, url string, fee int) (*WeChatOrder, error) {
	var ts = time.Now()
	var o * WeChatOrder = &WeChatOrder{
			WeChat        : w,
			Body          : desc,
			Detail        : detail,
			Attach        : additional,
			OrderNo       : orderId,
			TotalFee      : fee,
			Ip            : ip,
			Ts            : &ts,
			NotifiedUrl   : url,
			TradeType     : "JSAPI",
			Limit         : "",
			User          : user,
			}

	return o, nil
}
