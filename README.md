# CryptoHack Bot - Automated Forex/Crypto Trading System

## Project Overview

CryptoHack Bot is a comprehensive automated trading system designed for cryptocurrency and forex markets. The system provides real-time price monitoring, automated trading strategies, profit prediction, and user-friendly dashboard for managing trades and strategies.

### Key Features

- **Real-time Price Monitoring**: Live price charts with WebSocket streaming
- **Automated Trading Strategies**:
  - Grid Trading: Places buy orders below current price and sell orders above
  - Dollar-Cost Averaging (DCA): Systematic buying at regular intervals
- **Multi-Exchange Support**: Binance and Solana blockchain integration
- **User Authentication**: JWT-based secure authentication system
- **Trade Management**: Track positions, profit/loss, take profit, and stop loss
- **Signal Generation**: AI-powered trading signals with confidence scores
- **WebSocket Streaming**: Real-time price updates and notifications
- **Docker Containerization**: Easy deployment with docker-compose

## System Architecture

The system follows a microservices-like architecture with separate backend and frontend components:

### Backend (Go/Fiber)
- **Framework**: Go with Fiber web framework
- **Database**: PostgreSQL with GORM
- **Authentication**: JWT tokens with bcrypt password hashing
- **Dependency Injection**: Uber Fx for clean architecture
- **Real-time Communication**: WebSocket for price streaming
- **Exchange Integrations**: Binance API and Solana Web3.js
- **Trading Strategies**: Modular strategy implementations (Grid, DCA)
- **Worker Service**: Background analysis and automated trading

### Frontend (Node.js/Express)
- **Server**: Express.js proxy server
- **Templates**: EJS templating engine
- **Real-time Charts**: Chart.js for price visualization
- **Authentication**: JWT token management in localStorage
- **API Proxy**: Proxies all backend API calls
- **WebSocket Client**: Real-time price updates

### Database Schema
- **Users**: Authentication and exchange API keys
- **Trades**: Executed trades with profit/loss tracking
- **Signals**: Generated trading signals with confidence scores

### Deployment
- **Containerization**: Docker with multi-stage builds
- **Orchestration**: Docker Compose for local development
- **Networks**: Isolated container networking

## How Backend and Frontend Work Together

### Communication Flow

1. **User Authentication**:
   - Frontend collects login/register data
   - Proxies to backend `/api/auth/login` or `/api/auth/register`
   - Backend validates credentials, generates JWT token
   - Token stored in frontend localStorage for subsequent requests

2. **API Requests**:
   - Frontend makes requests to `/api/*` endpoints
   - Express proxy forwards to backend `http://localhost:3000`
   - Backend processes with JWT middleware for protected routes

3. **Real-time Price Streaming**:
   - Frontend establishes WebSocket connection to `/api/ws/price`
   - Backend streams live price data from exchanges
   - Chart.js updates in real-time with price and signal overlays

4. **Trading Operations**:
   - User selects strategy and parameters via frontend
   - Frontend calls backend prediction endpoints
   - Backend analyzes market data and generates signals
   - Worker service can execute automated trades based on signals

### Data Flow Architecture

```
User Browser (Port 8080)
    ↓ (HTTP/WebSocket)
Express Proxy Server
    ↓ (HTTP/WebSocket proxy)
Go Fiber API Server (Port 3000)
    ↓ (Database queries)
PostgreSQL Database
    ↓ (Exchange API calls)
Binance/Solana APIs
```

## Getting Started

### Prerequisites
- Docker and Docker Compose
- Git

### Environment Variables
Create a `.env` file in the backend directory:
```env
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=forexbot
PORT=3000
JWT_SECRET=your-jwt-secret-here
ALPHA_VANTAGE_API_KEY=your-api-key
BINANCE_API_KEY=your-binance-key
BINANCE_SECRET=your-binance-secret
```

### Installation and Setup

1. **Clone the repository**:
```bash
git clone <repository-url>
cd forex_bot
```

2. **Start the services**:
```bash
cd backend
docker-compose up --build
```

