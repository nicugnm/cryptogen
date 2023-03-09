package main

import (
	"cryptogen-retrieve/service"
	"time"
)

func main() {
	importService := service.NewService()

	// every 2 seconds make a request in order to get data
	ticker := time.NewTicker(2 * time.Second)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				importService.StartImportService()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	<-quit
}
