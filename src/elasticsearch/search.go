package elasticsearch

import (
	"lib"
	"fmt"
	"encoding/json"
	"gopkg.in/olivere/elastic.v3"
	"strings"
	"time"
)

type Ad struct {
	Id json.Number
	Style json.Number
	Title string
	Img string
	Intro string
	List_images string
	Info_badge string
	Org_url  string
}

func Esearch(page , size int) ([]Ad ,error){

	list := make([]Ad , 0)
	start := page*size

	searchResult , err := lib.Eclient.Search().
		Index("other_card").Type("ad").From(start).Size(size).Sort("pub_time",false).Pretty(true).Do()

	if err != nil {
		return list ,err
	}

	// Here's how you iterate through results with full control over each step.
	if searchResult.Hits.TotalHits > 0 {
		//fmt.Printf("Found a total of %d Joking\n", searchResult.Hits.TotalHits)
		//fmt.Printf("Found a maxscore of %d Joking\n", searchResult.Hits.MaxScore)
		//var jok map[string]string
		// Iterate through results
		for _, hit := range searchResult.Hits.Hits {
			// hit.Index contains the name of the index

			// Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
			var t Ad
			err := json.Unmarshal(*hit.Source, &t)
			if err != nil {
				// Deserialization failed
				panic(err)
			}

			// Work with Joking
			list = append(list , t)
		}
	} else {
		// No hits
		fmt.Print("Found no Joking\n")
	}

	return list ,nil
}

type Joking struct {
	Id string
	User string
	Content string
	Time json.Number
}

func EsearchForStream(userid string,page , size int) ([]Joking ,error){

	list := make([]Joking , 0)
	start := page*size

	time := time.Now().Format("2006-01-02 15:04:05")

	fmt.Println(time)

	query := elastic.NewBoolQuery()
	//querymatch := elastic.NewMatchPhraseQuery("user","匿名")
	//query = query.Should(querymatch)

	value , _ := lib.Rclient.HGet("Userhistory",userid).Result()
	ids := make([]string,0)
	ids = strings.Split(value,",")
	for _ , v := range ids {
		queryIdTerm := elastic.NewTermQuery("id",v)
		query = query.MustNot(queryIdTerm)
	}


	searchResult , err := lib.Eclient.Search().Query(query).
		Index("crawler").Type("crawler").From(start).Size(size).Sort("time",false).Pretty(true).Do()

	if err != nil {
		return list ,err
	}

	// Here's how you iterate through results with full control over each step.
	if searchResult.Hits.TotalHits > 0 {
		fmt.Printf("Found a total of %d Joking\n", searchResult.Hits.TotalHits)
		//fmt.Printf("Found a maxscore of %d Joking\n", searchResult.Hits.MaxScore)
		//var jok map[string]string
		// Iterate through results
		for _, hit := range searchResult.Hits.Hits {
			// hit.Index contains the name of the index

			// Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
			var t Joking
			err := json.Unmarshal(*hit.Source, &t)
			if err != nil {
				// Deserialization failed
				panic(err)
			}

			value = value+","+t.Id

			// Work with Joking
			list = append(list , t)
		}
	} else {
		// No hits
		fmt.Print("Found no Joking\n")
	}

	lib.Rclient.HSet("Userhistory",userid,value)

	return list ,nil
}