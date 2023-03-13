package repositories

import (
	"cryptogen-retrieve/domain"
)

type CryptoMetadataStorage interface {
	SaveDataList(metadata []domain.CryptoDataMetadata) error
}
