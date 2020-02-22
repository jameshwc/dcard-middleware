package main

import (
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/go-redis/redis"
)

const dbConfigFile = "db.toml"

// This project uses PostgreSQL.

type database struct {
	Host, Name, Password string
	Port                 int
}

func (db *database) Init() {
	if _, err := toml.DecodeFile("db.toml", db); err != nil {
		log.Fatal("error when reading db.toml", err)
	}
	addr := fmt.Sprintf("%s:%d", db.Host, db.Port)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: db.Password, // no password set
		DB:       0,           // use default DB
	})
	_, err := client.Ping().Result()
	if err != nil {
		log.Fatal("Redis server error:", err)
	}
}
