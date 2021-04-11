package usecases

import "adinunno.fr/twitter-to-telegram/src/core/domain"

type StatusRepo interface {
	SaveStatus(status *domain.Status)
	GetStatus(status *domain.Status) *domain.Status //TODO: replace parameter
}

type InstructionRepo interface {
	CreateInstruction(instruction *domain.Instruction)
	GetInstruction(instruction *domain.Instruction) *domain.Instruction //TODO: replace parameter with struct that only takes chat and user
	HasStopInstruction(instruction *domain.Instruction) bool
	HasLimitInstruction(instruction *domain.Instruction) int
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
