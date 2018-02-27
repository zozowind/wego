package core

import (
	"fmt"
	"math/rand"
	"time"
)

type CacheServer interface {
	Get() (string, error)            //获取token
	Set(string, time.Duration) error //设置token
	Lock() error                     //锁
	Unlock()                         //解锁
}

type CacheTokenServer struct {
	TokenFunc           func() (*AccessToken, error)
	CacheServer         CacheServer
	RefreshTokenReqChan chan *AccessToken // chan
}

func abs(x time.Duration) time.Duration {
	if x >= 0 {
		return x
	}
	return -x
}

func retryToken(attempts int, sleep time.Duration, fn func() (string, error)) (string, error) {
	token, err := fn()
	attempts--
	if nil != err && attempts > 0 {
		time.Sleep(sleep)
		return retryToken(attempts, 2*sleep, fn)
	}
	return token, err
}

func NewCacheTokenServer(cacheServer CacheServer, tokenFunc func() (*AccessToken, error)) *CacheTokenServer {
	srv := &CacheTokenServer{
		TokenFunc:           tokenFunc,
		CacheServer:         cacheServer,
		RefreshTokenReqChan: make(chan *AccessToken),
	}

	go srv.tokenUpdateDaemon(time.Hour * 2 * time.Duration(100+rand.Int63n(200)))
	return srv
}

func (this *CacheTokenServer) tokenUpdateDaemon(initTickDuration time.Duration) {
	tickDuration := initTickDuration

NEW_TICK_DURATION:
	ticker := time.NewTicker(tickDuration)
	for {
		select {
		case accessToken := <-this.RefreshTokenReqChan:
			tickDuration = time.Duration(accessToken.ExpiresIn) * time.Second
			ticker.Stop()
			goto NEW_TICK_DURATION
		case <-ticker.C:
			err := this.CacheServer.Lock()
			if nil == err {
				accessToken, err := this.TokenFunc()
				if nil == err {
					newTickDuration := time.Duration(accessToken.ExpiresIn) * time.Second
					if abs(tickDuration-newTickDuration) > time.Second*5 {
						tickDuration = newTickDuration
						ticker.Stop()
						goto NEW_TICK_DURATION
					}
				}
			}
			this.CacheServer.Unlock()
		}
	}
}

func (this *CacheTokenServer) Token() (string, error) {
	token, err := this.CacheServer.Get()
	if nil != err || token == "" {
		return this.RefreshToken()
	}
	return token, err
}

func (this *CacheTokenServer) RefreshToken() (string, error) {
	err := this.CacheServer.Lock()
	if nil != err {
		//retry
		this.CacheServer.Unlock()
		return retryToken(3, 300*time.Millisecond, this.CacheServer.Get)
	}
	accessToken, err := this.TokenFunc()
	if nil != err {
		return "", err
	}
	return accessToken.Token, fmt.Errorf("code %d, msg %s", accessToken.ErrCode, accessToken.ErrMsg)
}
