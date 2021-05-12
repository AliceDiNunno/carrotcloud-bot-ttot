package sqlite

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"log"
)

func CreateDB() *gorm.DB {
	//TODO: move this to configuration
	db, err := gorm.Open(sqlite.Open("./TTOT.db"), &gorm.Config{})

	if err != nil {
		log.Fatal("Unable to open database: " + err.Error() + "\n")
	}

	db.AutoMigrate(&Status{}, &Instruction{})

	return db
}
