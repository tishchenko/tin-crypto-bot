package main

import "github.com/tishchenko/tin-crypto-bot/config"

func main() {
	conf := config.NewConfig()
	println(conf)
	logic := config.NewLogic()
	println(logic)
}
