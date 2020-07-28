package pay

import (
	"encoding/json"
	"fmt"

	"github.com/zozowind/wego/libs/errmsg"
)

//StoreInfo 店铺信息
type StoreInfo struct {
	ID       string `json:"id"`        //SZTX001	门店编号，由商户自定义
	Name     string `json:"name"`      //腾讯大厦腾大餐厅	门店名称 ，由商户自定义
	AreaCode string `json:"area_code"` //440305	门店所在地行政区划码，详细见
	Address  string `json:"address"`   //科技园中一路腾讯大厦	门店详细地址 ，由商户自定义
}

//OrderUnifiedSceneInfo 场景信息
type OrderUnifiedSceneInfo struct {
	StoreInfo StoreInfo `json:"store_info"`
}

//GoodDetail 商品详情
type GoodDetail struct {
	GoodsID      string `json:"goods_id"`       //商品编码	goods_id	是	String(32)	商品编码	由半角的大小写字母、数字、中划线、下划线中的一种或几种组成
	WxpayGoodsID string `json:"wxpay_goods_id"` //微信侧商品编码	wxpay_goods_id	否	String(32)	1001	微信支付定义的统一商品编号（没有可不传）
	GoodsName    string `json:"goods_name"`     //商品名称	goods_name	否	String(256)	iPhone6s 16G	商品的实际名称
	Quantity     int64  `json:"quantity"`       //商品数量	quantity	是	int	1	用户购买的数量
	Price        int64  `json:"price"`          //商品单价	price	是	int	528800	单位为：分。如果商户有优惠，需传输商户优惠后的单价(例如：用户对一笔100元的订单使用了商场发的纸质优惠券100-50，则活动商品的单价应为原单价-50)
}

//DetailOrigin 订单详情
type DetailOrigin struct {
	CostPrice   int64         `json:"cost_price"`   //订单原价	cost_price	否	int	608800	1.商户侧一张小票订单可能被分多次支付，订单原价用于记录整张小票的交易金额。
	ReceiptID   string        `json:"receipt_id"`   //商品小票ID	receipt_id	否	String(32)	wx123	商家小票ID
	GoodsDetail []*GoodDetail `json:"goods_detail"` //单品列表	goods_detail	是	String	示例见下文	单品信息，使用Json数组格式提交
}

//OrderUnifiedReq 统一下单接口参数
type OrderUnifiedReq struct {
	// 必选参数
	CommonReq
	DeviceInfo string `xml:"device_info"` // 选填， 终端设备号(门店号或收银设备ID)，注意：PC网页或公众号内支付请传"WEB"
	//Sign           string `xml:"sign"`             //签名，最后添加到请求数据中，不列入结构体
	Body            string                 `xml:"body"` // 必填， 商品简单描述，该字段须严格按照规范传递，商家名称-销售商品类目
	DetailOrigin    *DetailOrigin          `xml:"-"`
	Detail          string                 `xml:"detail"`           // 选填， 商品详细描述，对于使用单品优惠的商户，该字段必须按照规范上传，详见“单品优惠参数说明”
	Attach          string                 `xml:"attach"`           // 选填， 附加数据，在查询API和支付通知中原样返回，该字段主要用于商户携带订单的自定义数据
	OutTradeNo      string                 `xml:"out_trade_no"`     // 必填， 商户系统内部的订单号,32个字符内、可包含字母, 其他说明见商户订单号
	FeeType         string                 `xml:"fee_type"`         // 选填， 符合ISO 4217标准的三位字母代码，默认人民币：CNY，其他值列表详见货币类型
	TotalFee        int64                  `xml:"total_fee"`        // 必填， 订单总金额，单位为分，详见支付金额
	SpbillCreateIP  string                 `xml:"spbill_create_ip"` // 必填， APP和网页支付提交用户端ip，Native支付填调用微信支付API的机器IP。
	TimeStart       string                 `xml:"time_start"`       // 选填， 订单生成时间，格式为yyyyMMddHHmmss，如2009年12月25日9点10分10秒表示为20091225091010。其他详见时间规则
	TimeExpire      string                 `xml:"time_expire"`      // 选填， 订单失效时间，格式为yyyyMMddHHmmss，如2009年12月27日9点10分10秒表示为20091227091010。其他详见时间规则
	GoodsTag        string                 `xml:"goods_tag"`        // 选填， 商品标记，代金券或立减优惠功能的参数，说明详见代金券或立减优惠
	NotifyURL       string                 `xml:"notify_url"`       // 必填， 接收微信支付异步通知回调地址，通知url必须为直接可访问的url，不能携带参数。
	TradeType       string                 `xml:"trade_type"`       // 必填， 取值如下：JSAPI，NATIVE，APP，详细说明见参数规定,这里取JSAPI
	ProductID       string                 `xml:"product_id"`       // 选填，32位trade_type=NATIVE时，此参数必传。此参数为二维码中包含的商品ID，商户自行定义。
	LimitPay        string                 `xml:"limit_pay"`        // 选填， no_credit--指定不能使用信用卡支付
	OpenID          string                 `xml:"openid"`           // 选填， trade_type=JSAPI，此参数必传，用户在商户appid下的唯一标识。
	Receipt         string                 `xml:"receipt"`          // 选填， Y，传入Y时，支付成功消息和支付详情页将出现开票入口。需要在微信支付商户平台或微信公众平台开通电子发票功能，传此字段才可生效
	SceneInfoOrigin *OrderUnifiedSceneInfo `xml:"-"`                // 场景信息原始数据
	SceneInfo       string                 `xml:"scene_info"`       // 选填，该字段常用于线下活动时的场景信息上报，支持上报实际门店信息，商户也可以按需求自己上报相关信息。该字段为JSON对象数据
	ProfitSharing   string                 `xml:"profit_sharing"`   // 选填， 是否分账，Y需要分账，N不分账，字母要求大写，默认认不分账
}

