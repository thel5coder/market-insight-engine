package interfaces

import (
	"context"
	"market-insight-engine/internal/domain/model"
)

type CryptoClient interface {
	FetchPrice(ctx context.Context, symbol string) (*model.MarketData, error)
}
