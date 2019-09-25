package main

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

const (
	tokenCacheKey   = "wx_token_cache_%s"
	tokenLockKey    = "wx_token_lock_%s"
	tokenExpiration = 10 * time.Second

	ticketCacheKey   = "wx_ticket_cache_%s"
	ticketLockKey    = "wx_ticket_lock_%s"
	ticketExpiration = 10 * time.Second

	tempTokenKey        = "wx_temp_token:%s"
	TempTokenExpiration = 10 * time.Minute
)

//RedisCacheServer 基于redis的缓存server
type RedisCacheServer struct {
	Client     *redis.Client
	CacheKey   string
	LockKey    string
	Expiration time.Duration
}

//TempToken 临时Token
type TempToken struct {
	AccessToken string `json:"accessToken"`
	OpenID      string `json:"openId"`
	AppType     string `json:"appType"`
	UnionID     string `json:"unionId"`
}

//NewRedisCacheServer 新建redisCacheServer
func NewRedisCacheServer(client *redis.Client, cacheKey string, lockKey string, expiration time.Duration) *RedisCacheServer {
	return &RedisCacheServer{
		Client:     client,
		CacheKey:   cacheKey,
		LockKey:    lockKey,
		Expiration: expiration,
	}
}

//Lock 锁定
func (rcs *RedisCacheServer) Lock() error {
	r := rcs.Client.SetNX(rcs.LockKey, "locked", rcs.Expiration)
	t, err := r.Result()
	if nil == err && !t {
		err = fmt.Errorf("set lock error")
	}
	return err
}

//Unlock 解锁
func (rcs *RedisCacheServer) Unlock() {
	_ = rcs.Client.Del(rcs.LockKey)
	return
}

//Get 获取
func (rcs *RedisCacheServer) Get() (string, error) {
	t := rcs.Client.Get(rcs.CacheKey)
	return t.Result()
}

//Set 设置
func (rcs *RedisCacheServer) Set(value string, expiration time.Duration) error {
	r := rcs.Client.Set(rcs.CacheKey, value, expiration)
	str, err := r.Result()
	if nil == err && str != "OK" {
		err = fmt.Errorf("set error: %s", err.Error())
	}
	return err
}
