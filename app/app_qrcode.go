package app

import "github.com/zozowind/wego/core"

//接口A: 适用于需要的码数量较少的业务场景 接口地址：
//https://api.weixin.qq.com/wxa/getwxacode?access_token=ACCESS_TOKEN
const (
	qrcodeURLA = core.WxAPIURL + "/wxa/getwxacode?access_token=%s"
	qrcodeURLB = core.WxAPIURL + "/wxa/getwxacodeunlimit?access_token=%s"
	qrcodeURLC = core.WxAPIURL + "/wxa/createwxaqrcode?access_token=%s"
)

//Qrcode qrcode interface
type Qrcode interface {
	URL() string
}

//RGB rgb color struct in qrcode
type RGB struct {
	R string `json:"r"`
	G string `json:"g"`
	B string `json:"b"`
}

//QrcodeA qrcode of type a
type QrcodeA struct {
	Path      string `json:"path"`
	Width     uint   `json:"width"`
	AutoColor bool   `json:"auto_color"`
	LineColor *RGB   `json:"line_color"`
}

//URL get qrcode url of type a
func (qr *QrcodeA) URL() string {
	return core.WxAPIURL + "/wxa/getwxacode?access_token=%s"
}

//QrcodeB qrcode of type b
type QrcodeB struct {
	Scene     string `json:"scene"`
	Page      string `json:"page"`
	Width     uint   `json:"width"`
	AutoColor bool   `json:"auto_color"`
	LineColor *RGB   `json:"line_color"`
}

//URL get qrcode url of type B
func (qr *QrcodeB) URL() string {
	return core.WxAPIURL + "/wxa/getwxacodeunlimit?access_token=%s"
}

//QrcodeC qrcode of type c
type QrcodeC struct {
	Path  string `json:"path"`
	Width uint   `json:"width"`
}

//URL get qrcode url of type c
func (qr *QrcodeC) URL() string {
	return core.WxAPIURL + "/wxa/createwxaqrcode?access_token=%s"
}

//GetQrcode get wechat app qrcode
func (client *WeAppClient) GetQrcode(qrcode Qrcode) ([]byte, error) {
	return client.PostResponseWithToken(qrcode.URL(), qrcode)
}
