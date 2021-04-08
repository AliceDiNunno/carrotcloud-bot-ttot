package src

import (
	"fmt"
	"github.com/AliceDiNunno/TwitterToTelegram/models"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/jinzhu/gorm"
	"gopkg.in/tucnak/telebot.v2"
	"strconv"
)

func writeHeader(sender *twitter.User, originalMessage *telebot.Message, lastMessage string) string {
	header := "@"+originalMessage.Sender.Username + ":\n"

	if lastMessage != "" {
		header = header + lastMessage + "\nPar "
	}

	header = header + sender.Name + " (@" + sender.ScreenName + "): \n"

	return header
}

func send(bot *telebot.Bot, tweets []twitter.Tweet, initialMessage *telebot.Message, lastMessage string) {
	message := ""
	headerWritten := false
	previewEnabled := false
	messageToReplyTo := initialMessage

	sendMessage := func() {
		if headerWritten == false {
			message = writeHeader(tweets[0].User, initialMessage, lastMessage) + "" + message
			headerWritten = true
		}
		sentMsg, err := bot.Send(initialMessage.Chat, message, &telebot.SendOptions {
			DisableWebPagePreview: !previewEnabled,
			ReplyTo: messageToReplyTo,
		})
		if err != nil {
			//Handle Error
		} else {
			messageToReplyTo = sentMsg
		}

		message = ""
		previewEnabled = false
	}

	for _, tweet := range tweets {
		print(tweet.FullText)
		//print("NEW TWEET \n")
		if tweet.ExtendedEntities != nil && len(tweet.ExtendedEntities.Media) > 0 {
			if len(message) > 0 {
				sendMessage()
			}
			previewEnabled = true
			message = tweet.FullText
			sendMessage()
		} else {
			if len(message + "\n\n" + tweet.FullText) > 4096 {
				sendMessage()
			}
			message = message + "\n\n" + tweet.FullText
			//print("APPEND MESSAGE \n" + message )
		}
	}

	if message != "" {
		sendMessage()
	}
}

func prepareToSend(db *gorm.DB, bot *telebot.Bot, m *telebot.Message, tweets []twitter.Tweet, lastMessage string) {
	var instruction models.TweetInstruction
	var limit int64
	limit = -1
	db.Where(&models.TweetInstruction{GroupId: m.Chat.ID, SenderId: m.Sender.ID}).Where("Date BETWEEN ? AND ?", m.Unixtime - 10, m.Unixtime + 10).First(&instruction)
	if instruction.Instruction != "" {
		stop := false
		if instruction.Instruction == "stop" {
			registerTweetStatus(db, m.ID, false, "La récupération de ce tweet à été annulée par l'utilisateur")
			stop = true
		} else {
			count, err := strconv.ParseInt(instruction.Instruction, 10, 64)
			if err != nil {
				_, _ = bot.Send(m.Chat, "Impossible de récuperer la limite")
			} else {
				limit = count
			}
		}
		db.Delete(&instruction)
		if stop {
			return
		}
	}

	if limit > 0 {
		limit = int64(len(tweets)) - limit

		for ; limit > 0; limit-- {
			tweets = tweets[:len(tweets)-1]
		}
	}

	for _, tweet := range tweets {
		fmt.Printf(tweet.FullText + "\n")
	}

	send(bot, tweets, m, lastMessage)
	registerTweetStatus(db, m.ID, true, "Ok")
}

