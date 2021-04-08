package src

import (
	"adinunno.fr/twitter-to-telegram/src/adapters/persistence/postgres"
	"adinunno.fr/twitter-to-telegram/src/infra"
	"github.com/jinzhu/gorm"
	"gopkg.in/tucnak/telebot.v2"
	"strings"
	"time"
)

func BotMiddleware(upd *telebot.Update) bool {
	if upd.Message == nil {
		return false
	} else if strings.HasSuffix(strings.ToLower(upd.Message.Sender.Username), "bot") {
		return false
	} else if time.Now().Unix()-upd.Message.Unixtime > 10 {
		return false
	}
	return true
}

func Init(bot *telebot.Bot, db *gorm.DB) {
	if (!db.HasTable(&postgres.TweetInstruction{})) {
		db.CreateTable(&postgres.TweetInstruction{})
	}
	if (!db.HasTable(&postgres.TweetRegistered{})) {
		db.CreateTable(&postgres.TweetRegistered{})
	}

	bot.Handle("/why", func(m *telebot.Message) {
		whyCommand(bot, db, m)
	})

	bot.Handle("/retry", func(m *telebot.Message) {
		retryCommand(bot, db, m)
	})

	bot.Handle("/limit", func(m *telebot.Message) {
		limitCommand(bot, db, m)
	})

	bot.Handle("/stop", func(m *telebot.Message) {
		stopCommand(bot, db, m)
	})

	bot.Handle(telebot.OnText, func(m *telebot.Message) {
		onText(bot, db, m)
	})
}

func SetupBot(conf infra.TelegramConf) {
	b, err := telebot.NewBot(telebot.Settings{
		Token: conf.BotToken,
		URL:   "https://api.telegram.org",
	})

	if err != nil {
		print("Telegram error")
		return
	}

	var poller telebot.Poller = &telebot.LongPoller{Timeout: 10 * time.Second}
	poller = telebot.NewMiddlewarePoller(poller, func(upd *telebot.Update) bool {
		return BotMiddleware(upd)
	})

	b.Poller = poller

	Init(b, botDatabase)

	b.Start()
}
