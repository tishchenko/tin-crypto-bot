package bot

import "github.com/tishchenko/tin-crypto-bot/config"

type CryptoBot struct {
	Logic config.Logic
}

func NewCryptoBot(conf *config.Logic) *CryptoBot {
	bot := &CryptoBot{}
	bot.Logic = *conf

	return bot
}

func (bot *CryptoBot) Run() {

}