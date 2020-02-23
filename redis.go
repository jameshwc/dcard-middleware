package main

import (
	"fmt"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/go-redis/redis/v7"
)

const dbConfigFile = "db.toml"

// This project uses Redis.

type Database interface {
	// Connect to the database server
	Init() error
	// Check if the IP is in database and whether it's forbidden or not
	Find(string) (bool, bool)
	// return X-RateLimit-Remaining and X-RateLimit-Reset
	GetKey(string) (int, string, error)
	// If IP is not found in database, then create one
	SetKey(string) error
	// Increment the visit counter of the IP, return X-RateLimit-Remaining
	IncrementVisitByIP(string) error
}

type redisServer struct {
	Config struct {
		Host, Name, Password string
		Port                 int
	}
	client *redis.Client
}

func (db *redisServer) Init() error {
	if _, err := toml.DecodeFile("db.toml", db); err != nil {
		return err
	}
	addr := fmt.Sprintf("%s:%d", db.Config.Host, db.Config.Port)
	db.client = redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   1, // use default DB
	})
	_, err := db.client.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}

func (db *redisServer) Find(ipaddr string) (bool, bool) {
	count, err := db.client.Get(ipaddr).Int()
	if err != nil && err != redis.Nil {
		fmt.Print(count, err)
		return false, false
	}
	if count >= maxIPInAnHour {
		return true, true
	}
	return true, false
}

func (db *redisServer) GetKey(ipaddr string) (int, string, error) {
	res, err := db.client.Get(ipaddr).Int()
	if err != nil {
		return 0, "", err
	}
	return res, db.client.TTL(ipaddr).Val().String(), nil
}

func (db *redisServer) SetKey(ipaddr string) error {
	err := db.client.Set(ipaddr, 1, time.Hour).Err()
	if err != nil {
		return err
	}
	return nil
}

func (db *redisServer) IncrementVisitByIP(ipaddr string) error {
	if _, err := db.client.Incr(ipaddr).Result(); err != nil {
		return err
	}
	return nil
}

func (db *redisServer) Reset() {
	db.client.FlushDB()
}
