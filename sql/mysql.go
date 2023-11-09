package sql

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
)

const (
	informationSchema = "/information_schema" // 数据库元数据库名
)

func NewMysql(url, dbName string) (*gorm.DB, error) {
	dsn := strings.TrimSuffix(url, "/"+dbName) + informationSchema
	return open(dsn)
}

func open(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
