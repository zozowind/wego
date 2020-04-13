package util

const (
	// SignTypeMD5 md5签名
	SignTypeMD5 = "MD5"
	// SignTypeHMACSHA256 HMAC-SHA256签名
	SignTypeHMACSHA256 = "HMAC-SHA256"

	// ProfitSharingYes 需要分账
	ProfitSharingYes = "Y"
	// ProfitSharingNo 不需要分账
	ProfitSharingNo = "N"

	// CodeSuccess 提交业务成功
	CodeSuccess = "SUCCESS"
	// CodeFail 提交业务失败
	CodeFail = "FAIL"

	// TradeTypeJSAPI JSAPI支付或者小程序支付
	TradeTypeJSAPI = "JSAPI"
	// TradeTypeNATIVE NATIVE支付（微信扫码支付）
	TradeTypeNATIVE = "NATIVE"
	// TradeTypeAPP APP支付（APP跳转微信支付）
	TradeTypeAPP = "APP"
	// TradeTypeH5 移动端在微信客户端外的网页支付
	TradeTypeH5 = "MWEB"

	// AccountTypeMerchantID 商户ID
	AccountTypeMerchantID = "MERCHANT_ID"
	// AccountTypePersonalWechatID 个人微信号
	AccountTypePersonalWechatID = "PERSONAL_WECHATID"
	// AccountTypePersonalOpenID 个人openid
	AccountTypePersonalOpenID = "PERSONAL_OPENID"

	// RelationTypeServiceProvider 服务商
	RelationTypeServiceProvider = "SERVICE_PROVIDER"
	// RelationTypeStore 门店
	RelationTypeStore = "STORE"
	// RelationTypeStaff 员工
	RelationTypeStaff = "STAFF"
	// RelationTypeStoreOwner 店主
	RelationTypeStoreOwner = "STORE_OWNER"
	// RelationTypePartner 合作伙伴
	RelationTypePartner = "PARTNER"
	// RelationTypeHeadQuarter 总部
	RelationTypeHeadQuarter = "HEADQUARTER"
	// RelationTypeHeadBrand 品牌方
	RelationTypeHeadBrand = "BRAND"
	// RelationTypeDistributor 分销商
	RelationTypeDistributor = "DISTRIBUTOR"
	// RelationTypeUser 用户
	RelationTypeUser = "USER"
	// RelationTypeSupplier 供应商
	RelationTypeSupplier = "SUPPLIER"
	// RelationTypeCustom 自定义
	RelationTypeCustom = "CUSTOM"
)
