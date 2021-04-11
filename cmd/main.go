package main

import (
	"adinunno.fr/twitter-to-telegram/src/adapters/gateway"
	"adinunno.fr/twitter-to-telegram/src/adapters/persistence/sqlite"
	"adinunno.fr/twitter-to-telegram/src/adapters/telegram"
	"adinunno.fr/twitter-to-telegram/src/core/usecases"
	"adinunno.fr/twitter-to-telegram/src/infra"
	"log"
)

func main() {
	infra.LoadEnv()

	db := sqlite.CreateDB()

	twitterConf := infra.LoadTwitterConfiguration()
	telegramConf := infra.LoadTelegramConfiguration()

	var statusRepo usecases.StatusRepo
	var instructionRepo usecases.InstructionRepo
	var twitterGateway usecases.TwitterGateway

	statusRepo = sqlite.StatusRepo{Db: db}
	instructionRepo = sqlite.InstructionRepo{Db: db}
	twitterGateway, err := gateway.NewTwitterGateway(twitterConf)

	if err != nil {
		log.Fatalln(err)
	}

	usecasesHandler := usecases.NewInteractor(statusRepo, instructionRepo, twitterGateway)

	bot := telegram.NewTelegramBot(telegramConf)
	commandHandler := telegram.NewCommandHandler(bot, usecasesHandler)
	commandHandler.SetCommands()
	bot.Start()
}
