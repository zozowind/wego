package app

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/zozowind/wego/core"
	"github.com/zozowind/wego/libs/errmsg"
	"github.com/zozowind/wego/util"
)

const (
	UnifiedOrderUrl = core.WxPayUrl + "/pay/unifiedorder"
)

type PayNotifyRequest struct {
	ReturnCode         string `xml:"return_code"`          //SUCCESS/FAIL 此字段是通信标识，非交易标识，交易是否成功需要查看result_code来判断
	ReturnMsg          string `xml:"return_msg"`           //返回信息，如非空，为错误原因
	Appid              string `xml:"appid"`                //微信分配的小程序ID
	MchId              string `xml:"mch_id"`               //微信支付分配的商户号
	DeviceInfo         string `xml:"device_info"`          //可选, 微信支付分配的终端设备号
	NonceStr           string `xml:"nonce_str"`            //随机字符串，不长于32位
	Sign               string `xml:"sign"`                 //签名
	SignType           string `xml:"sign_type"`            //签名类型，目前支持HMAC-SHA256和MD5，默认为MD5
	ResultCode         string `xml:"result_code"`          //SUCCESS/FAIL
	ErrCode            string `xml:"err_code"`             //错误返回的信息描述
	ErrCodeDes         string `xml:"err_code_des"`         //错误返回的信息描述
	OpenId             string `xml:"openid"`               //用户在商户appid下的唯一标识
	IsSubScribe        string `xml:"is_subscribe"`         //用户是否关注公众账号，Y-关注，N-未关注，仅在公众账号类型支付有效
	TradeType          string `xml:"trade_type"`           //JSAPI、NATIVE、APP, 这里应该是JSAPI
	BankType           string `xml:"bank_type"`            //银行类型，采用字符串类型的银行标识，银行类型见银行列表
	TotalFee           int64  `xml:"total_fee"`            //订单总金额，单位为分
	SettlementTotalFee int64  `xml:"settlement_total_fee"` //应结订单金额=订单金额-非充值代金券金额，应结订单金额<=订单金额
	FeeType            string `xml:"fee_type"`             //货币类型，符合ISO4217标准的三位字母代码，默认人民币：CNY
	CashFee            int64  `xml:"cash_fee"`             //现金支付金额订单现金支付金额
	CashFeeType        string `xml:"cash_fee_type"`        //货币类型，符合ISO4217标准的三位字母代码，默认人民币：CNY，
	CouponFee          int64  `xml:"coupon_fee"`           //代金券金额<=订单金额，订单金额-代金券金额=现金支付金额
	CouponCount        int64  `xml:"coupon_count"`         //代金券使用数量
	//同时支持4张代金券
	CouponFeeType0 string `xml:"coupon_type_0"`  //CASH--充值代金券 NO_CASH---非充值代金券
	CouponFee0     int64  `xml:"coupon_fee_0"`   //单个代金券支付金额
	CouponId0      string `xml:"coupon_id_0"`    //代金券ID
	CouponFeeType1 string `xml:"coupon_type_1"`  //CASH--充值代金券 NO_CASH---非充值代金券
	CouponFee1     int64  `xml:"coupon_fee_1"`   //单个代金券支付金额
	CouponId1      string `xml:"coupon_id_1"`    //代金券ID
	CouponFeeType2 string `xml:"coupon_type_2"`  //CASH--充值代金券 NO_CASH---非充值代金券
	CouponFee2     int64  `xml:"coupon_fee_2"`   //单个代金券支付金额
	CouponId2      string `xml:"coupon_id_2"`    //代金券ID
	CouponFeeType3 string `xml:"coupon_type_3"`  //CASH--充值代金券 NO_CASH---非充值代金券
	CouponFee3     int64  `xml:"coupon_fee_3"`   //单个代金券支付金额
	CouponId3      string `xml:"coupon_id_3"`    //代金券ID
	TransactionId  string `xml:"transaction_id"` //微信支付订单号
	OutTradeNo     string `xml:"out_trade_no"`   //商户系统的订单号，与请求一致
	Attach         string `xml:"attach"`         //商家数据包，原样返回
	TimeEnd        string `xml:"time_end"`       //支付完成时间，格式为yyyyMMddHHmmss
}

