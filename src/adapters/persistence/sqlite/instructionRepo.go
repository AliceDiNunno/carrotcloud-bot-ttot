package postgres

import (
	"github.com/jinzhu/gorm"
)

type InstructionRepo struct {
	Db *gorm.DB
}

type Instruction struct {
	gorm.Model
	SenderId    int
	GroupId     int64
	Date        int64
	Instruction string
}
