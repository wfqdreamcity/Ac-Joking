package elasticsearch

import (
	"lib"
	"fmt"
	"encoding/json"
	"gopkg.in/olivere/elastic.v3"
	"strings"
	//"time"
)

//新闻列表类型（不含content）
type news struct {
	Click_count json.Number `json:"click_count"`
	Comment_count json.Number `json:"comment_count"`
	Create_time json.Number	`json:"create_time"`
	Duration json.Number `json:"duration"`
	Entity_ids string `json:"entity_ids"`
	Entity_id string `json:"entity_id"`
	Entity_names string `json:"entity_names"`
	Entity_name string `json:"entity_name"`
	First_class_ids string `json:"first_class_ids"`
	Id string `json:"id"`
	Is_hot string `json:"is_hot"`
	List_images string `json:"list_images"`
	List_images_style json.Number `json:"list_images_style"`
	New_source string `json:"new_source"`
	Org_url string `json:"org_url"`
	Pub_time json.Number `json:"pub_time"`
	Title string `json:"title"`
	Image_list []string `json:"image_list"`
}

//新闻详情类型（含content）
type new struct {
	Click_count json.Number `json:"click_count"`
	Comment_count json.Number `json:"comment_count"`
	Create_time json.Number	`json:"create_time"`
	Duration json.Number `json:"duration"`
	Entity_ids string `json:"entity_ids"`
	Entity_id string `json:"entity_id"`
	Entity_names string `json:"entity_names"`
	Entity_name string `json:"entity_name"`
	First_class_ids string `json:"first_class_ids"`
	Id string `json:"id"`
	Is_hot string `json:"is_hot"`
	List_images string `json:"list_images"`
	List_images_style json.Number `json:"list_images_style"`
	New_source string `json:"new_source"`
	Org_url string `json:"org_url"`
	Pub_time json.Number `json:"pub_time"`
	Title string `json:"title"`
	Content string `json:"content"`
}


/*
*获取feed流信息
*start int 数据开始位置
*size int 每页显示数据条数
*UserId string 用户id
*news_type string 新闻类型：news_video 短视频 ，news_image 图集
*/
func GetListByEsearch(start int , size int ,UserId string , news_type string) ([]news ,error){
	list := make([]news , 0)

	query := elastic.NewBoolQuery()
	//querymatch := elastic.NewMatchPhraseQuery("user","匿名")
	//query = query.Should(querymatch)

	//获取用户浏览记录
	value , _ := lib.Rclient.HGet("Userhistory",UserId).Result()
	ids := make([]string,0)
	ids = strings.Split(value,",")
	for _ , v := range ids {
		queryIdTerm := elastic.NewTermQuery("id",v)
		query = query.MustNot(queryIdTerm)
	}

	searchResult , err := lib.Eclient.Search().
		Index("nm*").Type(news_type).Query(query).From(start).Size(size).Sort("pub_time",false).Pretty(true).Do()

	if err != nil {
		return list ,err
	}

	// Here's how you iterate through results with full control over each step.
	if searchResult.Hits.TotalHits > 0 {

		for _, hit := range searchResult.Hits.Hits {

			var t news

			err := json.Unmarshal(*hit.Source, &t)
			if err != nil {
				panic(err)
			}

			value = t.Id+","+value

			list = append(list , t)
		}
	} else {

		fmt.Print("Found no News\n")
	}

	//处理用户浏览记录
	handTheValue(value,UserId)

	return list ,nil

}

//处理用户浏览记录
func handTheValue(value ,userid string){
	idsArray :=make([]string , 0)

	historyNum :=300

	idsArray = strings.Split(value , ",")

	if len(idsArray) > historyNum {
		newsids := idsArray[0:historyNum]

		var idsString string

		for _ , v := range newsids {
			idsString = idsString+","+v
		}

		lib.Rclient.HSet("Userhistory",userid,idsString)
	}else {
		lib.Rclient.HSet("Userhistory",userid,value)
	}
}

//通过id获取新闻内容
func GetNewsDetailById(id string) *news{

	var new news
	queryTerm := elastic.NewTermQuery("id",id)

	searchResult , err := lib.Eclient.Search().Index("nm*").Type("news").Query(queryTerm).Pretty(true).Do()

	if err != nil {
		panic(err)
	}

	if searchResult.Hits.TotalHits > 0 {

		arr := searchResult.Hits.Hits


		err := json.Unmarshal(*arr[0].Source, &new)
		if err != nil {
			panic(err)
		}


	}


	return &new
}