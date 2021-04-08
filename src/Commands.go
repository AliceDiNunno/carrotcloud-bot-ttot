package src

import (
	"adinunno.fr/twitter-to-telegram/src/adapters/persistence/postgres"
	"fmt"
	"github.com/jinzhu/gorm"
	"gopkg.in/tucnak/telebot.v2"
	"strings"
	"time"
)

var lastMessages = []telegramLastMessage{}

func whyCommand(bot *telebot.Bot, db *gorm.DB, m *telebot.Message) {
	if m.ReplyTo == nil {
		_, _ = bot.Send(m.Chat, "Merci de répondre à un message contenant un tweet")
		return
	}
	id := m.ReplyTo.ID
	var tweet postgres.TweetRegistered
	if db.Where(&postgres.TweetRegistered{MessageId: id}).First(&tweet).RecordNotFound() {
		_, _ = bot.Send(m.Chat, "Aucun enregistrement trouvé pour ce message")
		return
	}

	status := "échec"
	if tweet.FetchSuccess {
		status = "réussite"
	}
	_, _ = bot.Send(m.Chat, "Détail: "+tweet.FetchStatus+" ("+status+")")
}

func stopCommand(bot *telebot.Bot, db *gorm.DB, m *telebot.Message) {
	instruction := postgres.TweetInstruction{
		Date:        m.Unixtime,
		SenderId:    m.Sender.ID,
		GroupId:     m.Chat.ID,
		Instruction: "stop", //todo: change to enum like
	}

	db.Create(&instruction) //Todo: check if existing
	//bot.Delete(m)
}

func limitCommand(bot *telebot.Bot, db *gorm.DB, m *telebot.Message) {
	print(m.Text + "\n")
	message := strings.Split(m.Text, " ")
	if len(message) > 1 {
		/*count, err := strconv.ParseInt(message[1], 10, 32)
		if (err != nil) {
			bot.Send(m.Chat, "Impossible de récuperer la limite")
			return
		}*/
		instruction := postgres.TweetInstruction{
			Date:        m.Unixtime,
			SenderId:    m.Sender.ID,
			GroupId:     m.Chat.ID,
			Instruction: message[1], //todo: change to enum like
			//todo: check if limit is bigger than 0
		}

		db.Create(&instruction) //Todo: check if existing
	} else {
		_, _ = bot.Send(m.Chat, "Merci de préciser jusqu'à combien de tweet vous voulez limiter\nPar exemple: `/limit 2`")
	}
	//bot.Delete(m)
}

func retryCommand(bot *telebot.Bot, db *gorm.DB, m *telebot.Message) {
	if m.ReplyTo == nil {
		_, _ = bot.Send(m.Chat, "Merci de répondre à un message contenant un tweet")
		return
	}
	id := hasTweet(db, m.ReplyTo)
	if id == nil {
		_, _ = bot.Send(m.Chat, "Aucun tweet n'a été trouvé dans ce message")
		return
	}
	tweets := fetchTweets(db, m, twitterClient, *id)

	if len(tweets) < 2 {
		registerTweetStatus(db, m.ID, false, "Ce n'est pas un thread ou aucune réponse n'a pu etre trouvée")
		return
	}

	prepareToSend(db, bot, m.ReplyTo, tweets, "")
}

func onText(bot *telebot.Bot, db *gorm.DB, m *telebot.Message) {
	id := hasTweet(db, m)
	if id != nil {
		tweets := fetchTweets(db, m, twitterClient, *id)

		fmt.Printf("Len: %d\n", len(tweets))

		if len(tweets) < 2 {
			registerTweetStatus(db, m.ID, false, "Ce n'est pas un thread ou aucune réponse n'a pu etre trouvée")
			return
		}

		lastMsg := getLastMessage(lastMessages, m.Chat.ID, m.Sender)

		if !(lastMsg.time-10 <= m.Unixtime && lastMsg.time+10 >= m.Unixtime) {
			lastMsg.message = ""
		}

		url := getTweetUrl(m)
		messageContent := strings.TrimSpace(strings.Replace(m.Text, url, "", -1))
		lastMessageContent := lastMsg.message
		userMessage := lastMessageContent

		if userMessage != "" {
			userMessage = userMessage + "\n"
		}
		userMessage = userMessage + url + "\n"
		if messageContent != "" {
			userMessage = userMessage + messageContent + "\n"
		}

		prepareToSend(db, bot, m, tweets, userMessage)

		if lastMsg.message != "" {
			lastMessages = addLastMessage(lastMessages, telegramLastMessage{
				time:    time.Now().Unix(),
				chatId:  m.Chat.ID,
				user:    m.Sender,
				message: "",
			})
		}
	} else {
		lastMessages = addLastMessage(lastMessages, telegramLastMessage{
			time:    time.Now().Unix(),
			chatId:  m.Chat.ID,
			user:    m.Sender,
			message: m.Text,
		})
	}
}
