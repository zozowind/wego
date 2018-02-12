package app

import "github.com/zozowind/wego/core"

//接口A: 适用于需要的码数量较少的业务场景 接口地址：
//https://api.weixin.qq.com/wxa/getwxacode?access_token=ACCESS_TOKEN
const (
	QrcodeUrlA = core.WxApiUrl + "/wxa/getwxacode?access_token=%s"
	QrcodeUrlB = core.WxApiUrl + "/wxa/getwxacodeunlimit?access_token=%s"
	QrcodeUrlC = core.WxApiUrl + "/wxa/createwxaqrcode?access_token=%s"
)

type Qrcode interface {
	Url() string
}

type RGB struct {
	R string `json:"r"`
	G string `json:"g"`
	B string `json:"b"`
}

type QrcodeA struct {
	Path      string `json:"path"`
	Width     uint   `json:"width"`
	AutoColor bool   `json:"auto_color"`
	LineColor *RGB   `json:"line_color"`
}

func (this *QrcodeA) Url() string {
	return core.WxApiUrl + "/wxa/getwxacode?access_token=%s"
}

type QrcodeB struct {
	Scene     string `json:"scene"`
	Path      string `json:"path"`
	Width     uint   `json:"width"`
	AutoColor bool   `json:"auto_color"`
	LineColor *RGB   `json:"line_color"`
}

func (this *QrcodeB) Url() string {
	return core.WxApiUrl + "/wxa/getwxacodeunlimit?access_token=%s"
}

type QrcodeC struct {
	Path  string `json:"path"`
	Width uint   `json:"width"`
}

func (this *QrcodeC) Url() string {
	return core.WxApiUrl + "/wxa/createwxaqrcode?access_token=%s"
}

func (this *WeAppClient) GetQrcode(qrcode Qrcode) ([]byte, error) {
	return this.Base.PostWithToken(qrcode.Url(), qrcode)
}
