package lib

import (
	"net/http"
	"io/ioutil"
)

//通过dataserver获取数据
func HbaseGet(method string,para map[string]string) ([]byte , error){
	var url string
	url = hbasehost+method+"?token="+dataToken+"&"
	var body []byte

	for i ,v := range para {
		url = url+i+"="+v+"&"
	}

	resp ,  err := HttpGet(url)
	if err != nil {
		return body ,err
	}

	body, err = ioutil.ReadAll(resp.Body)

	return body , err

}
//通过cfa获取数据
func CfaGet(method string,para map[string]string) ([]byte , error){
	var url string
	url = cfaHost+method+"?token="+dataToken+"&"
	var body []byte

	for i ,v := range para {
		url = url+i+"="+v+"&"
	}

	resp ,  err := HttpGet(url)
	if err != nil {
		return body ,err
	}

	body, err = ioutil.ReadAll(resp.Body)

	return body , err

}

func HttpGet(url string) (*http.Response ,error){
	var resp *http.Response
	var err error

	resp , err =http.Get(url)

	return resp, err

}
