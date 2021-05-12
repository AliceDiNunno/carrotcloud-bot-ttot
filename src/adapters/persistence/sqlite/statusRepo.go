package sqlite

import (
	"adinunno.fr/twitter-to-telegram/src/core/domain"
	"github.com/davecgh/go-spew/spew"
	"gorm.io/gorm"
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

func (s StatusRepo) GetStatus(status *domain.MessageMetadata) *domain.Status {
	var response Status

	s.Db.Where("chat_id = ? AND message_id = ?", status.Conversation, status.Id).First(&response)
	return statusToDomain(&response)
}

func statusFromDomain(status *domain.Status) Status {
	return Status{
		ChatId:       int64(status.MetaData.Conversation),
		MessageId:    int(status.MetaData.Id),
		FetchSuccess: status.DidSucceed,
		Details:      status.AdditionnalDetails,
	}
}

func statusToDomain(status *Status) *domain.Status {
	return &domain.Status{
		MetaData: domain.MessageMetadata{
			Id:           domain.ID(status.MessageId),
			Conversation: domain.Chat(status.ChatId),
			Sender:       0,
			SentDate:     0,
		},

		DidSucceed:         status.FetchSuccess,
		AdditionnalDetails: status.Details,
	}
}
