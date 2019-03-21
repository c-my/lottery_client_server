package datasource

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	databasePath = "./data/data.db"
	// DB is the only database
	DB = newDatabase()
)

// newDatabase returns a database
func newDatabase() *gorm.DB {
	if !isDatabaseExists(databasePath) {
		os.Create(databasePath)
	}
	DB, err := gorm.Open("sqlite3", databasePath)
	if err != nil {
		panic(fmt.Sprint("get error when connectiong to database: ", err))
	}
	return DB
}

func isDatabaseExists(path string) bool {
	s, err := os.Stat(path)
	// file not exists
	if (err != nil && os.IsNotExist(err)) || s.IsDir() {
		return false
	}
	return true
}
