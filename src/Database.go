package src

import (
	"github.com/jinzhu/gorm"
	"log"
)

var botDatabase *gorm.DB

func OpenDatabase()  {
	db, err := gorm.Open("sqlite3", "./TTOT.db")

	if err != nil {
		log.Fatal("Unable to open database: " + err.Error() + "\n")
	}

	botDatabase = db
}
