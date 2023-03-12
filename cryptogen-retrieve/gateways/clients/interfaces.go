package clients

import (
	"sync"
)

type CryptoRequests interface {
	RequestCryptoTypes(nb chan NonBlocking)
	RequestCryptoDetails(nb chan NonBlocking) []*interface{}

	SaveDataToRepository(nb chan NonBlocking, wg *sync.WaitGroup)
	SaveDataToFile(nb chan NonBlocking, wg *sync.WaitGroup)
}
