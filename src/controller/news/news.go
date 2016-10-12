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
func GetNewsList(rw http.ResponseWriter ,r *http.Request){

	lib.Success(rw ,"这是一个新闻列表页")
}

//获取置顶新闻接口
func GetTopNewList(rw http.ResponseWriter , r *http.Request){

	//para , ok := lib.CheckParameter(rw , r,"userId")
	//if !ok {
	//	return
	//}

	NewsId , err := lib.Rclient.Keys("is_top_new:*").Result()
	if err != nil {
		lib.Error(rw , "oop,get top new fail !"+err.Error())
		return
	}
	var key string
	//data := make(map[string]interface{})

	for _ , v := range NewsId {
		key = v
	}

	str , err := lib.Rclient.Get(key).Bytes()
	if err != nil {
		lib.Error(rw , "oop ,获取json数据失败 ！"+err.Error())
		return
	}

	var news interface{}

	err = json.Unmarshal(str,&news)
	if err != nil {
		lib.Error(rw , "oop ,解析json错误 !"+err.Error())
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

	//卡片
	list := make(map[string]interface{})
	list["style"] = 13
	list["entity_id"] = newArray["entity_ids"]
	list["entity_name"] = newArray["entity_names"]
	//list["create_time"] = newArray["create_time"]
	list["img"] = newArray["img"]
	list["data"] = content



	lib.Success(rw,list)
}

//获取视频列表
func GetFeedVideo(rw http.ResponseWriter ,r *http.Request){

	page :=  1
	size := 10
	userId :="0"

	r.ParseForm()
	if len(r.Form["page"]) > 0 {
		page  , _ = strconv.Atoi(r.Form["page"][0])
	}
	if len(r.Form["size"]) > 0 {
		size , _ = strconv.Atoi(r.Form["size"][0])
	}
	if len(r.Form["userId"]) > 0 {
		userId = r.Form["userId"][0]
	}

	list , err := elasticsearch.Esearch(page , size , userId)

	if err != nil {
		lib.Error(rw , err.Error())
		return
	}

	lib.Success(rw ,list)
}

//获取新闻正文接口
func GetNewsContent(rw http.ResponseWriter ,r *http.Request){

	para , ok := lib.CheckParameter(rw , r, "newsId","userId","is_wifi")
	if !ok {
		return
	}

	data := make(map[string]interface{})

	content := couchbase.GetContent(para["newsId"])

	newscontent := content.Content

	//如果不是是wifi环境下使用小图
	if para["is_wifi"] == "no" {
		//获取图片url
		reg, _ := regexp.Compile(`src=['|"]{1}(.*?)['|"]{1}`)
		newurls := reg.FindAllStringSubmatch(content.Content, -1)
		urls := make(map[string]string)
		for _, v := range newurls {
			urls[v[1]] = v[1] + "/thumbnail/!12p/sourceurl=" + v[1]
		}
		//替换原链接
		for i , v := range urls {
			newscontent = strings.Replace(newscontent,i,v,1)
		}
	}

	data["id"] = content.Id
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