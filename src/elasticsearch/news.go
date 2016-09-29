package elasticsearch

import (
	"lib"
	"fmt"
	"encoding/json"
	"gopkg.in/olivere/elastic.v3"
	"strings"
	//"time"
)

type news struct {
	Click_count json.Number `json:"click_count"`
	Comment_count json.Number `json:"comment_count"`
	Create_time json.Number	`json:"create_time"`
	Duration json.Number `json:"duration"`
	Entity_ids string `json:"entity_ids"`
	Entity_names string `json:"entity_names"`
	First_class_ids string `json:"first_class_ids"`
	Id string `json:"id"`
	Is_hot string `json:"is_hot"`
	List_images string `json:"list_images"`
	List_images_style json.Number `json:"list_images_style"`
	New_source string `json:"new_source"`
	Org_url string `json:"org_url"`
	Pub_time json.Number `json:"pub_time"`
	Title string `json:"title"`
}

func Esearch(page int , size int ,userid string) ([]news ,error){

	list := make([]news , 0)
	start := page*size

	query := elastic.NewBoolQuery()
	//querymatch := elastic.NewMatchPhraseQuery("user","匿名")
	//query = query.Should(querymatch)

	//获取用户浏览记录
	value , _ := lib.Rclient.HGet("Userhistory",userid).Result()
	ids := make([]string,0)
	ids = strings.Split(value,",")
	for _ , v := range ids {
		queryIdTerm := elastic.NewTermQuery("id",v)
		query = query.MustNot(queryIdTerm)
	}

	searchResult , err := lib.Eclient.Search().
		Index("nm*").Type("news_video").Query(query).From(start).Size(size).Sort("pub_time",false).Pretty(true).Do()

	if err != nil {
		return list ,err
	}

	// Here's how you iterate through results with full control over each step.
	if searchResult.Hits.TotalHits > 0 {
		//fmt.Printf("Found a total of %d Joking\n", searchResult.Hits.TotalHits)
		//fmt.Printf("Found a maxscore of %d Joking\n", searchResult.Hits.MaxScore)

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

		fmt.Print("Found no Joking\n")
	}

	//处理用户浏览记录
	handTheValue(value,userid)

	return list ,nil
}

//type Joking struct {
//	Id string
//	User string
//	Content string
//	Time json.Number
//}
//
//func EsearchForStream(userid string,page , size int) ([]Joking ,error){
//
//	list := make([]Joking , 0)
//	start := page*size
//
//	time := time.Now().Format("2006-01-02 15:04:05")
//
//	fmt.Println(time)
//
//	query := elastic.NewBoolQuery()
//	//querymatch := elastic.NewMatchPhraseQuery("user","匿名")
//	//query = query.Should(querymatch)
//
//	//获取用户浏览记录
//	value , _ := lib.Rclient.HGet("Userhistory",userid).Result()
//	ids := make([]string,0)
//	ids = strings.Split(value,",")
//	for _ , v := range ids {
//		queryIdTerm := elastic.NewTermQuery("id",v)
//		query = query.MustNot(queryIdTerm)
//	}
//
//	//时间过滤
//	queryFilter := elastic.NewFilterAggregation()
//
//	query = query.Filter(queryFilter)
//
//
//	searchResult , err := lib.Eclient.Search().Query(query).
//		Index("crawler").Type("crawler").From(start).Size(size).Sort("time",false).Pretty(true).Do()
//
//	if err != nil {
//		return list ,err
//	}
//
//	// Here's how you iterate through results with full control over each step.
//	if searchResult.Hits.TotalHits > 0 {
//		fmt.Printf("Found a total of %d Joking\n", searchResult.Hits.TotalHits)
//		//fmt.Printf("Found a maxscore of %d Joking\n", searchResult.Hits.MaxScore)
//		//var jok map[string]string
//		// Iterate through results
//		for _, hit := range searchResult.Hits.Hits {
//			// hit.Index contains the name of the index
//
//			// Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
//			var t Joking
//			err := json.Unmarshal(*hit.Source, &t)
//			if err != nil {
//				// Deserialization failed
//				panic(err)
//			}
//
//			value = t.Id+","+value
//
//			// Work with Joking
//			list = append(list , t)
//		}
//	} else {
//		// No hits
//		fmt.Print("Found no Joking\n")
//	}
//
//	//处理用户浏览记录
//	handTheValue(value,userid)
//
//	return list ,nil
//}

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