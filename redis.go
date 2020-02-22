package main

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/go-redis/redis"
)

const dbConfigFile = "db.toml"

// This project uses Redis.

type Database interface {
	// Connect to the database server
	Init() error
	// Check if the IP is in database and whether it's forbidden or not
	Find(string) (bool, bool)
	// return X-RateLimit-Remaining and X-RateLimit-Reset
	GetKeys(string) (int, string, error)
	// If IP is not found in database, then create one
	SetKeys(string) error
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
		Addr:     addr,
		Password: db.Config.Password, // no password set
		DB:       0,                  // use default DB
	})
	_, err := db.client.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}

func (db *redisServer) Find(ipaddr string) (bool, bool) {
	count, err := db.client.Get(ipaddr).Int()
	if err != nil {
		return false, false
	}
	if count > 1000 {
		return true, true
	}
	return true, false
}

func (db *redisServer) GetKey(ipaddr string) (int, string, error) {
	res, err := db.client.Get(ipaddr).Int()
	if err != nil {
		return 0, "", err
	}
	return res, db.client.TTL(ipaddr).String(), nil
}

func (db *redisServer) SetKey(ipaddr string) error {
	err := db.client.Set(ipaddr, 1, 3600).Err()
	if err != nil {
		return err
	}
	return nil
}

func (db *redisServer) IncrementVisitByIP(ipaddr string) error {
	if err := db.client.Incr(ipaddr).Err(); err != nil {
		return err
	}
	return nil
}
