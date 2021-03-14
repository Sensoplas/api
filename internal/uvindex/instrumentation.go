package uvindex

import (
	"context"
	"net/http"

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
