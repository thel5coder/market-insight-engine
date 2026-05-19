package interfaces

import (
	"context"
	"market-insight-engine/internal/domain/model"
)

type CryptoClient interface {
	FetchPrice(ctx context.Context, symbol string) (*model.MarketData, error)
}

type NewsClient interface {
	FetchLatestNews(ctx context.Context, category string) ([]*model.News, error)
	FetchEconomicCalendar(ctx context.Context, from, to string) ([]*model.EconomicEvent, error)
}

type ForexClient interface {
	FetchForexPrice(ctx context.Context, symbol string) (*model.MarketData, error)
}
