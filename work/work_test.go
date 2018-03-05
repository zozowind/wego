package work

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/zozowind/wego/core"
	"github.com/zozowind/wego/libs/errmsg"
)

var wework *WeWorkClient

func init() {
	wework = &WeWorkClient{}
	wework.AppID = "wwe1721397bbdb0977"
	wework.AppSecret = "sk7bXtdhhZ-d0QarAsMutT856LxoZh-2NmPV1449kNg"
	wework.AgentID = 1000002
	cacheServer := &core.DefaultMemoryCacheServer{}
	wework.TokenServer = core.NewCacheTokenServer(cacheServer, wework.RequestToken)
}

func Test_SendMessage(t *testing.T) {
	Convey("发送文本消息", t, func() {
		message := &TextMessage{
			Text: InMessageContent{
				Content: "hello world",
			},
		}
		res, err := wework.SendMessage(MessageTargetAll, []string{}, 0, message)
		if nil != err {
			t.Errorf("%s", err.(*errmsg.ErrMsg).Detail)
		} else {
			fmt.Printf("%#v", res)
		}
	})
}
