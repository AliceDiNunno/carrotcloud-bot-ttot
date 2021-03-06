package usecases

import (
	"adinunno.fr/twitter-to-telegram/src/core/domain"
	"errors"
)

func (i interactor) NewMessageReceived(date int64, chat domain.Chat, sender domain.User, id domain.TweetId) (domain.MessageList, error) {
	if (i.instructionRepo.HasStopInstruction(&domain.MessageMetadata{
		SentDate:     domain.Date(date),
		Conversation: chat,
		Sender:       sender,
	})) {
		return nil, errors.New("user has aborted the process using /stop")
	}
	tweet, err := i.twitterGateway.GetTweet(id)

	if err != nil {
		return nil, err
	}

	if !tweet.IsWithinTimeLimit() {
		return nil, errors.New("tweet is too old (above 15 days limit)")
	}

	tweets := domain.TweetList{}

	replies, err := i.twitterGateway.GetFullThread(tweet)

	tweets = append(tweets, replies...)

	if err != nil {
		return nil, err
	}

	if len(tweets) < 2 {
		return nil, errors.New("no thread has been found")
	}

	limit := i.instructionRepo.HasLimitInstruction(&domain.MessageMetadata{
		SentDate:     domain.Date(date),
		Conversation: chat,
		Sender:       sender,
	})

	var messages domain.MessageList
	for i, tweet := range tweets {
		if limit > 0 && i > limit-1 {
			break
		}
		messages = append(messages, &domain.Message{
			Metadata: domain.MessageMetadata{
				Conversation: chat,
			},

			Text: tweet.Message,
		})
	}

	return messages, nil
}

func (i interactor) FindTweetStatus(status domain.MessageMetadata) (*domain.Message, error) {
	dbstatus := i.statusRepo.GetStatus(&status)

	//TODO: this is the values returned when no entry was found. UPDATE to gorm2 required and fetch errors
	if dbstatus.MetaData.Conversation == 0 && dbstatus.MetaData.Sender == 0 {
		return &domain.Message{
			Metadata: domain.MessageMetadata{
				Conversation: status.Conversation,
			},
			Text: "this tweet is not in my database #sadface",
		}, nil
	}

	if dbstatus != nil {
		return &domain.Message{
			Metadata: domain.MessageMetadata{
				Conversation: status.Conversation,
			},
			Text: dbstatus.AdditionnalDetails,
		}, nil
	}

	return nil, nil
}
