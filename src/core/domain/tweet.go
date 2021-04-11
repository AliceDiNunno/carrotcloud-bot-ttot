package domain

import (
	"fmt"
	"time"
)

type TweetId int64

type Tweet struct {
	CreationTime time.Time

	Id       TweetId
	UserName string
	Message  string
}

func (t Tweet) IsWithinTimeLimit() bool { //Tweets can live up to 15 days before twitter do not provide the API with its replies
	fifteenDays := 15 * 24 * 60 * 60

	if (time.Now().Unix() - t.CreationTime.Unix()) > int64(fifteenDays) {
		return false
	}

	return true
}

func (t Tweet) Url() string {
	return fmt.Sprintf("https://twitter.com/%s/status/%d", t.UserName, t.Id)
}
