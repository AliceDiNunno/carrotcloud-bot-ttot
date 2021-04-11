package telegram

import (
	"adinunno.fr/twitter-to-telegram/src/core/domain"
	"gopkg.in/tucnak/telebot.v2"
)

func (r RoutesHandler) reply(messages domain.MessageList) {
	for _, message := range messages {
		r.bot.Send(&telebot.Chat{
			ID: int64(message.Recipient),
		}, message.Text)
	}
}
