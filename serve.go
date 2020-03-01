package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

func init() {
	f, err := os.OpenFile("limit.log", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal("error: create log file")
	}
	log.SetOutput(f)
}
func setupRouter(db Database) *gin.Engine {
	r := gin.Default()
	r.GET("/", limitVisit(db), hello)
	return r
}
func main() {
	var db redisServer
	if err := db.Init(1000, 3600); err != nil {
		log.Fatal("redis server", err)
	}
	db.Reset()
	r := setupRouter(&db)
	r.Run()
}
