package lib

import (
	"encoding/json"
	"io"
	"net/http"
	"errors"
)

//格式化成功输出
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
//格式化错误输出
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
//token 验证（显示处理）
func CheckToken(rw http.ResponseWriter ,r *http.Request) bool{

	var tokens string
	r.ParseForm()
	if len(r.Form["token"]) >0 {
		tokens = r.Form["token"][0]
		_ , err := checkToken(tokens)

		if err != nil {
			Error(rw , err.Error())
			return false
		}

		return true

	}else{
		Error(rw , "请输入token值")
		return false
	}
}

//token 验证（逻辑处理）
func checkToken(token string) (bool , error){


	//redis token检测
	value , err := Rclient.Get("app_token_check:"+token).Result()
	if err != nil {
		return false , errors.New("token is out of expire or get the wrong token:"+err.Error())
	}

	if value != token {
		return false , errors.New("token is out of expire or get the wrong token")
	}

	return true , nil



}

//验证请求参数
//验证输入的参数是否是否缺失，
//agrs 为必要参数的字符串类型
func CheckParameter(rw http.ResponseWriter,r *http.Request,agrs ...interface{}) (map[string]string ,bool){
	para := make(map[string]string)

	r.ParseForm()
	for _, v := range agrs {
		if key , ok := v.(string);ok {
			if len(r.Form[key]) >0 {
				para[key] =r.Form[key][0]
			}else {
				Error(rw ,"oop "+key+" is required !!")
				return para , false
			}
		}

	}
	return para , true
}