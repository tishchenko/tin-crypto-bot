package config

type Logic struct {
	Strategies Strategies `json:"strategies"`
}

type Strategies struct {
	Signals *[]Signal `json:"signals, omitempty"`
}

//var strategies = []string{"signals"}

type Strategy struct {
	Market string  `json:"market"`
	Pair   string  `json:"pair"`
	Count  float32 `json:"count"`
	Active *bool   `json:"active, omitempty"`
}

type Signal struct {
	Strategy
	BuyLevels  *[]Level `json:"buyLevels"`
	SellLevels *[]Level `json:"sellLevels"`
	StopLoss   *Level   `json:"sellLevels, omitempty"`
}

type Level struct {
	Price    float32 `json:"price"`
	Percent  float32 `json:"percent"`
	StopLoss *Level  `json:"sellLevels, omitempty"`
}