//OrderUnifiedRsp 统一下单接口返回数据
type OrderUnifiedRsp struct {
	CommonRsp
	//以下字段在return_code为SUCCESS的时候有返回
	DeviceInfo string `xml:"device_info"` // 终端设备号(门店号或收银设备ID)，注意：PC网页或公众号内支付请传"WEB"
	//下字段在return_code 和result_code都为SUCCESS的时候有返回
	PrepareID string `xml:"prepay_id"`
	TradeType string `xml:"trade_type"`
	CodeURL   string `xml:"code_url"` //trade_type=NATIVE时有返回，此url用于生成支付二维码，然后提供给用户进行扫码支付。
}

//OrderUnified 统一下单接口
func (cli *Client) OrderUnified(req *OrderUnifiedReq) (rsp *OrderUnifiedRsp, err error) {
	if nil != req.SceneInfoOrigin {
		var scene []byte
		scene, err = json.Marshal(req.SceneInfoOrigin)
		if nil != err {
			return
		}
		req.SceneInfo = string(scene)
	}

	if nil != req.DetailOrigin {
		var detail []byte
		detail, err = json.Marshal(req.DetailOrigin)
		if nil != err {
			return
		}
		req.SceneInfo = string(detail)
	}

	rsp = &OrderUnifiedRsp{}
	err = cli.request(orderUnifiedURL, req, rsp)
	if nil != err {
		err = errmsg.GetError(errPayRequestReq, err.Error())
		return
	}
	// 校验 trade_type
	if rsp.TradeType != req.TradeType {
		err = errmsg.GetError(errPayOrderRsp, fmt.Sprintf("trade_type mismatch, have: %s, want: %s", rsp.TradeType, req.TradeType))
	}
	return
}

//OrderQueryReq 订单查询接口参数
type OrderQueryReq struct {
	CommonReq
	TransactionID string `xml:"transaction_id"` //微信的订单号，建议优先使用
	OutTradeNo    string `xml:"out_trade_no"`   //商户系统内部订单号，要求32个字符内，只能是数字、大小写字母_-|*@ ，且在同一个商户号下唯一
}

//Coupon 优惠券
type Coupon struct {
	CouponType string //CASH--充值代金券 NO_CASH---非充值代金券
	CouponFee  int64  //单个代金券支付金额
	CouponID   string //代金券ID
}

