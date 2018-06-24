package market

import "github.com/tishchenko/tin-crypto-bot/config"

type Huobi struct {
}

func NewHuobi(conf *config.MarketConfig) *Huobi {
	c := &Huobi{}
	/*key := conf.Key
	secret := conf.Secret
	c.Client = huobi.NewClient(key, secret)*/
	return c
}

func (market *Huobi) CreateOrder(symbol string, orderType OrderType, quantity float64, price float64) (string, error) {
	return "", nil
}

func (market *Huobi) CancelOrder(orderId string) error {
	return nil
}