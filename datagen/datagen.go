package datagen

import (
	"context"
)

// DataGenerator represents a random data generator for API testing purposes
type DataGenerator interface {
	// Generate generates the data
	Generate(
		ctx context.Context,
		opts Options,
		stats StatisticsWriter,
	) (Graph, error)
}
