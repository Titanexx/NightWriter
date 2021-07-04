package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Enviromnent struct {
	Port          string `json:"port"`
	AccessSecret  string `json:"access_secret"`
	RefreshSecret string `json:"refresh_secret"`
	// CasbinWatcherEnable bool   `json:"casbin_watcher_enable"`
	Db    DbConfiguration
	Redis RedisConfiguration
}

type RedisConfiguration struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
}

type DbConfiguration struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Dbname   string `json:"dbname"`
}

var env *Enviromnent

func LoadEnv() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	}

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redis := RedisConfiguration{Host: redisHost, Port: redisPort, Password: redisPassword}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	db := DbConfiguration{Host: dbHost, Port: dbPort, User: dbUser, Password: dbPass, Dbname: dbName}

	port := os.Getenv("APP_PORT")
	accessSecret := os.Getenv("ACCESS_SECRET")
	refreshSecret := os.Getenv("REFRESH_SECRET")

	env = &Enviromnent{
		Port:          port,
		AccessSecret:  accessSecret,
		RefreshSecret: refreshSecret,
		// CasbinWatcherEnable: watcherEnable,
		Redis: redis,
		Db:    db,
	}

	log.Print("Env loaded")
}

func GetEnv() *Enviromnent {
	if env == nil {
		LoadEnv()
	}
	return env
}
