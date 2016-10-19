package lib

import (
	"net/http"
	"io/ioutil"
)

func HbaseGet(method string,para map[string]string) ([]byte , error){
	var url string
	url = hbasehost+method+"?token="+dataToken+"&"

	for i ,v := range para {
		url = url+i+"="+v+"&"
	}

	resp ,  err := HttpGet(url)

	body, err := ioutil.ReadAll(resp.Body)

	return body , err

}

func HttpGet(url string) (*http.Response ,error){
	var resp *http.Response
	var err error

	resp , err =http.Get(url)

	return resp, err

}
