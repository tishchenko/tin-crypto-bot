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

type OrderType string

const (
	OrderTypeBuyLimit  OrderType = "BUY-LIMIT"
	OrderTypeSellLimit OrderType = "SELL-LIMIT"
	OrderTypeBuyMarket  OrderType = "BUY-MARKET"
	OrderTypeSellMarket OrderType = "SELL-MARKET"
)

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
			handler(binanceKlineConv(event.Kline))
		}
		market.binance.SetKlinesHandler(symbol, interval, wsKlineHandler)
	}
	if marketName == "bitfinex" && market.bitfinex != nil {

	}
	if marketName == "cobinhood" && market.cobinhood != nil {

	}
	if marketName == "huobi" && market.huobi != nil {

	}
}

func (market *Market) CreateOrder(marketName string, symbol string, orderType OrderType, quantity float64, price float64) (string, error) {
	if marketName == "binance" && market.binance != nil {
		return market.binance.CreateOrder(symbol, orderType, quantity, price)
	}
	if marketName == "bitfinex" && market.bitfinex != nil {
		return market.bitfinex.CreateOrder(symbol, orderType, quantity, price)
	}
	if marketName == "cobinhood" && market.cobinhood != nil {
		return market.cobinhood.CreateOrder(symbol, orderType, quantity, price)
	}
	if marketName == "huobi" && market.huobi != nil {
		return market.huobi.CreateOrder(symbol, orderType, quantity, price)
	}
	return "", nil
}

func (market *Market) CancelOrder(marketName string, orderId string) error {
	if marketName == "binance" && market.binance != nil {
		return market.binance.CancelOrder(orderId)
	}
	if marketName == "bitfinex" && market.bitfinex != nil {
		return market.bitfinex.CancelOrder(orderId)
	}
	if marketName == "cobinhood" && market.cobinhood != nil {
		return market.cobinhood.CancelOrder(orderId)
	}
	if marketName == "huobi" && market.huobi != nil {
		return market.huobi.CancelOrder(orderId)
	}
	return nil
}

func binanceKlineConv(kline binance.WsKline) *Kline {
	k := &Kline{
		StartTime:    kline.StartTime,
		EndTime:      kline.EndTime,
		Symbol:       kline.Symbol,
		Interval:     kline.Interval,
		FirstTradeID: kline.FirstTradeID,
		LastTradeID:  kline.LastTradeID,
		Open:         kline.Open,
		Close:        kline.Close,
		High:         kline.High,
		Low:          kline.Low,
		Volume:       kline.Volume,
		IsFinal:      kline.IsFinal,
	}
	return k
}
