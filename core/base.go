package core

import (
	"net/http"
)

const (
	// WeMedia 公众号
	WeMedia = "MEDIA"
	// WeApp 小程序
	WeApp = "APP"
	// WeGame 小游戏
	WeGame = "GAME"
	// WeWork 企业微信内应用
	WeWork = "WORK"
)

// WeClient wechat client interface
type WeClient interface {
}

//WeBase wechat client base struct
type WeBase struct {
	AppID       string
	AppSecret   string
	PayID       string //支付账号，一般为商户账号
	PayKey      string //支付key
	TokenServer TokenServer
	HTTPClient  *http.Client
}