type PayNotifyResponse struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
}

func (this *WeAppClient) GetPayNotifyRequest(w http.ResponseWriter, r *http.Request) (*PayNotifyRequest, *errmsg.ErrMsg) {
	var err *errmsg.ErrMsg

	request := &PayNotifyRequest{}
	body, e := ioutil.ReadAll(r.Body)
	if e != nil {
		return request, errmsg.GetError(errPayNotifyData, err.Error())
	}

	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	e = xml.Unmarshal(body, request)
	if e != nil {
		return request, errmsg.GetError(errPayNotifyData, err.Error())
	}

	if request.ReturnCode == "SUCCESS" {
		sign := request.Sign
		request.Sign = ""
		payNotifyData, err := util.StructToUrlValue(request, "xml")
		if nil != err {
			return request, errmsg.GetError(errPayNotifyData, err.Error())
		}
		ok := util.CheckSignMd5(payNotifyData.Encode(), this.Base.PayKey, sign)
		if !ok {
			err = errmsg.GetError(errPayNotifySignCheck, fmt.Sprintf("response: %#v, signStr: %s", request, payNotifyData))
		}
	} else {
		err = errmsg.GetError(errPayNotifyResult, fmt.Sprintf("code: %s, message: %s", request.ReturnCode, request.ReturnMsg))
	}
	return request, err
}

type UnifiedOrderRequest struct {
	// 必选参数
	AppId          string `xml:"appid"`            //微信分配的小程序ID
	MchId          string `xml:"mch_id"`           //微信支付分配的商户号
	NonceStr       string `xml:"nonce_str"`        // 随机字符串，不长于32位。NOTE: 如果为空则系统会自动生成一个随机字符串。
	Body           string `xml:"body"`             // 商品简单描述，该字段须严格按照规范传递，商家名称-销售商品类目
	OutTradeNo     string `xml:"out_trade_no"`     // 商户系统内部的订单号,32个字符内、可包含字母, 其他说明见商户订单号
	TotalFee       int64  `xml:"total_fee"`        // 订单总金额，单位为分，详见支付金额
	SpbillCreateIP string `xml:"spbill_create_ip"` // APP和网页支付提交用户端ip，Native支付填调用微信支付API的机器IP。
	NotifyURL      string `xml:"notify_url"`       // 接收微信支付异步通知回调地址，通知url必须为直接可访问的url，不能携带参数。
	TradeType      string `xml:"trade_type"`       // 取值如下：JSAPI，NATIVE，APP，详细说明见参数规定,这里取JSAPI
	OpenId         string `xml:"openid"`           // rade_type=JSAPI，此参数必传，用户在商户appid下的唯一标识。
	//Sign           string `xml:"sign"`             //签名
	// 可选参数
	DeviceInfo string    `xml:"device_info"` // 终端设备号(门店号或收银设备ID)，注意：PC网页或公众号内支付请传"WEB"
	SignType   string    `xml:"sign_type"`   // 签名类型，默认为MD5，支持HMAC-SHA256和MD5。
	Detail     string    `xml:"detail"`      // 单品优惠字段(暂未上线)
	Attach     string    `xml:"attach"`      // 附加数据，在查询API和支付通知中原样返回，该字段主要用于商户携带订单的自定义数据
	FeeType    string    `xml:"fee_type"`    // 符合ISO 4217标准的三位字母代码，默认人民币：CNY，其他值列表详见货币类型
	TimeStart  time.Time `xml:"time_start"`  // 订单生成时间，格式为yyyyMMddHHmmss，如2009年12月25日9点10分10秒表示为20091225091010。其他详见时间规则
	TimeExpire time.Time `xml:"time_expire"` // 订单失效时间，格式为yyyyMMddHHmmss，如2009年12月27日9点10分10秒表示为20091227091010。其他详见时间规则
	GoodsTag   string    `xml:"goods_tag"`   // 商品标记，代金券或立减优惠功能的参数，说明详见代金券或立减优惠
	LimitPay   string    `xml:"limit_pay"`   // no_credit--指定不能使用信用卡支付
}

