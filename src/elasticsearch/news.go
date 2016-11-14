package elasticsearch

import (
	"lib"
	"fmt"
	"encoding/json"
	"gopkg.in/olivere/elastic.v3"
	"strings"
	"errors"
	"redis"
	"datatype"
)

/*
*获取feed流信息
*start int 数据开始位置
*size int 每页显示数据条数
*UserId string 用户id
*news_type string 新闻类型：news_video 短视频 ，news_image 图集
*/
func GetListByEsearch(start int , size int ,UserId string , news_type string ){
	//判断管道是否为空（防止高并发时产异步加载产生大量请求）
	lib.GetEsSource(lib.Chs)

	list := make([]datatype.News , 0)

	query := elastic.NewBoolQuery()

	//获取用户浏览记录
	ids , _ := redis.GetRelationByObjaIdAndFlagFromRedis(UserId , "01" , 0 , 600)
	newHistoryIds := make([]string,0)
	for _ , v := range ids {
		queryIdTerm := elastic.NewTermQuery("id",v)
		query = query.MustNot(queryIdTerm)
	}

	searchResult , err := lib.Eclient.Search().
		Index("nm*").Type(news_type).Query(query).From(start).Size(size).Sort("create_time",false).Pretty(true).Do()

	ok , _ := lib.CheckError(err)
	if !ok {
		panic(err)
	}

	//释放es请求
	lib.LockSource(lib.Chs)

	// Here's how you iterate through results with full control over each step.
	if searchResult.Hits.TotalHits > 0 {

		for _, hit := range searchResult.Hits.Hits {

			var t datatype.News

			err := json.Unmarshal(*hit.Source, &t)
			if err != nil {
				panic(err)
			}

			newHistoryIds = append(newHistoryIds , t.Id)

			list = append(list , t)
		}
	} else {

		fmt.Print("Found no News\n")
	}

	//处理用户浏览记录(异步处理)
	go redis.SetUserHistoryIntoRedis(UserId , newHistoryIds)

	//处理预加载队列
	newsList , err := json.Marshal(list)
	ok , _ = lib.CheckError(err)
	if !ok {
		panic(err)
	}
	redis.BulkGetNewsInfoForFeedList(UserId , news_type ,newsList)

}

/*
*图集相关推荐接口
*newsId 新闻id
*title 新闻标题
*返回值 新闻数组分片 和 error
*/
func GetRelateByMulitKeyWord(newsId , news_type string ,Mulit map[string]interface{} ,start int , size int) ([]datatype.News ,error) {
	list := make([]datatype.News , 0)

	query := elastic.NewBoolQuery()

	if Mulit["title"] != nil {
		queryMatch := elastic.NewMatchQuery("title",Mulit["title"])
		query = query.Must(queryMatch)
	}


	queryTerm := elastic.NewTermsQuery("id",newsId)
	query = query.MustNot(queryTerm)

	//time := time.Now().AddDate(0,0,-3).Unix()*1000
	//queryTime := elastic.NewRangeQuery("pub_time").Gte(time)
	//query = query.Must(queryTime)

	searchResult , err := lib.Eclient.Search().
		Index("nm*").Type(news_type).Query(query).From(start).Size(size).Pretty(true).Do()

	if err != nil {
		return list ,err
	}

	if searchResult.Hits.TotalHits > 0 {

		for _, hit := range searchResult.Hits.Hits {

			var t datatype.News

			err := json.Unmarshal(*hit.Source, &t)
			if err != nil {
				panic(err)
			}

			list = append(list , t)
		}
	} else {

		fmt.Print("Found no News\n")
	}

	return list ,nil
}

/*
*图集相关推荐接口
*newsId 新闻id
*title 新闻标题
*返回值 新闻数组分片 和 error
*/
func GetRelateImageByTitle(newsId , title string) ([]datatype.News ,error) {
	list := make([]datatype.News , 0)

	query := elastic.NewBoolQuery()

	queryMatch := elastic.NewMatchQuery("title",title)
	query = query.Must(queryMatch)

	queryTerm := elastic.NewTermsQuery("id",newsId)
	query = query.MustNot(queryTerm)

	//time := time.Now().AddDate(0,0,-3).Unix()*1000
	//queryTime := elastic.NewRangeQuery("pub_time").Gte(time)
	//query = query.Must(queryTime)

	searchResult , err := lib.Eclient.Search().
		Index("nm*").Type("news_image").Query(query).From(0).Size(6).Pretty(true).Do()

	if err != nil {
		return list ,err
	}

	if searchResult.Hits.TotalHits > 0 {

		for _, hit := range searchResult.Hits.Hits {

			var t datatype.News

			err := json.Unmarshal(*hit.Source, &t)
			if err != nil {
				panic(err)
			}

			list = append(list , t)
		}
	} else {

		fmt.Print("Found no News\n")
	}

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
func GetNewsDetailById(id string) *datatype.News{

	var new datatype.News
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

//通过id获取新闻内容
func GetNewsDetailByIdContent(id string) (map[string]interface{} , error){

	new := make(map[string]interface{})
	queryTerm := elastic.NewTermQuery("id",id)

	searchResult , err := lib.Eclient.Search().Index("nm*").Query(queryTerm).Pretty(true).Do()

	if err != nil {
		return new , err
	}

	if searchResult.Hits.TotalHits > 0 {

		arr := searchResult.Hits.Hits

		err := json.Unmarshal(*arr[0].Source, &new)
		if err != nil {
			return new , err
		}

		return new , err

	}else{
		return new , errors.New("oop ,no maping !!!")
	}
}

/*通过id数组获取新闻列表
*ids id数组分片
*newsType Es中对应的type
*/
func GetNewsListByIdArray(ids []string , newsType string) ([]datatype.News , error){
	list := make([]datatype.News , 0)

	if len(ids) < 1 {
		return list , nil
	}

	query := elastic.NewBoolQuery()
	for _ , v := range ids {
		queryTerm := elastic.NewTermsQuery("id", v)
		query = query.Should(queryTerm)
	}
	searchResult , err := lib.Eclient.Search().
		Index("nm*").Type(newsType).Query(query).Sort("pub_time",false).Pretty(true).Do()

	if err != nil {
		return list ,err
	}

	if searchResult.Hits.TotalHits > 0 {

		for _, hit := range searchResult.Hits.Hits {

			var t datatype.News

			err := json.Unmarshal(*hit.Source, &t)
			if err != nil {
				panic(err)
			}

			list = append(list , t)
		}
	} else {

		fmt.Print("Found no News\n")
	}

	return list ,nil

}