package lib

import (
	"encoding/json"
	"io"
	"net/http"
	"errors"
	"strconv"
	"fmt"
	"time"
)

//格式化成功输出
func Success(rw http.ResponseWriter , args ...interface{}){

	var data map[string]interface{}
	var list string

	data = make(map[string]interface{})

	data["statusCode"] =200
	data["msg"] ="ok"
	data["data"]=args[0]

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

	data["statusCode"] =400
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
	if r.Method == "GET" {
		for _, v := range agrs {
			if key, ok := v.(string); ok {
				if len(r.Form[key]) > 0 {
					para[key] = r.Form[key][0]
				} else {
					Error(rw, "oop " + key + " is required !!")
					return para, false
				}
			}

		}
	}else if r.Method =="POST" {
		for _, v := range agrs {
			if key, ok := v.(string); ok {
				if len(r.Form[key]) > 0 {
					para[key] = r.Form[key][0]
				} else {
					Error(rw, "oop " + key + " is required !!")
					return para, false
				}
			}

		}
	}

	return para , true
}

//默认分页处理
func GetPageAndSize(r *http.Request) (int , int){

	page := 0
	size := 5

	r.ParseForm()
	if len(r.Form["page"]) > 0 {
		page , _ = strconv.Atoi(r.Form["page"][0])
		page = page - 1
	}
	if len(r.Form["size"]) > 0 {
		size , _ = strconv.Atoi(r.Form["size"][0])
	}

	start := page * size

	return start , size

}

//验证错误(提取公共方法预留日志处理)
func CheckError(err error) (bool , error) {
	if err != nil {
		time := time.Now().Format("2006-01-02 15:04:05")
		fmt.Println(time+": ",err)
		return false , err
	}else {
		return true  , err
	}
}

//时间格式化输出
func TimeFormat(pubTime int64) string {

	var timeFormat string
	currentTime := time.Now().Unix()
	pubTime = pubTime /1000
	tmp := currentTime - pubTime
	if tmp < 60 {
		timeFormat ="刚刚"
	}else if tmp < 3600 {
		tp := tmp/60
		tpt := strconv.FormatInt(tp,10)
		timeFormat =tpt+"分钟前"
	}else if tmp < 86400 {
		tp := tmp/3600
		tpt := strconv.FormatInt(tp,10)
		timeFormat = tpt+"小时前"
	}else if tmp < 259200 {//3天内
		tp := tmp/86400
		tpt := strconv.FormatInt(tp,10)
		timeFormat = tpt+"天前"
	}else{
		date := time.Unix(pubTime,0)
		timeFormat = date.Format("01-02 15:04")
	}

	return timeFormat
}

//获取cfa地址
func GetCfaAddressByUserId(userId string) string{
	url := cfaHost

	return url
}

//统一输出函数避免每次测试输出是都得引入fmt
func Println(args ...interface{})  {
	fmt.Println(args)
}