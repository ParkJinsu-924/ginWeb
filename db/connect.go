package db

import (
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type dbType int

const (
	MainDB dbType = iota

	Max
)

var dbMap map[dbType]*gorm.DB

func GetDB(dbType dbType) *gorm.DB {
	return dbMap[dbType]
}

type dbConfig struct {
	FileName string
	Models   []any
}

func InitializeDB() {
	dbMap = make(map[dbType]*gorm.DB)

	dbTypeAndNameMap := make(map[dbType]dbConfig)
	{ // db mapping
		dbTypeAndNameMap[MainDB] = dbConfig{
			FileName: "ncc.db",
			Models: []any{
				&User{}, // users
				&Post{}, // posts
			},
		}
	}

	for dbType, dbConfig := range dbTypeAndNameMap {
		db, err := gorm.Open(sqlite.Open(dbConfig.FileName), &gorm.Config{})
		if err != nil || db == nil {
			log.Fatal("DB Connection Failure: ", dbConfig.FileName)
		}

		err = db.AutoMigrate(dbConfig.Models...)
		if err != nil {
			log.Fatal(err)
		}

		dbMap[dbType] = db
	}
}
