package main

import (
	"github.com/tishchenko/tin-crypto-bot/config"
	"github.com/tishchenko/tin-crypto-bot/bot"
)

func main() {
	conf := config.NewConfig()
	println(conf)
	logic := config.NewLogic()
	println(logic)

	telBot := bot.NewTelegramBot(conf)
	telBot.Run()

	cryptoBot := bot.NewCryptoBot(logic)
	cryptoBot.Run()
}
