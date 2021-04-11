package usecases

import "adinunno.fr/twitter-to-telegram/src/core/domain"

func (i interactor) RegisterThreadStatus(status *domain.Status) {
	i.statusRepo.SaveStatus(status)
}
