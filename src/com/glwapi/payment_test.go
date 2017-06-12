
package glwapi_test

import (
	"testing"
	"com/glwapi"
//	"bytes"
//	"strings"
//	"fmt"
	"encoding/xml"
)


func TestWeChatOrderXml(t * testing.T) {
	g := glwapi.D()
	g.WeChat.AppId ="1234"
	g.WeChat.Secret ="1234"
	g.WeChat.MchId ="1390613402"
	o, _ :=g.WeChat.CreateOrder1(nil, "1235", "aa","bb", "cc", "192.168.1.1", "http://test", 1234)
	s, e :=glwapi.EncodeCreateOrderXml(o)
	if e != nil {
		t.Errorf("%s\n", e.Error())
	} else {
		t.Logf("Test ok \n xml ===> \n%s\n", s)
	}
}

func TestQueryOrderXml(t * testing.T) {
	g := glwapi.D()
	g.WeChat.AppId ="1234"
	g.WeChat.Secret ="1234"
	g.WeChat.MchId ="1390613402"
	o, _ :=g.WeChat.CreateOrder1(nil, "", "aa","bb", "cc", "192.168.1.1", "http://test", 1234)
	s, e :=glwapi.EncodeQueryOrderXml(o)
	if e == nil {
		t.Errorf("Should Ok, because no OrderNo")
	} else {
		t.Logf("Test ok \n xml ===> \n%s\n", s)
	}

	o.OrderNo = "12344444"
	s, e =glwapi.EncodeCloseOrderXml(o)
	if e != nil {
		t.Errorf("Test failed : %s", e.Error())
	} else {
		t.Logf("Test ok for normal \n xml ===> \n%s\n", s)
	}
}


type xmlData struct {
	AppId   string  `xml:"appid"`
	MchId   string  `xml:"mch_id"`
	Nonce   string  `xml:"nonce_str"`
	OrderNo   string  `xml:"out_trade_no"`
}

func TestCloseOrderXml(t * testing.T) {
	g := glwapi.D()
	g.WeChat.AppId ="1234"
	g.WeChat.Secret ="1234"
	g.WeChat.MchId ="1390613402"
	o, _ :=g.WeChat.CreateOrder1(nil, "", "aa","bb", "cc", "192.168.1.1", "http://test", 1234)
	s, e :=glwapi.EncodeCloseOrderXml(o)
	if e == nil {
		t.Logf("xml ===> \n%s\n", s)
		t.Errorf("Should no ok, because no OrdeNo")
	} else {
		t.Logf("Test ok  for No order no")
	}

	o.OrderNo = "12344444"
	s, e =glwapi.EncodeCloseOrderXml(o)
	if e != nil {
		t.Errorf("Test failed : %s", e.Error())
	} else {
		xd	 := xmlData{}
		xml.Unmarshal([]byte(s), &xd)
		if xd.OrderNo != "12344444" {
			t.Errorf("Test failed for orderno mismatch :%s\n   %s\n", s, xd.OrderNo)
			return
		}	
		t.Logf("Test ok for normal ")
	}
}

func TestRequestRefundXml(t * testing.T) {
	g := glwapi.D()
	g.WeChat.AppId ="1234"
	g.WeChat.Secret ="1234"
	g.WeChat.MchId ="1390613402"
	o, _ :=g.WeChat.CreateOrder1(nil, "", "aa","bb", "cc", "192.168.1.1", "http://test", 1234)
	s, e :=glwapi.EncodeRequestRefundXml(o)
	if e == nil {
		t.Logf("xml ===> \n%s\n", s)
		t.Errorf("Should no ok, because no OrdeNo")
	} else {
		t.Logf("Test ok  for No orderno\n xml ===> \n%s\n  %s", s, e.Error())
	}

	o.OrderNo = "12344444"
	s, e =glwapi.EncodeRequestRefundXml(o)
	if e != nil {
		t.Errorf("Test failed : %s", e.Error())
	} else {
		xd	 := xmlData{}
		xml.Unmarshal([]byte(s), &xd)
		if xd.OrderNo != "12344444" {
			t.Errorf("Test failed for orderno mismatch :%s\n   %s\n", s, xd.OrderNo)
			return
		}
		t.Logf("Test ok for normal ")
	}
}

