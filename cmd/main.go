package main

import (
	"adinunno.fr/twitter-to-telegram/src"
	"adinunno.fr/twitter-to-telegram/src/infra"
)

func main() {
	infra.LoadEnv()

	twitterConf := infra.LoadTwitterConfiguration()
	telegramConf := infra.LoadTelegramConfiguration()

	src.OpenDatabase()
	src.CreateTwitterClient(twitterConf)
	src.SetupBot(telegramConf)
}
