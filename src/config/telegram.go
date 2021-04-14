package config

type TelegramConf struct {
	BotToken string
}

func LoadTelegramConfiguration() TelegramConf {
	return TelegramConf{RequireEnvString("TELEGRAM_BOT_TOKEN")}
}
