package couchbase

import (
	"lib"
	"net/http"
	"gopkg.in/couchbase/gocb.v1"
)

type Comment struct {
	Id string `json:"id"`
	News_id string `json:"news_id"`
	UserId string `json:"userId"`
	Up_count int `json:"up_count"`
	Comment string `json:"comment"`
	Type int `json:"type"`
	Post_time string `json:"post_time"`
}

var bucket *gocb.Bucket

func init(){
	buc , err := lib.OpenBucket("comments")
	if err != nil {
		panic(err)
	}

	bucket =buc
}

func GetComments(rw http.ResponseWriter ,r *http.Request){

	var com Comment
	var key string

	r.ParseForm()

	if len(r.Form["key"]) >0 {
		key = r.Form["key"][0]
	}else{
		lib.Error(rw , "请输入key值")
		return
	}
	bucket.Get(key,&com)

	lib.Success(rw ,&com)

}

//获取新闻评论
func GetCommentsByNewId(news_id string,ctype string,start int,size int) []interface{}{

	list := make([]interface{},0)

	var query *gocb.N1qlQuery
	// Use query
	if ctype == "new" {
		//最新评论
		query = gocb.NewN1qlQuery("SELECT * FROM comments WHERE news_id = $1 and up_count <= 3 ORDER BY post_time desc LIMIT $2 OFFSET $3")
	}else if ctype == "hot"{
		//最热评论
		query = gocb.NewN1qlQuery("SELECT * FROM comments WHERE news_id = $1 and up_count > 3 ORDER BY post_time DESC LIMIT $2 OFFSET $3")
	}

	rows, err := bucket.ExecuteN1qlQuery(query, []interface{}{news_id,size,start})
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	row := make(map[string]interface{})
	for rows.Next(&row) {

		list = append(list,row["comments"])
	}


	return list
}

//获取新闻评论总数
func GetCommentsCountById(news_id string) float64{
	var query *gocb.N1qlQuery

	query = gocb.NewN1qlQuery("SELECT COUNT(id) AS num From comments WHERE news_id = $1")

	rows, err := bucket.ExecuteN1qlQuery(query, []interface{}{news_id})
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var num float64
	row := make(map[string]interface{})
	for rows.Next(&row) {
		num = row["num"].(float64)
	}

	return num
}