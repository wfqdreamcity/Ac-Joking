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
	"controller/user"
	"math/rand"
	"redis"
	"fmt"
)


//新闻热点新闻列表页
func GetFeedNew(rw http.ResponseWriter ,r *http.Request){

	para , ok := lib.CheckParameter(rw ,r,"userId")
	if !ok {
		return
	}

	news := make([]map[string]interface{} , 0)
	var err error
	start , size := lib.GetPageAndSize(r)
	start = 0
	//获取数据
	paraMap := make(map[string]string)
	paraMap["user_id"] = para["userId"]
	paraMap["from"] = strconv.Itoa(start)
	paraMap["size"] = strconv.Itoa(size)
	result , err := lib.CfaGet("hot",paraMap)
	ok , _ = lib.CheckError(err)

	if ok {
		resultTemp := make(map[string]interface{})
		json.Unmarshal(result ,&resultTemp)

		if resultTemp["statusCode"] == "200" {
			resultTempNew := resultTemp["data"].(map[string]interface{})
			list := resultTempNew["data"].([]interface{})
			for _ , v := range list {
				newTemp := v.(map[string]interface{})
				//添加对应最热评论
				c := make([]interface{} , 3)
				comment := make([]interface{} , 0)
				comment = couchbase.GetCommentsByNewId(newTemp["id"].(string),"hot",0,3)
				copy(c, comment)
				//如果少于三条最新评论补全
				if len(comment) < 3 {
					num := 3 - len(comment)
					commentTemp := couchbase.GetCommentsByNewId(newTemp["id"].(string),"new",0,num)
					for _ , k := range commentTemp {
						comment = append(comment , k)
					}
				}
				////将评论组装成指定格式
				Fcomment := make([]map[string]interface{} , 0)
				for _ , k := range comment {
					commentT := make(map[string]interface{})
					vT := k.(map[string]interface{})
					commentT["name"] = vT["user_name"]
					commentT["content"] = vT["comment"]
					Fcomment = append(Fcomment , commentT)
				}
				newTemp["comment"] = Fcomment
				news = append(news , newTemp)
			}

			news , _ = getUserAndNewsRelate(para["userId"],news , "11")

		}else{
			fmt.Println(resultTemp["message"])
		}
	}

	lib.Success(rw , news)
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

	news := make([]map[string]interface{} , 0)
	var err error
	start , size := lib.GetPageAndSize(r)
	start = 0
	//获取数据
	//检查redis数据是否存在
	number := redis.CheckRedisBucketIsAvalible(para["userId"],"news_video")
	if number <= size {
		//es中 写入搜索数据 并写入redis 异步提前加载
		go elasticsearch.GetListByEsearch(0 ,300 ,para["userId"],"news_video")
	}
	//redis 获取数据
	ok , list := redis.GetNewsListFormRedis(para["userId"] ,"news_video" , start , size)

	if ok {

		for _ , v := range list {
			new := make(map[string]interface{})

			//通过id随机生成点击数
			rs := []rune(v.Id)
			var strTemp string
			for _ , v := range rs[4:20] {
				strTemp = strTemp+string(v)
			}
			dim , _ := strconv.ParseInt(strTemp, 10, 64)
			r := rand.New(rand.NewSource(dim))
			clickCount :=strconv.Itoa(r.Intn(100))
			new["id"] = v.Id
			new["click_count"]=clickCount+"万"
			new["tags"] = v.Tags
			new["comment_count"] = v.Comment_count
			new["type"] = "video"
			tempTime , err := v.Create_time.Float64()
			lib.CheckError(err)
			new["pub_time"] = lib.TimeFormat(int64(tempTime))
			new["duration"] = v.Duration
			new["title"] = v.Title
			new["news_source"] = v.News_source
			new["list_images"] = v.List_images
			new["landing_param"] = v.Org_url
			new["org_url"] = v.Org_url
			new["list_images_style"] = v.List_images_style
			new["link_type"] = "native"
			new["news_title"] = v.Title
			news = append(news , new)
		}
		//判断各个图集是否已经收藏
		news , err = getUserAndNewsRelate(para["userId"] , news ,"1")
		ok , _ = lib.CheckError(err)
	}

	lib.Success(rw ,news)
}

