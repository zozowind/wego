package app

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/zozowind/wego/core"
)

// 直播相关接口

const (
	createRoomURL        = core.WxAPIURL + "/wxaapi/broadcast/room/create?access_token=%s"       //1.创建直播间
	getLiveInfoURL       = core.WxAPIURL + "/wxa/business/getliveinfo?access_token=%s"           //2.获取直播房间列表, 获取直播间回放
	addRoomGoodsURL      = core.WxAPIURL + "/wxaapi/broadcast/room/addgoods?access_token=%s"     //3.直播间导入商品
	addGoodsURL          = core.WxAPIURL + "/wxaapi/broadcast/goods/add?access_token=%s"         //1.商品添加并提审
	resetAuditGoodsURL   = core.WxAPIURL + "/wxaapi/broadcast/goods/resetaudit?access_token=%s"  //2.撤回审核
	auditGoodsURL        = core.WxAPIURL + "/wxaapi/broadcast/goods/audit?access_token=%s"       //3.重新提交审核
	deleteGoodsURL       = core.WxAPIURL + "/wxaapi/broadcast/goods/delete?access_token=%s"      //4.删除商品
	updateGoodsURL       = core.WxAPIURL + "/wxaapi/broadcast/goods/update?access_token=%s"      //5.更新商品
	getGoodsWarehouseURL = core.WxAPIURL + "/wxa/business/getgoodswarehouse?access_token=%s"     //6.获取商品状态
	getGoodsApprovedURL  = core.WxAPIURL + "/wxaapi/broadcast/goods/getapproved?access_token=%s" //7.获取商品列表

	//LiveTypePushFlow 推流直播
	LiveTypePushFlow = 1
	//LiveTypeMobile 手机直播
	LiveTypeMobile = 0
	//LiveScreenPortrait 竖屏
	LiveScreenPortrait = 0
	//LiveScreenLandscape 横屏
	LiveScreenLandscape = 1
	//LiveTrueFlag 是标记
	LiveTrueFlag = 1
	//LiveFalseFlag 否标记
	LiveFalseFlag = 0
	//actionGetReplay 获取回放action
	actionGetReplay = "get_replay"

	//LiveStatusDoing 直播中
	LiveStatusDoing = 101
	//LiveStatusNotStart 未开始
	LiveStatusNotStart = 102
	//LiveStatusOver 已结束
	LiveStatusOver = 103
	//LiveStatusForbidden 禁播
	LiveStatusForbidden = 104
	//LiveStatusPause 暂停
	LiveStatusPause = 105
	//LiveStatusException 异常
	LiveStatusException = 106
	//LiveStatusExpire 已过期
	LiveStatusExpire = 107

	//PriceTypeOne 一口价
	PriceTypeOne = 1
	//PriceTypeRange 价格区间
	PriceTypeRange = 2
	//PriceTypeDiscount 折扣价
	PriceTypeDiscount = 3

	//AuditStatusInit 未审核
	AuditStatusInit = 0
	//AuditStatusDoing 审核中
	AuditStatusDoing = 1
	//AuditStatusSuccess 审核成功
	AuditStatusSuccess = 2
	//AuditStatusFail 审核失败
	AuditStatusFail = 3

	//ThirdPartyTagAddFromAPI 从api添加商品
	ThirdPartyTagAddFromAPI = 2
)

//CreateRoomParam 创建直播间参数
type CreateRoomParam struct {
	Name         string `json:"name"`         //  "测试直播房间1",   房间名字
	CoverImg     string `json:"coverImg"`     //  "",    通过 uploadfile 上传，填写 mediaID
	StartTime    int64  `json:"startTime"`    //  1588237130,    开始时间
	EndTime      int64  `json:"endTime"`      //  1588237130 ,  结束时间
	AnchorName   string `json:"anchorName"`   //  "zefzhang1",   主播昵称
	AnchorWechat string `json:"anchorWechat"` //  "WxgQiao_04",   主播微信号
	ShareImg     string `json:"shareImg"`     //  "" ,  通过 uploadfile 上传，填写 mediaID
	Type         int    `json:"type"`         //  1 ,  直播类型，1 推流 0 手机直播
	ScreenType   int    `json:"screenType"`   //  0,   1：横屏 0：竖屏
	CloseLike    int    `json:"closeLike"`    //  0 ,  是否 关闭点赞 1 关闭
	CloseGoods   int    `json:"closeGoods"`   //  0,  是否 关闭商品货架，1：关闭
	CloseComment int    `json:"closeComment"` //  0  是否开启评论，1：关闭
}