//OrderQueryRsp 订单查询接口返回
type OrderQueryRsp struct {
	CommonRsp
	DeviceInfo         string    `xml:"device_info"`          // 设备号	device_info	否	String(32)	013467007045764	微信支付分配的终端设备号
	OpenID             string    `xml:"openid"`               // 用户标识	openid	是	String(128)	oUpF8uMuAJO_M2pxb1Q9zNjWeS6o	用户在商户appid下的唯一标识
	IsSubscribe        string    `xml:"is_subscribe"`         //  是否关注公众账号	is_subscribe	是	String(1)	Y	用户是否关注公众账号，Y-关注，N-未关注
	TradeType          string    `xml:"trade_type"`           // 交易类型	trade_type	是	String(16)	JSAPI	调用接口提交的交易类型，取值如下：JSAPI，NATIVE，APP，MICROPAY，详细说明见参数规定
	TradeState         string    `xml:"trade_state"`          // 交易状态	trade_state	是	String(32)	SUCCESS
	BankType           string    `xml:"bank_type"`            //付款银行	bank_type	是	String(16)	CMC	银行类型，采用字符串类型的银行标识
	TotalFee           int64     `xml:"total_fee"`            // 标价金额	total_fee	是	Int	100	订单总金额，单位为分
	SettlementTotalFee int64     `xml:"settlement_total_fee"` // 应结订单金额	settlement_total_fee	否	Int	100	当订单使用了免充值型优惠券后返回该参数，应结订单金额=订单金额-免充值优惠券金额。
	FeeType            string    `xml:"fee_type"`             // 标价币种	fee_type	否	String(8)	CNY	货币类型，符合ISO 4217标准的三位字母代码，默认人民币：CNY，其他值列表详见货币类型
	CashFee            int64     `xml:"cash_fee"`             // 现金支付金额	cash_fee	是	Int	100	现金支付金额订单现金支付金额，详见支付金额
	CashFeeType        string    `xml:"cash_fee_type"`        // 现金支付币种	cash_fee_type	否	String(16)	CNY	货币类型，符合ISO 4217标准的三位字母代码，默认人民币：CNY，其他值列表详见货币类型
	CouponFee          int64     `xml:"coupon_fee"`           // 代金券金额	coupon_fee	否	Int	100	“代金券”金额<=订单金额，订单金额-“代金券”金额=现金支付金额，详见支付金额
	CouponCount        int64     `xml:"coupon_count"`         // 代金券使用数量	coupon_count	否	Int	1	代金券使用数量
	Coupons            []*Coupon `xml:"-"`                    //代金券
	TransactionID      string    `xml:"transaction_id"`       //微信支付订单号	transaction_id	是	String(32)	1009660380201506130728806387	微信支付订单号
	OutTradeNo         string    `xml:"out_trade_no"`         //商户订单号	out_trade_no	是	String(32)	20150806125346	商户系统内部订单号，要求32个字符内，只能是数字、大小写字母_-|*@ ，且在同一个商户号下唯一。
	Attach             string    `xml:"attach"`               //附加数据	attach	否	String(128)	深圳分店	附加数据，原样返回
	TimeEnd            string    `xml:"time_end"`             //支付完成时间	time_end	是	String(14)	20141030133525	订单支付时间，格式为yyyyMMddHHmmss
	TradeStateDesc     string    `xml:"trade_state_desc"`     //交易状态描述	trade_state_desc	是	String(256)	支付失败，请重新下单支付	对当前查询订单状态的描述和下一步操作的指引
}

//OrderQuery 订单查询接口
func (cli *Client) OrderQuery(req *OrderQueryReq) (rsp *OrderQueryRsp, err error) {
	rsp = &OrderQueryRsp{}
	err = cli.request(orderQueryURL, req, rsp)
	if nil != err {
		err = errmsg.GetError(errPayRequestReq, err.Error())
	}
	//@todo 需要遍历coupons

	return
}

//OrderCloseReq 订单关闭接口参数
type OrderCloseReq struct {
	CommonReq
	OutTradeNo string `xml:"out_trade_no"` //商户系统内部订单号，要求32个字符内，只能是数字、大小写字母_-|*@ ，且在同一个商户号下唯一
}

//OrderClose 订单关闭接口
func (cli *Client) OrderClose(req *OrderCloseReq) (rsp *CommonRsp, err error) {
	rsp = &CommonRsp{}
	err = cli.request(orderCloseURL, req, rsp)
	if nil != err {
		err = errmsg.GetError(errPayRequestReq, err.Error())
	}
	return
}

