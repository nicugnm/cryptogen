package main

import (
	"context"
	"cryptogen-retrieve/service"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	importService := service.NewService()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// create a signal channel to listen for SIGINT and SIGTERM signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-ticker.C:
			// create a context with a 30-second timeout
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			// defer the cancel function to ensure it gets called
			defer cancel()

			// call your function using the context
			err := executeFunction(ctx, importService)
			if err != nil {
				fmt.Println("Error executing function:", err)
			}
		case sig := <-sigCh:
			// received a signal, stop the ticker and quit the program
			fmt.Println("Received signal:", sig)
			ticker.Stop()
			os.Exit(0)
		}
	}
}

func executeFunction(ctx context.Context, importService *service.CryptoService) error {
	// simulate a long-running operation
	time.Sleep(20 * time.Second)

	// check if the context has been cancelled
	select {
	case <-ctx.Done():
		// return an error if the context was cancelled
		return ctx.Err()
	default:
		// continue with the function if the context is still valid
		importService.StartImportService()
	}

	fmt.Println("Function execution complete.")
	return nil
}
