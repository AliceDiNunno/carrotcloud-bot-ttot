package TToT

import (
	"github.com/AliceDiNunno/TwitterToTelegram/models"
	"github.com/jinzhu/gorm"
	"gopkg.in/tucnak/telebot.v2"
	"os"
	"strings"
	"time"
)

func BotMiddleware(upd *telebot.Update) bool {
	if upd.Message == nil {
		return false
	} else if strings.HasSuffix(strings.ToLower(upd.Message.Sender.Username), "bot") {
		return false
	} else if time.Now().Unix() - upd.Message.Unixtime > 10 {
		return false
	}
	return true
}

func Init(bot *telebot.Bot, db *gorm.DB) {
	if (!db.HasTable(&models.TweetInstruction{})) {
		db.CreateTable(&models.TweetInstruction{})
	}
	if (!db.HasTable(&models.TweetRegistered{})) {
		db.CreateTable(&models.TweetRegistered{})
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



func SetupBot() {
	b, err := telebot.NewBot(telebot.Settings{
		Token:  os.Getenv("telegram_bot_key"),
		URL: "https://api.telegram.org",
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
