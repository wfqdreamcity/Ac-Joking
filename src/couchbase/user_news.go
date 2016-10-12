package couchbase

import (
	"lib"
	"gopkg.in/couchbase/gocb.v1"
)

type UserNews struct {
	Relation string `json:"relation"`
	First_class_id string `json:"first_class_id"`
	End_id string `json:"end_id"`
	Start_id string `json:"start_id"`
}

var bucketRelationNews *gocb.Bucket

func init(){
	buc , err := lib.OpenBucket("user_news")
	if err != nil {
		panic(err)
	}

	bucketRelationNews =buc
}

//获取新闻评论
func GetRelationNews(userId , entityId string) bool{

	var query *gocb.N1qlQuery

	//最热评论
	query = gocb.NewN1qlQuery("SELECT count(*) FROM user_news WHERE start_id = $1 and end_id = $2 and relation=11")


	rows, err := bucketRelationNews.ExecuteN1qlQuery(query, []interface{}{userId,entityId})
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	row := make(map[string]interface{})
	for rows.Next(&row) {
		if row["count(*)"] != nil {
			return true
		}
	}

	return false
}
