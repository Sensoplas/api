package trend

import (
	"context"

	"cloud.google.com/go/firestore"
)

type ExposureEntry struct {
	Time     int64   `json:"time"`
	Exposure float32 `json:"exposure"`
}

type Service interface {
	GetTimeSeries(ctx context.Context, userID string, start, end int64) ([]ExposureEntry, error)
}

func NewFirestoreTrendingService(fs *firestore.Client) *FirestoreTrendingService {
	return &FirestoreTrendingService{fs}
}

type FirestoreTrendingService struct {
	fs *firestore.Client
}

func (s *FirestoreTrendingService) GetTimeSeries(ctx context.Context, userID string, start, end int64) ([]ExposureEntry, error) {
	docs, err := s.fs.Collection("exposure").Doc("time-series").Collection(userID).Where("time", ">=", start).Where("time", "<=", end).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	timeSeries := make([]ExposureEntry, len(docs))

	for i, doc := range docs {
		timeSeries[i] = ExposureEntry{doc.Data()["time"].(int64), float32(doc.Data()["exposure"].(float64))}
	}

	return timeSeries, nil
}
