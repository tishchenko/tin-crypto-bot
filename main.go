package main

import (
	"github.com/tishchenko/tin-crypto-bot/config"
	"github.com/tishchenko/tin-crypto-bot/bot"
	"github.com/tishchenko/tin-crypto-bot/logger"
)

func main() {
	conf := config.NewConfig()
	println(conf)
	logic := config.NewLogic()
	println(logic)

	telBot := bot.NewTelegramBot(conf)

	cryptoLogger := logger.NewCryptoLogger()
	cryptoBot := bot.NewCryptoBot(&conf.Markets, logic, cryptoLogger)

	botAdapter := bot.NewBotAdapter(cryptoBot, telBot)
	botAdapter.Run()
}
