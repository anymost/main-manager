package main

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/gpmgo/gopm/modules/log"
	"io/ioutil"
	"strconv"
)

var redisConfig *Config
var redisClient *redis.Client

type Config struct {
	Host string `json:"host"`
	Port int `json:"port"`
	Password string `json:"password"`
}

func readConfig() *Config  {
	log.Info("start reading config")
	if redisConfig == nil {
		data, _ := ioutil.ReadFile("./config/config.json")
		json.Unmarshal(data, &redisConfig)
	}
	log.Info("reading config complete")
	return redisConfig
}


func StartRedis() *redis.Client  {
	log.Info("start connecting redis")
	if redisClient == nil {
		config := readConfig()
		address := config.Host + ":" + strconv.Itoa(config.Port)
		redisClient = redis.NewClient(&redis.Options{
			Addr:     address,
			Password: config.Password,
			DB:       0,
		})
	}
	log.Info("connecting redis complete")
	return redisClient
}