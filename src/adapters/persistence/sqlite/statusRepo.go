package sqlite

import (
	"adinunno.fr/twitter-to-telegram/src/core/domain"
	"github.com/jinzhu/gorm"
)

type StatusRepo struct {
	Db *gorm.DB
}

type Status struct {
	gorm.Model
	ChatId       int64
	MessageId    int
	FetchSuccess bool
	Details      string
}

func (s StatusRepo) SaveStatus(status *domain.Status) {
	request := statusFromDomain(status)

	s.Db.Create(&request)
}

func (s StatusRepo) GetStatus(status *domain.Status) *domain.Status {
	request := statusFromDomain(status)
	var response Status

	s.Db.Where("chat_id = ? AND message_id = ?", request.ChatId, request.MessageId).First(&response)
	return statusToDomain(&response)
}

func statusFromDomain(status *domain.Status) Status {
	return Status{
		ChatId:       int64(status.Recipient),
		MessageId:    int(status.Sender),
		FetchSuccess: status.DidSucceed,
		Details:      status.AdditionnalDetails,
	}
}

func statusToDomain(status *Status) *domain.Status {
	return &domain.Status{
		Recipient:          domain.Chat(status.ChatId),
		Sender:             domain.User(status.MessageId),
		DidSucceed:         status.FetchSuccess,
		AdditionnalDetails: status.Details,
	}
}
