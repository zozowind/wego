package app

import (
	"github.com/zozowind/wego/core"
)

type WeAppClient struct {
	Base *core.WeBase
}

type WeAppConfig struct {
}

func InitWeApp(base *core.WeBase) (*WeAppClient, error) {
	client := &WeAppClient{
		Base: base,
	}
	return client, nil
}
