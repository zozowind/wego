package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/go-redis/redis"
	"github.com/zozowind/wego/core"
)

const (
	tokenHost  = "http://127.0.0.1:7000"
	getPath    = "/token/get"
	verifyPath = "/token/verify"
)

var redisAddr = "127.0.0.1:6379"

type tokenClient struct {
	Name        string
	CacheServer core.TokenServer
}

func (c *tokenClient) getAndVerify() {
	for {
		n := rand.Intn(10) + 1
		time.Sleep(time.Duration(n) * time.Second)
		token, err := c.CacheServer.Token()
		if nil != err {
			fmt.Println(c.Name, " 获取token错误 ", err.Error())
			return
		}
		param := url.Values{}
		param.Set("token", token)
		rsp, err := http.Get(tokenHost + verifyPath + "?" + param.Encode())
		if nil != err {
			fmt.Println(c.Name, " 请求token验证错误 ", err.Error())
			return
		}
		data, err := ioutil.ReadAll(rsp.Body)
		if nil != err {
			fmt.Println(c.Name, " rsp读取错误 ", err.Error())
			return
		}
		fmt.Println(c.Name, string(data))
	}
}

func startClient(n int) (err error) {
	for i := 0; i < n; i++ {
		newClient(i)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("111"))
	})
	return http.ListenAndServe(":7001", nil)
}

func (c *tokenClient) tokenFunc() (token *core.AccessToken, err error) {
	param := url.Values{}
	param.Set("client", c.Name)
	rsp, err := http.Get(tokenHost + getPath + "?" + param.Encode())
	if nil != err {
		fmt.Println("[err]请求token错误", err.Error())
		return
	}
	defer rsp.Body.Close()
	data, err := ioutil.ReadAll(rsp.Body)
	if nil != err {
		fmt.Println("[err]rsp读取错误", err.Error())
		return
	}
	token = &core.AccessToken{}
	err = json.Unmarshal(data, token)
	if nil != err {
		fmt.Println("[err]解析错误", err.Error())
	}

	switch {
	case token.ExpiresIn > 60*3:
		token.ExpiresIn -= 60
	case token.ExpiresIn > 30:
		token.ExpiresIn -= 10
	default:
		err = fmt.Errorf("expires_in too small: %d", token.ExpiresIn)
		return token, err
	}

	return
}

func newClient(i int) (err error) {
	rc := redis.NewClient(&redis.Options{
		Addr:        redisAddr,
		Password:    "",
		DB:          0,
		PoolSize:    3,
		DialTimeout: 60 * time.Second,
	})

	redisCacheServer := &RedisCacheServer{
		Client:     rc,
		CacheKey:   fmt.Sprintf(tokenCacheKey, "testappid"),
		LockKey:    fmt.Sprintf(tokenLockKey, "testappid"),
		Expiration: tokenExpiration,
	}

	if nil != err {
		return
	}

	c := &tokenClient{
		Name: fmt.Sprintf("客户端-%d", i),
	}
	c.CacheServer = core.NewCacheTokenServer(redisCacheServer, c.tokenFunc)
	go c.getAndVerify()
	return
}
