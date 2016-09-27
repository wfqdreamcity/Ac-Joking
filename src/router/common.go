/*
  公共路由
*/
package router

import (
	"net/http"
	"controller/token"
	"entity"
	"couchbase"
)
func HandleFuncRouter(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(pattern , handler)
}

func CommonRouter()  {

	http.HandleFunc("/gettoken",token.GetToken)
	//http.HandleFunc("/checktoken",token.CheckToken)
}



func EntityRouter(){
	HandleFuncRouter("/getentitylist",entity.GetEntityList)
	HandleFuncRouter("/getentityes",entity.IndexEsearch)
}

func CommentRouter(){
	//获取评论
	http.HandleFunc("/getcommentsbyid",couchbase.GetComments)
	http.HandleFunc("/getcommentsbynewid",couchbase.GetCommentsByNewId)

	//http.HandleFunc("/",entity.GetStream)
}