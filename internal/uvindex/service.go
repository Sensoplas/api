package uvindex

import (
	"context"
	"math/rand"
)

type Service interface {
	Compute(context.Context, SensorData) (float32, error)
}

type RNGService struct{}

func (s *RNGService) Compute(ctx context.Context, data SensorData) (float32, error) {
	// this is not meant to be cryptographically secure
	//nolint: gosec
	return rand.Float32() * 25, nil
}
