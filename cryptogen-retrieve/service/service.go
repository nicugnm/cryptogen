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

	log.Printf("Started import data service")

	s.ImportData()
}

func (s *CryptoService) ImportData() {
	start := time.Now()

	nb := make(chan clients.NonBlocking, clients.NUM_REQUESTS)
	wg := &sync.WaitGroup{}

	requests := clients.ClientsRequests()

	for i := 0; i < clients.NUM_REQUESTS; i++ {
		wg.Add(1)

		go requests.RequestCryptoTypes(nb)
	}

	go requests.SaveDataToFile(nb, wg)
	go requests.SaveDataToRepository(nb, wg)

	wg.Wait()

	elapsed := time.Since(start)

	log.Infof("Duration of request: %s", elapsed)
}
