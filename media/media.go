package media

import (
	"github.com/zozowind/wego/core"
)

//WeAppClient wechat app client struct

type WeMediaClient struct {
	core.WeBase
	TicketServer  TicketServer
	MessageConfig *core.MessageConfig
}
