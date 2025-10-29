package api

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/ratheeshkumar25/forex_bot/backend/pkg/middleware"
	"github.com/ratheeshkumar25/forex_bot/backend/pkg/model"
	"github.com/ratheeshkumar25/forex_bot/backend/pkg/repository"
	"golang.org/x/crypto/bcrypt"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	userRepo  *repository.UserRepository
	jwtSecret string
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(userRepo *repository.UserRepository, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

// RegisterRequest represents the registration request
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest represents the login request
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Register handles user registration
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validate input
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Username, email, and password are required"})
	}

	    // Check if user already exists
		existingUser, err := h.userRepo.GetUserByUsername(req.Username)
		if err == nil && existingUser != nil {
			return c.Status(409).JSON(fiber.Map{"error": "Username already exists"})
		}
	
	    existingUser, err = h.userRepo.GetUserByEmail(req.Email)
	    if err == nil && existingUser != nil {
	        return c.Status(409).JSON(fiber.Map{"error": "Email already exists"})
	    }
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	// Create user
	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := h.userRepo.CreateUser(user); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create user in database"})
	}

	// Generate JWT token
	token, err := h.generateToken(user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return c.JSON(fiber.Map{
		"message": "User registered successfully",
		"token":   token,
		"user": fiber.Map{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// Login handles user login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	    // Get user by username
		user, err := h.userRepo.GetUserByUsername(req.Username)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "User not found"})
		}
	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Generate JWT token
	token, err := h.generateToken(user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return c.JSON(fiber.Map{
		"message": "Login successful",
		"token":   token,
		"user": fiber.Map{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// GetProfile handles getting user profile
func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	user, err := h.userRepo.GetUserByID(userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(fiber.Map{
		"user": fiber.Map{
			"id":               user.ID,
			"username":         user.Username,
			"email":            user.Email,
			"has_binance_keys": user.BinanceAPIKey != "" && user.BinanceSecretKey != "",
			"has_solana_key":   user.SolanaPrivateKey != "",
			"created_at":       user.CreatedAt,
			"updated_at":       user.UpdatedAt,
		},
	})
}

// UpdateExchangeKeys handles updating exchange API keys
func (h *AuthHandler) UpdateExchangeKeys(c *fiber.Ctx) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req struct {
		BinanceAPIKey    string `json:"binance_api_key,omitempty"`
		BinanceSecretKey string `json:"binance_secret_key,omitempty"`
		SolanaPrivateKey string `json:"solana_private_key,omitempty"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	user, err := h.userRepo.GetUserByID(userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	// Update keys if provided
	if req.BinanceAPIKey != "" {
		user.BinanceAPIKey = req.BinanceAPIKey
	}
	if req.BinanceSecretKey != "" {
		user.BinanceSecretKey = req.BinanceSecretKey
	}
	if req.SolanaPrivateKey != "" {
		user.SolanaPrivateKey = req.SolanaPrivateKey
	}

	user.UpdatedAt = time.Now()

	if err := h.userRepo.UpdateUser(user); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update exchange keys"})
	}

	return c.JSON(fiber.Map{"message": "Exchange keys updated successfully"})
}

// generateToken generates a JWT token for the user
func (h *AuthHandler) generateToken(user *model.User) (string, error) {
	claims := middleware.JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.jwtSecret))
}
