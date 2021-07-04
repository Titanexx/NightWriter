package config

import (
	"log"
	
	"gogomddoc/middleware/auth"

	"github.com/go-redis/redis/v7"
)

type RedisStruct struct {
	Client     *redis.Client
	Auth *auth.RedisAuthService
}

var Redis RedisStruct

func ConnectRedis() *redis.Client {
	configuration := GetEnv().Redis

	Redis.Client = redis.NewClient(&redis.Options{
		Addr:     configuration.Host  + ":" + configuration.Port,
		Password: configuration.Password,
		DB:       0,
	})

	_, err := Redis.Client.Ping().Result()
	if err != nil {
		log.Fatalf("Failed to connect to redis: %s",err)
	} else {
		log.Print("Success to connect to redis.")
	}

	Redis.Auth = auth.NewAuthService(Redis.Client )

	return Redis.Client 
}

