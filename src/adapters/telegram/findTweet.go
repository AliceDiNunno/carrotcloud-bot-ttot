package telegram

import (
	"adinunno.fr/twitter-to-telegram/src/core/domain"
	"errors"
	"regexp"
	"strconv"
)

func findTweet(text string) (domain.TweetId, error) {
	r := regexp.MustCompile(`https:\/\/twitter.com\/[a-zA-Z0-9_]*\/status\/([0-9]*)`)
	match := r.FindStringSubmatch(text)
	if len(match) > 1 {
		id := match[1]

		id64, err := strconv.ParseInt(id, 10, 64)
		if err != nil { //conversion failed so we have a bad id on our hands
			//here I'm doing the unusual "bypass an error to create a new one" because I wanted to be a little bit more explicit
			//feel free to open an issue and tell my that I'm wrong with a better way to handle this. I'll listen :)
			return -1, errors.New("malformed or inexistant tweet id")
		}
		return domain.TweetId(id64), nil
	}
	return -1, errors.New("no twitter url found")
}
