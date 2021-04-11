package sqlite

import (
	"github.com/jinzhu/gorm"
)

type TweetRepo struct {
	Db *gorm.DB
}

type Tweet struct {
	gorm.Model
	MessageId    int
	FetchSuccess bool
	FetchStatus  string
}
