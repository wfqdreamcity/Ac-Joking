package lib

import (
      "gopkg.in/olivere/elastic.v3"
	//"time"
	//"log"
	//"os"
	"fmt"
	"time"
	"log"
	"os"
)

var Eclient *elastic.Client

func init(){
	// Create a client
	cl, err := elastic.NewClient(
		elastic.SetURL(elastichost),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetMaxRetries(5),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)))

	if err != nil {
		// Handle error
		panic(err)
	}

	fmt.Println("elasticsearch is ok!")

	Eclient = cl
}