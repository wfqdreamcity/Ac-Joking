package lib

import (
	"gopkg.in/redis.v4"
	"fmt"
)

var Rclient *redis.Client

func init(){

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		DB:       0,  // use default DB
	})

	ping, err := client.Ping().Result()

	if err != nil {
		panic(err)
	}
	fmt.Println("redis is ok!",ping)

	Rclient = client
}
