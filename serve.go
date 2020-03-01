package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	var db redisServer
	if err := db.Init(10, 10); err != nil {
		log.Fatal("redis server", err)
	}
	db.Reset()
	r := gin.Default()
	r.GET("/", limitVisit(hello, &db))
	r.Run(":8001")
}
