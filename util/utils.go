package util

import (
	"reflect"
	"bytes"
	"k8s.io/apimachinery/pkg/util/yaml"
	"os"
	"github.com/jonluo94/cool/log"
)

const (
	Separator = "---"
)

var logger = log.GetLogger("xorm", log.ERROR)

// struct 转 map
func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

//Yamls 转 Jsons
func Yamls2Jsons(yamlBytes []byte) [][]byte {
	jsons := make([][]byte, 0)
	yamls := bytes.Split(yamlBytes, []byte(Separator))
	for _, v := range yamls {
		if len(v) == 0 {
			continue
		}
		obj, err := yaml.ToJSON(v)
		if err != nil {
			logger.Error(err.Error())
		}
		jsons = append(jsons, obj)
	}

	return jsons
}

//判断文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//创建文件夹
func CreatedDir(dir string) {
	exist, err := PathExists(dir)
	if err != nil {
		logger.Error("get dir error![%v]\n", err)
		return
	}

	if exist {
		logger.Info("has dir![%v]\n", dir)
	} else {
		logger.Info("no dir![%v]\n", dir)
		// 创建文件夹
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			logger.Error("mkdir failed![%v]\n", err)
		} else {
			logger.Info("mkdir success!\n")
		}
	}
}