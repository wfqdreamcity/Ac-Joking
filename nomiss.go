package main

import (
	"log"
	"net/http"
	"lib"
	"entity"
	//"couchbase"
)



func main() {
	http.HandleFunc("/test",gethot)

	//实体相关api
	http.HandleFunc("/getentitylist",entity.GetEntityList)
	//http.HandleFunc("/getentityes",entity.IndexEsearch)

	//获取评论
	//http.HandleFunc("/getcommentsbyid",couchbase.GetComments)
	//http.HandleFunc("/getcommentsbynewid",couchbase.GetCommentsByNewId)

	http.HandleFunc("/",entity.GetStream)

	http.HandleFunc("/gettoken",lib.GetToken)
	http.HandleFunc("/checktoken",lib.CheckToken)
	err := http.ListenAndServe(":8888" , nil)
	if err != nil {
		log.Fatal("Listening fail!")
	}

}


func gethot(rw http.ResponseWriter,r *http.Request) {

	message := "欢迎使用nomiss golang api"

	lib.Error(rw , message)
}

