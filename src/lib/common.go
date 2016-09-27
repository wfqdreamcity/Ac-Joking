package lib

import (
	"encoding/json"
	"io"
	"net/http"
)


func Success(rw http.ResponseWriter , args ...interface{}){

	var data map[string]interface{}
	var list string

	data = make(map[string]interface{})

	data["code"] =200
	data["msg"] ="ok"
	data["data"]=args

	result , err := json.Marshal(data)

	if err != nil {
		io.WriteString(rw ,"格式化数据错误")
	}

	list = string(result)

	io.WriteString(rw ,list)

}

func Error(rw http.ResponseWriter , args ...interface{}){
	var data map[string]interface{}
	var list string

	data = make(map[string]interface{})

	data["code"] =400
	data["msg"] =args
	data["data"]=nil

	result , err := json.Marshal(data)

	if err != nil {
		io.WriteString(rw ,"格式化数据错误")
	}

	list = string(result)

	io.WriteString(rw ,list)
}