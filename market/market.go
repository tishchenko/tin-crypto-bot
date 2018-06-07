package market

import "github.com/tishchenko/tin-crypto-bot/config"

type Market struct {
	Name      string
	bitfinex  *Bitfinex
	binance   *Binance
	cobinhood *Cobinhood
	huobi     *Huobi
}

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

func (market *Market) Klines() {
	// TODO Demo
	if market.binance != nil {
		market.binance.Klines()
	}
}
