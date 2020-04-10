package main

import (
	"github.com/Sab94/udspubsub"
	"os"
	"os/signal"
)

func main() {
	sigchan := make(chan os.Signal)
	signal.Notify(sigchan, os.Interrupt)

	udspubsub.Listener("/tmp/udspubsub.sock")
	<- sigchan
	return
}
