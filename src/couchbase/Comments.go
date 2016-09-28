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

func GetCommentsByNewId(rw http.ResponseWriter, r *http.Request){

	var news_id string

	list := make([]interface{},0)

	r.ParseForm()
	if len(r.Form["news_id"]) > 0 {
		news_id = r.Form["news_id"][0]
	}else {
		lib.Error(rw , "请输入news_id值")
		return
	}

	// Use query
	query := gocb.NewN1qlQuery("SELECT * FROM comments WHERE news_id = $1")
	rows, _ := bucket.ExecuteN1qlQuery(query, []interface{}{news_id})
	defer rows.Close()

	var row interface{}
	for rows.Next(&row) {
		list = append(list,row)
	}


	lib.Success(rw , list)
}