//RefundReq 退款接口参数
type RefundReq struct {
	CommonReq
	TransactionID string `xml:"transaction_id"`  //微信的订单号，建议优先使用
	OutTradeNo    string `xml:"out_trade_no"`    //商户系统内部订单号，要求32个字符内，只能是数字、大小写字母_-|*@ ，且在同一个商户号下唯一
	OutRefundNo   string `xml:"out_refund_no"`   //是	String(64)	1217752501201407033233368018	商户系统内部的退款单号，商户系统内部唯一，只能是数字、大小写字母_-|*@ ，同一退款单号多次请求只退一笔。
	TotalFee      int64  `xml:"total_fee"`       // 订单金额	total_fee	是	Int	100	订单总金额，单位为分，只能为整数，详见支付金额
	RefundFee     int64  `xml:"refund_fee"`      // 退款金额	refund_fee	是	Int	100	退款总金额，订单总金额，单位为分，只能为整数，详见支付金额
	RefundFeeType string `xml:"refund_fee_type"` //退款货币种类	refund_fee_type	否	String(8)	CNY	退款货币类型，需与支付一致，或者不填。符合ISO 4217标准的三位字母代码，默认人民币：CNY，其他值列表详见货币类型
	RefundDesc    string `xml:"refund_desc"`     //退款原因	refund_desc	否	String(80)	商品已售完
	RefundAccount string `xml:"refund_account"`  //退款资金来源	refund_account	否	String(30)	REFUND_SOURCE_RECHARGE_FUNDS
	NotifyURL     string `xml:"notify_url"`      //退款结果通知url	notify_url	否	String(256)	https://weixin.qq.com/notify/
}

//RefundRsp 退款接口返回
type RefundRsp struct {
	CommonRsp
	TransactionID       string    `xml:"transaction_id"`        //微信订单号	transaction_id	是	String(32)	4007752501201407033233368018	微信订单号
	OutTradeNo          string    `xml:"out_trade_no"`          //商户订单号	out_trade_no	是	String(32)	33368018	商户系统内部订单号，要求32个字符内，只能是数字、大小写字母_-|*@ ，且在同一个商户号下唯一。
	OutRefundNo         string    `xml:"out_refund_no"`         //商户退款单号	out_refund_no	是	String(64)	121775250	商户系统内部的退款单号，商户系统内部唯一，只能是数字、大小写字母_-|*@ ，同一退款单号多次请求只退一笔。
	RefundID            string    `xml:"refund_id"`             //微信退款单号	refund_id	是	String(32)	2007752501201407033233368018	微信退款单号
	RefundFee           int64     `xml:"refund_fee"`            //退款金额	refund_fee	是	Int	100	退款总金额,单位为分,可以做部分退款
	SettlementRefundFee int64     `xml:"settlement_refund_fee"` //应结退款金额	settlement_refund_fee	否	Int	100	去掉非充值代金券退款金额后的退款金额，退款金额=申请退款金额-非充值代金券退款金额，退款金额<=申请退款金额
	TotalFee            int64     `xml:"total_fee"`             //标价金额	total_fee	是	Int	100	订单总金额，单位为分，只能为整数，详见支付金额
	SettlementTotalFee  int64     `xml:"settlement_total_fee"`  //应结订单金额	settlement_total_fee	否	Int	100	去掉非充值代金券金额后的订单总金额，应结订单金额=订单金额-非充值代金券金额，应结订单金额<=订单金额。
	FeeType             string    `xml:"fee_type"`              //标价币种	fee_type	否	String(8)	CNY	订单金额货币类型，符合ISO 4217标准的三位字母代码，默认人民币：CNY，其他值列表详见货币类型
	CashFee             int64     `xml:"cash_fee"`              //现金支付金额	cash_fee	是	Int	100	现金支付金额，单位为分，只能为整数，详见支付金额
	CashFeeType         string    `xml:"cash_fee_type"`         //现金支付币种	cash_fee_type	否	String(16)	CNY	货币类型，符合ISO 4217标准的三位字母代码，默认人民币：CNY，其他值列表详见货币类型
	CashRefundFee       int64     `xml:"cash_refund_fee"`       //现金退款金额	cash_refund_fee	否	Int	100	现金退款金额，单位为分，只能为整数，详见支付金额
	Coupons             []*Coupon `xml:"-"`                     //代金券
	CouponRefundFee     int64     `xml:"coupon_refund_fee"`     // Int	100	代金券退款金额<=退款金额，退款金额-代金券或立减优惠退款金额为现金，说明详见代金券或立减优惠
	CouponRefundCount   int64     `xml:"coupon_refund_count"`   // Int	1	退款代金券使用数量
}

