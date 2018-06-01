package bot

import (
	"github.com/tishchenko/tin-crypto-bot/config"
	"net/http"
	"golang.org/x/net/proxy"
	"log"
	"net"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"context"
	"reflect"
)

const version = "0.1.0 alpha"

type TelegramBot struct {
	Bot    *tgbotapi.BotAPI
	States map[int64]State
	MesChan chan string
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
	bot.MesChan = make(chan string)

	log.Printf("Authorized on account %s", bot.Bot.Self.UserName)

	return bot
}

func (bot *TelegramBot) Run() {
	//go bot.poll(m)
	go bot.sendBroadcastMessage(bot.MesChan)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.Bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {

			switch update.Message.Command() {
			case "start":
				bot.States[update.Message.Chat.ID] = State{}
				bot.sendMessage(update.Message.Chat.ID, "Привет, "+update.Message.From.FirstName+" \xE2\x9C\x8C\n Я запущен.")
			case "stop":
				delete(bot.States, update.Message.Chat.ID)
				bot.sendMessage(update.Message.Chat.ID, "Ну, ты это, зови если что...")
				bot.sendMessage(update.Message.Chat.ID, "\xF0\x9F\x92\xA4")
				/*case "help":
					help := "Доступны команды:\n" +
						"/start - запуск мониторинга состояния очередей\n" +
						"/stop - приостановка мониторинга состояния очередей\n" +
						"/health - проверка бота на работоспособность\n" +
						"/queue (queueName [queueType]) - вывод статистики по очереди queueName (<b>" + strings.Join(queueNames, ", ") + "</b>) за последнее время; queueType - если не задан, то выводится статистика по normal queue, если <b>EQ</b>, то выводится статистика по exception queue\n"
					//"/rules - выводит правила оповещения, указанные в настройках"
					bot.sendMessage(update.Message.Chat.ID, help)
				case "health":
					bot.printHealth(update.Message.Chat.ID)
				case "queue":
					bot.printQueueStat(update.Message.Chat.ID, update.Message.CommandArguments())
				case "rules":
					bot.printRules(update.Message.Chat.ID)*/
			}
		}

	}
}

func (bot *TelegramBot) sendBroadcastMessage(m chan string) {
	for {
		message := <-m
		for chatID, state := range bot.States {
			bot.sendMessage(chatID, message)
			println(chatID, ": ", reflect.ValueOf(&state).String())
		}

		print("+")
	}
}

func (bot *TelegramBot) sendMessage(chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ParseMode = tgbotapi.ModeHTML
	bot.Bot.Send(msg)
}

func (bot *TelegramBot) persistStates() {

}

func (bot *TelegramBot) initStates() {
	bot.States = map[int64]State{}

}
