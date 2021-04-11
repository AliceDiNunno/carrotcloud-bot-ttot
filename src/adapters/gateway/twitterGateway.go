package gateway

import (
	"adinunno.fr/twitter-to-telegram/src/core/domain"
	"adinunno.fr/twitter-to-telegram/src/core/usecases"
	"adinunno.fr/twitter-to-telegram/src/infra"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"sort"
	"time"
)

type twitterGateway struct {
	twitterClient *twitter.Client
}

func tweetToDomain(tweet *twitter.Tweet) (*domain.Tweet, error) {
	creationTime, err := time.Parse("Mon Jan 02 15:04:05 -0700 2006", tweet.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &domain.Tweet{
		CreationTime: creationTime,

		Id:       domain.TweetId(tweet.ID),
		UserName: tweet.User.ScreenName,
		Message:  tweet.FullText,
	}, nil
}

func (t twitterGateway) GetTweet(id domain.TweetId) (*domain.Tweet, error) {
	statusLookupParams := &twitter.StatusShowParams{
		TweetMode: "extended",
	}

	tweet, _, err := t.twitterClient.Statuses.Show(int64(id), statusLookupParams)

	if err != nil {
		return nil, err
	}

	return tweetToDomain(tweet)
}

func (t twitterGateway) GetFullThread(tweet *domain.Tweet) (domain.TweetList, error) {
	searchTweetParams := &twitter.SearchTweetParams{
		Query:      "from:" + tweet.UserName + " to:" + tweet.UserName,
		SinceID:    int64(tweet.Id),
		ResultType: "recent",
		TweetMode:  "extended",
		Count:      1000,
	}

	replies := domain.TweetList{tweet}

	tweets, _, err := t.twitterClient.Search.Tweets(searchTweetParams)

	if err != nil {
		return nil, err
	}

	statuses := tweets.Statuses
	sort.Slice(statuses, func(i, j int) bool {
		return statuses[i].ID < statuses[j].ID
	})

	for _, reply := range statuses {
		if replies.ContainsId(domain.TweetId(reply.InReplyToStatusID)) &&
			!replies.ContainsId(domain.TweetId(reply.ID)) {
			domainTweet, err := tweetToDomain(&reply)

			if err != nil {
				return nil, err
			}

			replies = append(replies, domainTweet)
		}
	}
	return replies, nil
}

func NewTwitterGateway(conf infra.TwitterConf) (usecases.TwitterGateway, error) {
	consumerKey := conf.ApiConsumerKey
	consumerSecret := conf.ApiConsumerSecret
	accessToken := conf.UserAccesToken
	accessSecret := conf.UserAccessSecret

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	twitterClient := twitter.NewClient(httpClient)

	return twitterGateway{
		twitterClient: twitterClient,
	}, nil
}