//Refund 退款接口
func (cli *Client) Refund(req *RefundReq) (rsp *RefundRsp, err error) {
	rsp = &RefundRsp{}
	err = cli.request(refundURL, req, rsp)
	if nil != err {
		err = errmsg.GetError(errPayRequestReq, err.Error())
	}
	return
}

//RefundQueryReq 退款查询接口参数
type RefundQueryReq struct {
	CommonReq
	TransactionID string `xml:"transaction_id"` //微信订单号	transaction_id	四选一	String(32)	1217752501201407033233368018	微信订单号查询的优先级是： refund_id > out_refund_no > transaction_id > out_trade_no
	OutTradeNo    string `xml:"out_trade_no"`   //商户订单号	out_trade_no	String(32)	1217752501201407033233368018	商户系统内部订单号，要求32个字符内，只能是数字、大小写字母_-|*@ ，且在同一个商户号下唯一。
	OutRefundNo   string `xml:"out_refund_no"`  //商户退款单号	out_refund_no	String(64)	1217752501201407033233368018	商户系统内部的退款单号，商户系统内部唯一，只能是数字、大小写字母_-|*@ ，同一退款单号多次请求只退一笔。
	RefundID      string `xml:"refund_id"`      //微信退款单号	refund_id	String(32)	1217752501201407033233368018
	Offset        int64  `xml:"offset"`         //偏移量	offset	否	Int	15	偏移量，当部分退款次数超过10次时可使用，表示返回的查询结果从这个偏移量开始取记录
}

//RefundInfo 退款详情
type RefundInfo struct {
	OutRefundNo         string //商户退款单号	out_refund_no_$n	是	String(64)	1217752501201407033233368018	商户系统内部的退款单号，商户系统内部唯一，只能是数字、大小写字母_-|*@ ，同一退款单号多次请求只退一笔。
	RefundID            string //微信退款单号	refund_id_$n	是	String(32)	1217752501201407033233368018	微信退款单号
	RefundChannel       string //退款渠道	refund_channel_$n	否	String(16)	ORIGINAL
	RefundFee           int64  //申请退款金额	refund_fee_$n	是	Int	100	退款总金额,单位为分,可以做部分退款
	SettlementRefundFee int64  //退款金额	settlement_refund_fee_$n	否	Int	100	退款金额=申请退款金额-非充值代金券退款金额，退款金额<=申请退款金额
	CouponRefundFee     int64  //总代金券退款金额	coupon_refund_fee_$n	否	Int	100	代金券退款金额<=退款金额，退款金额-代金券或立减优惠退款金额为现金，说明详见代金券或立减优惠
	CouponRefundCount   int64  //退款代金券使用数量	coupon_refund_count_$n	否	Int	1	退款代金券使用数量 ,$n为下标,从0开始编号
	RefundStatus        string //退款状态 refund_status_$n	是	String(16)	SUCCESS
	RefundAccount       string //退款资金来源	refund_account_$n	否	String(30)	REFUND_SOURCE_RECHARGE_FUNDS
	RefundRecvAccount   string //退款入账账户	refund_recv_accout_$n	是	String(64)	招商银行信用卡0403	取当前退款单的退款入账方
	RefundSuccessTime   string //否	String(20)	2016-07-25 15:26:26	退款成功时间，当退款状态为退款成功时有返回。$n为下标，从0开始编号。
	Coupons             []*Coupon
}

//RefundQueryRsp 退款查询接口返回
type RefundQueryRsp struct {
	CommonRsp
	TotalRefundCount   int64         `xml:"total_refund_count"`   //订单总退款次数	total_refund_count	否	Int	35	订单总共已发生的部分退款次数，当请求参数传入offset后有返回
	TransactionID      string        `xml:"transaction_id"`       //微信订单号	transaction_id	是	String(32)	1217752501201407033233368018	微信订单号
	OutTradeNo         string        `xml:"out_trade_no"`         //商户订单号	out_trade_no	是	String(32)	1217752501201407033233368018	商户系统内部订单号，要求32个字符内，只能是数字、大小写字母_-|*@ ，且在同一个商户号下唯一。
	TotalFee           int64         `xml:"total_fee"`            //订单金额	total_fee	是	Int	100	订单总金额，单位为分，只能为整数，详见支付金额
	SettlementTotalFee int64         `xml:"settlement_total_fee"` //应结订单金额	settlement_total_fee	否	Int	100	当订单使用了免充值型优惠券后返回该参数，应结订单金额=订单金额-免充值优惠券金额。
	FeeType            string        `xml:"fee_type"`             //货币种类	fee_type	否	String(8)	CNY	订单金额货币类型，符合ISO 4217标准的三位字母代码，默认人民币：CNY，其他值列表详见货币类型
	CashFee            int64         `xml:"cash_fee"`             //现金支付金额	cash_fee	是	Int	100	现金支付金额，单位为分，只能为整数，详见支付金额
	RefundCount        int64         `xml:"refund_count"`         //退款笔数	refund_count	是	Int	1	当前返回退款笔数
	RefundList         []*RefundInfo `xml:"-"`
}

