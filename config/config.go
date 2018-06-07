package config

/**
Config file template:
{
  "token": "",
  "proxy": {
    "ip": "",
    "port": "",
    "user": "",
    "password": ""
  },
  "markets": {
    "marketName": {
		"key": "",
		"secret": ""
	}
  }
}
 */

import (
	"os"
	"encoding/json"
	"github.com/tishchenko/tin-crypto-bot/utils"
	"log"
)

const (
	defConfigFileName = "config.json"
)

type Config struct {
	TelApiToken string                  `json:"token"`
	Proxy       *Proxy                  `json:"proxy,omitempty"`
	Markets     map[string]MarketConfig `json:"markets"`
}

type Proxy struct {
	IpAddr   string `json:"ip"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}

var marketNames = []string{"bitfinex", "binance", "cobinhood", "huobi"}

type MarketConfig struct {
	Key    string `json:"key"`
	Secret string `json:"secret"`
}

func NewConfig() *Config {
	return NewConfigWithCustomFile("")
}

func NewConfigWithCustomFile(fileName string) *Config {

	if fileName == "" {
		fileName = defConfigFileName
	}

	c := &Config{}

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalln("Can't open configuration file!")
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&c)
	if err != nil {
		log.Fatalln("Configuration file is wrong!")
	}

	if c.Markets == nil || len(c.Markets) < 1 {
		log.Fatalln("Can't find markets configuration!")
	}

	c.validate()

	return c
}

func (c *Config) validate() {
	for market, _ := range c.Markets {
		if !utils.StringInSlice(market, marketNames) {
			log.Println("Unknown market \"" + market + "\" in configuration file.")
			delete(c.Markets, market)
		}
	}
}
