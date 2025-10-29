package model

import "time"

type Trade struct {
	Symbol      string
	EntryPrice  float64
	TakeProfit  float64
	StopLoss    float64
	OpenTime    time.Time
	MaxDuration time.Duration
	IsOpen      bool
	Size        float64 // Amount of asset to buy/sell
	RiskPercent float64 // % of capital risked on this trade
}
