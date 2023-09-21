package main

import (
	"github.com/gin-gonic/gin"
	sql "server/internal/pkg/orm"
	rd "server/internal/pkg/redis"
	rt "server/route"
)

func main() {
	r := gin.Default()
	route := new(rt.Routes)
	go route.InitRouter(r)
	go rd.InitRedis()
	go sql.InitDB()
	r.Run(":8080")

}
