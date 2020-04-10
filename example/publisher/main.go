package main

import (
	"github.com/Sab94/udspubsub"
	"log"
)

func main(){
	client, err := udspubsub.NewPubSubClient("/tmp/udspubsub.sock")
	if err != nil {
		log.Println("Dialer err : ", err.Error())
		return
	}
	log.Println(client.ID)
	err = client.Publish("channelOne", []byte(`{ "votes": { "option_A": "3" } }`))
	if err != nil {
		log.Println("Publish err : ", err.Error())
		return
	}
	return
}
