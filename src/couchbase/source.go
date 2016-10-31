package couchbase

import (
	"lib"
	"gopkg.in/couchbase/gocb.v1"
	"encoding/json"
)

type Source struct {
	Id           string `json:"id"`
	M_c          json.Number `json:"m_c"`
	D_c          json.Number `json:"d_c"`
	Tl           string `json:"tl"`
	U_c          json.Number `json:"u_c"`
	C_c          json.Number `json:"c_c"`
	Entity_names []string `json:"entity_names"`
	Cl_c         json.Number `json:"cl_c"`
	F_c          json.Number `json:"f_c"`
	P_d          string `json:"p_d"`
}

var bucketSource *gocb.Bucket

func init(){
	buc , err := lib.OpenBucket("source_other")
	if err != nil {
		panic(err)
	}

	bucketSource =buc
}

func GetSource(key string) *Source{

	var com Source

	bucketSource.Get(key,&com)

	return &com

}
