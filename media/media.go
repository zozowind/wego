package media

import (
	"github.com/zozowind/wego/core"
	"github.com/zozowind/wego/media/message"
)

//WeMediaClient wechat media client struct
type WeMediaClient struct {
	core.WeBase
	TicketServer  TicketServer
	MessageConfig *core.MessageConfig

	eventHandler map[message.EventType]func(*message.MixMessage) (*message.Reply, error)
}
