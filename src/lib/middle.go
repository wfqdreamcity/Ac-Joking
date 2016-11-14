package lib

import (
	"net/http"
	"time"
	"fmt"
)

type Log struct {
	router string
	endNao int64
	startNao int64
	times int64
	timems int64
	timenas int64

}

func (f *Log)HandleTime(){
	mutime := f.endNao - f.startNao
	f.times = mutime / 1000000000
	f.timems = mutime /1000000
	f.timenas = mutime
}

var logInfo Log
//中间层验证 在此注册方法
func MidleCheck(w http.ResponseWriter ,r *http.Request) bool {

	logInfo.router = r.RequestURI
	logInfo.startNao = time.Now().UnixNano()

	return  true
	//token验证
	if(CheckToken(w ,r)){
		return true
	}else{
		return false
	}
}

//函数出口
func LastHandle(w http.ResponseWriter ,r *http.Request){

	logInfo.endNao = time.Now().UnixNano()
	current := time.Now().Format("2006-01-02 15:04:05")
	logInfo.HandleTime()
	fmt.Println(current , logInfo)

}