package TToT

import (
	"github.com/jinzhu/gorm"
	"gopkg.in/tucnak/telebot.v2"
	"regexp"
	"strconv"
)


func hasTweet(db *gorm.DB, m *telebot.Message) *int64 {
	r := regexp.MustCompile(`https:\/\/twitter.com\/[a-zA-Z0-9_]*\/status\/([0-9]*)`)
	match := r.FindStringSubmatch(m.Text)
	if len(match) > 1 {
		id := match[1]

		id64, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			registerTweetStatus(db, m.ID, false, "Id malformé ou inexistant")
			return nil
		}
		return &id64
	}
	return nil
}

func addLastMessage(lastMessages []telegramLastMessage, message telegramLastMessage) []telegramLastMessage {
	for idx, lmessage := range lastMessages {
		if lmessage.chatId == message.chatId && lmessage.user.ID == message.user.ID {
			lastMessages[idx].message = message.message
			return lastMessages
		}
	}

	return append(lastMessages, message)
}

func getLastMessage(lastMessages []telegramLastMessage, chatId int64, userId *telebot.User) telegramLastMessage {
	for _, lastMessage := range lastMessages {
		if lastMessage.user.ID == userId.ID && lastMessage.chatId == chatId {
			return lastMessage
		}
	}

	return telegramLastMessage{
		user: userId,
		chatId: chatId,
		message: "",
	}
}