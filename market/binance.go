package market

import (
	//https://github.com/adshao/go-binance
	"github.com/adshao/go-binance"
	"github.com/tishchenko/tin-crypto-bot/config"
	"log"
	"context"
	"strconv"
	"strings"
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
		doneC, _, err := binance.WsKlineServe(
			market.convSymbol(symbol),
			interval,
			handler,
			errHandler,
		)
		if err != nil {
			log.Println(err)
			return
		}
		<-doneC
	}()
}

func (market *Binance) CreateOrder(symbol string, orderType OrderType, quantity float64, price float64) (string, error) {
	sideType1, orderType1 := market.convOrderType(orderType)
	order, err := market.Client.NewCreateOrderService().
		Symbol(market.convSymbol(symbol)).
		Side(sideType1).
		Type(orderType1).
	//TimeInForce(binance.TimeInForceGTC).
		Quantity(market.convFloatNum(quantity)).
		Price(market.convFloatNum(price)).
		Do(context.Background())
	if err != nil {
		log.Println(err)
		return "", err
	}
	log.Println(order)

	return market.convOrderId(order.OrderID), nil
}

func (market *Binance) CancelOrder(orderId string) error {
	_, err := market.Client.NewCancelOrderService().
	//Symbol("BNBETH"). //TODO
		OrderID(4432844).Do(context.Background())
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (market *Binance) convSymbol(symbol string) string {
	symbols := strings.Split(symbol, "-")
	symbol = strings.Join(symbols, "")
	return symbol
}

func (market *Binance) convOrderId(orderId int64) string {
	return strconv.Itoa(int(orderId))
}

func (market *Binance) convFloatNum(num float64) string {
	return strconv.FormatFloat(num, 'f', 8, 64)
}

func (market *Binance) convOrderType(orderType OrderType) (binance.SideType, binance.OrderType) {
	var sideTypeRes binance.SideType
	var orderTypeRes binance.OrderType
	orderTypes := strings.Split(string(orderType), "-")
	if len(orderTypes) < 2 {
		return sideTypeRes, orderTypeRes
	}
	switch orderTypes[0] {
	case "BUY":
		sideTypeRes = binance.SideTypeBuy
	case "SELL":
		sideTypeRes = binance.SideTypeSell
	}
	switch orderTypes[1] {
	case "LIMIT":
		orderTypeRes = binance.OrderTypeLimit
	case "MARKET":
		orderTypeRes = binance.OrderTypeMarket
	}
	return sideTypeRes, orderTypeRes
}
