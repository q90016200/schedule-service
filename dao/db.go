package dao

import (
	"gorm.io/gorm"
	"os"
	"scheduleService/dao/mysql"
)

func Load() *gorm.DB {
	mysqlConfig := mysql.Config{
		UserName: os.Getenv("MYSQL_USER"),
		PassWord: os.Getenv("MYSQL_PASSWORD"),
		Host:     os.Getenv("MYSQL_HOST"),
		Port:     os.Getenv("MYSQL_PORT"),
		DataBase: os.Getenv("MYSQL_DATABASE"),
	}

	dbConn := mysqlConfig.Conn()
	return dbConn
}
