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

func (market *Binance) Klines() {
	/*klines, err := market.Client.NewKlinesService().Symbol("LTCBTC").
		Interval("15m").Do(context.Background())
	if err != nil {
		return
	}
	for _, k := range klines {
		fmt.Println(k)
	}*/

	errHandler := func(err error) {
		log.Println(err)
	}
	wsKlineHandler := func(event *binance.WsKlineEvent) {
		log.Println(event.Kline)
	}
	doneC, _, err := binance.WsKlineServe("LTCBTC", "1h", wsKlineHandler, errHandler)
	if err != nil {
		log.Println(err)
		return
	}
	<-doneC
}
