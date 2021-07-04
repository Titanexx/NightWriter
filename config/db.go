package config

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDb() {
	dbEnv := GetEnv().Db
	dsn := "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable "
	dsn = fmt.Sprintf(dsn, dbEnv.Host, dbEnv.Port, dbEnv.User, dbEnv.Password, dbEnv.Dbname)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect to database: %s", err)
	} else {
		log.Print("Success to connect to database.")
	}
	
	DB = db

	// db.Debug()// &models.Policy{}, &models.UserRole{},) //database migration
}