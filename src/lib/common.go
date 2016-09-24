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

func GetToken(rw http.ResponseWriter , r *http.Request){
	r.ParseForm()

	//var token []map[string]string
	var appKey string
	var secretkey string

	if len(r.Form["appKey"]) == 0 {
		Error(rw,"参数appKey缺失")
		return
	}
	if len(r.Form["secretkey"]) ==0 {
		Error(rw,"参数secretkey缺失")
		return
	}

	appKey = r.Form["appKey"][0]
	secretkey = r.Form["secretkey"][0]

	list :=make([]interface{},0)

	list = append(list ,"通过"+appKey+secretkey+"获取token")

	Success(rw,list)
}

func CheckToken(rw http.ResponseWriter , r *http.Request){
	r.ParseForm()

	var token string

	if len(r.Form["token"]) ==0 {
		Error(rw ,"参数token 为必填字段")
		return
	}

	token = r.Form["token"][0]
	//token check
	err := Rclient.Set("token", token, 0).Err()
	if err != nil {
		panic(err)
	}

	value , ok := Rclient.Get("token").Result()
	if ok != nil {
		Error(rw , "读取redis 数据错误"+err.Error())
		return
	}

	list :=make([]interface{},0)

	list = append(list , "token 验证成功",token,value)

	Success(rw ,list)

}