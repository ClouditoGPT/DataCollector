package collector

import (
	"DataCollector/internal/models"
	"context"
)

type Collector interface {
	Name() string

	Collect(
		ctx context.Context,
	) (
		<-chan models.Document,
		error,
	)
}