3. **Access the application**:
- Frontend: http://localhost:8080
- Backend API: http://localhost:3000

### Development Setup

For development without Docker:

1. **Backend**:
```bash
cd backend
go mod download
go run cmd/main.go
```

2. **Frontend**:
```bash
cd frontend
npm install
npm start
```

## API Documentation

### Authentication Endpoints

#### POST `/api/auth/register`
Register a new user account.

**Request Body**:
```json
{
  "username": "string",
  "email": "string",
  "password": "string"
}
```

**Response**:
```json
{
  "message": "User registered successfully",
  "token": "jwt-token",
  "user": {
    "id": 1,
    "username": "string",
    "email": "string"
  }
}
```

#### POST `/api/auth/login`
Authenticate user and get JWT token.

**Request Body**:
```json
{
  "username": "string",
  "password": "string"
}
```

#### GET `/api/auth/profile`
Get authenticated user profile (requires JWT).

### Trading Endpoints

#### GET `/api/price/:exchange`
Get current price for a symbol.

**Parameters**:
- `exchange`: binance
- `symbol`: BTCUSDT

#### GET `/api/predict/:strategy`
Predict profit for a trading strategy.

**Parameters**:
- `strategy`: grid | dca
- `symbol`: BTCUSDT
- `investment`: 1000
- `timeframe`: short | long

#### GET `/api/signals/:strategy`
Get trading signals for a strategy.

#### GET `/api/trades`
Get user's trade history (requires JWT).

### WebSocket Endpoints

#### `/api/ws/price`
Real-time price streaming.

**Query Parameters**:
- `symbol`: BTCUSDT

**Message Format**:
```json
{
  "price": 45000.50,
  "timestamp": "2024-01-01T12:00:00Z"
}
```

## Trading Strategies

### Grid Trading
Places multiple buy orders below the current price and sell orders above it. The bot takes profit when price hits sell levels and cuts losses when it drops below buy levels.

### Dollar-Cost Averaging (DCA)
Systematically buys assets at regular intervals or when price drops significantly. Takes profit at predetermined levels above the average purchase price.

## Security Features

- JWT-based authentication
- Password hashing with bcrypt
- CORS protection
- Input validation and sanitization
- Secure API key storage
- Rate limiting (configurable)

## Monitoring and Logging

- Structured logging with Zap
- Request/response logging
- Error tracking and reporting
- Performance metrics (configurable)

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For support and questions:
- Create an issue in the repository
- Check the documentation
- Review the code comments

---

## System Architecture Diagram

Below is a textual representation of the system architecture. For a visual PNG diagram, use tools like draw.io or Lucidchart to create the following structure:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   User Browser  │    │ Express Proxy   │    │   Go Fiber API  │
│   (Port 8080)   │◄──►│   Server        │◄──►│   Server         │
│                 │    │ (Port 8080)    │    │   (Port 3000)    │
│ - EJS Templates │    │                 │    │                 │
│ - Chart.js      │    │ - API Proxy     │    │ - REST API      │
│ - WebSocket     │    │ - Static Files  │    │ - WebSocket     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  PostgreSQL DB  │    │  Worker Service │    │ Exchange APIs   │
│                 │    │                 │    │                 │
│ - Users         │    │ - Analysis      │    │ - Binance       │
│ - Trades        │    │ - Automation    │    │ - Solana        │
│ - Signals       │    │ - Scheduling    │    │ - Price Data    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Component Descriptions:

1. **User Browser**: Client-side interface with real-time charts and trading controls
2. **Express Proxy Server**: Serves frontend assets and proxies API requests to backend
3. **Go Fiber API Server**: Core business logic, authentication, trading strategies
4. **PostgreSQL Database**: Persistent storage for users, trades, and signals
5. **Worker Service**: Background processing for market analysis and automated trading
6. **Exchange APIs**: External cryptocurrency exchanges for price data and trading

### Data Flow:
- Solid lines: HTTP/WebSocket communication
- Dashed lines: Database queries and external API calls
- Arrows indicate request/response direction
