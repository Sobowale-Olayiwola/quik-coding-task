package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"quik/domain"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DataSources struct {
	MySQLDB         *gorm.DB
	RedisInMemoryDB *redis.Client
}

// InitDS establishes connections to fields in dataSources
func initDS() (*DataSources, error) {
	log.Printf("Initializing data sources\n")
	// Initialize MySQLDB connection
	dsn := os.Getenv("DB_USERNAME") + ":" + os.Getenv("DB_PASSWORD") + "@tcp" + "(" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + ")/" + os.Getenv("DB_NAME") + "?" + "charset=utf8mb4&parseTime=True&loc=Local"
	mysql.Open(dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&domain.Player{}, &domain.Wallet{})

	//Initalize RedisDB connection
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_CONNECTION_URI"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	ctx := context.TODO()
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Failed to connect to Redis", err.Error())
	}

	return &DataSources{
		MySQLDB:         db,
		RedisInMemoryDB: rdb,
	}, nil
}
