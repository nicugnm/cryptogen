package clients

import (
	"sync"
)

type CryptoRequests interface {
	RequestCryptoTypes(nb chan CryptoTypeChannel, wg *sync.WaitGroup)
	RequestCryptoDetails(nbType chan CryptoTypeChannel, nbMetadata chan CryptoMetadataChannel, wg *sync.WaitGroup)

	SaveDataToRepository(nb chan CryptoMetadataChannel, wg *sync.WaitGroup)
	SaveTypeToFile(nb chan CryptoTypeChannel, wg *sync.WaitGroup)
	SaveMetadataToFile(nb chan CryptoMetadataChannel, wg *sync.WaitGroup)
}
