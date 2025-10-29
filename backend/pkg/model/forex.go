package model

type ForexData struct {
	Timestamp string
	Price     float64
}

type Prediction struct {
	Pair       string  `json:"pair"`
	Signal     string  `json:"signal"` // "buy", "sell", "hold"
	Confidence float64 `json:"confidence"`
	Price      float64 `json:"price"`
	Reason     string  `json:"reason,omitempty"`
}
