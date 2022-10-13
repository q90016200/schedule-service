package mysql

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	UserName string
	PassWord string
	Host     string
	Port     string
	DataBase string
}

func (config *Config) Conn() (db *gorm.DB, err error) {
	dsn := config.UserName + ":" + config.PassWord + "@tcp(" + config.Host + ":" + config.Port + ")/" + config.DataBase + "?charset=utf8mb4&parseTime=true&loc=Local"

	fmt.Println(dsn)

	//db, err = gorm.Open(mysql.New(mysql.Config{
	//	DSN: dsn,
	//}))

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	return
}
