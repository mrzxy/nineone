package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
	"nineone/config"
)
import "github.com/go-redis/redis/v7"
var (
	db *gorm.DB
)

func InitDB() {
	var err error
	db, err = gorm.Open("mysql", config.Conf.Server.MySqlUrl)
	if err != nil {
		logrus.Fatal("db connect error ", err)
	}
	if config.Conf.Server.Debug {
		db.LogMode(true)
	}

	InitRedis()
}
var redisClient *redis.Client
func InitRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     config.Conf.Redis.Addr,
		Password: config.Conf.Redis.Password,
		DB:       config.Conf.Redis.DB,
	})

	_, err := redisClient.Ping().Result()
	if err != nil {
		logrus.Fatal("redis connect error ", err)
	}
}
func Redis() *redis.Client {
	return redisClient
}
func CloseRedis() {
	redisClient.Close()
}
func DB() *gorm.DB {
	return db
}

func CloseDB() error {
	return db.Close()
}
