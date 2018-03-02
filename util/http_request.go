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
func HTTPXMLPost(client *http.Client, url string, params url.Values) ([]byte, error) {
	request := URLValueToXML(params)
	response := []byte{}
	var reader io.Reader
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
