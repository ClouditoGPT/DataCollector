package main

import (
	"DataCollector/internal/crawler"
	"DataCollector/internal/sources/wikipedia"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	seeds := []string{"ایران"}
	ctx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nShutting down...")
		cancel()
	}()

	f := wikipedia.NewClient("https://fa.wikipedia.org/w/api.php")
	source := crawler.NewCollector(
		f,
		crawler.WithSeeds(seeds),
	)

	fmt.Printf("Starting crawler: %s\n", f.Name())

	docs, err := source.Collect(ctx)
	if err != nil {
		panic(err)
	}

	for doc := range docs {
		fmt.Printf("Crawled: %s (lang=%s)\n", doc.Title, doc.Language)
	}

	fmt.Println("Done!")
}