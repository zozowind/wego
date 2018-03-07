package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// DefaultHTTPClient default http client
var DefaultHTTPClient *http.Client

// DefaultMediaHTTPClient default http client for request media
var DefaultMediaHTTPClient *http.Client

func init() {
	shortTimeClient := *http.DefaultClient
	shortTimeClient.Timeout = time.Second * 5
	DefaultHTTPClient = &shortTimeClient

	longTimeclient := *http.DefaultClient
	longTimeclient.Timeout = time.Second * 60
	DefaultMediaHTTPClient = &longTimeclient
}

// HTTPJsonPost http post request with json in body
func HTTPJsonPost(httpClient *http.Client, url string, param interface{}) ([]byte, error) {
	data := []byte{}
	if httpClient == nil {
		httpClient = DefaultHTTPClient
	}

	//请求参数
	body, err := json.Marshal(param)
	if err != nil {
		return data, err
	}

	//请求
	httpResp, err := httpClient.Post(url, "application/json; charset=utf-8", bytes.NewReader(body))
	if err != nil {
		return data, err
	}

	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return data, fmt.Errorf("http.Status: %s", httpResp.Status)
	}
	data, err = ioutil.ReadAll(httpResp.Body)
	return data, err
}

// HTTPGet http get request
func HTTPGet(httpClient *http.Client, url string) ([]byte, error) {
	data := []byte{}
	if httpClient == nil {
		httpClient = DefaultHTTPClient
	}

	//请求
	httpResp, err := httpClient.Get(url)
	if err != nil {
		return data, err
	}

	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return data, fmt.Errorf("http.Status: %s", httpResp.Status)
	}

	data, err = ioutil.ReadAll(httpResp.Body)
	return data, err
}

// HTTPXMLPost http post request with XML in body
func HTTPXMLPost(httpClient *http.Client, url string, params url.Values) ([]byte, error) {
	if httpClient == nil {
		httpClient = DefaultHTTPClient
	}
	request := URLValueToXML(params)
	response := []byte{}
	var reader io.Reader
	if params != nil {
		reader = strings.NewReader(request)
	}
	resp, err := httpClient.Post(url, "text/xml", reader)
	if err != nil {
		return response, err
	}

	response, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return response, err
	}
	return response, nil
}

//RequestFile upload media file strcut
type RequestFile struct {
	Name string
	Data *bytes.Buffer
}

func (f *RequestFile) Read(p []byte) (n int, err error) {
	return f.Data.Read(p)
}

//HTTPFormPost http request using form post
func HTTPFormPost(httpClient *http.Client, url string, params url.Values, files map[string][]*RequestFile) ([]byte, error) {
	if httpClient == nil {
		httpClient = DefaultHTTPClient
	}
	//create an empty form
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	//add string params

	for key, s := range params {
		for _, v := range s {
			_ = bodyWriter.WriteField(key, v)
		}
	}

	//add file to upload
	for key, fs := range files {
		for _, f := range fs {
			//create file field
			fileWriter, err := bodyWriter.CreateFormFile(key, f.Name)
			if nil != err {
				return nil, err
			}
			//copy filedata to form
			_, err = io.Copy(fileWriter, f)
			if err != nil {
				return nil, err
			}
		}
	}

	// get upload content-type like multipart/form-data; boundary=...
	contentType := bodyWriter.FormDataContentType()

	// close bodyWriter now, not in deferr, it will add close tag to body
	bodyWriter.Close()

	response := []byte{}
	// post to server
	resp, err := httpClient.Post(url, contentType, bodyBuf)
	if err != nil {
		return response, err
	}

	response, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return response, err
	}
	return response, nil
}
