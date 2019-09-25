package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/zozowind/wego/core"
)

var (
	seed     int64
	duration int64
	startAt  int64
	ticker   *time.Ticker
)

func init() {
	duration = 300
	ticker = time.NewTicker(time.Duration(duration) * time.Second)
}

func getTokenHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Form.Get("client"), "请求token")
	token := serverGetToken()
	data, _ := json.Marshal(token)
	w.Write(data)
}

func verifyTokenHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if nil != err {
		w.Write([]byte(err.Error()))
		return
	}
	token := r.Form.Get("token")
	vToken := fmt.Sprintf("%d", seed)
	if token == vToken {
		w.Write([]byte(fmt.Sprintf("token %s 验证正确", token)))
	} else {
		w.Write([]byte(fmt.Sprintf("token %s 验证错误，应该是 %s", token, vToken)))
	}
}

func tickerGen(t *time.Ticker) {
	for {
		<-t.C
		genToken()
	}
}

func startServer() (err error) {
	genToken()
	go tickerGen(ticker)
	http.HandleFunc("/token/get", getTokenHandler)
	http.HandleFunc("/token/verify", verifyTokenHandler)
	return http.ListenAndServe(":7000", nil)
}

func genToken() {
	seed++
	startAt = time.Now().Unix()
	fmt.Println("生成新的Token ", seed, " 当前时间戳 ", startAt)
}

func serverGetToken() (token *core.AccessToken) {
	ticker.Stop()
	genToken()
	ticker = time.NewTicker(time.Duration(duration) * time.Second)
	go tickerGen(ticker)
	return &core.AccessToken{
		Token:     fmt.Sprintf("%d", seed),
		ExpiresIn: duration - time.Now().Unix() + startAt,
	}
}
