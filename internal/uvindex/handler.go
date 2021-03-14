package uvindex

import (
	httptransport "github.com/go-kit/kit/transport/http"
	"go.uber.org/zap"
)

func MakeHTTPHandler(svc Service, logger *zap.Logger) *httptransport.Server {
	service := svc
	endpoint := MakeUVIComputeEndpoint(service)
	endpoint = endpointLoggingMiddleware(logger)(endpoint)

	return httptransport.NewServer(
		endpoint,
		httpDecoderMiddleware(logger)(DecodeUVIRequestHTTP),
		EncodeResponse,
	)
}
