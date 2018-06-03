package config

import (
	"os"
	"log"
	"encoding/json"
	"strings"
	"github.com/tishchenko/tin-crypto-bot/utils"
	"math"
	"github.com/satori/go.uuid"
)

const (
	defLogicFileName = "logic.json"
)

type Logic struct {
	FileName   string      `json:"-"`
	Strategies *Strategies `json:"strategies"`
}

type Strategies struct {
	Signals *[]Signal `json:"signals, omitempty"`
}

type Strategy struct {
	Id       string  `json:"id"`
	Market   string  `json:"market"`
	Pair     string  `json:"pair"`
	TradeCap float32 `json:"tradeCap"`
	Active   *bool   `json:"active, omitempty"`
}

type Signal struct {
	Strategy
	BuyLevels  *[]Level `json:"buyLevels"`
	SellLevels *[]Level `json:"sellLevels"`
	StopLoss   *Level   `json:"stopLoss, omitempty"`
}

type Level struct {
	Price    float32 `json:"price"`
	Percent  float32 `json:"percent"`
	StopLoss *Level  `json:"stopLoss, omitempty"`
}

func NewLogic() *Logic {
	return NewLogicWithCustomFile("")
}

func NewLogicWithCustomFile(fileName string) *Logic {

	logic := &Logic{}
	logic.FileName = fileName

	logic.Reload()
	logic.validate()

	return logic
}

func (logic *Logic) Reload() {
	if logic.FileName == "" {
		logic.FileName = defLogicFileName
	}

	file, err := os.Open(logic.FileName)
	if err != nil {
		log.Fatalln("Can't open logic file!")
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&logic)
	if err != nil {
		log.Fatalln("Logic file is wrong!")
	}
}

func (logic *Logic) validate() {
	if logic.Strategies == nil {
		log.Fatalln("Can't find \"Strategies\" block!")
	}
	if logic.Strategies.Signals != nil {
		for _, signal := range *logic.Strategies.Signals {
			logic.validateSignal(signal)
		}
	}
}

func (logic *Logic) validateSignal(signal Signal) bool {
	if !logic.validateSignalMarket(signal) {
		return false
	}

	if !logic.validateSignalPair(signal) {
		return false
	}

	if !logic.validateSignalId(signal) {
		return false
	}

	if !logic.validateSignalTradeCap(signal) {
		return false
	}

	if !logic.validateSignalBuyLevels(signal) {
		return false
	}

	if !logic.validateSignalSellLevels(signal) {
		return false
	}

	return true
}

func (logic *Logic) validateSignalMarket(signal Signal) bool {
	if !utils.StringInSlice(signal.Market, marketNames) {
		log.Printf("Unknown market \"%s\"! Supported: %s.\n", signal.Market, strings.Join(marketNames, ", "))
		return false
	}
	return true
}

func (logic *Logic) validateSignalPair(signal Signal) bool {
	if len(signal.Pair) < 5 || !strings.Contains(signal.Pair, "-") {
		log.Printf("Pair \"%s\" have wrong format! Example: BTC-USDT.\n", signal.Pair)
		return false
	}
	return true
}

func (logic *Logic) validateSignalId(signal Signal) bool {
	if signal.Id == "" {
		for {
			signal.Id = uuid.Must(uuid.NewV4()).String()
			if logic.validateSignalIdUnique(signal.Id) {
				break
			}
		}
		log.Printf("Generated new id \"%s\" for trade strategy \"%s\" \"%s\".", signal.Id, signal.Pair, signal.Market)

	}
	return true
}

func (logic *Logic) validateSignalIdUnique(id string) bool {
	// TODO
	return true
}

func (logic *Logic) validateSignalTradeCap(signal Signal) bool {
	if signal.TradeCap <= 0 {
		log.Printf("Trade cap must be greater than zero. Signal \"%s\"!\n", signal.Id)
		return false
	}
	return true
}

func (logic *Logic) validateSignalBuyLevels(signal Signal) bool {
	var percentSum float32 = 0
	var maxBuyPrice float32 = 0

	if signal.BuyLevels != nil && len(*signal.BuyLevels) > 0 {
		for _, level := range *signal.BuyLevels {
			percentSum += level.Percent
			if level.Price > maxBuyPrice {
				maxBuyPrice = level.Price
			}
		}
	} else {
		log.Printf("Must be at least one buy level for signal \"%s\"\n", signal.Id)
		return false
	}

	return true
}

func (logic *Logic) validateSignalSellLevels(signal Signal) bool {
	var percentSum float32 = 0
	var minSellPrice float32 = math.MaxFloat32

	if signal.SellLevels != nil && len(*signal.SellLevels) > 0 {
		percentSum = 0
		for _, level := range *signal.SellLevels {
			percentSum += level.Percent
			if level.Price < minSellPrice {
				minSellPrice = level.Price
			}
		}
	} else {
		log.Printf("Must be at least one sell level for signal \"%s\"\n", signal.Id)
		return false
	}

	return true
}
