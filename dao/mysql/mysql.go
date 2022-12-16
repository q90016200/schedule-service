package mysql

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DBConn DBConn
var DBConn *gorm.DB

type Config struct {
	UserName string
	PassWord string
	Host     string
	Port     string
	DataBase string
}

func (config *Config) Conn() *gorm.DB {
	dsn := config.UserName + ":" + config.PassWord + "@tcp(" + config.Host + ":" + config.Port + ")/" + config.DataBase + "?charset=utf8mb4&parseTime=true&loc=Local"

	//fmt.Println(dsn)

	//db, err = gorm.Open(mysql.New(mysql.Config{
	//	DSN: dsn,
	//}))

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err) // 連不到就直接panic裡服務重起再連
	}

	return db
}