//CreateRoomResponse 创建直播间结果
type CreateRoomResponse struct {
	core.WxErrorResponse
	RoomID int64 `json:"roomId"` //房间ID
}

//CreateRoom 创建直播间
func (client *WeAppClient) CreateRoom(param *CreateRoomParam) (res *CreateRoomResponse, err error) {
	res = &CreateRoomResponse{}

	data, err := client.PostWithToken(createRoomURL, param)
	if nil != err {
		return
	}
	err = json.Unmarshal(data, res)
	if nil != err {
		return
	}
	err = res.Check()
	return
}

//GetRoomListParam 获取直播间列表参数
type GetRoomListParam struct {
	Start int64 `json:"start"` // 0 起始拉取房间，start = 0 表示从第 1 个房间开始拉取
	Limit int64 `json:"limit"` // 10 每次拉取的个数上限，不要设置过大，建议 100 以内
}

//Good 商品信息
type Good struct {
	CoverImg string  `json:"cover_img"` //商品封面图链接
	URL      string  `json:"url"`       //商品小程序路径
	Price    float64 `json:"price"`     //商品价格
	Name     string  `json:"name"`      //商品名称
}

//Room 房间信息
type Room struct {
	Name       string  `json:"name"`   // "name":"直播房间名"
	RoomID     int64   `json:"roomid"` // "roomid": 1,
	CoverImg   string  `json:"cover_img"`
	ShareImg   string  `json:"share_img"`
	LiveStatus int     `json:"live_status"`
	StartTime  int64   `json:"start_time"`
	EndTime    int64   `json:"end_time"`
	AnchorName string  `json:"anchor_name"`
	Goods      []*Good `json:"goods"` // 商品详细
	Total      int64   `json:"total"` // 总数
}

//GetRoomListResponse 获取直播间列表结果
type GetRoomListResponse struct {
	core.WxErrorResponse
	RoomInfo []*Room `json:"room_info"`
}

//GetRoomList 获取直播间列表
func (client *WeAppClient) GetRoomList(param *GetRoomListParam) (res *GetRoomListResponse, err error) {
	res = &GetRoomListResponse{}

	data, err := client.PostWithToken(getLiveInfoURL, param)
	if nil != err {
		return
	}
	err = json.Unmarshal(data, res)
	if nil != err {
		return
	}
	err = res.Check()
	return
}

//GetRoomPlaybackParam 获取直播间回放参数
type GetRoomPlaybackParam struct {
	Action string `json:"action"` //默认值 get_replay
	RoomID int64  `json:"room_id"`
	Start  int64  `json:"start"` // 0 起始拉取房间，start = 0 表示从第 1 个房间开始拉取
	Limit  int64  `json:"limit"` // 10 每次拉取的个数上限，不要设置过大，建议 100 以内
}

//Replay 回放
type Replay struct {
	ExpireTime string `json:"expire_time"` //回放视频url过期时间
	CreateTime string `json:"create_time"` //回放视频创建时间
	MediaURL   string `json:"media_url"`   //回放视频链接
	Total      int64  `json:"total"`       //回放视频片段个数
}

//GetRoomPlaybackResponse 获取直播间列表结果
type GetRoomPlaybackResponse struct {
	core.WxErrorResponse
	LiveReplay []*Replay `json:"live_replay"` // 回放列表
	Total      int64     `json:"total"`       // 总数
}

