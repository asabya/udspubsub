package main

import (
	"github.com/Sab94/udspubsub"
	"log"
	"os"
	"os/signal"
)

func main(){
	sigchan := make(chan os.Signal)
	signal.Notify(sigchan, os.Interrupt)

	client, err := udspubsub.NewPubSubClient("/tmp/udspubsub.sock")
	if err != nil {
		log.Println("Dialer err : ", err.Error())
		return
	}
	log.Println(client.ID)
	onMessage, err := client.Subscribe("channelOne")
	if err != nil {
		log.Println("Subscription err : ", err.Error())
		return
	}
	go func() {
		for {
			message := <- onMessage
			log.Println(message)
		}
	}()
	<- sigchan
	return
}
