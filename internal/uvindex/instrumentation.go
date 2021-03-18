package uvindex

import (
	"context"
	"net/http"
	"time"

	"github.com/Sensoplas/api/internal/cont"
	"github.com/Sensoplas/api/internal/logging"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"go.uber.org/zap"
)

func endpointLoggingMiddleware(logger *zap.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			epLogger := logging.WithLayerName(logger, "endpoint")
			name, err := cont.GetOperationName(ctx)
			if err == nil {
				epLogger = logging.WithOperationName(epLogger, name)
			}
			epLogger.Debug("calling endpoint")

			res, err := next(ctx, request)

			if err != nil {
				epLogger.Error("endpoint returned error", zap.Error(err))
			} else {
				epLogger.Debug("called endpoint")
			}

			return res, err
		}
	}
}

func httpDecoderMiddleware(logger *zap.Logger) func(httptransport.DecodeRequestFunc) httptransport.DecodeRequestFunc {
	return func(drf httptransport.DecodeRequestFunc) httptransport.DecodeRequestFunc {
		return func(c context.Context, r *http.Request) (request interface{}, err error) {
			decoderLogger := logging.WithLayerName(logger, "transport")

			op, err := cont.GetOperationName(c)
			if err == nil {
				decoderLogger = logging.WithOperationName(decoderLogger, op)
			}

			decoderLogger.Debug("calling decoder")

			req, err := drf(c, r)
			if err != nil {
				decoderLogger.Error("decoder returned error", zap.Error(err))
			} else {
				decoderLogger.Debug("called decoder")
			}

			return req, err
		}
	}
}

func WithServiceLogger(logger *zap.Logger, next Service) *LoggingService {
	return &LoggingService{logger, next}
}

var _ Service = &LoggingService{}

type LoggingService struct {
	logger *zap.Logger
	next   Service
}

func (s *LoggingService) Compute(c context.Context, d string, lat, long float32) (output float32, err error) {
	defer func(begin time.Time) {
		fields := []zap.Field{
			zap.String("method", "compute"),
			zap.Float32("output", output),
			zap.Duration("took", time.Since(begin)),
		}
		id, err := cont.GetIDTokenClaims(c)
		if err == nil {
			fields = append(fields, zap.String("userID", id.Subject))
		}
		s.logger.Info("", fields...)
	}(time.Now())

	output, err = s.next.Compute(c, d, lat, long)
	return
}
