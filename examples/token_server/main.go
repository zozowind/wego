package main

import (
	"flag"
)

//开发一个token_server server 和 客户端 用于测试access_token 锁

var (
	isClient  = flag.Bool("c", false, "is client")
	isServer  = flag.Bool("s", false, "is server")
	clientNum = flag.Int("n", 1, "client num")
)

func main() {
	var err error
	flag.Parse()
	defer func() {
		if nil != err {
			panic(err.Error)
		}
	}()
	if *isClient {
		if *clientNum > 0 {
			err = startClient(*clientNum)
			if nil != err {
				return
			}
		}
	}
	if *isServer {
		err = startServer()
		if nil != err {
			return
		}
	}
}
