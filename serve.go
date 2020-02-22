package main

import (
	"log"
	"net/http"
)

func main() {
	var db redisServer
	if err := db.Init(); err != nil {
		log.Fatal("redis server", err)
	}
	http.HandleFunc("/", limitVisit(hello, &db))
	if err := http.ListenAndServe(":8001", nil); err != nil {
		log.Fatal(err)
	}
}
