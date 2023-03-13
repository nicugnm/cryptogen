package service

import (
	"cryptogen-retrieve/gateways/clients"
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

	log.Printf("Start importing data..")

	s.ImportData()
}

func (s *CryptoService) ImportData() {
	start := time.Now()

	nbType := make(chan clients.CryptoTypeChannel, 1)
	nbMetadata := make(chan clients.CryptoMetadataChannel, 1)

	wg := &sync.WaitGroup{}

	requests := clients.ClientsRequests()

	wg.Add(1)

	go requests.RequestCryptoTypes(nbType, wg)

	// we need to wait all goroutines to finish in order to have the correct data
	// the actual data are the symbols from the cryptocurrency and then use them for the second request to obtain all the data of the cryptocurrencies
	wg.Wait()

	wg.Add(1)
	go requests.RequestCryptoDetails(nbType, nbMetadata, wg)

	// we need to wait all goroutines to finish in order to have data published to the channels
	// the actual data will be saved in files and in redis
	wg.Wait()

	//fmt.Println("Start save to file")

	wg.Add(1)

	//go requests.SaveTypeToFile(nbType, wg)
	//go requests.SaveMetadataToFile(nbMetadata, wg)

	go requests.SaveDataToRepository(nbMetadata, wg)

	wg.Wait()

	elapsed := time.Since(start)

	log.Infof("Duration of request: %s", elapsed)
}
