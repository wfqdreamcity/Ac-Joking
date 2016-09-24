package lib

import (
	"gopkg.in/couchbase/gocb.v1"
)

var cluster *gocb.Cluster

const couchbasehost  = "couchbase://tapi01.nomiss.hb02.allydata.cn"

func init__(){
	clu, _ := gocb.Connect(couchbasehost)
	cluster = clu
}

func OpenBucket(selectbucket string) (*gocb.Bucket ,error){

	bucket, err := cluster.OpenBucket(selectbucket,"")
	if err != nil {
		panic(err)
	}

	return bucket,nil

}