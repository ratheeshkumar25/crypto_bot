package repository

import (
	"database/sql"

	"github.com/ratheeshkumar25/forex_bot/backend/pkg/model"
)

// SignalRepository handles database operations for signals
type SignalRepository struct {
	db *sql.DB
}

// NewSignalRepository creates a new signal repository
func NewSignalRepository(db *sql.DB) *SignalRepository {
	return &SignalRepository{db: db}
}

// CreateSignal creates a new signal
func (r *SignalRepository) CreateSignal(signal *model.Signal) error {
	query := `INSERT INTO signals (symbol, strategy, signal_type, price, take_profit, stop_loss, confidence, created_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	return r.db.QueryRow(query, signal.Symbol, signal.Strategy, signal.Type, signal.Price, signal.TakeProfit, signal.StopLoss, signal.Confidence, signal.CreatedAt).Scan(&signal.ID)
}

// GetSignalsBySymbol retrieves signals for a symbol
func (r *SignalRepository) GetSignalsBySymbol(symbol string) ([]*model.Signal, error) {
	query := `SELECT id, symbol, strategy, signal_type, price, take_profit, stop_loss, confidence, created_at
	          FROM signals WHERE symbol = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(query, symbol)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var signals []*model.Signal
	for rows.Next() {
		signal := &model.Signal{}
		err := rows.Scan(&signal.ID, &signal.Symbol, &signal.Strategy, &signal.Type, &signal.Price, &signal.TakeProfit, &signal.StopLoss, &signal.Confidence, &signal.CreatedAt)
		if err != nil {
			return nil, err
		}
		signals = append(signals, signal)
	}
	return signals, nil
}

// DeleteSignal deletes a signal
func (r *SignalRepository) DeleteSignal(id int) error {
	query := `DELETE FROM signals WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
