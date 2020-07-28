package pay

import (
	"fmt"
	"testing"
	"time"
)

var cli *Client

func init() {
	var err error
	cli, err = InitClient("your_appid", "your_appkey", true)
	if nil != err {
		panic(err.Error())
	}
}

func TestOrderUnified(t *testing.T) {
	req := &OrderUnifiedReq{}
	req.AppID = "your_open_id"
	req.DeviceInfo = "TestDeviceInfo"
	req.Body = "支付测试"
	req.Attach = "附加数据"
	req.OutTradeNo = fmt.Sprintf("%d", time.Now().Nanosecond())
	req.TotalFee = 101 //沙箱必须为101
	req.SpbillCreateIP = "your_ip"
	req.TimeStart = time.Now().Format("20060102150405")
	req.TimeExpire = time.Now().Add(15 * time.Minute).Format("20060102150405")
	req.NotifyURL = "your_notify_url"
	req.TradeType = TradeTypeJSAPI
	req.OpenID = "your_open_id"

	rsp, err := cli.OrderUnified(req)
	if nil != err {
		t.Errorf(err.Error())
		return
	}
	fmt.Printf("%#v", rsp)
	return
}
