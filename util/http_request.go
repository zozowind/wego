package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var DefaultHttpClient *http.Client
var DefaultMediaHttpClient *http.Client

func init() {
	shortTimeClient := *http.DefaultClient
	shortTimeClient.Timeout = time.Second * 5
	DefaultHttpClient = &shortTimeClient

	longTimeclient := *http.DefaultClient
	longTimeclient.Timeout = time.Second * 60
	DefaultMediaHttpClient = &longTimeclient
}

type WeBaseResponse struct {
	Code    int    `json:"errcode"`
	Message string `json:"errmsg"`
}

func HttpJsonPost(httpClient *http.Client, url string, param interface{}) ([]byte, error) {
	data := []byte{}
	if httpClient == nil {
		httpClient = DefaultHttpClient
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

func HttpGet(httpClient *http.Client, url string) ([]byte, error) {
	data := []byte{}
	if httpClient == nil {
		httpClient = DefaultHttpClient
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

func HttpXMLPost(client *http.Client, url string, params url.Values) ([]byte, error) {
	request := UrlValueToXml(params)
	response := []byte{}
	var reader io.Reader = nil
	if params != nil {
		reader = strings.NewReader(request)
	}
	resp, err := client.Post(url, "text/xml", reader)
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
