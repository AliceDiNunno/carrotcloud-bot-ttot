package sqlite

import (
	"adinunno.fr/twitter-to-telegram/src/core/domain"
	"github.com/jinzhu/gorm"
)

type InstructionRepo struct {
	Db *gorm.DB
}

type Instruction struct {
	gorm.Model
	Date                 int64
	SenderId             int
	GroupId              int64
	Instruction          string
	InstructionParameter int //this is used for "/limit" not "/stop", subject to changed if new commands where to arrive
}

func (s InstructionRepo) CreateInstruction(instruction *domain.Instruction) {
	request := instructionFromDomain(instruction)

	s.Db.Create(&request)
}

func (s InstructionRepo) GetInstruction(instruction *domain.Instruction) *domain.Instruction {
	request := instructionFromDomain(instruction)
	var response Instruction

	s.Db.Where("sender_id = ? AND group_id = ? AND date BETWEEN ? AND ?", request.SenderId, request.GroupId, instruction.Date-10, instruction.Date+10).First(&response)
	return instructionToDomain(&response)
}

func (s InstructionRepo) HasStopInstruction(instruction *domain.Instruction) bool {
	return s.GetInstruction(instruction).Instruction == domain.StopInstruction
}

func (s InstructionRepo) HasLimitInstruction(instruction *domain.Instruction) int {
	limitInstruction := s.GetInstruction(instruction)
	if limitInstruction.Instruction == domain.LimitInstruction {
		return limitInstruction.Parameter
	}
	return -1
}

func instructionFromDomain(instruction *domain.Instruction) *Instruction {
	return &Instruction{
		Date:                 instruction.Date,
		SenderId:             int(instruction.User),
		GroupId:              int64(instruction.Recipient),
		Instruction:          string(instruction.Instruction),
		InstructionParameter: instruction.Parameter,
	}
}

func instructionToDomain(instruction *Instruction) *domain.Instruction {
	return &domain.Instruction{
		Date:        instruction.Date,
		Recipient:   domain.Chat(instruction.GroupId),
		User:        domain.User(instruction.SenderId),
		Instruction: domain.InstructionType(instruction.Instruction),
		Parameter:   instruction.InstructionParameter,
	}
}
