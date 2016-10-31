package news

import (
	"net/http"
	"lib"
	"couchbase"
	"elasticsearch"
	"regexp"
	"strings"
	"encoding/json"
	"strconv"
)


//新闻列表页
func GetFeedNew(rw http.ResponseWriter ,r *http.Request){

	para , ok := lib.CheckParameter(rw ,r,"userId")
	if !ok {
		return
	}

	start , size := lib.GetPageAndSize(r)

	list , err := elasticsearch.GetListByEsearch(start ,size ,para["userId"],"news")

	ok , _ = lib.CheckError(err)

	lib.Success(rw , list)
}

//获取置顶新闻接口
func GetTopNewList(rw http.ResponseWriter , r *http.Request){

	para , ok := lib.CheckParameter(rw , r,"userId")
	if !ok {
		return
	}

	err_result := ""
	//redis 查找，查看是否含有置顶新闻（正常情况下就一条数据）
	NewsId , err := lib.Rclient.Keys("is_top_new:*").Result()
	ok , _ = lib.CheckError(err)
	if !ok {
		lib.Success(rw ,err_result)
		return
	}

	var key string
	for _ , v := range NewsId {
		key = v
	}

	str , err := lib.Rclient.Get(key).Bytes()
	ok , _ = lib.CheckError(err)
	if !ok {
		lib.Success(rw , err_result)
		return
	}

	var news interface{}

	err = json.Unmarshal(str,&news)
	ok , _ = lib.CheckError(err)
	if !ok {
		lib.Success(rw , err_result)
		return
	}

	newArray := news.(map[string]interface{})

	//新闻详情
	content := make(map[string]interface{})
	content["id"] = newArray["id"]
	content["news_title"] = newArray["title"]
	content["list_images_style"] = newArray["list_images_style"]
	content["list_images"] = newArray["list_images"]
	content["news_source"] = newArray["news_source"]
	content["link_type"] = "native"
	id := newArray["id"]
	newid := id.(string)
	num_int := couchbase.GetCommentsCountById(newid)
	content["count"] = strconv.Itoa(int(num_int))

	//添加一层数组结果适配客户端
	result :=  make([]interface{},0)

	result = append(result,content)

	//卡片
	list := make(map[string]interface{})
	list["style"] = 13
	list["entity_id"] = newArray["entity_ids"]
	list["entity_name"] = newArray["entity_names"]
	//list["create_time"] = newArray["create_time"]
	list["img"] = newArray["img"]
	//查看当前用户是否关注实体
	if ok := couchbase.GetRelationEntity(para["userId"],newArray["entity_ids"].(string)) ; ok {
		list["is_followed"] = "1";
	}else{
		list["is_followed"] = "0";
	}
	list["data"] = result



	lib.Success(rw,list)
}

//获取视频列表
func GetFeedVideo(rw http.ResponseWriter ,r *http.Request){

	para , ok := lib.CheckParameter(rw ,r,"userId")
	if !ok {
		return
	}

	start , size := lib.GetPageAndSize(r)

	list , err := elasticsearch.GetListByEsearch(start , size , para["userId"],"news_video")

	ok , _ = lib.CheckError(err)

	lib.Success(rw ,list)
}

//获取图集列表
func GetFeedImage(rw http.ResponseWriter , r *http.Request){
	para , ok := lib.CheckParameter(rw ,r,"userId")
	if !ok {
		return
	}

	start , size := lib.GetPageAndSize(r)

	list , err := elasticsearch.GetListByEsearch(start ,size ,para["userId"],"news_image")

	ok , _ = lib.CheckError(err)

	lib.Success(rw , list)
}



//获取新闻正文接口
func GetNewsContent(rw http.ResponseWriter ,r *http.Request){

	para , ok := lib.CheckParameter(rw , r, "newsId","userId","is_wifi")
	if !ok {
		return
	}

	data := make(map[string]interface{})

	//content := couchbase.GetContent(para["newsId"])
	paras := make(map[string]string)
	paras["newsId"] = para["newsId"]
	content , err := lib.HbaseGet("/news/getContentById",paras)

	ok , _ = lib.CheckError(err)
	if !ok {
		return
	}

	new := make(map[string]interface{})
	json.Unmarshal(content,&new)

	if new["statusCode"] != "200" {
		lib.Error(rw , "hbase获取数据 失败！")
		return
	}

	newdetail := new["response"].(map[string]interface{})
	newcontent := newdetail["data"].(string)
	newscontent :=newcontent

	//如果不是是wifi环境下使用小图
	if para["is_wifi"] != "1" {
		//获取图片url
		reg, _ := regexp.Compile(`src=['|"]{1}(.*?)['|"]{1}`)
		newurls := reg.FindAllStringSubmatch(newcontent, -1)
		urls := make(map[string]string)
		for _, v := range newurls {
			urls[v[1]] = v[1] + "/thumbnail/!12p/sourceurl=" + v[1]
		}
		//替换原链接
		for i , v := range urls {
			newscontent = strings.Replace(newscontent,i,v,1)
		}
	}

	data["id"] = para["newsId"]
	data["n_c"] =newscontent

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
	data["is_wifi"] = para["is_wifi"]




	lib.Success(rw ,data)


}

//获取新闻属性信息（除正文es获取）
func GetNewsAttributeByIdWithoutContent(rw http.ResponseWriter ,r *http.Request){
	para , ok := lib.CheckParameter(rw , r, "newsId","userId")
	if !ok {
		return
	}

	new := elasticsearch.GetNewsDetailById(para["newsId"])

	lib.Success(rw,new)
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