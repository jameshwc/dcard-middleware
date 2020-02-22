package main

import (
	"log"
	"net/http"
)

func main() {
	var db database
	db.Init()
	http.HandleFunc("/", hello)
	if err := http.ListenAndServe(":8001", nil); err != nil {
		log.Fatal(err)
	}
}
