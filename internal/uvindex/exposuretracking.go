package uvindex

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/Sensoplas/api/internal/cont"
	"go.uber.org/zap"
)

// might abstract data access away one day :)

func WithExposureTracking(s Service, firestore *firestore.Client, l *zap.Logger) *ExposureTrackingService {
	return &ExposureTrackingService{
		s, firestore, l.With(zap.String("co", "exposure-tracker-middleware")),
	}
}

var _ Service = &ExposureTrackingService{}

type ExposureTrackingService struct {
	service   Service
	firestore *firestore.Client
	logger    *zap.Logger
}

func (s *ExposureTrackingService) Compute(ctx context.Context, data string, lat, long float32) (float32, error) {
	pred, err := s.service.Compute(ctx, data, lat, long)

	if err == nil {
		token, tokenErr := cont.GetIDTokenClaims(ctx)
		if err != nil {
			s.logger.Error("unable to get userID to write exposure tracking data", zap.Error(tokenErr))
		}

		if tokenErr == nil {
			_, res, writeErr := s.firestore.Collection("exposure").Doc("time-series").Collection(token.Subject).Add(ctx, map[string]interface{}{
				"time":     time.Now().Unix(),
				"exposure": pred,
			})

			s.logger.Info("wrote to firestore", zap.String("time", res.UpdateTime.String()))

			if writeErr != nil {
				s.logger.Error("unable to write data to firestore")
			}
		}
	}

	return pred, err
}
