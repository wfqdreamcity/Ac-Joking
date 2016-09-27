package main

import (
	"log"
	"net/http"
	"router"
	"lib"
)

func main() {
	http.HandleFunc("/",index)

	//获取评论
	router.CommentRouter()

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

