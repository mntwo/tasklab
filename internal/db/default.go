package db

import "gorm.io/gorm"

var dbs map[string]*gorm.DB

func GetDB(dbName string) *gorm.DB {
	return dbs[dbName]
}
