package usecases

import "adinunno.fr/twitter-to-telegram/src/core/domain"

func (i interactor) LimitNextThread(date int64, chat domain.Chat, sender domain.User, limit int) error {
	i.instructionRepo.CreateInstruction(&domain.Instruction{
		Metadata: domain.MessageMetadata{
			Id:           0,
			Conversation: chat,
			Sender:       sender,
			SentDate:     domain.Date(date),
		},
		Instruction: domain.LimitInstruction,
		Parameter:   limit,
	})

	return nil
}

func (i interactor) StopNextThread(date int64, chat domain.Chat, sender domain.User) error {
	i.instructionRepo.CreateInstruction(&domain.Instruction{
		Metadata: domain.MessageMetadata{
			Id:           0,
			Conversation: chat,
			Sender:       sender,
			SentDate:     domain.Date(date),
		},
		Instruction: domain.StopInstruction,
		Parameter:   -1, //Unused for /stop
	})

	return nil
}
