package lib

import (
      "gopkg.in/olivere/elastic.v3"
	"fmt"
	"time"
)

var Eclient *elastic.Client

//异步加载是锁定es请求
var EsChannel chan int

func init(){
	// Create a client
	cl, err := elastic.NewClient(
		elastic.SetURL(elastichost),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetMaxRetries(5))
		//elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		//elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)))

	if err != nil {
		// Handle error
		panic(err)
	}

	fmt.Println("elasticsearch is ok!")

	Eclient = cl

	//Create a channel
	EsChannel = make(chan int)
}