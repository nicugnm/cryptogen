package domain

type CryptoMetadataStorage interface {
	SaveDataList(metadata []*CryptoMetadata) error
}
