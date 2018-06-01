package bot

type BotAdapter struct {
	cryptoBot   *CryptoBot
	telegramBot *TelegramBot
	//cryptoBotChannels  map[string][]chan CryptoBotEvent
	//telegramBotChannel chan string
}

func NewBotAdapter(cryptoBot *CryptoBot, telegramBot *TelegramBot) *BotAdapter {
	bot := &BotAdapter{}
	//bot.cryptoBotChannels = cryptoBot.consumers
	//bot.telegramBotChannel = telegramBot.MesChan
	bot.cryptoBot = cryptoBot
	bot.telegramBot = telegramBot
	return bot
}

func (bot *BotAdapter) Run() {
	var e chan CryptoBotEvent

	if bot.cryptoBot != nil {
		bot.cryptoBot.AddConsumer("TelBot", e)
		bot.cryptoBot.Run()
	}
	if bot.telegramBot != nil {
		bot.telegramBot.Run()
	}

	// TODO https://flaviocopes.com/golang-event-listeners/
	go func(event chan CryptoBotEvent) {
		ev := <- e
		bot.telegramBot.MesChan <- "Event " + string(ev.EventType)
	}(e)
}