//GetRoomPlayback 获取直播间回放
func (client *WeAppClient) GetRoomPlayback(param *GetRoomPlaybackParam) (res *GetRoomPlaybackResponse, err error) {
	res = &GetRoomPlaybackResponse{}

	data, err := client.PostWithToken(getLiveInfoURL, param)
	if nil != err {
		return
	}
	err = json.Unmarshal(data, res)
	if nil != err {
		return
	}
	err = res.Check()
	return
}

//AddRoomGoodsParam 直播间导入商品参数
type AddRoomGoodsParam struct {
	IDs    []int64 `json:"ids"` // 数组列表，可传入多个，里面填写 商品 ID
	RoomID int64   `json:"roomId"`
}

//AddRoomGoods 直播间导入商品
func (client *WeAppClient) AddRoomGoods(param *AddRoomGoodsParam) (res *core.WxErrorResponse, err error) {
	res = &core.WxErrorResponse{}

	data, err := client.PostWithToken(addGoodsURL, param)
	if nil != err {
		return
	}
	err = json.Unmarshal(data, res)
	if nil != err {
		return
	}
	err = res.Check()
	return
}

//GoodsInfo 商品信息
type GoodsInfo struct {
	GoodsID       int64   `json:"goodsId"`     //商品ID 更新时必传
	CoverImgURL   string  `json:"coverImgUrl"` //填入mediaID（mediaID获取后，三天内有效）；max 300*300
	Name          string  `json:"name"`        //商品名称，最长14个汉字，1个汉字相当于2个字符
	PriceType     int     `json:"priceType"`   //价格类型，1：一口价（只需要传入price，price2不传） 2：价格区间（price字段为左边界，price2字段为右边界，price和price2必传） 3：显示折扣价（price字段为原价，price2字段为现价， price和price2必传）
	Price         float64 `json:"price"`       //数字，最多保留两位小数，单位元
	Price2        float64 `json:"price2"`      //数字，最多保留两位小数，单位元
	URL           string  `json:"url"`         //商品详情页的小程序路径，路径参数存在 url 的，该参数的值需要进行 encode 处理再填入
	ThirdPartyTag int     `json:"thirdPartyTag"`
}

//GoodsItem 商品信息返回数据
type GoodsItem struct {
	GoodsID       int64   `json:"goods_id"` //商品ID
	CoverImgURL   string  `json:"cover_img_url"`
	Name          string  `json:"name"`
	PriceType     int     `json:"price_type"`
	Price         float64 `json:"price"`
	Price2        float64 `json:"price2"`
	URL           string  `json:"url"`
	AuditStatus   int     `json:"audit_status"`
	ThirdPartyTag int     `json:"third_party_tag"`
}

//AddGoodsParam 商品添加并提审参数
type AddGoodsParam struct {
	GoodsInfo *GoodsInfo `json:"goodsInfo"`
}

//AddGoodsResponse 商品添加并提审结果
type AddGoodsResponse struct {
	core.WxErrorResponse
	GoodID  int64 `json:"goodsId"` //商品ID
	AuditID int64 `json:"auditId"` //审核单ID
}

//AddGoods 商品添加并提审
func (client *WeAppClient) AddGoods(param *AddGoodsParam) (res *AddGoodsResponse, err error) {
	res = &AddGoodsResponse{}

	data, err := client.PostWithToken(addRoomGoodsURL, param)
	if nil != err {
		return
	}
	err = json.Unmarshal(data, res)
	if nil != err {
		return
	}
	err = res.Check()
	return
}

//ResetAuditGoodsParam 撤回提审参数
type ResetAuditGoodsParam struct {
	GoodsID int64 `json:"goodsId"`
	AuditID int64 `json:"auditId"` //审核单ID
}

//ResetAuditGoods 撤回提审
func (client *WeAppClient) ResetAuditGoods(param *ResetAuditGoodsParam) (res *core.WxErrorResponse, err error) {
	res = &core.WxErrorResponse{}

	data, err := client.PostWithToken(resetAuditGoodsURL, param)
	if nil != err {
		return
	}
	err = json.Unmarshal(data, res)
	if nil != err {
		return
	}
	err = res.Check()
	return
}

