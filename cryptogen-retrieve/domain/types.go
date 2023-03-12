package domain

import (
	"time"
)

type CryptoTypeMetadata struct {
	Id                 int64     `json:"id"`
	Name               string    `json:"name"`
	Symbol             string    `json:"symbol"`
	Rank               int64     `json:"rank"`
	LastHistoricalData time.Time `json:"last_historical_data"`
}

type CryptoDataMetadata struct {
	TotalSuply int64 `json:"total_supply"`
	MaxSuply   int64 `json:"max_supply"`
	Quote      Quote `json:"quote"`
}

type Quote struct {
	USD MoneyInformation `json:"USD"`
}

type MoneyInformation struct {
	Price                    float64   `json:"price"`
	Volume                   float64   `json:"volume_24h"`
	VolumeChange24h          float64   `json:"volume_change_24h"`
	PercentChange1h          float64   `json:"percent_change_1h"`
	PercentChange24h         float64   `json:"percent_change_24h"`
	PercentChange7d          float64   `json:"percent_change_7d"`
	PercentChange30d         float64   `json:"percent_change_30d"`
	MarketingCap             float64   `json:"marketing_cap"`
	MarketingCapDominance    float64   `json:"marketing_cap_dominance"`
	FullyDilutedMarketingCap float64   `json:"fully_diluted_market_cap"`
	LastUpdated              time.Time `json:"last_updated"`
}