//RefundQuery 退款查询接口
func (cli *Client) RefundQuery(req *RefundQueryReq) (rsp *RefundQueryRsp, err error) {
	rsp = &RefundQueryRsp{}
	err = cli.request(refundQueryURL, req, rsp)
	if nil != err {
		err = errmsg.GetError(errPayRequestReq, err.Error())
	}
	//@todo 需要遍历refund和coupon
	return
}

//BillDownloadReq 下载交易账单接口参数
type BillDownloadReq struct {
	CommonReq
	BillDate string `xml:"bill_date"` //对账单日期	bill_date	是	String(8)	20140603	下载对账单的日期，格式：20140603
	BillType string `xml:"bill_type"` //账单类型	bill_type	否	String(8)	ALL
	TarType  string `xml:"tar_type"`  //压缩账单	tar_type	否	String	GZIP	非必传参数，固定值：GZIP，返回格式为.gzip的压缩包账单。不传则默认为数据流形式。
}

//BillDownload 下载交易账单接口
func (cli *Client) BillDownload(req *BillDownloadReq) (err error) {
	//这里需要处理文本
	return
}

//FundFlowDownloadReq 下载资金账单接口参数
type FundFlowDownloadReq struct {
	CommonReq
	BillDate    string `xml:"bill_date"`    //资金账单日期	bill_date	是	String(8)	20140603	下载对账单的日期，格式：20140603
	AccountType string `xml:"account_type"` //资金账户类型	account_type	是	String(8)	Basic
	TarType     string `xml:"tar_type"`     //压缩账单	tar_type	否	String(8)	GZIP	非必传参数，固定值：GZIP，返回格式为.gzip的压缩包账单。不传则默认为数据流形式。
}

//FundFlowDownload 下载资金账单接口
func (cli *Client) FundFlowDownload(req *FundFlowDownloadReq) (err error) {
	//这里需要处理文本
	return
}

//ItilReportReq 交易保障接口参数
type ItilReportReq struct {
	CommonReq
	DeviceInfo   string `xml:"device_info"`   //设备号	device_info	否	String(32)	013467007045764	微信支付分配的终端设备号，商户自定义
	InterfaceURL string `xml:"interface_url"` //接口URL	interface_url	是	String(127)
	ExecuteTime  int64  `xml:"execute_time"`  //接口耗时	execute_time_	是	Int	1000
	ReturnCode   string `xml:"return_code"`   //返回状态码	return_code	是	String(16)	SUCCESS
	ReturnMsg    string `xml:"return_msg"`    //返回信息	return_msg	否	String(128)	签名失败
	ResultCode   string `xml:"result_code"`   //业务结果	result_code	是	String(16)	SUCCESS
	ErrCode      string `xml:"err_code"`      //错误代码	err_code	否	String(32)	SYSTEMERROR
	ErrCodeDes   string `xml:"err_code_des"`  //错误代码描述	err_code_des	否	String(128)	系统错误
	OutTradeNo   string `xml:"out_trade_no"`  //商户订单号	out_trade_no	否	String(32)	1217752501201407033233368018	商户系统内部的订单号,商户可以在上报时提供相关商户订单号方便微信支付更好的提高服务质量。
	UserIP       string `xml:"user_ip"`       //访问接口IP	user_ip	是	String(16)	8.8.8.8	发起接口调用时的机器IP
	Time         string `xml:"time"`          //商户上报时间	time	否	String(14)	20091227091010
}

//ItilReport 交易保障接口
func (cli *Client) ItilReport(req *ItilReportReq) (rsp *CommonRsp, err error) {
	rsp = &CommonRsp{}
	err = cli.request(itilReportURL, req, rsp)
	if nil != err {
		err = errmsg.GetError(errPayRequestReq, err.Error())
	}
	return
}
