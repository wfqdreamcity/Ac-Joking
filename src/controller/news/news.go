package news

import (
	"net/http"
	"lib"
	"couchbase"
	"elasticsearch"
)



func GetNewsList(rw http.ResponseWriter ,r *http.Request){

	lib.Success(rw ,"这是一个新闻列表页")
}

//获取新闻正文接口
func GetNewsContent(rw http.ResponseWriter ,r *http.Request){

	para , ok := lib.CheckParameter(rw , r, "newsId","userId")
	if !ok {
		return
	}

	data := make(map[string]interface{})

	content := couchbase.GetContent(para["newsId"])

	data["id"] = content.Id
	data["n_c"] = content.Content

	source_other := couchbase.GetSource(para["newsId"])

	data["c_c"] = source_other.C_c
	data["m_c"] = source_other.M_c
	data["d_c"] = source_other.D_c
	data["tl"] = source_other.Tl
	data["u_c"] = source_other.U_c
	data["entity_names"] = source_other.Entity_names
	data["cl_c"] = source_other.Cl_c
	data["f_c"] = source_other.F_c
	data["p_d"] = source_other.P_d

	news := elasticsearch.GetNewsDetailById(para["newsId"])

	data["n_s"] = news




	lib.Success(rw ,data)


}

//获取新闻评论
func GetNewsComments(rw http.ResponseWriter ,r *http.Request){

	para , ok := lib.CheckParameter(rw , r ,"newsId","type")
	if !ok {
		return
	}

	start , size := lib.GetPageAndSize(r)

	comments := couchbase.GetCommentsByNewId(para["newsId"],para["type"],start,size)

	lib.Success(rw ,comments)
}