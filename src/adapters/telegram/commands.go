package telegram

import (
	"adinunno.fr/twitter-to-telegram/src/core/domain"
	"adinunno.fr/twitter-to-telegram/src/core/usecases"
	"errors"
	"gopkg.in/tucnak/telebot.v2"
	"strconv"
	"strings"
)

type RoutesHandler struct {
	bot      *telebot.Bot
	usecases usecases.Usecases
}

func (r RoutesHandler) WhyCommand(m *telebot.Message) (domain.MessageList, error) {
	if m.ReplyTo == nil {
		return domain.MessageList{
			{
				Recipient: domain.Chat(m.Chat.ID),
				Text:      "Please use this command while replying",
			},
		}, errors.New("command was used without replying to a message")
	}

	reply, err := r.usecases.FindTweetStatus(domain.Status{
		Recipient: domain.Chat(m.Chat.ID),
		Sender:    domain.User(m.ReplyTo.ID), //TODO: this should not be sender but messageId or smthg
	})

	if err != nil {
		return nil, err
	}

	return domain.MessageList{
		reply,
	}, nil
}

func (r RoutesHandler) LimitCommand(m *telebot.Message) (domain.MessageList, error) {
	arguments := strings.Split(m.Text, " ")
	if len(arguments) > 1 {
		result, err := strconv.Atoi(arguments[1])

		if err != nil {
			return domain.MessageList{
				{
					Recipient: domain.Chat(m.Chat.ID),
					Text:      "Please use /limit with a number. For example /limit 2",
				},
			}, errors.New("/limit used with an invalid argument")
		}

		err = r.usecases.LimitNextThread(m.Unixtime, domain.Chat(m.Chat.ID), domain.User(m.Sender.ID), result)

		return nil, err //Command is silently processed if no err
	} else {
		return domain.MessageList{
			{
				Recipient: domain.Chat(m.Chat.ID),
				Text:      "Please use /limit with a number. For example /limit 2",
			},
		}, errors.New("/limit used without arguments")
	}
}

func (r RoutesHandler) StopCommand(m *telebot.Message) (domain.MessageList, error) {
	err := r.usecases.StopNextThread(m.Unixtime, domain.Chat(m.Chat.ID), domain.User(m.Sender.ID))
	return nil, err //Command is silently processed if no err
}

func (r RoutesHandler) RetryCommand(m *telebot.Message) (domain.MessageList, error) {
	if m.ReplyTo == nil {
		return domain.MessageList{
			{
				Recipient: domain.Chat(m.Chat.ID),
				Text:      "Please use this command while replying",
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

	messages, err := r.usecases.NewMessageReceived(m.Unixtime, domain.Chat(m.Chat.ID), domain.User(m.Sender.ID), id)

	if err != nil {
		return nil, err
	}

	return messages, nil
}
