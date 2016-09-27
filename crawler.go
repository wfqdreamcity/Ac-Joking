package main

import (
	"net/http"
	"fmt"
	"time"
	//"bufio"
	//"os"
	//"strings"
)

func main() {

	url := "http://laravel.com/crawlerstart"

	var crawler_url string

	var ch chan int

	number :=1

	//for {
	//
	//	fmt.Println(`Enter following commands to control the crawler:  url website`)
	//
	//	r := bufio.NewReader(os.Stdin)
	//
	//	rawline ,_ ,_ := r.ReadLine()
	//
	//	line := string(rawline)
	//
	//	common := strings.Split(line," ")
	//
	//	if common[0] =="url" {
	//		crawler_url = common[1]
	//		break
	//	}else {
	//		crawler_url ="http://www.qiushibaike.com/"
	//	}

	crawler_url ="http://www.qiushibaike.com/"

	url = url+"?url="+crawler_url

	for i:= 0 ; i< number ; i++{

		go crawlerSingle(url)

	}
	<-ch



}

func crawlerMuilt(url string){
	i := 0
	for {
		_ ,err :=http.Get(url)

		if err != nil {
			fmt.Println("get require fail!")
		}

		time := time.Now().Format("2006-01-02 15:04:05")
		i++
		fmt.Println(time,"==>time : ",i)
	}
}

func crawlerSingle(url string){
	for {
		_ , err := http.Get(url)

		if err != nil {
			panic(err)
		}
	}
}
