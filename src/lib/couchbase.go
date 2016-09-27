package lib

import (
	"gopkg.in/couchbase/gocb.v1"
	"fmt"
)

var cluster *gocb.Cluster

func init(){
	clu, _ := gocb.Connect(couchbasehost)
	cluster = clu
}

func OpenBucket(selectbucket string) (*gocb.Bucket ,error){

	bucket, err := cluster.OpenBucket(selectbucket,"")
	if err != nil {
		panic(err)
	}

	fmt.Println("Couchbase is ok !")

	return bucket,nil

}