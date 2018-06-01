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

	cryptoBot := bot.NewCryptoBot(logic)

	botAdapter := bot.NewBotAdapter(cryptoBot, telBot)
	botAdapter.Run()
}
