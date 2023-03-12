package clients

import (
	"cryptogen-retrieve/domain"
	"net/http"
)

type CryptoTypeRequest struct {
	Status interface{}                  `json:"status"`
	Data   []*domain.CryptoTypeMetadata `json:"data"`
}

type CryptoMetadataRequest struct {
	Status interface{}                          `json:"status"`
	Data   map[string]domain.CryptoDataMetadata `json:"data"`
}

type CryptoTypeChannel struct {
	Response *http.Response
	Error    error
}

type CryptoMetadataChannel struct {
	Response *http.Response
	Error    error
}
