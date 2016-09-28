package main

import (
	"log"
	"net/http"
	"router"
	"lib"
)


var ServeMuxCus = lib.DefaultServeMux

func main() {


	lib.HandleFunc("/",index)

	//引入路由
	router.Router()

	err := http.ListenAndServe(":8888" , ServeMuxCus)
	if err != nil {
		log.Fatal("Listening fail port 8888 !")
	}

}


func index(rw http.ResponseWriter,r *http.Request) {

	message := "欢迎使用nomiss golang api"

	lib.Success(rw , message)
}
