package lib

import (
	"gopkg.in/couchbase/gocb.v1"
	"fmt"
)

var cluster *gocb.Cluster

func init(){
	clu, err := gocb.Connect(couchbasehost)
	if err != nil {
		panic(err)
	}

	cluster = clu

	fmt.Println("Couchbase is ok !")

}

func OpenBucket(selectbucket string) (*gocb.Bucket ,error){

	bucket, err := cluster.OpenBucket(selectbucket,"")
	if err != nil {
		panic(err)
	}

	return bucket,nil

}