package main

import (
	"github.com/chdlvy/go-chatFrame"
	"github.com/chdlvy/go-chatFrame/internal/msggateway"
	"log"
	"time"
)

func main() {
	chatFrame.Default()
	err := msggateway.RunWsAndServer(8080, 5*time.Second, 1000)
	if err != nil {
		log.SetFlags(log.Llongfile)
		log.Println(err)
	}
}
