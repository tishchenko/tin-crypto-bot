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