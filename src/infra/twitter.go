package infra

type TwitterConf struct {
	ApiConsumerKey    string
	ApiConsumerSecret string

	UserAccesToken   string
	UserAccessSecret string
}

func LoadTwitterConfiguration() TwitterConf {
	return TwitterConf{
		RequireEnvString("TWITTER_CONSUMER_KEY"),
		RequireEnvString("TWITTER_CONSUMER_SECRET"),
		RequireEnvString("TWITTER_ACCESS_TOKEN"),
		RequireEnvString("TWITTER_ACCESS_SECRET"),
	}
}
