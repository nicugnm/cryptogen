package clients

import (
	"net/http"
	"time"
)

type CryptoMetadata struct {
	Id                 int64     `json:"id"`
	Name               string    `json:"name"`
	Symbol             string    `json:"symbol"`
	Rank               int64     `json:"rank"`
	LastHistoricalData time.Time `json:"last_historical_data"`
}

type NonBlocking struct {
	Response *http.Response
	Error    error
}
