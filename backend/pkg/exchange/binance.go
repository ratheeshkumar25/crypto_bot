package exchange

import (
	"context"
	"strconv"

	"github.com/adshao/go-binance/v2"
)

// BinanceExchange implements the Exchange interface for Binance
type BinanceExchange struct {
	client *binance.Client
}

// NewBinanceExchange creates a new Binance exchange instance
func NewBinanceExchange(apiKey, secret string) Exchange {
	client := binance.NewClient(apiKey, secret)
	return &BinanceExchange{client: client}
}

// GetPrice retrieves the current price for a symbol
func (b *BinanceExchange) GetPrice(ctx context.Context, symbol string) (float64, error) {
	prices, err := b.client.NewListPricesService().Symbol(symbol).Do(ctx)
	if err != nil {
		return 0, err
	}
	if len(prices) == 0 {
		return 0, nil // or error
	}
	price, err := strconv.ParseFloat(prices[0].Price, 64)
	if err != nil {
		return 0, err
	}
	return price, nil
}

// PlaceOrder places an order on Binance
func (b *BinanceExchange) PlaceOrder(ctx context.Context, symbol string, side string, quantity float64, price float64) error {
	order := b.client.NewCreateOrderService().Symbol(symbol).
		Side(binance.SideType(side)).Type(binance.OrderTypeLimit).
		TimeInForce(binance.TimeInForceTypeGTC).Quantity(strconv.FormatFloat(quantity, 'f', -1, 64)).
		Price(strconv.FormatFloat(price, 'f', -1, 64))
	_, err := order.Do(ctx)
	return err
}

// GetVolume retrieves the trading volume for a symbol over a timeframe
func (b *BinanceExchange) GetVolume(ctx context.Context, symbol string, timeframe string) (float64, error) {
	// Use 24hr ticker statistics for volume
	stats, err := b.client.NewListPriceChangeStatsService().Symbol(symbol).Do(ctx)
	if err != nil {
		return 0, err
	}
	if len(stats) == 0 {
		return 0, nil
	}
	volume, err := strconv.ParseFloat(stats[0].Volume, 64)
	if err != nil {
		return 0, err
	}
	return volume, nil
}

// GetBalance retrieves the balance for an asset
func (b *BinanceExchange) GetBalance(ctx context.Context, asset string) (float64, error) {
	account, err := b.client.NewGetAccountService().Do(ctx)
	if err != nil {
		return 0, err
	}
	for _, balance := range account.Balances {
		if balance.Asset == asset {
			free, err := strconv.ParseFloat(balance.Free, 64)
			if err != nil {
				return 0, err
			}
			return free, nil
		}
	}
	return 0, nil
}
