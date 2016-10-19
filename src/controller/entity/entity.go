package entity

import (
	"lib"
	"net/http"
	"strconv"
	"elasticsearch"
	"couchbase"
)

type Entityinfo struct {
	Id string
	NAME string
	NICKNAMES string
	IMG string
}

//查看当前用户是否已经关注实体
func GetRelationForEntityId(rw http.ResponseWriter, r *http.Request){

	para , ok := lib.CheckParameter(rw,r,"userId","entityId")
	if !ok {
		return
	}

	ok = couchbase.GetRelationEntity(para["userId"],para["entityId"])

	lib.Success(rw , ok)
}

func GetEntityList(rw http.ResponseWriter , r *http.Request){

	start , size := lib.GetPageAndSize(r)

	rows, err := lib.DB.Query("select id,name,nicknames,img from entity limit ?,?",start, size)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	list := make([]Entityinfo  , 0)

	var personinfo Entityinfo


	var id string
	var name string
	var nicknames string
	var img string

	for rows.Next() {
		if err := rows.Scan(&id,&name,&nicknames,&img); err != nil {
			panic(err)
		}
		personinfo = Entityinfo{id,name,nicknames,img}
		list = append(list , personinfo)
	}

	lib.Success(rw , list)
}

func IndexEsearch(rw http.ResponseWriter ,r *http.Request){

	page :=  1
	size := 10
	userId :="0"

	r.ParseForm()
	if len(r.Form["page"]) > 0 {
		page  , _ = strconv.Atoi(r.Form["page"][0])
	}
	if len(r.Form["size"]) > 0 {
		size , _ = strconv.Atoi(r.Form["size"][0])
	}
	if len(r.Form["userId"]) > 0 {
		userId = r.Form["userId"][0]
	}

	list , err := elasticsearch.Esearch(page , size , userId)

	if err != nil {
		lib.Error(rw , err.Error())
		return
	}

	lib.Success(rw ,list)
}
//func GetStream(rw http.ResponseWriter ,r *http.Request){
//
//	page :=  1
//	size := 10
//	userId :="0"
//
//	r.ParseForm()
//	if len(r.Form["page"]) > 0 {
//		page  , _ = strconv.Atoi(r.Form["page"][0])
//	}
//	if len(r.Form["size"]) > 0 {
//		size , _ = strconv.Atoi(r.Form["size"][0])
//	}
//	if len(r.Form["userId"]) > 0 {
//		userId = r.Form["userId"][0]
//	}
//
//	list , err := elasticsearch.EsearchForStream(userId , page , size)
//
//	if err != nil {
//		lib.Error(rw , err.Error())
//		return
//	}
//
//	lib.Success(rw ,list)
//}
