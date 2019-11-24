package models

import (
	"github.com/jinzhu/gorm"
)

type TweetInstruction struct {
	gorm.Model
	SenderId int
	GroupId int64
	Date int64
	Instruction string
}
