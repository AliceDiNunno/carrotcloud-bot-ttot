package telegram

import (
	"adinunno.fr/twitter-to-telegram/src/config"
	"gopkg.in/tucnak/telebot.v2"
	"log"
	"strings"
	"time"
)

func shouldUpdate(newMessage *telebot.Update) bool {
	if newMessage.Message == nil { //if text message is empty we abort
		return false
	} else if strings.HasSuffix(strings.ToLower(newMessage.Message.Sender.Username), "bot") { //Or if the message is sent from a bot
		return false
	} else if time.Now().Unix()-newMessage.Message.Unixtime > 10 { //Or if the message has been sent 10mS ago
		//TODO: ^ move this to configuration ^
		return false
	}
	return true
}

func NewTelegramBot(conf config.TelegramConf) *telebot.Bot {
	bot, err := telebot.NewBot(telebot.Settings{
		Token: conf.BotToken,
		URL:   "https://api.telegram.org",
	})

	if err != nil {
		log.Fatalln("Unable to start telegram bot: ", err.Error())
	}

	var poller telebot.Poller = &telebot.LongPoller{Timeout: 10 * time.Second} //TODO: move this to configuration
	//TODO: move timeout to configuration

	poller = telebot.NewMiddlewarePoller(poller, func(upd *telebot.Update) bool {
		return shouldUpdate(upd)
	})

	bot.Poller = poller

	return bot
}
