package uvindex

import (
	"context"
	"math/rand"

	"github.com/briandowns/openweathermap"
)

type Service interface {
	Compute(context.Context, string, float32, float32) (float32, error)
}

type RNGService struct{}

func (s *RNGService) Compute(ctx context.Context, data SensorData, lat, long float32) (float32, error) {
	// this is not meant to be cryptographically secure
	//nolint: gosec
	return rand.Float32() * 25, nil
}

type LocationService struct {
	WeatherAPI *openweathermap.UV
}

func (s *LocationService) Compute(ctx context.Context, data string, lat, long float32) (float32, error) {
	err := s.WeatherAPI.Current(&openweathermap.Coordinates{Latitude: float64(lat), Longitude: float64(long)})
	if err != nil {
		return 0, err
	}

	return float32(s.WeatherAPI.Value), nil
}
