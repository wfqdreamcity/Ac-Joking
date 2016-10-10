package couchbase

import (
	"lib"
	"gopkg.in/couchbase/gocb.v1"
)

type Relation struct {
	Relation string `json:"relation"`
	First_class_id string `json:"first_class_id"`
	End_id string `json:"end_id"`
	Start_id string `json:"start_id"`
}

var bucketRelationEntity *gocb.Bucket

func init(){
	buc , err := lib.OpenBucket("user_entity")
	if err != nil {
		panic(err)
	}

	bucketRelationEntity =buc
}

//获取新闻评论
func GetRelationEntity(userId , entityId string) bool{

	var query *gocb.N1qlQuery

	//最热评论
	query = gocb.NewN1qlQuery("SELECT count(*) FROM user_entity WHERE start_id = $1 and end_id = $2 and relation=11")


	rows, err := bucket.ExecuteN1qlQuery(query, []interface{}{userId,entityId})
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
