package repository

import (
	"database/sql"

	"github.com/ratheeshkumar25/forex_bot/backend/pkg/model"
)

// TradeRepository handles database operations for trades
type TradeRepository struct {
	db *sql.DB
}

// NewTradeRepository creates a new trade repository
func NewTradeRepository(db *sql.DB) *TradeRepository {
	return &TradeRepository{db: db}
}

// CreateTrade creates a new trade
func (r *TradeRepository) CreateTrade(trade *model.DBTrade) error {
	query := `INSERT INTO trades (user_id, symbol, side, quantity, price, strategy, profit_loss, take_profit, stop_loss, status, executed_at, closed_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id`
	return r.db.QueryRow(query, trade.UserID, trade.Symbol, trade.Side, trade.Quantity, trade.Price, trade.Strategy, trade.ProfitLoss, trade.TakeProfit, trade.StopLoss, trade.Status, trade.ExecutedAt, trade.ClosedAt).Scan(&trade.ID)
}

// GetTradeByID retrieves a trade by ID
func (r *TradeRepository) GetTradeByID(id int) (*model.DBTrade, error) {
	trade := &model.DBTrade{}
	query := `SELECT id, user_id, symbol, side, quantity, price, strategy, profit_loss, take_profit, stop_loss, status, executed_at, closed_at
	          FROM trades WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&trade.ID, &trade.UserID, &trade.Symbol, &trade.Side, &trade.Quantity, &trade.Price, &trade.Strategy, &trade.ProfitLoss, &trade.TakeProfit, &trade.StopLoss, &trade.Status, &trade.ExecutedAt, &trade.ClosedAt)
	if err != nil {
		return nil, err
	}
	return trade, nil
}

// GetTradesByUserID retrieves all trades for a user
func (r *TradeRepository) GetTradesByUserID(userID int) ([]*model.DBTrade, error) {
	query := `SELECT id, user_id, symbol, side, quantity, price, strategy, profit_loss, take_profit, stop_loss, status, executed_at, closed_at
	          FROM trades WHERE user_id = $1 ORDER BY executed_at DESC`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trades []*model.DBTrade
	for rows.Next() {
		trade := &model.DBTrade{}
		err := rows.Scan(&trade.ID, &trade.UserID, &trade.Symbol, &trade.Side, &trade.Quantity, &trade.Price, &trade.Strategy, &trade.ProfitLoss, &trade.TakeProfit, &trade.StopLoss, &trade.Status, &trade.ExecutedAt, &trade.ClosedAt)
		if err != nil {
			return nil, err
		}
		trades = append(trades, trade)
	}
	return trades, nil
}

// GetOpenTradesByUserID retrieves open trades for a user
func (r *TradeRepository) GetOpenTradesByUserID(userID int) ([]*model.DBTrade, error) {
	query := `SELECT id, user_id, symbol, side, quantity, price, strategy, profit_loss, take_profit, stop_loss, status, executed_at, closed_at
	          FROM trades WHERE user_id = $1 AND status = 'OPEN' ORDER BY executed_at DESC`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trades []*model.DBTrade
	for rows.Next() {
		trade := &model.DBTrade{}
		err := rows.Scan(&trade.ID, &trade.UserID, &trade.Symbol, &trade.Side, &trade.Quantity, &trade.Price, &trade.Strategy, &trade.ProfitLoss, &trade.TakeProfit, &trade.StopLoss, &trade.Status, &trade.ExecutedAt, &trade.ClosedAt)
		if err != nil {
			return nil, err
		}
		trades = append(trades, trade)
	}
	return trades, nil
}

// UpdateTrade updates a trade
func (r *TradeRepository) UpdateTrade(trade *model.DBTrade) error {
	query := `UPDATE trades SET profit_loss = $1, status = $2, closed_at = $3 WHERE id = $4`
	_, err := r.db.Exec(query, trade.ProfitLoss, trade.Status, trade.ClosedAt, trade.ID)
	return err
}

// DeleteTrade deletes a trade
func (r *TradeRepository) DeleteTrade(id int) error {
	query := `DELETE FROM trades WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
