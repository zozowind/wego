package media

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/zozowind/wego/core"
	"github.com/zozowind/wego/util"
)

const (
	//TempQrcodeIntActionName 临时二维码
	TempQrcodeIntActionName = "QR_SCENE"
	//TempQrcodeStrActionName 临时二维码
	TempQrcodeStrActionName = "QR_STR_SCENE"
	//LimitQrcodeIntActionName 永久二维码
	LimitQrcodeIntActionName = "QR_LIMIT_SCENE"
	//LimitQrcodeStrActionName 永久二维码
	LimitQrcodeStrActionName = "QR_LIMIT_STR_SCENE"

	qrcodeURL = core.WxAPIURL + "/cgi-bin/qrcode/create"
)

//QrcodeScene 二维码Scene
type QrcodeScene struct {
	SceneID  int64  `json:"scene_id"`
	SceneStr string `json:"scene_str"`
}

//QrcodeActionInfo 二维码ActionInfo
type QrcodeActionInfo struct {
	Scene QrcodeScene `json:"scene"`
}

//QrcodeReq 请求
type QrcodeReq struct {
	ExpireSeconds int64            `json:"expire_seconds"`
	ActionName    string           `json:"action_name"`
	ActionInfo    QrcodeActionInfo `json:"action_info""`
}

//QrcodeRsp 返回
type QrcodeRsp struct {
	Ticket        string `json:"ticket"`
	ExpireSeconds int64  `json:"expire_seconds"`
	URL           string `json:"url"`
}

// QrcodeTicket 获取二维码ticket
func (wm *WeMediaClient) QrcodeTicket(req *QrcodeReq) (qrcode *QrcodeRsp, err error) {
	token, err := wm.TokenServer.Token()
	if nil != err {
		return
	}
	params := url.Values{}
	params.Set("access_token", token)

	data, err := util.HTTPJsonPost(nil, qrcodeURL+"?"+params.Encode(), req)
	rsp := &core.WxErrorResponse{}
	err = json.Unmarshal(data, rsp)
	if nil != err {
		return
	}
	err = rsp.Check()
	if nil != err {
		return
	}
	qrcode = &QrcodeRsp{}
	err = json.Unmarshal(data, qrcode)
	return
}

//QrcodeURL 使用ticket换取二维码
func QrcodeURL(ticket string) string {
	ticket = url.QueryEscape(ticket)
	return fmt.Sprintf("https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket=%s", ticket)
}
