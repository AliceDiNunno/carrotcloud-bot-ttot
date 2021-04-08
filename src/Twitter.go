package src

import (
	"adinunno.fr/twitter-to-telegram/src/adapters/persistence/postgres"
	"adinunno.fr/twitter-to-telegram/src/infra"
	"fmt"
	"github.com/dghubble/oauth1"
	"gopkg.in/tucnak/telebot.v2"
	"regexp"
	"sort"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/dghubble/go-twitter/twitter"
	_ "github.com/mattn/go-sqlite3"
)

var twitterClient *twitter.Client

type telegramLastMessage struct {
	time    int64
	chatId  int64
	user    *telebot.User
	message string
}

func getTweetUrl(m *telebot.Message) string {
	r := regexp.MustCompile(`https:\/\/twitter.com\/([a-zA-Z0-9_]*)\/status\/([0-9]*)`)
	match := r.FindStringSubmatch(m.Text)
	if len(match) > 2 {
		return "https://twitter.com/" + match[1] + "/status/" + match[2]
	}
	return ""
}

func fetchTweets(db *gorm.DB, m *telebot.Message, client *twitter.Client, id int64) []twitter.Tweet {
	tweetList := []twitter.Tweet{}

	statusLookupParams := &twitter.StatusShowParams{
		TweetMode: "extended",
	}
	tweet, _, err := client.Statuses.Show(id, statusLookupParams)

	if err != nil {
		print(err.Error())
		registerTweetStatus(db, m.ID, false, "Erreur twitter: "+err.Error())
		return []twitter.Tweet{}
	}

	tweettime, err := time.Parse("Mon Jan 02 15:04:05 -0700 2006", tweet.CreatedAt)
	if err != nil {
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

	searchTweetParams := &twitter.SearchTweetParams{
		Query:      "from:" + tweet.User.ScreenName + " to:" + tweet.User.ScreenName,
		SinceID:    tweet.ID,
		ResultType: "recent",
		TweetMode:  "extended",
		Count:      1000,
	}
	tweets, _, _ := client.Search.Tweets(searchTweetParams)

	statuses := tweets.Statuses
	sort.Slice(statuses, func(i, j int) bool {
		return statuses[i].ID < statuses[j].ID
	})

	var knownids []string
	knownids = append(knownids, tweet.IDStr)

	for _, twt := range statuses {
		if contains(knownids, twt.InReplyToStatusIDStr) {
			knownids = append(knownids, twt.IDStr)
			tweetList = append(tweetList, twt)
		}
	}

	return tweetList
}

func registerTweetStatus(db *gorm.DB, id int, success bool, reason string) {
	print("Tweet status: " + reason)

	db.Save(&postgres.TweetRegistered{
		MessageId:    id,
		FetchSuccess: success,
		FetchStatus:  reason,
	})
}

func CreateTwitterClient(conf infra.TwitterConf) {
	consumerKey := conf.ApiConsumerKey
	consumerSecret := conf.ApiConsumerSecret
	accessToken := conf.UserAccesToken
	accessSecret := conf.UserAccessSecret

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	twitterClient = twitter.NewClient(httpClient)
}

//todo bug %  (!o(missing)) in tweet
// "%of" instead of "% of"
