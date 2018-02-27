package core

import (
	"net/http"
)

const WeMedia = "MEDIA" //公众号
const WeApp = "APP"     //小程序
const WeGame = "GAME"   //小游戏

type WeClient interface {
}

type WeBase struct {
	AppId       string
	AppSecret   string
	PayId       string //支付账号，一般为商户账号
	PayKey      string //支付key
	TokenServer TokenServer
	HttpClient  *http.Client
}
