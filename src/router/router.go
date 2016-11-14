/*
  公共路由
*/
package router

import (
	"lib"
	"couchbase"
	"controller/token"
	"controller/news"
	"controller/entity"
	"controller/user"
	"controller/rela"
)

func Router()  {

	//获取token
	lib.HandleFunc("/gettoken",token.GetToken)

	//feed流相关
	lib.HandleFuncMiddle("/getfeedhot",news.GetFeedNew)
	lib.HandleFuncMiddle("/gettopnew",news.GetTopNewList)
	lib.HandleFuncMiddle("/getfeedvideo",news.GetFeedVideo)
	lib.HandleFuncMiddle("/getfeedimage",news.GetFeedImage)

	//获取文章详情页
	lib.HandleFuncMiddle("/getnewscontent",news.GetNewsContent) //获取文章正文
	lib.HandleFuncMiddle("/getnewattributewithoutcontent",news.GetNewsAttributeByIdWithoutContent)//获取文章相关属性不包含正文
	lib.HandleFuncMiddle("/getrelateimagebytitle",news.GetRelateImageByTitle)//图集推荐

	//实体相关
	lib.HandleFuncMiddle("/getrelationforentityid",entity.GetRelationForEntityId)
	lib.HandleFuncMiddle("/getentitylist",entity.GetEntityList)


	//获取评论
	lib.HandleFuncMiddle("/getcommentsbyid",couchbase.GetComments)
	lib.HandleFuncMiddle("/getnewscomments",news.GetNewsComments)
	lib.HandleFuncMiddle("/getnewscommentscount",news.GetNewsCommentsCount)


	//对象关系相关接口（实体订阅，新闻收藏 等相关接口）
	lib.HandleFuncMiddle("/setrelation",rela.SetRelationUserIdAndNewsId) //设置对象之间的关系

	//用户收藏相关接口
	lib.HandleFuncMiddle("/getusercollectlist",user.GetUserCollectionList)
}