package domain

type CryptoMetadataStorage interface {
	SaveData(metadata *CryptoMetadata) error
	SaveDataList(metadata []*CryptoMetadata) error
}
