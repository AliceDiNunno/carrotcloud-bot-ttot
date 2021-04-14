package usecases

import "adinunno.fr/twitter-to-telegram/src/core/domain"

type StatusRepo interface {
	SaveStatus(status *domain.Status)
	GetStatus(status *domain.MessageMetadata) *domain.Status
}

type InstructionRepo interface {
	CreateInstruction(instruction *domain.Instruction)
	GetInstruction(instruction *domain.MessageMetadata) *domain.Instruction
	HasStopInstruction(instruction *domain.MessageMetadata) bool
	HasLimitInstruction(instruction *domain.MessageMetadata) int
}

type interactor struct {
	statusRepo      StatusRepo
	instructionRepo InstructionRepo
	twitterGateway  TwitterGateway
}

func NewInteractor(sR StatusRepo, iR InstructionRepo, tG TwitterGateway) interactor {
	return interactor{
		statusRepo:      sR,
		instructionRepo: iR,
		twitterGateway:  tG,
	}
}
