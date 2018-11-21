package main

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"io/ioutil"
	"log"
	"strconv"
)

var redisConfig *Config
var redisClient *redis.Client

type Config struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
}

func readConfig() {
	log.Println("start reading config")
	if redisConfig == nil {
		data, _ := ioutil.ReadFile("./config/config.json")
		json.Unmarshal(data, &redisConfig)
	}
	log.Println("reading config complete")
}

func initRedis() {
	if redisConfig == nil {
		readConfig()
	}
	log.Println("start connecting redis")
	if redisClient == nil {
		address := redisConfig.Host + ":" + strconv.Itoa(redisConfig.Port)
		redisClient = redis.NewClient(&redis.Options{
			Addr:     address,
			Password: redisConfig.Password,
			DB:       0,
		})
	}
	log.Println("connecting redis complete")
}
