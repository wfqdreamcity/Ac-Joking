package couchbase

import (
	"lib"
	"gopkg.in/couchbase/gocb.v1"
)

type Content struct {
	Id string `json:"id"`
	Content string `json:"content"`
}

var bucketContent *gocb.Bucket

func init(){
	buc , err := lib.OpenBucket("content")
	if err != nil {
		panic(err)
	}

	bucketContent =buc
}

func GetContent(key string) *Content{

	var com Content

	bucketContent.Get(key,&com)

	return &com

}
