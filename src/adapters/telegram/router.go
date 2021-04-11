package telegram

import (
	"adinunno.fr/twitter-to-telegram/src/core/domain"
	"gopkg.in/tucnak/telebot.v2"
)

type endpointHandler func(*telebot.Message) (domain.MessageList, error)

func (r RoutesHandler) handle(bot *telebot.Bot, endpoint string, handler endpointHandler) {
	bot.Handle(endpoint, func(message *telebot.Message) {
		messages, err := handler(message)

		if err != nil {

		} else {

		}

		messages = r.usecases.FormatMessages(messages)
		r.reply(messages)
	})
}

func (r RoutesHandler) SetCommands() {
	r.handle(r.bot, "/why", r.WhyCommand)
	r.handle(r.bot, "/retry", r.RetryCommand)
	r.handle(r.bot, "/limit", r.WhyCommand)
	r.handle(r.bot, "/stop", r.WhyCommand)
	r.handle(r.bot, telebot.OnText, r.NewMessage)
}
