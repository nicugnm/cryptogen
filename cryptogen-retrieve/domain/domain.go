package domain

import "context"

type CryptoMetadataStorage interface {
	SaveFile(ctx context.Context, metadata *CryptoMetadata) error
	RetrieveFile(ctx context.Context, sha256 string) (*CryptoMetadata, error)
}
