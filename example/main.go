package main

import (
	"fmt"
	"github.com/chdlvy/go-chatFrame"
	"github.com/chdlvy/go-chatFrame/pkg/common/config"
	"github.com/chdlvy/go-chatFrame/router"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"path/filepath"
	"time"
)

func main() {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("无法获取当前工作目录：", err)
		return
	}
	absolutePath := filepath.Join(currentDir, "/config/")
	config.InitConfig("config.yaml", absolutePath)

	go func() {
		r := gin.Default()
		route := &router.Route{}
		route.InitRouter(r)
		if err := r.Run(":8081"); err != nil {
			log.Println(err)
		}
	}()

	if err := chatFrame.RunWsServer(9090, 5*time.Second, 1000); err != nil {
		log.SetFlags(log.Llongfile)
		log.Println(err)
	}
}
