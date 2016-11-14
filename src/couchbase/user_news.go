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
func GetRelationNews(userId , newsId ,relateType string) ([]map[string]string , error){

	var query *gocb.N1qlQuery
	relate := make([]map[string]string,0)

	//最热评论
	query = gocb.NewN1qlQuery("SELECT end_id,start_id FROM user_news WHERE start_id = $1 and end_id in ($2) and look_relation=$3")

	rows, err := bucketRelationNews.ExecuteN1qlQuery(query, []interface{}{userId,newsId,relateType})
	if err != nil {
		return relate , err
	}
	defer rows.Close()

	row := make(map[string]string)
	for rows.Next(&row) {

		relate = append(relate , row)
	}

	return relate , err
}
