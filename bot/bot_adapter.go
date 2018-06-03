package bot

type BotAdapter struct {
	cryptoBot   *CryptoBot
	telegramBot *TelegramBot
}

func NewBotAdapter(cryptoBot *CryptoBot, telegramBot *TelegramBot) *BotAdapter {
	bot := &BotAdapter{}
	bot.cryptoBot = cryptoBot
	bot.telegramBot = telegramBot
	return bot
}

func (bot *BotAdapter) Run() {
	e := make(chan CryptoBotEvent)

	go func() {
		for {
			ev := <-e
			bot.telegramBot.MesChan <- "Event " + string(ev.EventType)
		}
	}()

	if bot.cryptoBot != nil {
		bot.cryptoBot.AddConsumer("TelBot", e)
		bot.cryptoBot.Run()
	}
	if bot.telegramBot != nil {
		bot.telegramBot.Run()
	}
}
