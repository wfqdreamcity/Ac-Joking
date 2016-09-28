/*
  公共路由
*/
package router

import (
	"net/http"
	"controller/token"
	"entity"
	"couchbase"
	"lib"
)
func HandleFuncRouter(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(pattern , handler)
}

func Router()  {

	//获取token
	lib.HandleFunc("/gettoken",token.GetToken)



	lib.HandleFuncMiddle("/getentitylist",entity.GetEntityList)
	lib.HandleFuncMiddle("/getentityes",entity.IndexEsearch)


	//获取评论
	lib.HandleFuncMiddle("/getcommentsbyid",couchbase.GetComments)
	lib.HandleFuncMiddle("/getcommentsbynewid",couchbase.GetCommentsByNewId)

	//http.HandleFunc("/",entity.GetStream)
}