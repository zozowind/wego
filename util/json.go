package util

import (
	"bytes"
	"encoding/json"
)

//JSONMarshal 对json中特殊字符不处理
func JSONMarshal(param interface{}) (data []byte, err error) {
	body := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(body)
	jsonEncoder.SetEscapeHTML(false)
	err = jsonEncoder.Encode(param)
	if err != nil {
		return
	}
	data = body.Bytes()
	return
}
