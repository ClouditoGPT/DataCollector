package main

import (
	"DataCollector/internal/crawler"
	"DataCollector/internal/sources/htmlcrawler"
	"DataCollector/internal/sources/wikipedia"
	"encoding/json"
	"os"
	"time"
)

type SourceConfig struct {
	Name        string   `json:"name"`
	URL         string   `json:"url"`
	Seeds       []string `json:"seeds"`
	Language    string   `json:"language"`
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

func NewSourceFetcher(name, baseURL string) crawler.SourceFetcher {
	switch name {
	case "wikipedia":
		return wikipedia.NewClient(baseURL)
	case "html":
		return htmlcrawler.NewClient(baseURL)
	}
	return nil
}

func NewSourceCollector(f crawler.SourceFetcher, cfg SourceConfig) *crawler.Collector {
	return crawler.NewCollector(
		f,
		crawler.WithSeeds(cfg.Seeds),
		crawler.WithWorkers(cfg.Workers),
		crawler.WithRateDelay(time.Duration(cfg.RateDelayMs)*time.Millisecond),
		crawler.WithQueuePath("./data/"+cfg.Name+"_queue.json"),
		crawler.WithVisitedPath("./data/"+cfg.Name+"_visited.json"),
	)
}