//AuditGoodsParam 重新提交审核参数
type AuditGoodsParam struct {
	GoodsID int64 `json:"goodsId"`
}

//AuditGoodsResponse 重新提交审核结果
type AuditGoodsResponse struct {
	core.WxErrorResponse
	AuditID int64 `json:"auditId"` //审核单ID
}

//AuditGoods 重新提交审核
func (client *WeAppClient) AuditGoods(param *AuditGoodsParam) (res *AuditGoodsResponse, err error) {
	res = &AuditGoodsResponse{}

	data, err := client.PostWithToken(auditGoodsURL, param)
	if nil != err {
		return
	}
	err = json.Unmarshal(data, res)
	if nil != err {
		return
	}
	err = res.Check()
	return
}

//DeleteGoodsParam 删除商品参数
type DeleteGoodsParam struct {
	GoodsID int64 `json:"goodsId"`
}

//DeleteGoods 删除商品
func (client *WeAppClient) DeleteGoods(param *DeleteGoodsParam) (res *core.WxErrorResponse, err error) {
	res = &core.WxErrorResponse{}

	data, err := client.PostWithToken(deleteGoodsURL, param)
	if nil != err {
		return
	}
	err = json.Unmarshal(data, res)
	if nil != err {
		return
	}
	err = res.Check()
	return
}

//UpdateGoodsParam 更新商品参数
type UpdateGoodsParam struct {
	GoodsInfo *GoodsInfo `json:"goodsInfo"`
}

//UpdateGoods 更新商品
func (client *WeAppClient) UpdateGoods(param *UpdateGoodsParam) (res *core.WxErrorResponse, err error) {
	res = &core.WxErrorResponse{}

	data, err := client.PostWithToken(updateGoodsURL, param)
	if nil != err {
		return
	}
	err = json.Unmarshal(data, res)
	if nil != err {
		return
	}
	err = res.Check()
	return
}

//GetGoodsWarehouseParam 获取商品状态参数
type GetGoodsWarehouseParam struct {
	GoodsIDs []int64 `json:"goods_ids"`
}

//GetGoodsWarehouseResponse 获取商品状态结果
type GetGoodsWarehouseResponse struct {
	core.WxErrorResponse
	Goods []*GoodsItem `json:"goods"`
	Total int64        `json:"total"` // 总数
}

//GetGoodsWarehouse 获取商品状态
func (client *WeAppClient) GetGoodsWarehouse(param *GetGoodsWarehouseParam) (res *GetGoodsWarehouseResponse, err error) {
	res = &GetGoodsWarehouseResponse{}

	data, err := client.PostWithToken(getGoodsWarehouseURL, param)
	if nil != err {
		return
	}
	err = json.Unmarshal(data, res)
	if nil != err {
		return
	}
	err = res.Check()
	return
}

//GetGoodsApprovedParam 获取商品列表参数
type GetGoodsApprovedParam struct {
	Offset int64
	Limit  int
	Status int
}

//GetGoodsApprovedResponse 获取商品列表结果
type GetGoodsApprovedResponse struct {
	core.WxErrorResponse
	Goods []*GoodsInfo `json:"goods"`
	Total int64        `json:"total"`
}

//GetGoodsApproved 获取商品列表
func (client *WeAppClient) GetGoodsApproved(param *GetGoodsApprovedParam) (res *GetGoodsApprovedResponse, err error) {
	res = &GetGoodsApprovedResponse{}

	v := url.Values{}
	v.Set("offset", strconv.Itoa(int(param.Offset)))
	if param.Limit < 1 {
		param.Limit = 30
	} else if param.Limit > 100 {
		param.Limit = 100
	}

	v.Set("limit", strconv.Itoa(param.Limit))
	v.Set("status", strconv.Itoa(param.Status))
	data, err := client.GetResponseWithToken(getGoodsWarehouseURL, v)
	if nil != err {
		return
	}
	err = json.Unmarshal(data, res)
	if nil != err {
		return
	}
	err = res.Check()
	return
}
