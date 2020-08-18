package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/zozowind/wego/util"
)

// 直播相关接口

//注意点：
// 1、临时素材media_id是可复用的。
// 2、媒体文件在微信后台保存时间为3天，即3天后media_id失效。
// 3、上传临时素材的格式、大小限制与公众平台官网一致。
// 图片（image）: 10M，支持PNG\JPEG\JPG\GIF格式
// 语音（voice）：2M，播放长度不超过60s，支持AMR\MP3格式
// 视频（video）：10MB，支持MP4格式
// 缩略图（thumb）：64KB，支持JPG格式
const (
	uploadMediaURL = WxAPIURL + "/cgi-bin/media/upload?access_token=%s&type=%s" //临时素材上传
	getMediaURL    = WxAPIURL + "/cgi-bin/media/get?access_token=%s"            //临时素材获取
	//MediaTypeImage 临时素材图片
	MediaTypeImage = "image"
	// MediaTypeVoice 临时素材语音
	MediaTypeVoice = "voice"
	// MediaTypeVideo 临时素材视频
	MediaTypeVideo = "video"
	//MediaTypeThumb 临时素材缩略图
	MediaTypeThumb = "thumb"
)

//UploadMediaParam 上传临时素材
type UploadMediaParam struct {
	Type string
	Name string
	Data []byte
}

//UploadMediaResponse 获取临时素材结果
type UploadMediaResponse struct {
	WxErrorResponse
	Type      string `json:"type"`
	MediaID   string `json:"media_id"`
	CreatedAt int64  `json:"created_at"`
}

//UploadMedia 上传临时素材
func (client *WeBase) UploadMedia(param *UploadMediaParam) (res *UploadMediaResponse, err error) {
	res = &UploadMediaResponse{}
	data := []byte{}
	token, err := client.Token()
	if nil != err {
		return
	}
	f := &util.RequestFile{
		Name: param.Name,
		Data: bytes.NewBuffer(param.Data),
	}
	url := fmt.Sprintf(uploadMediaURL, token, param.Type)
	files := map[string][]*util.RequestFile{
		"media": []*util.RequestFile{
			f,
		},
	}
	data, err = util.HTTPFormPost(client.HTTPClient, url, nil, files)
	if nil != err {
		return
	}
	err = json.Unmarshal(data, res)
	if nil != err {
		return
	}
	err = res.Check()
	return
}

//GetMediaResponse 获取临时素材结果
type GetMediaResponse struct {
	WxErrorResponse
	VideoURL string `json:"video_url"`
}

//GetMedia 获取临时素材
func (client *WeBase) GetMedia(mediaID string) (res *GetMediaResponse, err error) {
	res = &GetMediaResponse{}

	v := url.Values{}
	v.Set("media_id", mediaID)
	data, err := client.GetResponseWithToken(getMediaURL, v)
	if nil != err {
		return
	}
	err = json.Unmarshal(data, res)
	if nil != err {
		return
	}
	err = res.Check()
	return
}