type UnifiedOrderResponse struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	//以下字段在return_code为SUCCESS的时候有返回
	AppId      string `xml:"appid"`       //微信分配的小程序ID
	MchId      string `xml:"mch_id"`      //微信支付分配的商户号
	DeviceInfo string `xml:"device_info"` // 终端设备号(门店号或收银设备ID)，注意：PC网页或公众号内支付请传"WEB"
	NonceStr   string `xml:"nonce_str"`   // 随机字符串，不长于32位。NOTE: 如果为空则系统会自动生成一个随机字符串。
	Sign       string `xml:"sign"`        //签名
	ResultCode string `xml:"result_code"`
	ErrCode    string `xml:"err_code"`
	ErrCodeDes string `xml:"err_code_des"`
	//下字段在return_code 和result_code都为SUCCESS的时候有返回
	PrepareId string `xml:"prepay_id"`
	TradeType string `xml:"trade_type"`
}

type PayPackage struct {
	AppId     string `json:"appId"`
	NonceStr  string `json:"nonceStr"`
	Package   string `json:"package"`
	SignType  string `json:"signType"`
	TimeStamp string `json:"timeStamp"`
	PaySign   string `json:"paySign"`
}

func (this *WeAppClient) GeneratePayPackage(request *UnifiedOrderRequest) (*PayPackage, *errmsg.ErrMsg) {
	request.AppId = this.Base.AppId
	request.MchId = this.Base.PayId
	request.NonceStr = util.RandString(32)
	request.TradeType = "JSAPI"
	request.SignType = "MD5" //暂时只支持MD5
	params, err := util.StructToUrlValue(request, "xml")
	if nil != err {
		return nil, errmsg.GetError(errUnifiedOrderReq, err.Error())
	}
	params.Add("sign", util.SignMd5(params.Encode(), this.Base.PayKey))
	data, err := util.HttpXMLPost(this.Base.HttpClient, UnifiedOrderUrl, params)
	if nil != err {
		return nil, errmsg.GetError(errUnifiedOrderRsp, err.Error())
	}

	response := &UnifiedOrderResponse{}
	err = xml.Unmarshal(data, &response)
	if nil != err {
		return nil, errmsg.GetError(errUnifiedOrderRsp, err.Error())
	}

	//校验结果
	if response.ReturnCode != "SUCCESS" {
		return nil, errmsg.GetError(errUnifiedOrderResult, fmt.Sprintf("code: %s, message: %s", response.ReturnCode, response.ReturnMsg))
	}

	if response.ResultCode != "SUCCESS" {
		return nil, errmsg.GetError(errUnifiedOrderResult, fmt.Sprintf("result_code:%s, err_code: %s, err_code_des: %s", response.ResultCode, response.ErrCode, response.ErrCodeDes))
	}

	// 校验 trade_type
	if response.TradeType != request.TradeType {
		return nil, errmsg.GetError(errUnifiedOrderResult, fmt.Sprintf("trade_type mismatch, have: %s, want: %s", response.TradeType, request.TradeType))
	}

	//校验签名
	sign := response.Sign
	response.Sign = ""
	resData, err := util.StructToUrlValue(response, "xml")
	if nil != err {
		return nil, errmsg.GetError(errUnifiedOrderResult, err.Error())
	}
	ok := util.CheckSignMd5(resData.Encode(), this.Base.PayKey, sign)
	if !ok {
		return nil, errmsg.GetError(errUnifiedOrderResult, fmt.Sprintf("sign invalid, response: %#v, signStr: %s", response, resData))
	}

	//组织支付包
	payPackage := &PayPackage{
		AppId:     this.Base.AppId,
		NonceStr:  util.RandString(32),
		Package:   fmt.Sprintf("prepay_id=%s", response.PrepareId),
		SignType:  "MD5",
		TimeStamp: fmt.Sprintf("%d", time.Now().Unix()),
	}
	params, err = util.StructToUrlValue(payPackage, "json")
	if nil != err {
		return payPackage, errmsg.GetError(errPayPackage, err.Error())
	}
	payPackage.PaySign = util.SignMd5(params.Encode(), this.Base.PayKey)
	return payPackage, nil
}
