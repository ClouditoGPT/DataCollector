package main

import (
	"DataCollector/internal/crawler"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type SourceConfig struct {
	Name        string   `json:"name"`
	URL         string   `json:"url"`
	Seeds       []string `json:"seeds"`
	Workers     int      `json:"workers"`
	RateDelayMs int      `json:"rate_delay_ms"`
}

type Config struct {
	Sources []SourceConfig `json:"sources"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func runSource(ctx context.Context, src SourceConfig) {
	f := crawler.NewFetcher(src.Name, src.URL)
	if f == nil {
		fmt.Printf("Unknown source type: %s\n", src.Name)
		return
	}

	source := crawler.NewCollector(
		f,
		crawler.WithSeeds(src.Seeds),
		crawler.WithWorkers(src.Workers),
		crawler.WithRateDelay(time.Duration(src.RateDelayMs)*time.Millisecond),
	)

	fmt.Printf("Starting crawler: %s (workers=%d, seeds=%d)\n", src.Name, src.Workers, len(src.Seeds))

	docs, err := source.Collect(ctx)
	if err != nil {
		fmt.Printf("Crawler error: %v\n", err)
		return
	}

	for doc := range docs {
		fmt.Printf("Crawled: %s (lang=%s)\n", doc.Title, doc.Language)
	}

	state := source.GetState()
	fmt.Printf("Source %s done - visited: %d, errors: %d\n", src.Name, state.GetVisited(), state.GetErrors())
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nShutting down...")
		cancel()
	}()

	cfg, err := LoadConfig("./configs/sources.json")
	if err != nil {
		fmt.Printf("Config error: %v, using defaults\n", err)
		cfg = &Config{
			Sources: []SourceConfig{
				{Name: "html", URL: "https://en.wikipedia.org", Seeds: []string{"https://en.wikipedia.org/wiki/Iran"}, Workers: 5, RateDelayMs: 500},
			},
		}
	}

	for _, src := range cfg.Sources {
		runSource(ctx, src)
	}

	fmt.Println("All sources completed!")
}