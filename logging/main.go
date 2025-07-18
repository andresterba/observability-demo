package main

import (
	"fmt"
	"time"

	"go.uber.org/zap"
)

const url = "http://my.site.example.com"

func main() {
	logger, _ := zap.NewProduction()
	defer func() {
		err := logger.Sync()
		if err != nil {
			fmt.Printf("Error syncing logger: %v\n", err)
		}
	}()
	sugar := logger.Sugar()

	fmt.Println("Structured logging:")
	sugar.Infow("failed to fetch URL",
		// Structured context as loosely typed key-value pairs.
		"url", url,
		"attempt", 3,
		"backoff", time.Second,
	)

	fmt.Println("test")

	sugar.Infof("Failed to fetch URL: %s", url)

	fmt.Println("\n\nPlain logging:")

	fmt.Printf("[%s] failed to fetch URL: %s\n", time.Now(), url)
}
