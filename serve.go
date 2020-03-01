package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func setupRouter(db Database) *gin.Engine {
	r := gin.Default()
	r.GET("/", limitVisit(db), hello)
	return r
}
func main() {
	var db redisServer
	if err := db.Init(10, 10); err != nil {
		log.Fatal("redis server", err)
	}
	db.Reset()
	r := setupRouter(&db)
	r.Run()
}
