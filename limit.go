package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// const maxIPInAnHour int = 1000

func limitVisit(db Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
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
			log.WithFields(log.Fields{
				"ip": ip,
			}).Info("Someone's ip has been forbidden")
			c.AbortWithStatus(429)
		} else {
			c.Status(200)
		}
		c.Next()
	}
}
