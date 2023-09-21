package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"server/msggateway"
	"server/pkg/chatlog"
	"server/pkg/common/config"
	"server/pkg/common/db"
	"server/pkg/friend"
	"server/pkg/group"
	"server/route"
	"time"
)

func main() {
	err := config.InitConfig("config.yaml", "/config/")
	if err != nil {
		log.Println("config init failed", err)
	}
	go func() {
		if err := group.Start(); err != nil {
			log.Println(err)
			return
		}
		if err := db.StartUserServer(); err != nil {
			log.Println(err)
			return
		}
		if err := friend.StartFriendServer(); err != nil {
			log.Println(err)
			return
		}
		if err := chatlog.StartChatLogServer(); err != nil {
			log.Println(err)
			return
		}

		r := gin.Default()
		tm := &route.TestMode{}
		tm.InitTestRouter(r)
		if err := r.Run(":8081"); err != nil {
			log.Println(err)
		}
		//_, err = db.NewGormDB()
		//if err != nil {
		//	log.Println("init gorm failed", err)
		//}
	}()

	msggateway.ReSetLongConnConf(8080, 5*time.Second, 1000)
	err = msggateway.RunWsAndServer()
	if err != nil {
		log.SetFlags(log.Llongfile)
		log.Println(err)
	}
}
