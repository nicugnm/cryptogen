package clients

import (
	"cryptogen-retrieve/domain"
	"net/http"
)

type CryptoRequest struct {
	Status interface{}              `json:"status"`
	Data   []*domain.CryptoMetadata `json:"data"`
}

type NonBlocking struct {
	Response *http.Response
	Error    error
}
