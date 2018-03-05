package work

import (
	"github.com/zozowind/wego/core"
)

const (
	// WxWorkAPIURL wechat work api ur
	WxWorkAPIURL = "https://qyapi.weixin.qq.com"
)

//WeWorkClient wechat app client struct
type WeWorkClient struct {
	core.WeBase
	AgentID int
}
