package domain

type TweetList []*Tweet

func (tl TweetList) Contains(tweet *Tweet) bool {
	return tl.ContainsId(tweet.Id)
}

func (tl TweetList) ContainsId(id TweetId) bool {
	for _, entry := range tl {
		if entry.Id == id {
			return true
		}
	}
	return false
}
