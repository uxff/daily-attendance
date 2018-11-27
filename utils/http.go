package utils

import (
	"io/ioutil"
	"net/http"
	"time"
	"net"
	"bytes"
	"os"
	"mime/multipart"
	"path/filepath"
	"io"
)

var clnt *http.Client

// todo
func init() {
	tr := &http.Transport{
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     60 * time.Second,
		Dial: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
	}
	clnt = &http.Client{
		Transport: tr,
		Timeout:   5 * time.Second,
	}
}

func httpReq(req *http.Request, headers map[string]string) (resp *http.Response, err error) {
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err = clnt.Do(req)
	if err != nil {
		return nil, err
	}

	return
}

func HttpGet(url string, headers map[string]string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return httpReq(req, headers)
}

func HttpPost(url string, headers map[string]string, data []byte) (resp *http.Response, err error) {
	body := bytes.NewReader(data)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	return httpReq(req, headers)
}


// 完整的HTTP GET
func HttpGetBody(url string, headers map[string]string) ([]byte, error) {
	resp, err := HttpGet(url, headers)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// 完整的HTTP POST
func HttpPostBody(url string, headers map[string]string, data []byte) ([]byte, error) {
	resp, err := HttpPost(url, headers, data)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}



// http上传文件
func PostFile(fullpath string, targetUrl string) ([]byte, error) {
	// 打开文件句柄操作
	fh, err := os.Open(fullpath)
	if err != nil {
		return nil, err
	}
	defer fh.Close()

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	// 关键的一步操作
	filename := filepath.Base(fullpath)
	fileWriter, err := bodyWriter.CreateFormFile("uploadfile", filename)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return nil, err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := clnt.Post(targetUrl, contentType, bodyBuf)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

