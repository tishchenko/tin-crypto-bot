package market

import (
	//https://github.com/adshao/go-binance
	"github.com/adshao/go-binance"
	"github.com/tishchenko/tin-crypto-bot/config"
	"fmt"
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

	wsKlineHandler := func(event *binance.WsKlineEvent) {
		fmt.Println(event.Kline)
	}
	doneC, err := binance.WsKlineServe("LTCBTC", "1m", wsKlineHandler)
	if err != nil {
		fmt.Println(err)
		return
	}
	<-doneC
}
