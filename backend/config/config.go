package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	AlphaVantageAPIKey string
	BinanceAPIKey      string
	BinanceSecret      string
	Port               string
	JWTSecret          string
	// Database configuration
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	// Add more exchange keys as needed, e.g., CoinbaseAPIKey, etc.
}

// NewConfig creates a new Config struct from environment variables.
// This function will be used by fx as a dependency provider.
func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, will rely on environment variables")
	}

	apiKey := os.Getenv("ALPHA_VANTAGE_API_KEY")
	if apiKey == "" {
		log.Println("WARNING: ALPHA_VANTAGE_API_KEY is not set. The application might not work correctly.")
	}

	binanceAPIKey := os.Getenv("BINANCE_API_KEY")
	binanceSecret := os.Getenv("BINANCE_SECRET")
	if binanceAPIKey == "" || binanceSecret == "" {
		// Use sample test keys for demonstration (replace with real keys for production)
		log.Println("WARNING: Using sample Binance test keys. Replace with real API keys for production use.")
		binanceAPIKey = "test_api_key"    // Sample key - replace with real
		binanceSecret = "test_secret_key" // Sample secret - replace with real
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key"
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "password"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "forexbot"
	}

	return &Config{
		AlphaVantageAPIKey: apiKey,
		BinanceAPIKey:      binanceAPIKey,
		BinanceSecret:      binanceSecret,
		Port:               port,
		JWTSecret:          jwtSecret,
		DBHost:             dbHost,
		DBPort:             dbPort,
		DBUser:             dbUser,
		DBPassword:         dbPassword,
		DBName:             dbName,
	}, nil
}
