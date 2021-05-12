package telegram

import (
	"adinunno.fr/twitter-to-telegram/src/core/domain"
	"adinunno.fr/twitter-to-telegram/src/core/usecases"
	"gopkg.in/tucnak/telebot.v2"
)

type endpointHandler func(*telebot.Message) (domain.MessageList, error)

func (r RoutesHandler) handle(bot *telebot.Bot, endpoint string, handler endpointHandler) {
	bot.Handle(endpoint, func(message *telebot.Message) {
		messages, err := handler(message)

		var details = ""

		if err != nil {
			details = err.Error()
		}

		r.usecases.RegisterThreadStatus(&domain.Status{
			MetaData: domain.MessageMetadata{
				Id:           domain.ID(message.ID),
				Conversation: domain.Chat(message.Chat.ID),
				Sender:       domain.User(message.Sender.ID),
			},
			DidSucceed:         err != nil,
			AdditionnalDetails: details,
		})

		messages = r.usecases.FormatMessages(messages)
		r.reply(messages)
	})
}

func (r RoutesHandler) SetCommands() {
	r.handle(r.bot, "/why", r.WhyCommand)
	r.handle(r.bot, "/retry", r.RetryCommand)
	r.handle(r.bot, "/limit", r.LimitCommand)
	r.handle(r.bot, "/stop", r.StopCommand)
	r.handle(r.bot, telebot.OnText, r.NewMessage)
}

func NewCommandHandler(bot *telebot.Bot, usecasesHandler usecases.Usecases) RoutesHandler {
	return RoutesHandler{
		bot:      bot,
		usecases: usecasesHandler,
	}
}
