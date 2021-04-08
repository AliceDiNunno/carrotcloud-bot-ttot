package postgres

import (
	"github.com/jinzhu/gorm"
)

type TweetRegistered struct {
	gorm.Model
	MessageId    int
	FetchSuccess bool
	FetchStatus  string
}
