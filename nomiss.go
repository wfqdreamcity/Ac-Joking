package main

import (
	"log"
	"net/http"
	"lib"
	"couchbase"
	"router"
)

func main() {
	http.HandleFunc("/",index)



	//获取评论
	http.HandleFunc("/getcommentsbyid",couchbase.GetComments)
	http.HandleFunc("/getcommentsbynewid",couchbase.GetCommentsByNewId)

	//http.HandleFunc("/",entity.GetStream)





	//公共路由
	router.CommonRouter()

	//实体api
	router.EntityRouter()




	err := http.ListenAndServe(":8888" , nil)
	if err != nil {
		log.Fatal("Listening fail port 8888 !")
	}

}


func index(rw http.ResponseWriter,r *http.Request) {

	message := "欢迎使用nomiss golang api"

	lib.Success(rw , message)
}

