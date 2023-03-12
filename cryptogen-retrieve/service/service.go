package service

import (
	"cryptogen-retrieve/gateways/clients"
	"fmt"
	runtime "github.com/banzaicloud/logrus-runtime-formatter"
	log "github.com/sirupsen/logrus"
	"os"
	"sync"
	"time"
)

type CryptoService struct {
}

func NewService() *CryptoService {
	return &CryptoService{}
}

var _ CryptoMetadataService = (*CryptoService)(nil)

func (s *CryptoService) StartImportService() {
	formatter := runtime.Formatter{ChildFormatter: &log.TextFormatter{
		FullTimestamp:          true,
		DisableLevelTruncation: true,
	}}
	formatter.Line = true
	log.SetFormatter(&formatter)
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	log.Printf("Started import data service")

	s.ImportData()
}

func (s *CryptoService) ImportData() {
	start := time.Now()

	nbType := make(chan clients.CryptoTypeChannel, clients.NUM_REQUESTS)
	nbMetadata := make(chan clients.CryptoMetadataChannel, clients.NUM_REQUESTS)

	wg := &sync.WaitGroup{}

	requests := clients.ClientsRequests()

	for i := 0; i < clients.NUM_REQUESTS; i++ {
		wg.Add(1)

		go requests.RequestCryptoTypes(nbType, wg)
	}

	// we need to wait all to all goroutines to finish in order to have the correct data
	// the actual data are the symbols from the cryptocurrency and then use them for the second request to obtain all the data of the cryptocurrencies
	wg.Wait()

	for i := 0; i < clients.NUM_REQUESTS; i++ {
		wg.Add(1)

		go requests.RequestCryptoDetails(nbType, nbMetadata, wg)
	}

	wg.Wait()

	fmt.Println("Start save to file")

	wg.Add(2)

	go requests.SaveTypeToFile(nbType, wg)
	go requests.SaveMetadataToFile(nbMetadata, wg)

	go requests.SaveDataToRepository(nbMetadata, wg)

	wg.Wait()

	elapsed := time.Since(start)

	log.Infof("Duration of request: %s", elapsed)
}
