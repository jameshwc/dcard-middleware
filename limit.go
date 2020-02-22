package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const maxIPInAnHour int = 10

func getIP(r *http.Request) string {
	ip := r.Header.Get("X-Real-IP")
	fmt.Println(ip, r.RemoteAddr)
	if ip == "" {
		return strings.Split(r.RemoteAddr, ":")[0]
	}
	return ip
}
func limitVisit(next http.HandlerFunc, db Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := getIP(r)
		fmt.Println(ip)
		exist, tooMany := db.Find(ip)
		if !exist {
			db.SetKey(ip)
		} else {
			db.IncrementVisitByIP(ip)
		}
		count, ttl, err := db.GetKey(ip)
		if err != nil {
			log.Fatal("get redis key", err)
		}
		remaining := maxIPInAnHour - count
		w.Header().Add("X-RateLimit-Remaining", strconv.Itoa(remaining))
		w.Header().Add("X-RateLimit-Reset", ttl)
		if tooMany {
			w.WriteHeader(429)
		} else {
			w.WriteHeader(200)
		}
		next.ServeHTTP(w, r)
	}
}
