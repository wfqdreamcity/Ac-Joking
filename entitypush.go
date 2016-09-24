package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"strconv"
)


func main(){

	domain := "http://admin.nomiss.com"

	url := domain+"/bulkget"

	fmt.Println(url)

	resp ,err := http.Get(url)

	if err != nil {
		fmt.Println("获取数据错误")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	res := make([]interface{},0)


	erre := json.Unmarshal(body , &res)

	if erre != nil {
		fmt.Println(erre.Error())
	}

	//新建管道，并发数
	num :=	15

	chs := make([]chan int ,num)

	for i := 0 ; i < num ; i++ {
		chs[i] = make(chan int)
	}

	for i , v := range res {

		a := v.(map[string]interface{})

		number := i%num

		key := strconv.Itoa(i)
		a_ := a["id"]
		id_int := a_.(float64)
		id := strconv.FormatFloat(id_int, 'g', 6, 64)

		go httpGet(key, id, chs[number],number)

		if number == num-1 {
			for _ , ch := range chs {
				fmt.Println("channel out :" ,<-ch)
			}
		}

	}


}

func httpGet(key string, id string, ch chan int,number int) {



	domain := "http://admin.nomiss.com"

	url := domain+"/bulkepush?key="+key+"&id="+id

	fmt.Println(url)

	resp ,err := http.Get(url)

	if err != nil {
		fmt.Println("获取数据错误")
	}

	defer resp.Body.Close()

	ch <- number
	fmt.Println("counting")

}