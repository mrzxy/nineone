package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
	"nineone/config"
)

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
}

func DB() *gorm.DB {
	return db
}

func CloseDB() error {
	return db.Close()
}
