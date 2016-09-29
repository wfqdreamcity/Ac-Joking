/*
  公共路由
*/
package router

import (
	"net/http"
	"couchbase"
	"lib"
	"controller/token"
	"controller/news"
	"controller/entity"
)
func HandleFuncRouter(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(pattern , handler)
}

func Router()  {

	//获取token
	lib.HandleFunc("/gettoken",token.GetToken)

	//获取文章详情页
	lib.HandleFuncMiddle("/getnewslist",news.GetNewsList)
	lib.HandleFuncMiddle("/getnewscontent",news.GetNewsContent)

	lib.HandleFuncMiddle("/getentitylist",entity.GetEntityList)
	lib.HandleFuncMiddle("/getentityes",entity.IndexEsearch)


	//获取评论
	lib.HandleFuncMiddle("/getcommentsbyid",couchbase.GetComments)
	lib.HandleFuncMiddle("/getnewscomments",news.GetNewsComments)

	//http.HandleFunc("/",entity.GetStream)
}