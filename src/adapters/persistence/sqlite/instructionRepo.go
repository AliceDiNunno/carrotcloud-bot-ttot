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

func (s InstructionRepo) GetInstruction(instruction *domain.MessageMetadata) *domain.Instruction {
	var response Instruction

	s.Db.Where("sender_id = ? AND group_id = ? AND date BETWEEN ? AND ?", instruction.Sender, instruction.Conversation, instruction.SentDate-10, instruction.SentDate+10).First(&response)
	return instructionToDomain(&response)
}

func (s InstructionRepo) HasStopInstruction(instruction *domain.MessageMetadata) bool {
	return s.GetInstruction(instruction).Instruction == domain.StopInstruction
}

func (s InstructionRepo) HasLimitInstruction(instruction *domain.MessageMetadata) int {
	limitInstruction := s.GetInstruction(instruction)
	if limitInstruction.Instruction == domain.LimitInstruction {
		return limitInstruction.Parameter
	}
	return -1
}

func instructionFromDomain(instruction *domain.Instruction) *Instruction {
	return &Instruction{
		Date:                 int64(instruction.Metadata.SentDate),
		SenderId:             int(instruction.Metadata.Sender),
		GroupId:              int64(instruction.Metadata.Conversation),
		Instruction:          string(instruction.Instruction),
		InstructionParameter: instruction.Parameter,
	}
}

func instructionToDomain(instruction *Instruction) *domain.Instruction {
	return &domain.Instruction{
		Metadata: domain.MessageMetadata{
			Id:           0,
			Conversation: domain.Chat(instruction.GroupId),
			Sender:       domain.User(instruction.SenderId),
			SentDate:     domain.Date(instruction.Date),
		},

		Instruction: domain.InstructionType(instruction.Instruction),
		Parameter:   instruction.InstructionParameter,
	}
}