func TestQueryRefundXml(t * testing.T) {
	g := glwapi.D()
	g.WeChat.AppId ="1234"
	g.WeChat.Secret ="1234"
	g.WeChat.MchId ="1390613402"
	o, _ :=g.WeChat.CreateOrder1(nil, "", "aa","bb", "cc", "192.168.1.1", "http://test", 1234)
	s, e :=glwapi.EncodeQueryRefundXml(o)
	if e == nil {
		t.Logf("xml ===> \n%s\n", s)
		t.Errorf("Should no ok, because no OrdeNo")
	} else {
		t.Logf("Test ok  for No order no\n xml ===> \n%s\n", s)
	}

	o.OrderNo = "12344444"
	s, e =glwapi.EncodeQueryRefundXml(o)
	if e != nil {
		t.Errorf("Test failed : %s", e.Error())
	} else {
		t.Logf("Test ok for normal \n xml ===> \n%s\n", s)
	}
}



func TestDecodePaymentResult(t * testing.T) {
	o := glwapi.WeChatOrder{}	
	testD := "<xml>"
   	testD +="<return_code><![CDATA[SUCCESS]]></return_code>"
    	testD +="  <return_msg><![CDATA[OK]]></return_msg>"
     	testD +=" <appid><![CDATA[wx2421b1c4370ec43b]]></appid>"
      	testD +="<mch_id><![CDATA[10000100]]></mch_id>"
      	testD +="<device_info><![CDATA[1000]]></device_info>"
      	testD +="<nonce_str><![CDATA[TN55wO9Pba5yENl8]]></nonce_str>"
      	testD +="<sign><![CDATA[BDF0099C15FF7BC6B1585FBB110AB635]]></sign>"
      	testD +="<result_code><![CDATA[SUCCESS]]></result_code>"
      	testD +="<openid><![CDATA[oUpF8uN95-Ptaags6E_roPHg7AG0]]></openid>"
      	testD +="<is_subscribe><![CDATA[Y]]></is_subscribe>"
      	testD +="<trade_type><![CDATA[MICROPAY]]></trade_type>"
      	testD +="<bank_type><![CDATA[CCB_DEBIT]]></bank_type>"
      	testD +="<total_fee>1</total_fee>"
      	testD +="<fee_type><![CDATA[CNY]]></fee_type>"
      	testD +="<transaction_id><![CDATA[1008450740201411110005820873]]></transaction_id>"
      	testD +="<out_trade_no><![CDATA[1415757673]]></out_trade_no>"
      	testD +="<attach><![CDATA[订单额外描述]]></attach>"
      	testD +="<time_end><![CDATA[20141111170043]]></time_end>"
      	testD +="<trade_state><![CDATA[SUCCESS]]></trade_state>"
   	testD +="<coupon_fee_0>0</coupon_fee_0> "
   	testD +="<coupon_fee_1><![CDATA[1]]></coupon_fee_1> "
   	testD +="<coupon_type_1>CASH</coupon_type_1> "
   	testD +="</xml>"
	t.Logf("%s", testD)
	e := glwapi.DecodePaymentResult(&o, []byte(testD))
	if e == nil {
		if o.OrderNo != "1415757673" || o.TradeState != glwapi.TRADE_STATE_SUCCESS  || len(o.Coupon) != 2 || o.Coupon[1].CouponType !="CASH" || o.Coupon[1].CouponFee != 1 || o.TransId != "1008450740201411110005820873"{
			t.Errorf("Extract Data mismatch %s", o)
		} else {
			t.Logf("Test ok")
		}
	} else {
		t.Errorf("Test Failed: "+ e.Error())
	}
}



