package httputil

import (

	"net/http"
	"strings"
	"io/ioutil"
	"net/url"
	"os"
	"bytes"
	"mime/multipart"
	"io"
	"path/filepath"
	"github.com/jonluo94/cool/log"
	"github.com/jonluo94/cool/json"
)

var logger = log.GetLogger("httputil", log.ERROR)

func PostLocalFile(fileParam,filename string, targetUrl string,params map[string]string) []byte  {

	file, err := os.Open(filename)
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(fileParam, filepath.Base(filename))
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	_, err = io.Copy(part, file)
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	for key, val := range params {
		err = writer.WriteField(key, val)
		if err != nil {
			logger.Error(err.Error())
			return nil
		}
	}
	err = writer.Close()
	if err != nil {
		logger.Error(err.Error())
		return nil
	}

	request, err := http.NewRequest("POST", targetUrl, body)
	request.Header.Add("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	defer resp.Body.Close()
	//响应
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Read failed:", err)
		return nil
	}
	//返回结果
	return response
}


func PostMultiFile(file multipart.File,fileParam,filename, targetUrl string,params map[string]string) []byte {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(fileParam, filename)
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	_, err = io.Copy(part, file)
	if err != nil {
		logger.Error(err.Error())
		return nil
	}

	for key, val := range params {
		err = writer.WriteField(key, val)
		if err != nil {
			logger.Error(err.Error())
			return nil
		}
	}

	err = writer.Close()
	if err != nil {
		logger.Error(err.Error())
		return nil
	}

	request, err := http.NewRequest("POST", targetUrl, body)
	request.Header.Add("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	defer resp.Body.Close()
	//响应
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Read failed:", err)
		return nil
	}

	//返回结果
	return response
}

func PostJson(uri string, jsons interface{}) []byte {

	jsonParam, errs := json.Marshal(jsons) //转换成JSON返回的是byte[]
	if errs != nil {
		logger.Error(errs.Error())
	}

	//发送请求
	req, err := http.NewRequest("POST", uri, strings.NewReader(string(jsonParam)))
	if err != nil {
		logger.Error(err.Error())
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	defer resp.Body.Close()
	//响应
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Read failed:", err)
	}

	//返回结果
	return response
}

func PostForm(uri string, paras map[string][]string) []byte {

	resp, err := http.PostForm(uri, url.Values(paras))
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err.Error())
	}
	return body

}

func Get(uri string) []byte {

	resp, err := http.Get(uri)
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err.Error())
	}

	return body

}
