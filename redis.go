package main

import (
	"fmt"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/go-redis/redis/v7"
)

const dbConfigFile = "db.toml"

// This project uses Redis.

type redisServer struct {
	Config struct {
		Host, Name, Password string
		Port                 int
	} `toml:"database"`
	client  *redis.Client
	maxIP   int
	timeout time.Duration
}

func (db *redisServer) Init(maxIP int, timeout int) error {
	if _, err := toml.DecodeFile("config.toml", db); err != nil {
		return err
	}
	addr := fmt.Sprintf("%s:%d", db.Config.Host, db.Config.Port)
	db.client = redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0, // use default DB
	})
	db.timeout = time.Duration(timeout) * time.Second
	db.maxIP = maxIP
	_, err := db.client.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}

func (db *redisServer) Find(ipaddr string) (existed bool, toomuch bool) {
	count, err := db.client.Get(ipaddr).Int()
	if err == redis.Nil {
		return false, false
	}
	if count >= db.maxIP {
		return true, true
	}
	return true, false
}

func (db *redisServer) GetKey(ipaddr string) (int, string, error) {
	res, err := db.client.Get(ipaddr).Int()
	if err != nil {
		return 0, "", err
	}
	remaining := db.maxIP - res
	if remaining < 0 {
		remaining = 0
	}
	return remaining, db.client.TTL(ipaddr).Val().String(), nil
}

func (db *redisServer) SetKey(ipaddr string) error {
	err := db.client.Set(ipaddr, 1, db.timeout).Err()
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
	db.client.FlushAll()
}
