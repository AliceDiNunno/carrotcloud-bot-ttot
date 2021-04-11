package telegram

import (
	"adinunno.fr/twitter-to-telegram/src/core/domain"
	"adinunno.fr/twitter-to-telegram/src/core/usecases"
	"errors"
	"gopkg.in/tucnak/telebot.v2"
)

type RoutesHandler struct {
	bot      *telebot.Bot
	usecases usecases.Usecases
}

func (r RoutesHandler) WhyCommand(m *telebot.Message) (domain.MessageList, error) {
	/*	if m.ReplyTo == nil {
			_, err := r.bot.Send(m.Chat, "Merci de répondre à un message contenant un tweet")

			if err != nil {
				println("Unable to send meesage:", err.Error())
			}

			return
		}
		id := m.ReplyTo.ID
		var tweet sqlite.Tweet
		if db.Where(&sqlite.Tweet{MessageId: id}).First(&tweet).RecordNotFound() {
			_, err := r.bot.Send(m.Chat, "Aucun enregistrement trouvé pour ce message")

			if err != nil {
				println("Unable to send meesage:", err.Error())
			}

			return
		}

		status := "échec"
		if tweet.FetchSuccess {
			status = "réussite"
		}
		_, err := r.bot.Send(m.Chat, "Détail: "+tweet.FetchStatus+" ("+status+")")

		if err != nil {
			println("Unable to send meesage:", err.Error())
		}
	*/
	r.usecases.WhyTweetNotWorking()
	return domain.MessageList{
		{
			RecipientId: m.Chat.ID,
			Text:        "This is not implemented",
		},
	}, errors.New("waiting for implementation")
}

func (r RoutesHandler) RetryCommand(m *telebot.Message) (domain.MessageList, error) {
	if m.ReplyTo == nil {
		return domain.MessageList{
			{
				RecipientId: m.Chat.ID,
				Text:        "Please use this command while replying",
			},
		}, errors.New("command was used without replying to a message")
	}

	return r.NewMessage(m.ReplyTo)
}

func (r RoutesHandler) NewMessage(m *telebot.Message) (domain.MessageList, error) {
	id, err := findTweet(m.Text)

	if err != nil {
		return nil, err
	}

	tweets, err := r.usecases.NewMessageReceived(m.Text, id)

	if err != nil {
		return nil, err
	}

	var messages domain.MessageList
	for _, tweet := range tweets {
		messages = append(messages, &domain.Message{
			RecipientId: m.Chat.ID,
			Text:        tweet.Message,
		})
	}
	return messages, nil
}

func NewCommandHandler(bot *telebot.Bot, usecasesHandler usecases.Usecases) RoutesHandler {
	return RoutesHandler{
		bot:      bot,
		usecases: usecasesHandler,
	}
}