func TestDecodeQueryOrderResp(t * testing.T) {
	o := glwapi.WeChatOrder{}	
	testD := "<xml>"
   	testD +="<return_code><![CDATA[SUCCESS]]></return_code>"
    	testD +="  <return_msg><![CDATA[OK]]></return_msg>"
     	testD +=" <appid><![CDATA[wx2421b1c4370ec43b]]></appid>"
      	testD +="<mch_id><![CDATA[10000100]]></mch_id>"
      	testD +="<device_info><![CDATA[1000]]></device_info>"
      	testD +="<nonce_str><![CDATA[TN55wO9Pba5yENl8]]></nonce_str>"
      	testD +="<sign><![CDATA[BDF0099C15FF7BC6B1585FBB110AB635]]></sign>"
      	testD +="<result_code><![CDATA[SUCCESS]]></result_code>"
      	testD +="<openid><![CDATA[oUpF8uN95-Ptaags6E_roPHg7AG0]]></openid>"
      	testD +="<is_subscribe><![CDATA[Y]]></is_subscribe>"
      	testD +="<trade_type><![CDATA[MICROPAY]]></trade_type>"
      	testD +="<bank_type><![CDATA[CCB_DEBIT]]></bank_type>"
      	testD +="<total_fee>1</total_fee>"
      	testD +="<fee_type><![CDATA[CNY]]></fee_type>"
      	testD +="<transaction_id><![CDATA[1008450740201411110005820873]]></transaction_id>"
      	testD +="<out_trade_no><![CDATA[1415757673]]></out_trade_no>"
      	testD +="<attach><![CDATA[订单额外描述]]></attach>"
      	testD +="<time_end><![CDATA[20141111170043]]></time_end>"
      	testD +="<trade_state><![CDATA[SUCCESS]]></trade_state>"
   	testD +="<coupon_fee_0>0</coupon_fee_0> "
   	testD +="<coupon_fee_1><![CDATA[1]]></coupon_fee_1> "
   	testD +="<coupon_id_1><![CDATA[1234]]></coupon_id_1> "
   	testD +="<coupon_id_0><![CDATA[01234]]></coupon_id_0> "
   	testD +="<coupon_type_1>CASH</coupon_type_1> "
   	testD +="</xml>"
	e := glwapi.DecodeQueryOrderResp(&o, []byte(testD))
	if e == nil {
		if o.OrderNo != "1415757673" ||
		 o.TradeState != glwapi.TRADE_STATE_SUCCESS  ||
		 o.TransId != "1008450740201411110005820873"  ||
		 len(o.Coupon) != 2 ||
		 o.Coupon[1].CouponType !="CASH" ||
		 o.Coupon[1].CouponFee != 1    ||
		 o.Coupon[1].CouponId !="1234" ||
		 o.Coupon[0].CouponId !="01234" {
			t.Errorf("Extract Data mismatch %s", o)
		} else {
			t.Logf("Test ok")
		}
	} else {
		t.Errorf("Test Failed: "+ e.Error())
	}
}


func TestDecodeCreateOrderResp(t * testing.T) {
	g := glwapi.D()

	var testData = "<xml> "
   	testData +="<return_code><![CDATA[SUCCESS]]></return_code> "
   	testData +="<return_msg><![CDATA[OK]]></return_msg> "
   	testData +="<appid><![CDATA[wx2421b1c4370ec43b]]></appid> "
   	testData +="<mch_id><![CDATA[10000100]]></mch_id> "
   	testData +="<nonce_str><![CDATA[IITRi8Iabbblz1Jc]]></nonce_str> "
   	testData +="<openid><![CDATA[oUpF8uMuAJO_M2pxb1Q9zNjWeS6o]]></openid> "
   	testData +="<sign><![CDATA[7921E432F65EB8ED0CE9755F0E86D72F]]></sign> "
   	testData +="<result_code><![CDATA[SUCCESS]]></result_code> "
   	testData +="<prepay_id><![CDATA[wx201411101639507cbf6ffd8b0779950874]]></prepay_id> "
   	testData +="<trade_type><![CDATA[JSAPI]]></trade_type> "
   	testData +="</xml>"

	o, _ :=g.WeChat.CreateOrder1(nil, "1235", "aa","bb", "cc", "192.168.1.1", "http://test", 1234)
	er := glwapi.DecodeCreateOrderResp(o, []byte(testData))
	if er != nil {
		t.Errorf("Test failed :%s\n", er.Error())
	} else {
		if o.PrepayId != "wx201411101639507cbf6ffd8b0779950874" {
			t.Errorf("Test failed prepayID mismatch :%s\n",o.PrepayId )
		} else {
			t.Logf("Test ok  PrepayId extract  matched%s", o.PrepayId)
		}
	}


	//Test failed data
	testData = "<xml> "
   	testData +="<return_code><![CDATA[ORDERCLOSED]]></return_code> "
   	testData +="<return_msg><![CDATA[OK]]></return_msg> "
   	testData +="<appid><![CDATA[wx2421b1c4370ec43b]]></appid> "
   	testData +="<mch_id><![CDATA[10000100]]></mch_id> "
   	testData +="<nonce_str><![CDATA[IITRi8Iabbblz1Jc]]></nonce_str> "
   	testData +="<openid><![CDATA[oUpF8uMuAJO_M2pxb1Q9zNjWeS6o]]></openid> "
   	testData +="<sign><![CDATA[7921E432F65EB8ED0CE9755F0E86D72F]]></sign> "
   	testData +="<result_code><![CDATA[SUCCESS]]></result_code> "
   	testData +="<prepay_id><![CDATA[wx201411101639507cbf6ffd8b0779950874]]></prepay_id> "
	o, _ =g.WeChat.CreateOrder1(nil, "1235", "aa","bb", "cc", "192.168.1.1", "http://test", 1234)
	er = glwapi.DecodeCreateOrderResp(o, []byte(testData))
	if er == nil {
		t.Errorf("Test failed: should not success ")
	} else {
		t.Logf("Test ok  for response failed ")

	}

	
}




