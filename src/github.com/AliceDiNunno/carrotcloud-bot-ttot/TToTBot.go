package TToT

import (
	"fmt"
	"github.com/dghubble/oauth1"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/dghubble/go-twitter/twitter"
	"gopkg.in/tucnak/telebot.v2"
	_ "github.com/mattn/go-sqlite3"
)

var engineDatabase *gorm.DB

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func getTweetUrl(m *telebot.Message) string {
	r := regexp.MustCompile(`https:\/\/twitter.com\/([a-zA-Z0-9_]*)\/status\/([0-9]*)`)
	match := r.FindStringSubmatch(m.Text)
	if len(match) > 2 {
		return "https://twitter.com/"+match[1]+"/status/"+match[2]
	}
	return ""
}

func fetchTweets(db *gorm.DB, m *telebot.Message, client *twitter.Client, id int64) []twitter.Tweet {
	tweetList := []twitter.Tweet{}

	statusLookupParams := &twitter.StatusShowParams{
		TweetMode: "extended",
	}
	tweet, _, err := client.Statuses.Show(id, statusLookupParams)

	if (err != nil) {
		print(err.Error())
		registerTweetStatus(db, m.ID, false, "Erreur twitter: " + err.Error())
		return []twitter.Tweet{}
	}

	tweettime, err := time.Parse("Mon Jan 02 15:04:05 -0700 2006", tweet.CreatedAt)
	if (err != nil) {
		print(err.Error() + "\n")
	} else {
		fifteenDays := 15 * 24 * 60 * 60

		fmt.Printf("TWEET TIME: %d\n", tweettime.Unix())

		if (time.Now().Unix() - tweettime.Unix()) > int64(fifteenDays) {
			registerTweetStatus(db, m.ID, false, "Le tweet est daté de plus de 15 jours")
			return []twitter.Tweet{}
		}
	}

	tweetList = append(tweetList, *tweet)

	searchTweetParams := &twitter.SearchTweetParams {
		Query:     "from:"+tweet.User.ScreenName + " to:"+tweet.User.ScreenName,
		SinceID:	tweet.ID,
		ResultType: "recent",
		TweetMode: "extended",
		Count:     1000,
	}
	tweets, _, _ := client.Search.Tweets(searchTweetParams)

	statuses := tweets.Statuses
	sort.Slice(statuses, func(i, j int) bool {
		return statuses[i].ID < statuses[j].ID
	})

	var knownids []string
	knownids = append(knownids, tweet.IDStr)

	for _, twt := range statuses {
		if (contains(knownids, twt.InReplyToStatusIDStr)) {
			knownids = append(knownids, twt.IDStr)
			tweetList = append(tweetList, twt)
		}
	}

	return tweetList
}

type TweetInstruction struct {
	gorm.Model
	SenderId int
	GroupId int64
	Date int64
	Instruction string
}

type TweetRegistered struct {
	gorm.Model
	MessageId int
	FetchSucess bool
	FetchStatus string
}

type telegramLastMessage struct {
	time int64
	chatId int64
	user *telebot.User
	message string
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
		if (lastMessage.user.ID == userId.ID && lastMessage.chatId == chatId) {
			return lastMessage
		}
	}

	return telegramLastMessage{
		user: userId,
		chatId: chatId,
		message: "",
	}
}

func writeHeader(sender *twitter.User, originalMessage *telebot.Message, lastMessage string) string {
	header := "@"+originalMessage.Sender.Username + ":\n"

	if (lastMessage != "") {
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
		if (headerWritten == false) {
			message = writeHeader(tweets[0].User, initialMessage, lastMessage) + "" + message
			headerWritten = true
		}
		sentMsg, err := bot.Send(initialMessage.Chat, message, &telebot.SendOptions {
			DisableWebPagePreview: !previewEnabled,
			ReplyTo: messageToReplyTo,
		})
		if (err != nil) {
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
			if (len(message) > 0) {
				sendMessage()
			}
			previewEnabled = true
			message = tweet.FullText
			sendMessage()
		} else {
			if (len(message + "\n\n" + tweet.FullText) > 4096) {
				sendMessage()
			}
			message = message + "\n\n" + tweet.FullText
			//print("APPEND MESSAGE \n" + message )
		}
	}

	if (message != "") {
		sendMessage()
	}
}

