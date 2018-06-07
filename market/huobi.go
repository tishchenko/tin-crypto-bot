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
