package work

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/zozowind/wego/libs/errmsg"

	"github.com/zozowind/wego/core"
	"github.com/zozowind/wego/util"
)

const (
	mediaUploadURL = WxWorkAPIURL + "/cgi-bin/media/upload?access_token=%s&type=%s"
	//MediaTypeImage image
	MediaTypeImage = "image"
	//MediaTypeVoice voice
	MediaTypeVoice = "voice"
	//MediaTypeVideo video
	MediaTypeVideo = "video"
	//MediaTypeFile file
	MediaTypeFile = "file"
)

//WxWorkUploadMediaResponse wx work upload media response
type WxWorkUploadMediaResponse struct {
	core.WxErrorResponse
	Type      string `json:"type"`
	MediaID   string `json:"media_id"`
	CreatedAt string `json:"created_at"`
}

// UploadLocalMedia upload media from local path
func (w *WeWorkClient) UploadLocalMedia(t string, filePath string) (*WxWorkUploadMediaResponse, error) {
	filePath = "error.go"
	file, err := os.Open(filePath)
	if nil != err {
		return nil, errmsg.GetError(errUploadMedia, err.Error())
	}
	f := &util.RequestFile{
		Name: path.Base(file.Name()),
		Data: &bytes.Buffer{},
	}
	_, err = io.Copy(f.Data, file)
	if nil != err {
		return nil, errmsg.GetError(errUploadMedia, err.Error())
	}
	res, err := w.uploadMedia(t, f)
	if nil != err {
		return nil, errmsg.GetError(errUploadMedia, err.Error())
	}
	return res, err
}

// UploadMemoryMedia upload media from memory
func (w *WeWorkClient) UploadMemoryMedia(t string, name string, data []byte) (*WxWorkUploadMediaResponse, error) {
	f := &util.RequestFile{
		Name: name,
		Data: bytes.NewBuffer(data),
	}

	res, err := w.uploadMedia(t, f)
	if nil != err {
		return nil, errmsg.GetError(errUploadMedia, err.Error())
	}
	return res, err
}

func (w *WeWorkClient) uploadMedia(t string, f *util.RequestFile) (*WxWorkUploadMediaResponse, error) {
	data := []byte{}
	token, err := w.Token()
	if nil != err {
		return nil, err
	}
	url := fmt.Sprintf(mediaUploadURL, token, t)
	files := map[string][]*util.RequestFile{
		"media": []*util.RequestFile{
			f,
		},
	}
	data, err = util.HTTPFormPost(w.HTTPClient, url, nil, files)

	if nil != err {
		return nil, err
	}

	res := &WxWorkUploadMediaResponse{}
	err = json.Unmarshal(data, res)
	if nil != err {
		return nil, err
	}
	return res, nil
}

func (w *WeWorkClient) GetMedia() {

}