func TestDecodeQueryRefundResponse(t * testing.T) {
	g := glwapi.D()
	var testData = "<xml> "
   	testData +="<coupon_refund_fee_0>0</coupon_refund_fee_0> "
   	testData +="<result_code>SUCCESS</result_code> "
   	testData +="<return_code><![CDATA[ORDERCLOSED]]></return_code> "
   	testData +="<return_msg><![CDATA[OK]]></return_msg> "
   	testData +="<appid><![CDATA[aaaaaaaaa]]></appid> "
   	testData +="<mch_id><![CDATA[123123]]></mch_id> "
   	testData +="<sign><![CDATA[dfwerwew23123]]></sign> "
   	testData +="<return_msg><![CDATA[OK]]></return_msg> "
   	testData +="<coupon_refund_fee_1><![CDATA[2333]]></coupon_refund_fee_1> "
   	testData +="<coupon_refund_fee_2><![CDATA[3111]]></coupon_refund_fee_2> "
   	testData +="<coupon_refund_fee_3><![CDATA[4233]]></coupon_refund_fee_3> "
   	testData +="<refund_status_3><![CDATA[REFUNDCLOSE]]></refund_status_3> "
   	testData +="<refund_success_time_2><![CDATA[REFUNDCLOSE]]></refund_success_time_2> "
   	testData +="<refund_account_1><![CDATA[12345]]></refund_account_1> "
   	testData +="<refund_recv_accout_0><![CDATA[2345]]></refund_recv_accout_0> "
   	testData +="<coupon_fee_0>0</coupon_fee_0> "
   	testData +="<coupon_fee_1><![CDATA[1]]></coupon_fee_1> "
   	testData +="<coupon_fee_2><![CDATA[2]]></coupon_fee_2> "
   	testData +="<coupon_fee_3><![CDATA[3]]></coupon_fee_3> "
   	testData +="<coupon_fee_4><![CDATA[4]]></coupon_fee_4> "
   	testData +="<coupon_fee_5><![CDATA[5]]></coupon_fee_5> "
   	testData +="<cash_fee><![CDATA[100]]></cash_fee> "
   	testData +="<cash_fee_type><![CDATA[bb]]></cash_fee_type> "
   	testData +="<transaction_id><![CDATA[xxxxxxxx]]></transaction_id> "
	testData +="</xml>"
	o, _ :=g.WeChat.CreateOrder1(nil, "1235", "aa","bb", "cc", "192.168.1.1", "http://test", 1234)
	glwapi.DecodeQueryRefundResp(o, []byte(testData))
	if o.Cash != 100 ||
		 o.CashType != "bb" ||
		 o.TransId != "xxxxxxxx" ||
		 len(o.Refund.CouponRefund) != 4 ||
		 len(o.Coupon) != 6  ||
		 o.Refund.CouponRefund[2].Fee != 3111 ||
		 o.Refund.CouponRefund[3].Status != "REFUNDCLOSE"  ||
		 o.Refund.CouponRefund[1].Account != "12345"  ||
		 o.Refund.CouponRefund[0].RAccount != "2345"  ||
		 o.Coupon[2].CouponFee != 2 {
		t.Errorf("Extract failed:%s", o)
		t.Errorf("Extract failed:%s", o.Refund.CouponRefund)
	} else {
		t.Logf("Test ok for Decode QueryFund")
	}
	
}


