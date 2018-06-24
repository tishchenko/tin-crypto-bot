package market

import "github.com/tishchenko/tin-crypto-bot/config"

type Bitfinex struct {

}

func NewBitfinex(conf *config.MarketConfig) *Bitfinex {
	c := &Bitfinex{}
	/*key := conf.Key
	secret := conf.Secret
	c.Client = bitfinex.NewClient(key, secret)*/
	return c
}

func (market *Bitfinex) CreateOrder(symbol string, orderType OrderType, quantity float64, price float64) (string, error) {
	return "", nil
}

func (market *Bitfinex) CancelOrder(orderId string) error {
	return nil
}