func prepareToSend(db *gorm.DB, bot *telebot.Bot, m *telebot.Message, tweets []twitter.Tweet, lastMessage string) {
	var instruction TweetInstruction
	var limit int64
	limit = -1
	db.Where(&TweetInstruction{GroupId: m.Chat.ID, SenderId: m.Sender.ID}).Where("Date BETWEEN ? AND ?", m.Unixtime - 10, m.Unixtime + 10).First(&instruction)
	if (instruction.Instruction != "") {
		stop := false
		if (instruction.Instruction == "stop") {
			registerTweetStatus(db, m.ID, false, "La récupération de ce tweet à été annulée par l'utilisateur")
			stop = true
		} else {
			count, err := strconv.ParseInt(instruction.Instruction, 10, 64)
			if (err != nil) {
				bot.Send(m.Chat, "Impossible de récuperer la limite")
			} else {
				limit = count
			}
		}
		db.Delete(&instruction)
		if (stop) {
			return
		}
	}

	if (limit > 0) {
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

func registerTweetStatus(db *gorm.DB, id int, success bool, reason string) {
	print("Tweet status: " + reason)

	db.Save(&TweetRegistered{
		MessageId: id,
		FetchSucess: success,
		FetchStatus: reason,
	})
}

func hasTweet(db *gorm.DB, m *telebot.Message) *int64 {
	r := regexp.MustCompile(`https:\/\/twitter.com\/[a-zA-Z0-9_]*\/status\/([0-9]*)`)
	match := r.FindStringSubmatch(m.Text)
	if len(match) > 1 {
		id := match[1]

		id64, err := strconv.ParseInt(id, 10, 64)
		if (err != nil) {
			registerTweetStatus(db, m.ID, false, "Id malformé ou inexistant")
			return nil
		}
		return &id64
	}
	return nil
}

type bot struct{}

var configuration = map[string]string{}

func ConsumeConfigurationItem(key, value string) {
	configuration[key] = value
}

func BotMiddleware(upd *telebot.Update) bool {
	if (upd.Message == nil) {
		return false
	} else if (strings.HasSuffix(strings.ToLower(upd.Message.Sender.Username), "bot")) {
		return false
	} else if (time.Now().Unix() - upd.Message.Unixtime > 10) {
		return false
	}
	return true
}

func Init(bot *telebot.Bot, db *gorm.DB) {
	client := CreateTwitterClient(configuration)
	lastMessages := []telegramLastMessage{}

	if (!db.HasTable(&TweetInstruction{})) {
		db.CreateTable(&TweetInstruction{})
	}
	if (!db.HasTable(&TweetRegistered{})) {
		db.CreateTable(&TweetRegistered{})
	}

	bot.Handle("/why", func(m *telebot.Message) {
		if (m.ReplyTo == nil) {
			bot.Send(m.Chat, "Merci de répondre à un message contenant un tweet")
			return
		}
		id := m.ReplyTo.ID
		var tweet TweetRegistered
		if(db.Where(&TweetRegistered{MessageId: id}).First(&tweet).RecordNotFound()) {
			bot.Send(m.Chat, "Aucun enregistrement trouvé pour ce message")
			return
		}

		status := "échec"
		if (tweet.FetchSucess) {
			status = "réussite"
		}
		bot.Send(m.Chat, "Détail: "  + tweet.FetchStatus + " (" + status + ")")
	})

	bot.Handle("/retry", func(m *telebot.Message) {
		if (m.ReplyTo == nil) {
			bot.Send(m.Chat, "Merci de répondre à un message contenant un tweet")
			return
		}
		id := hasTweet(db, m.ReplyTo)
		if (id == nil) {
			bot.Send(m.Chat, "Aucun tweet n'a été trouvé dans ce message")
			return
		}
		tweets := fetchTweets(db, m, client, *id)

		if (len(tweets) < 2) {
			registerTweetStatus(db, m.ID, false, "Ce n'est pas un thread ou aucune réponse n'a pu etre trouvée")
			return
		}

		prepareToSend(db, bot, m.ReplyTo, tweets, "")
	})

	bot.Handle("/test", func(m *telebot.Message) {
		bot.Send(m.Chat, "It works!")
	})

	bot.Handle("/limit", func(m *telebot.Message) {
		print(m.Text + "\n")
		message := strings.Split(m.Text, " ")
		if (len(message) > 1) {
			/*count, err := strconv.ParseInt(message[1], 10, 32)
			if (err != nil) {
				bot.Send(m.Chat, "Impossible de récuperer la limite")
				return
			}*/
			instruction := TweetInstruction{
				Date: m.Unixtime,
				SenderId: m.Sender.ID,
				GroupId: m.Chat.ID,
				Instruction: message[1], //todo: change to enum like
				//todo: check if limit is bigger than 0
			}

			db.Create(&instruction) //Todo: check if existing
		} else {
			bot.Send(m.Chat, "Merci de préciser jusqu'à combien de tweet vous voulez limiter\nPar exemple: `/limit 2`")
		}
		//bot.Delete(m)
	})

	bot.Handle("/stop", func(m *telebot.Message) {
		instruction := TweetInstruction{
			Date: m.Unixtime,
			SenderId: m.Sender.ID,
			GroupId: m.Chat.ID,
			Instruction: "stop", //todo: change to enum like
		}

		db.Create(&instruction) //Todo: check if existing
		//bot.Delete(m)
	})

	bot.Handle(telebot.OnText, func(m *telebot.Message) {
		id := hasTweet(db, m)
		if id != nil {
			tweets := fetchTweets(db, m, client, *id)

			fmt.Printf("Len: %d\n", len(tweets))

			if (len(tweets) < 2) {
				registerTweetStatus(db, m.ID, false, "Ce n'est pas un thread ou aucune réponse n'a pu etre trouvée")
				return
			}

			lastMsg := getLastMessage(lastMessages, m.Chat.ID, m.Sender)

			if (!(lastMsg.time - 10 <= m.Unixtime && lastMsg.time + 10 >= m.Unixtime)) {
				lastMsg.message = ""
			}

			url := getTweetUrl(m)
			messageContent := strings.TrimSpace(strings.Replace(m.Text, url,"", -1))
			lastMessageContent := lastMsg.message
			userMessage := lastMessageContent

			if (userMessage != "") {
				userMessage = userMessage + "\n"
			}
			userMessage = userMessage + url + "\n"
			if (messageContent != "") {
				userMessage = userMessage + messageContent + "\n"
			}

			prepareToSend(db, bot, m, tweets, userMessage)

			if (lastMsg.message != "") {
				lastMessages = addLastMessage(lastMessages, telegramLastMessage{
					time: time.Now().Unix(),
					chatId: m.Chat.ID,
					user: m.Sender,
					message: "",
				})
			}
		} else {
			lastMessages = addLastMessage(lastMessages, telegramLastMessage{
				time: time.Now().Unix(),
				chatId: m.Chat.ID,
				user: m.Sender,
				message: m.Text,
			})
		}
	})
	configuration = map[string]string{}
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

	Init(b, engineDatabase)

	b.Start()
}

func Database() *gorm.DB {
	db, err := gorm.Open("sqlite3", "./TTOT.db")
	db.LogMode(true)

	if (err != nil) {
		log.Fatal("Unable to open database: " + err.Error() + "\n")
	}

	return db
}

func CreateTwitterClient(configuration map[string]string) *twitter.Client {
	consumerKey := os.Getenv("twitter_consumer_key")
	consumerSecret := os.Getenv("twitter_consumer_secret")
	accessToken := os.Getenv("twitter_access_token")
	accessSecret := os.Getenv("twitter_access_secret")

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)

	return client
}

//todo bug %  (!o(missing)) in tweet
// "%of" instead of "% of"

func SetupDatabase(db *gorm.DB) {
	engineDatabase = db
}

func Start() {
	SetupDatabase(Database())
	SetupBot()
}