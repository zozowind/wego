package app

import (
	"github.com/zozowind/wego/core"
)

//WeAppClient wechat app client struct
type WeAppClient struct {
	Base *core.WeBase
}

// InitWeApp init a wechat app client
func InitWeApp(base *core.WeBase) (*WeAppClient, error) {
	client := &WeAppClient{
		Base: base,
	}
	return client, nil
}
