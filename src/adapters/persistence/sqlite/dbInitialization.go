package postgres

import (
	"github.com/jinzhu/gorm"
	"log"
)

func CreateDB() gorm.DB {
	db, err := gorm.Open("sqlite3", "./TTOT.db")

	if err != nil {
		log.Fatal("Unable to open database: " + err.Error() + "\n")
	}

	return db
}