func TestDecodeRequestRefundResp(t * testing.T) {

	testD := ""
	testD +="<xml>"
	testD +="  <return_code><![CDATA[SUCCESS]]></return_code>"
	testD +=" <return_msg><![CDATA[OK]]></return_msg>"
	testD +="  <appid><![CDATA[wx2421b1c4370ec43b]]></appid>"
 	testD +="  <mch_id><![CDATA[10000100]]></mch_id>"
	testD +="  <nonce_str><![CDATA[NfsMFbUFpdbEhPXP]]></nonce_str>"
	testD +="  <sign><![CDATA[B7274EB9F8925EB93100DD2085FA56C0]]></sign>"
	testD +="  <result_code><![CDATA[SUCCESS]]></result_code>"
	testD +="  <transaction_id><![CDATA[1008450740201411110005820873]]></transaction_id>"
	testD +="  <out_trade_no><![CDATA[1415757673]]></out_trade_no>"
	testD +="  <out_refund_no><![CDATA[1415701182]]></out_refund_no>"
	testD +="  <refund_id><![CDATA[2008450740201411110000174436]]></refund_id>"
	testD +="  <refund_channel><![CDATA[]]></refund_channel>"
	testD +="  <refund_fee>23</refund_fee> "
	testD +="</xml>"
	g := glwapi.D()
	o, _ :=g.WeChat.CreateOrder1(nil, "1235", "aa","bb", "cc", "192.168.1.1", "http://test", 1234)
	glwapi.DecodeRequestRefundResp(o, []byte(testD))
	if o.Refund.RefundId != "2008450740201411110000174436" || o.Refund.Fee != 23 || o.TransId != "1008450740201411110005820873"  {
		t.Errorf("Extract failed:%s", o)
	} else {
		t.Logf("Test ok for Decode RequestFundResp")
	}
}



func TestDecodeCloseOrderResp(t * testing.T) {
	testD := ""
	testD += "<xml>"
   	testD += "<return_code><![CDATA[SUCCESS]]></return_code>"
   	testD += "<return_msg><![CDATA[OK]]></return_msg>"
   	testD += "<appid><![CDATA[wx2421b1c4370ec43b]]></appid>"
   	testD += "<mch_id><![CDATA[10000100]]></mch_id>"
   	testD += "<nonce_str><![CDATA[BFK89FC6rxKCOjLX]]></nonce_str>"
   	testD += "<sign><![CDATA[72B321D92A7BFA0B2509F3D13C7B1631]]></sign>"
   	testD += "<result_code><![CDATA[SUCCESS]]></result_code>"
   	testD += "<result_msg><![CDATA[OK]]></result_msg>"
	testD += "</xml>"
	g :=glwapi.D()
	o, _ :=g.WeChat.CreateOrder1(nil, "1235", "aa","bb", "cc", "192.168.1.1", "http://test", 1234)
	e := glwapi.DecodeCloseOrderResp(o, []byte(testD))
	if e == nil {
		t.Logf("Test ok for success order")
	} else {
		t.Logf("Test failed :"+ e.Error())
	}

	testD = ""
	testD += "<xml>"
   	testD += "<return_code><![CDATA[FAIL]]></return_code>"
   	testD += "<return_msg><![CDATA[OK]]></return_msg>"
   	testD += "<appid><![CDATA[wx2421b1c4370ec43b]]></appid>"
   	testD += "<mch_id><![CDATA[10000100]]></mch_id>"
   	testD += "<nonce_str><![CDATA[BFK89FC6rxKCOjLX]]></nonce_str>"
   	testD += "<sign><![CDATA[72B321D92A7BFA0B2509F3D13C7B1631]]></sign>"
   	testD += "<err_code><![CDATA[SUCCESS]]></err_code>"
   	testD += "<err_msg><![CDATA[OK]]></err_msg>"
	testD += "</xml>"

	o, _ =g.WeChat.CreateOrder1(nil, "1235", "aa","bb", "cc", "192.168.1.1", "http://test", 1234)
	e = glwapi.DecodeCloseOrderResp(o, []byte(testD))
	if e == nil {
		t.Logf("Test faild for system error order")
	} else {
		t.Logf("Test ok for get error :"+ e.Error())
	}
}
