package usecases

import "adinunno.fr/twitter-to-telegram/src/core/domain"

type TwitterGateway interface {
	GetTweet(id domain.TweetId) (*domain.Tweet, error)
	GetFullThread(tweet *domain.Tweet) (domain.TweetList, error)
}

type Usecases interface {
	FindTweetStatus(metadata domain.MessageMetadata) (*domain.Message, error)
	NewMessageReceived(date int64, chat domain.Chat, sender domain.User, id domain.TweetId) (domain.MessageList, error)
	LimitNextThread(date int64, chat domain.Chat, sender domain.User, limit int) error
	StopNextThread(date int64, chat domain.Chat, sender domain.User) error

	RegisterThreadStatus(status *domain.Status)

	FormatMessages(list domain.MessageList) domain.MessageList
}
