package bot

import (
	"github.com/tishchenko/tin-crypto-bot/config"
	"github.com/tishchenko/tin-crypto-bot/market"
	"log"
	"github.com/tishchenko/tin-crypto-bot/logger"
)

type CryptoBot struct {
	Logic     config.Logic
	Markets   map[string]market.Market
	consumers map[string][]chan CryptoBotEvent
	logger    logger.CryptoLogger
}

type CryptoBotEventType string

const (
	SetLimitBuyOrder    CryptoBotEventType = "Установлен лимитный ордер на покупку"
	SetLimitSellOrder   CryptoBotEventType = "Установлен лимитный ордер на продажу"
	FillBuyOrder        CryptoBotEventType = "Выполнен лимитный ордер на покупку"
	FillSellOrder       CryptoBotEventType = "Выполнен лимитный ордер на продажу"
	FillMarketBuyOrder  CryptoBotEventType = "Выполнен рыночный ордер на покупку"
	FillMarketSellOrder CryptoBotEventType = "Выполнен рыночный ордер на продажу"
	PriceIsHigher       CryptoBotEventType = "Цена превысила отметку"
	PriceIsLower        CryptoBotEventType = "Цена опустилась ниже отметки"
	PriceJumpUp         CryptoBotEventType = "Скачёк цены вверх"
	PriceJumpDown       CryptoBotEventType = "Скачёк цены вниз"
)

type CryptoBotEvent struct {
	EventType CryptoBotEventType
}

func NewCryptoBot(marketsConf *map[string]config.MarketConfig, logic *config.Logic, logger *logger.CryptoLogger) *CryptoBot {
	bot := &CryptoBot{}

	bot.Markets = make(map[string]market.Market)
	for marketName, marketConf := range *marketsConf {
		bot.Markets[marketName] = *market.NewMarket(marketName, &marketConf)
	}

	bot.Logic = *logic
	bot.logger = *logger

	return bot
}

func (bot *CryptoBot) Run() {

	bot.signalStrategiesProcessing(bot.Logic.Strategies.Signals)

	// TODO Demo
	/*for _, market := range bot.Markets {
		market.SetKlinesHandler(
			market.Name,
			"LTCBTC",
			"1h",
			klineHandler,
		)
	}*/

	/*go func() {
		for {

			print("-")
			// TODO Demo
			bot.Emit("TelBot", CryptoBotEvent{SetLimitBuyOrder})
			time.Sleep(5 * time.Second)
		}
	}()*/
}

func (bot *CryptoBot) AddConsumer(e string, ch chan CryptoBotEvent) {
	if bot.consumers == nil {
		bot.consumers = make(map[string][]chan CryptoBotEvent)
	}
	if _, ok := bot.consumers[e]; ok {
		bot.consumers[e] = append(bot.consumers[e], ch)
	} else {
		bot.consumers[e] = []chan CryptoBotEvent{ch}
	}
}

func (bot *CryptoBot) RemoveConsumer(e string, ch chan CryptoBotEvent) {
	if _, ok := bot.consumers[e]; ok {
		for i := range bot.consumers[e] {
			if bot.consumers[e][i] == ch {
				bot.consumers[e] = append(bot.consumers[e][:i], bot.consumers[e][i+1:]...)
				break
			}
		}
	}
}

func (bot *CryptoBot) Emit(e string, response CryptoBotEvent) {
	if _, ok := bot.consumers[e]; ok {
		for _, handler := range bot.consumers[e] {
			go func(handler chan CryptoBotEvent) {
				handler <- response
			}(handler)
		}
	}
}

func (bot *CryptoBot) signalStrategiesProcessing(signals *[]config.Signal) {
	for _, signal := range *signals {
		if signal.Active != nil && !*signal.Active {
			continue
		}
		market, ok := bot.Markets[signal.Market]
		if !ok {
			log.Printf("Market %s for signal %s is not found.", signal.Market, signal.Id)
			continue
		}
		market.Name = signal.Market
		market.SetKlinesHandler(
			market.Name,
			signal.Pair,
			"1m",
			bot.cookSignalHandler(signal),
		)
	}
}

func (bot *CryptoBot) cookSignalHandler(signal config.Signal) market.KlineHandler {
	return func(kline *market.Kline) {
		// TODO Delete this Demo
		bot.Emit("TelBot", CryptoBotEvent{PriceJumpUp})

		log.Print(signal.Id)
		log.Print(signal.Hash())
		log.Println(kline.Close)
	}
}
