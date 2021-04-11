package usecases

import "adinunno.fr/twitter-to-telegram/src/core/domain"

func (i interactor) LimitNextThread(date int64, chat domain.Chat, sender domain.User, limit int) error {
	i.instructionRepo.CreateInstruction(&domain.Instruction{
		Date:        date,
		Recipient:   chat,
		User:        sender,
		Instruction: domain.LimitInstruction,
		Parameter:   limit,
	})

	return nil
}

func (i interactor) StopNextThread(date int64, chat domain.Chat, sender domain.User) error {
	i.instructionRepo.CreateInstruction(&domain.Instruction{
		Date:        date,
		Recipient:   chat,
		User:        sender,
		Instruction: domain.StopInstruction,
		Parameter:   -1, //Unused
	})

	return nil
}
