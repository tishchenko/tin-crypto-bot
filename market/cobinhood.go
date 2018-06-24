package market

import "github.com/tishchenko/tin-crypto-bot/config"

type Cobinhood struct {

}

func NewCobinhood(conf *config.MarketConfig) *Cobinhood {
	c := &Cobinhood{}
	/*key := conf.Key
	secret := conf.Secret
	c.Client = cobinhood.NewClient(key, secret)*/
	return c
}

func (market *Cobinhood) CreateOrder(symbol string, orderType OrderType, quantity float64, price float64) (string, error) {
	return "", nil
}

func (market *Cobinhood) CancelOrder(orderId string) error {
	return nil
}