package bot

import (
	"github.com/tishchenko/tin-crypto-bot/config"
	"time"
	"github.com/tishchenko/tin-crypto-bot/market"
	"log"
)

type CryptoBot struct {
	Logic     config.Logic
	Markets   map[string]market.Market
	consumers map[string][]chan CryptoBotEvent
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

func NewCryptoBot(marketsConf *map[string]config.MarketConfig, logic *config.Logic) *CryptoBot {
	bot := &CryptoBot{}

	bot.Markets = make(map[string]market.Market)
	for marketName, marketConf := range *marketsConf {
		bot.Markets[marketName] = *market.NewMarket(marketName, &marketConf)
	}

	bot.Logic = *logic

	return bot
}

func (bot *CryptoBot) Run() {
	klineHandler := func(kline *market.Kline) {
		bot.Emit("TelBot", CryptoBotEvent{PriceJumpUp})
		log.Println(*kline)
	}
	// TODO Demo
	for _, market := range bot.Markets {
		market.SetKlinesHandler(
			market.Name,
			"LTCBTC",
			"1h",
			klineHandler,
		)
	}

	go func() {
		for {

			print("-")
			// TODO Demo
			bot.Emit("TelBot", CryptoBotEvent{SetLimitBuyOrder})
			time.Sleep(5 * time.Second)
		}
	}()
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
