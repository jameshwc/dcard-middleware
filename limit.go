package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// const maxIPInAnHour int = 1000

func getIP(r *http.Request) string {
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		return strings.Split(r.RemoteAddr, ":")[0]
	}
	return ip
}
func limitVisit(db Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		fmt.Println("ip:", ip)
		exist, tooMany := db.Find(ip)
		if !exist {
			err := db.SetKey(ip)
			if err != nil {
				log.Fatal("Set redis key", err)
			}
		} else {
			err := db.IncrementVisitByIP(ip)
			if err != nil {
				log.Fatal("Increment redis key", err)
			}
		}
		remaining, ttl, err := db.GetKey(ip)
		if err != nil {
			log.Fatal("Get redis key", err)
		}
		c.Writer.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Writer.Header().Set("X-RateLimit-Reset", ttl)
		if tooMany {
			c.AbortWithStatus(429)
		} else {
			c.Status(200)
		}
		c.Next()
	}
}
