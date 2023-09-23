package main

import (
	"github.com/chdlvy/go-chatFrame"
	"log"
	"time"
)

func main() {

	chatFrame.Default()

	err := chatFrame.RunWsServer(8080, 5*time.Second, 1000)
	if err != nil {
		log.SetFlags(log.Llongfile)
		log.Println(err)
	}
}
