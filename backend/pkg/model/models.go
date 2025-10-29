package model

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID               int       `json:"id" db:"id"`
	Username         string    `json:"username" db:"username"`
	Email            string    `json:"email" db:"email"`
	PasswordHash     string    `json:"-" db:"password_hash"`
	BinanceAPIKey    string    `json:"-" db:"binance_api_key"`
	BinanceSecretKey string    `json:"-" db:"binance_secret_key"`
	SolanaPrivateKey string    `json:"-" db:"solana_private_key"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// DBTrade represents a trade executed by the bot (for database storage)
type DBTrade struct {
	ID         int        `json:"id" db:"id"`
	UserID     int        `json:"user_id" db:"user_id"`
	Symbol     string     `json:"symbol" db:"symbol"`
	Side       string     `json:"side" db:"side"` // BUY or SELL
	Quantity   float64    `json:"quantity" db:"quantity"`
	Price      float64    `json:"price" db:"price"`
	Strategy   string     `json:"strategy" db:"strategy"`
	ProfitLoss float64    `json:"profit_loss" db:"profit_loss"`
	TakeProfit float64    `json:"take_profit" db:"take_profit"`
	StopLoss   float64    `json:"stop_loss" db:"stop_loss"`
	Status     string     `json:"status" db:"status"` // OPEN, CLOSED, CANCELLED
	ExecutedAt time.Time  `json:"executed_at" db:"executed_at"`
	ClosedAt   *time.Time `json:"closed_at" db:"closed_at"`
}

// Signal represents a trading signal
type Signal struct {
	ID         int       `json:"id" db:"id"`
	Symbol     string    `json:"symbol" db:"symbol"`
	Strategy   string    `json:"strategy" db:"strategy"`
	Type       string    `json:"type" db:"type"` // BUY or SELL
	Price      float64   `json:"price" db:"price"`
	TakeProfit float64   `json:"take_profit" db:"take_profit"`
	StopLoss   float64   `json:"stop_loss" db:"stop_loss"`
	Confidence float64   `json:"confidence" db:"confidence"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}
