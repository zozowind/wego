package mobile

import (
	"github.com/zozowind/wego/core"
)

//WeMobileClient 微信移动应用
type WeMobileClient struct {
	core.WeBase
	MessageConfig *core.MessageConfig
}
