package market

import (
	//https://github.com/adshao/go-binance
	"github.com/adshao/go-binance"
	"github.com/tishchenko/tin-crypto-bot/config"
	"log"
)

type Binance struct {
	Client *binance.Client
}

func NewBinance(conf *config.MarketConfig) *Binance {
	c := &Binance{}

	key := conf.Key
	secret := conf.Secret
	c.Client = binance.NewClient(key, secret)
	c.Client.UserAgent = "TinCryptoBot"
	c.Client.Debug = true

	return c
}

func (market *Binance) SetKlinesHandler(symbol string, interval string, handler binance.WsKlineHandler) {
	errHandler := func(err error) {
		log.Println(err)
	}

	//doneC, _, err := binance.WsKlineServe("LTCBTC", "1h", handler, errHandler)
	go func() {
		doneC, _, err := binance.WsKlineServe(symbol, interval, handler, errHandler)
		if err != nil {
			log.Println(err)
			return
		}
		<-doneC
	}()
}
