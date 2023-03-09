package service

import (
	"cryptogen-retrieve/gateways/clients"
	runtime "github.com/banzaicloud/logrus-runtime-formatter"
	log "github.com/sirupsen/logrus"
	"os"
	"sync"
	"time"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) StartImportService() {
	formatter := runtime.Formatter{ChildFormatter: &log.TextFormatter{
		FullTimestamp:          true,
		DisableLevelTruncation: true,
	}}
	formatter.Line = true
	log.SetFormatter(&formatter)
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	log.Printf("Started import data service")

	ImportData()
}

func ImportData() {
	start := time.Now()

	nb := make(chan clients.NonBlocking, clients.NUM_REQUESTS)
	wg := &sync.WaitGroup{}

	for i := 0; i < clients.NUM_REQUESTS; i++ {
		wg.Add(1)
		go clients.Request(nb)
	}

	go clients.HandleResponse(nb, wg)

	wg.Wait()

	elapsed := time.Since(start)

	log.Infof("Duration of request: %s", elapsed)
}
