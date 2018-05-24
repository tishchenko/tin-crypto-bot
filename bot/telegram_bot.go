package bot

import (
	"github.com/tishchenko/tin-crypto-bot/config"
	"net/http"
	"golang.org/x/net/proxy"
	"log"
	"net"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"context"
)

const version = "0.1.0 alpha"

type TelegramBot struct {
	Bot    *tgbotapi.BotAPI
	States map[int64]State
}

type State struct {
}

type AlarmMesData struct {
}

func NewTelegramBot(conf *config.Config) *TelegramBot {
	bot := &TelegramBot{}
	bot.initStates()

	var err error
	var client *http.Client

	if conf == nil {
		panic("Для Telegram Bot не заданы настройки!")
	}

	if conf.Proxy != nil {
		auth := &proxy.Auth{conf.Proxy.User, conf.Proxy.Password}

		dialer, err := proxy.SOCKS5("tcp", conf.Proxy.IpAddr+":"+conf.Proxy.Port, auth, proxy.Direct)
		if err != nil {
			log.Panic(err)
		}
		transport := &http.Transport{}
		transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
		}
		client = &http.Client{Transport: transport}
	} else {
		client = &http.Client{}
	}

	bot.Bot, err = tgbotapi.NewBotAPIWithClient(conf.TelApiToken, client)
	if err != nil {
		log.Panic(err)
	}

	bot.Bot.Debug = true

	log.Printf("Authorized on account %s", bot.Bot.Self.UserName)

	return bot
}

func (bot *TelegramBot) Run() {

}

func (bot *TelegramBot) persistStates() {

}

func (bot *TelegramBot) initStates() {
	bot.States = map[int64]State{}

}