//获取图集列表
func GetFeedImage(rw http.ResponseWriter , r *http.Request){
	para , ok := lib.CheckParameter(rw ,r,"userId")
	if !ok {
		return
	}

	images := make([]map[string]interface{},0)
	var err error
	start , size := lib.GetPageAndSize(r)
	start = 0
	//获取数据
	//检查redis数据是否存在
	number := redis.CheckRedisBucketIsAvalible(para["userId"],"news_image")
	if number <= size {
		//es中 写入搜索数据 并写入redis(一次导入300条数据异步提前加载)
		go elasticsearch.GetListByEsearch(0 ,300 ,para["userId"],"news_image")

	}
	//redis 获取数据
	ok , list := redis.GetNewsListFormRedis(para["userId"] ,"news_image" , start , size)

	if ok {
		for i , v := range list {
			image := make(map[string]interface{})
			image["id"] = v.Id
			image["title"] = v.Title
			image["list_images"] = v.List_images
			image["list_images_style"] = v.List_images_style
			image["count"] = string(v.Comment_count)
			image["news_source"] = v.News_source
			tempTime , err := v.Create_time.Float64()
			lib.CheckError(err)
			image["pub_time"] = lib.TimeFormat(int64(tempTime))
			//将字符串转换成json
			image_list := make([]interface{},0)
			for _ , img := range v.Image_list {
				var imglist interface{}
				imgarr := []byte(img)
				json.Unmarshal(imgarr , &imglist)
				image_list = append(image_list , imglist)
			}
			image["image_list"] = image_list

			//配置展示类型style 12单图 14 3图一大二小 15 3图轮播
			strArray := strings.Split(v.List_images," ")
			count := len(strArray)
			if count == 1 {
				image["list_images"] = v.List_images
				image["style"] ="12"
			}else if count == 3 {
				num := i%3
				if num == 0 {
					//单图
					image["list_images"] = strArray[0]
					image["style"] ="12"
				}else if num == 1 {
					//三图一大二小
					image["list_images"] = v.List_images
					image["style"] ="14"
				}else{
					//三图轮播
					image["list_images"] = v.List_images
					image["style"] ="15"
				}
			}

			//只有封面图为一图或三图时加入召回列表
			if count == 1 || count == 3 {
				images = append(images , image)
			}

		}
		//判断各个图集是否已经收藏
		images , err = user.GetUserAndNewsRelate(para["userId"] , images ,"41")
		ok , _ = lib.CheckError(err)
	}

	lib.Success(rw , images)
}

//获取图集相关推荐
func GetRelateImageByTitle(rw http.ResponseWriter , r *http.Request){
	para , ok := lib.CheckParameter(rw ,r,"userId","title","newsId")
	if !ok {
		return
	}

	images := make([]map[string]interface{},0)

	list , err := elasticsearch.GetRelateImageByTitle(para["newsId"] , para["title"])

	ok , _ = lib.CheckError(err)

	if ok {
		for _ , v := range list {
			image := make(map[string]interface{})
			image["id"] = v.Id
			image["title"] = v.Title
			image["list_images"] = v.List_images
			image["list_images_style"] = v.List_images_style
			image["count"] = string(v.Comment_count)
			//将字符串转换成json
			image_list := make([]interface{},0)
			for _ , img := range v.Image_list {
				var imglist interface{}
				imgarr := []byte(img)
				json.Unmarshal(imgarr , &imglist)
				image_list = append(image_list , imglist)
			}
			image["image_list"] = image_list
			images = append(images , image)
		}
		//判断各个图集是否已经收藏
		images , err = user.GetUserAndNewsRelate(para["userId"] , images ,"41")
		ok , _ = lib.CheckError(err)
	}

	lib.Success(rw , images)
}


//获取新闻正文接口
func GetNewsContent(rw http.ResponseWriter ,r *http.Request){

	para , ok := lib.CheckParameter(rw , r, "newsId","userId","is_wifi")
	if !ok {
		return
	}

	data := make(map[string]interface{})

	//先去redis中获取
	var newcontent string
	newDetail , _ := redis.GetNewDetailByRedisWithContent(para["newsId"])
	if len(newDetail.Content) != 0 {
		newcontent = newDetail.Content
	}else{
		//es 中获取
		newDetail , err := elasticsearch.GetNewsDetailByIdContent(para["newsId"])
		ok , _ = lib.CheckError(err)
		if !ok {
			return
		}
		newcontent = newDetail["content"].(string)
	}

	newscontent :=newcontent

	//如果不是是wifi环境下使用小图
	if para["is_wifi"] != "1" {
		//获取图片url
		reg, _ := regexp.Compile(`src=['|"]{1}(.*?)['|"]{1}`)
		newurls := reg.FindAllStringSubmatch(newcontent, -1)
		urls := make(map[string]string)
		for _, v := range newurls {
			urls[v[1]] = v[1] + "/thumbnail/!36p/sourceurl=" + v[1]
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

//获取新闻评论总数
func GetNewsCommentsCount(rw http.ResponseWriter ,r *http.Request){

	para , ok := lib.CheckParameter(rw , r ,"newsId")
	if !ok {
		return
	}

	count := couchbase.GetCommentsCountById(para["newsId"])

	lib.Success(rw ,count)
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

//查询新闻是否已经收藏
func getUserAndNewsRelate(userId string ,list []map[string]interface{} ,relate_type string) ([]map[string]interface{} , error){
	newsIds :="0"

	newsList := make([]map[string]interface{} , 0)

	for _ , v := range list {
		newsIds = newsIds+","+v["id"].(string)
	}

	resp , err := couchbase.GetRelationNews(userId , newsIds ,relate_type)

	if err != nil {
		return newsList ,err
	}

	for _ , v := range list {
		new := make(map[string]interface{})

		new = v
		new["is_collected"] = "0"
		//判断是否已经收藏
		for _ , k := range resp {
			if k["end_id"] == v["id"] {
				new["is_collected"] = "1"
			}
		}
		newsList = append(newsList , new)
	}

	return newsList , err

}
