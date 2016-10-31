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

	//feed流相关
	lib.HandleFuncMiddle("/getnewslist",news.GetFeedNew)
	lib.HandleFuncMiddle("/gettopnew",news.GetTopNewList)
	lib.HandleFuncMiddle("/getfeedvideo",news.GetFeedVideo)
	lib.HandleFuncMiddle("/getfeedimage",news.GetFeedImage)

	//获取文章详情页
	lib.HandleFuncMiddle("/getnewscontent",news.GetNewsContent) //获取文章正文
	lib.HandleFuncMiddle("/getnewattributewithoutcontent",news.GetNewsAttributeByIdWithoutContent)//获取文章相关属性不包含正文

	//实体相关
	lib.HandleFuncMiddle("/getrelationforentityid",entity.GetRelationForEntityId)
	lib.HandleFuncMiddle("/getentitylist",entity.GetEntityList)


	//获取评论
	lib.HandleFuncMiddle("/getcommentsbyid",couchbase.GetComments)
	lib.HandleFuncMiddle("/getnewscomments",news.GetNewsComments)

	//http.HandleFunc("/",entity.GetStream)
}