package media

import (
	"math/rand"
	"time"

	"github.com/zozowind/wego/core"
)

//CacheJsTicketServer token server use cache server
type CacheJsTicketServer struct {
	TicketFunc           func() (*JsAPITicket, error)
	CacheServer          core.CacheServer
	RefreshTicketReqChan chan *JsAPITicket // chan
}

func abs(x time.Duration) time.Duration {
	if x >= 0 {
		return x
	}
	return -x
}

func retryTicket(attempts int, sleep time.Duration, fn func() (string, error)) (string, error) {
	token, err := fn()
	attempts--
	if nil != err && attempts > 0 {
		time.Sleep(sleep)
		return retryTicket(attempts, 2*sleep, fn)
	}
	return token, err
}

//NewCacheJsTicketServer get a new cache token server
func NewCacheJsTicketServer(cacheServer core.CacheServer, tokenFunc func() (*JsAPITicket, error)) *CacheJsTicketServer {
	srv := &CacheJsTicketServer{
		TicketFunc:           tokenFunc,
		CacheServer:          cacheServer,
		RefreshTicketReqChan: make(chan *JsAPITicket),
	}

	go srv.tokenUpdateDaemon(time.Hour * 2 * time.Duration(100+rand.Int63n(200)))
	return srv
}

func (cts *CacheJsTicketServer) tokenUpdateDaemon(initTickDuration time.Duration) {
	tickDuration := initTickDuration

NEW_TICK_DURATION:
	ticker := time.NewTicker(tickDuration)
	for {
		select {
		case jsAPITicket := <-cts.RefreshTicketReqChan:
			_ = cts.CacheServer.Set(jsAPITicket.Ticket, time.Duration(jsAPITicket.ExpiresIn)*time.Second)
			tickDuration = time.Duration(jsAPITicket.ExpiresIn) * time.Second
			ticker.Stop()
			goto NEW_TICK_DURATION
		case <-ticker.C:
			err := cts.CacheServer.Lock()
			if nil == err {
				jsAPITicket, err := cts.TicketFunc()
				if nil == err {
					_ = cts.CacheServer.Set(jsAPITicket.Ticket, time.Duration(jsAPITicket.ExpiresIn)*time.Second)
					newTickDuration := time.Duration(jsAPITicket.ExpiresIn) * time.Second
					if abs(tickDuration-newTickDuration) > time.Second*5 {
						tickDuration = newTickDuration
						ticker.Stop()
						goto NEW_TICK_DURATION
					}
				}
			}
			cts.CacheServer.Unlock()
		}
	}
}

//Ticket Token get ticket from CacheJsTicketServer
func (cts *CacheJsTicketServer) Ticket() (string, error) {
	ticket, err := cts.CacheServer.Get()
	if nil != err || ticket == "" {
		return cts.RefreshTicket()
	}
	return ticket, err
}

//RefreshTicket RefreshToken refresh ticket from CacheJsTicketServer
func (cts *CacheJsTicketServer) RefreshTicket() (string, error) {
	err := cts.CacheServer.Lock()
	if nil != err {
		//retry
		return retryTicket(3, 300*time.Millisecond, cts.CacheServer.Get)
	}
	jsAPITicket, err := cts.TicketFunc()
	if nil != err {
		return "", err
	}
	cts.RefreshTicketReqChan <- jsAPITicket
	return jsAPITicket.Ticket, nil
}
