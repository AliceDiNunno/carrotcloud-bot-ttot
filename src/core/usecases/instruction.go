package usecases

import "adinunno.fr/twitter-to-telegram/src/core/domain"

func (i interactor) LimitNextThread(chat domain.Chat, sender domain.User, limit int) error {
	panic("implement me")
}

func (i interactor) StopNextThread(chat domain.Chat, sender domain.User) error {
	panic("implement me")
}
