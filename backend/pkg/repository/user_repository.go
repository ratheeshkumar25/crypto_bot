package repository

import (
	"database/sql"

	"github.com/ratheeshkumar25/forex_bot/backend/pkg/model"
)

// UserRepository handles database operations for users
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser creates a new user
func (r *UserRepository) CreateUser(user *model.User) error {
	query := `INSERT INTO users (username, email, password_hash, binance_api_key, binance_secret_key, solana_private_key, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	return r.db.QueryRow(query, user.Username, user.Email, user.PasswordHash, user.BinanceAPIKey, user.BinanceSecretKey, user.SolanaPrivateKey, user.CreatedAt, user.UpdatedAt).Scan(&user.ID)
}

// GetUserByID retrieves a user by ID
func (r *UserRepository) GetUserByID(id int) (*model.User, error) {
	user := &model.User{}
	query := `SELECT id, username, email, password_hash, binance_api_key, binance_secret_key, solana_private_key, created_at, updated_at FROM users WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.BinanceAPIKey, &user.BinanceSecretKey, &user.SolanaPrivateKey, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByUsername retrieves a user by username
func (r *UserRepository) GetUserByUsername(username string) (*model.User, error) {
	user := &model.User{}
	query := `SELECT id, username, email, password_hash, binance_api_key, binance_secret_key, solana_private_key, created_at, updated_at FROM users WHERE username = $1`
	err := r.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.BinanceAPIKey, &user.BinanceSecretKey, &user.SolanaPrivateKey, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	user := &model.User{}
	query := `SELECT id, username, email, password_hash, binance_api_key, binance_secret_key, solana_private_key, created_at, updated_at FROM users WHERE email = $1`
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.BinanceAPIKey, &user.BinanceSecretKey, &user.SolanaPrivateKey, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUser updates a user
func (r *UserRepository) UpdateUser(user *model.User) error {
	query := `UPDATE users SET username = $1, email = $2, password_hash = $3, binance_api_key = $4, binance_secret_key = $5, solana_private_key = $6, updated_at = $7 WHERE id = $8`
	_, err := r.db.Exec(query, user.Username, user.Email, user.PasswordHash, user.BinanceAPIKey, user.BinanceSecretKey, user.SolanaPrivateKey, user.UpdatedAt, user.ID)
	return err
}

// DeleteUser deletes a user
func (r *UserRepository) DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
