package trend

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	chart "github.com/wcharczuk/go-chart/v2"
)

type ExposureEntry struct {
	Time     int64   `json:"time"`
	Exposure float32 `json:"exposure"`
}

type Service interface {
	GetTimeSeries(ctx context.Context, userID string, start, end int64, width, height int) (*chart.Chart, error)
}

func NewFirestoreTrendingService(fs *firestore.Client) *FirestoreTrendingService {
	return &FirestoreTrendingService{fs}
}

type FirestoreTrendingService struct {
	fs *firestore.Client
}

func (s *FirestoreTrendingService) GetTimeSeries(ctx context.Context, userID string, start, end int64, width, height int) (*chart.Chart, error) {
	docs, err := s.fs.Collection("exposure").Doc("time-series").Collection(userID).Where("time", ">=", start).Where("time", "<=", end).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	xVals := make([]time.Time, len(docs))
	yVals := make([]float64, len(docs))
	for i, doc := range docs {
		xVals[i] = time.Unix(doc.Data()["time"].(int64), 0)
		yVals[i] = doc.Data()["exposure"].(float64)
	}

	graph := &chart.Chart{
		Width:  width,
		Height: height,
		Series: []chart.Series{
			chart.TimeSeries{
				Name:  "UV Exposure",
				YAxis: chart.YAxisPrimary,
				Style: chart.Style{
					StrokeColor: chart.GetDefaultColor(0).WithAlpha(64),
					StrokeWidth: 1,
					FillColor:   chart.GetDefaultColor(0).WithAlpha(64),
				},
				XValues: xVals,
				YValues: yVals,
			},
		},
	}

	return graph, nil
}
