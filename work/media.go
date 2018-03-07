package work

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

func (w *WeWorkClient) UploadMedia() {

}

func (w *WeWorkClient) GetMedia() {

}
