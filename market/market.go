package market

import (
	"github.com/tishchenko/tin-crypto-bot/config"
	"github.com/adshao/go-binance"
)

type Market struct {
	Name      string
	bitfinex  *Bitfinex
	binance   *Binance
	cobinhood *Cobinhood
	huobi     *Huobi
}

type Kline struct {
	StartTime    int64
	EndTime      int64
	Symbol       string
	Interval     string
	FirstTradeID int64
	LastTradeID  int64
	Open         string
	Close        string
	High         string
	Low          string
	Volume       string
	IsFinal      bool
}

type KlineHandler func(kline *Kline)

func NewMarket(marketName string, conf *config.MarketConfig) *Market {
	m := &Market{}

	m.Name = marketName
	switch m.Name {
	case "binance":
		m.binance = NewBinance(conf)
	case "bitfinex":
		m.bitfinex = NewBitfinex(conf)
	case "cobinhood":
		m.cobinhood = NewCobinhood(conf)
	case "huobi":
		m.huobi = NewHuobi(conf)
	}

	return m
}

func (market *Market) SetKlinesHandler(marketName string, symbol string, interval string, handler KlineHandler) {

	if marketName == "binance" && market.binance != nil {
		wsKlineHandler := func(event *binance.WsKlineEvent) {
			handler(binanceKlineConverter(event.Kline))
		}
		market.binance.SetKlinesHandler(symbol, interval, wsKlineHandler)
	}
}

func binanceKlineConverter(kline binance.WsKline) *Kline {
	k := &Kline{
		StartTime: kline.StartTime,
		EndTime: kline.EndTime,
		Symbol: kline.Symbol,
		Interval: kline.Interval,
		FirstTradeID: kline.FirstTradeID,
		LastTradeID: kline.LastTradeID,
		Open: kline.Open,
		Close: kline.Close,
		High: kline.High,
		Low: kline.Low,
		Volume: kline.Volume,
		IsFinal: kline.IsFinal,
	}
	return k
}