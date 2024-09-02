package db

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMySQLStorage(cfg mysql.Config) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.New(cfg), &gorm.Config{})
	if err != nil {
		log.Println(err)
	}

	return db, nil
}
