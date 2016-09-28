package news

import (
	"net/http"
	"lib"
	"couchbase"
)



func GetNewsList(rw http.ResponseWriter ,r *http.Request){

	lib.Success(rw ,"这是一个新闻列表页")
}


func GetNewsDetail(rw http.ResponseWriter ,r *http.Request){

	para , ok := lib.CheckParameter(rw , r, "newsId","userId")
	if !ok {
		return
	}

	com := couchbase.GetContent(para["newsId"])

	lib.Success(rw ,com)


}