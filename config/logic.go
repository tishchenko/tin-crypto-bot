package config

import (
	"os"
	"log"
	"encoding/json"
	"strings"
	"github.com/tishchenko/tin-crypto-bot/utils"
	"math"
	"crypto/sha1"
	"encoding/hex"
	"strconv"
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

	signals := []Signal{}
	if logic.Strategies.Signals != nil {
		for _, signal := range *logic.Strategies.Signals {
			/*var active bool
			if signal.Active == nil {
				active = true
			} else {
				active = *signal.Active
			}
			active = active && logic.validateSignal(signal)
			signal.Active = &(active)*/

			if signal.validate() {
				signals = append(signals, signal)
			}
		}
		logic.Strategies.Signals = &signals
	}
}

// TODO Refactor this !
func (signal *Signal) validate() bool {
	if !signal.validateMarket() {
		return false
	}

	if !signal.validatePair() {
		return false
	}

	if !signal.validateId() {
		return false
	}

	if !signal.validateTradeCap() {
		return false
	}

	if !signal.validateBuyLevels() {
		return false
	}

	if !signal.validateSellLevels() {
		return false
	}

	return true
}

func (signal *Signal) validateMarket() bool {
	if !utils.StringInSlice(signal.Market, marketNames) {
		log.Printf("Unknown market \"%s\"! Supported: %s.\n", signal.Market, strings.Join(marketNames, ", "))
		return false
	}
	return true
}

func (signal *Signal) validatePair() bool {
	if len(signal.Pair) < 5 || !strings.Contains(signal.Pair, "-") {
		log.Printf("Pair \"%s\" have wrong format! Example: BTC-USDT.\n", signal.Pair)
		return false
	}
	return true
}

func (signal *Signal) validateId() bool {
	if signal.Id == "" {
		for {
			//signal.Id = uuid.Must(uuid.NewV4()).String()
			signal.Id = signal.Hash()
			if signal.validateIdUnique(signal.Id) {
				break
			}
		}
		log.Printf("Generated new id \"%s\" for trade strategy \"%s\" \"%s\".", signal.Id, signal.Pair, signal.Market)

	}
	return true
}

func (signal *Signal) validateIdUnique(id string) bool {
	// TODO
	return true
}

func (signal *Signal) validateTradeCap() bool {
	if signal.TradeCap <= 0 {
		log.Printf("Trade cap must be greater than zero. Signal \"%s\"!\n", signal.Id)
		return false
	}
	return true
}

func (signal *Signal) validateBuyLevels() bool {
	var percentSum float32 = 0
	var maxBuyPrice float32 = 0

	if signal.BuyLevels != nil && len(*signal.BuyLevels) > 0 {
		for _, level := range *signal.BuyLevels {
			percentSum += level.Percent
			if level.Price > maxBuyPrice {
				maxBuyPrice = level.Price
			}
			if percentSum > 100 {
				log.Printf("Buy levels quantity more than 100%% for signal \"%s\"\n", signal.Id)
				return false
			}
		}
	} else {
		log.Printf("Must be at least one buy level for signal \"%s\"\n", signal.Id)
		return false
	}

	return true
}

func (signal *Signal) validateSellLevels() bool {
	var percentSum float32 = 0
	var minSellPrice float32 = math.MaxFloat32

	if signal.SellLevels != nil && len(*signal.SellLevels) > 0 {
		percentSum = 0
		for _, level := range *signal.SellLevels {
			percentSum += level.Percent
			if level.Price < minSellPrice {
				minSellPrice = level.Price
			}
			if percentSum > 100 {
				log.Printf("Sell levels quantity more than 100%% for signal \"%s\"\n", signal.Id)
				return false
			}
		}
	} else {
		log.Printf("Must be at least one sell level for signal \"%s\"\n", signal.Id)
		return false
	}

	return true
}

func (signal *Signal) Hash() string {
	s := signal.Market + signal.Pair

	var level2string = func(level Level) string {
		var level2string = func(level Level) string {
			s := "percent:" + strconv.FormatFloat(float64(level.Percent), 'g', 1, 32) +
				"price:" + strconv.FormatFloat(float64(level.Price), 'g', 1, 32)
			return s
		}
		s := level2string(level)
		if level.StopLoss != nil {
			s += "stoploss:" + level2string(*level.StopLoss)
		}
		return s
	}

	s += "buyLevels:"
	for _, level := range *signal.BuyLevels {
		s += level2string(level)
	}
	s += "sellLevels:"
	for _, level := range *signal.SellLevels {
		s += level2string(level)
	}
	if signal.StopLoss != nil {
		s += "stoploss:" + level2string(*signal.StopLoss)
	}

	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
