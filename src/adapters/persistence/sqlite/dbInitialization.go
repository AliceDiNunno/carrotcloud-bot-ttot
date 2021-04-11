package sqlite

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
)

func CreateDB() *gorm.DB {
	//TODO: move this to configuration
	db, err := gorm.Open("sqlite3", "./TTOT.db")

	if err != nil {
		log.Fatal("Unable to open database: " + err.Error() + "\n")
	}

	db.AutoMigrate(&Status{}, &Instruction{})

	return db
}
