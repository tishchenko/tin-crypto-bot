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
  "markets": [
    {}
  ]
}
 */

import (
	"os"
	"encoding/json"
)

const (
	defConfigFileName = "config.json"
)

type Config struct {
	TelApiToken string         `json:"token"`
	Proxy       *Proxy         `json:"proxy,omitempty"`
	Markets     []MarketConfig `json:"markets"`
}

type Proxy struct {
	IpAddr   string `json:"ip"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type MarketConfig struct {

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
		panic("Config file is wrong!")
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&c)
	if err != nil {
		return c
	}

	if c.Markets == nil || len(c.Markets) < 1 {
		panic("Can't find markets configuration!")
	}

	return c
}