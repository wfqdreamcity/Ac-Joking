package lib

import "net/http"

//中间层验证 在此注册方法
func MidleCheck(w http.ResponseWriter ,r *http.Request) bool {
	return  true
	//token验证
	if(CheckToken(w ,r)){
		return true
	}else{
		return false
	}
}