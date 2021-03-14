package uvindex

import (
	"github.com/Sensoplas/api/models"
	httptransport "github.com/go-kit/kit/transport/http"
	"go.uber.org/zap"
)

func MakeHTTPHandler(svc Service, logger *zap.Logger) *httptransport.Server {
	service := svc
	endpoint := MakeUVIComputeEndpoint(service)
	endpoint = endpointLoggingMiddleware(logger)(endpoint)

	server := httptransport.NewServer(
		endpoint,
		httpDecoderMiddleware(logger)(DecodeUVIRequestHTTP),
		EncodeResponse,
		httptransport.ServerBefore(),
		httptransport.ServerErrorEncoder(models.EncodeErrorHTTP(logger)),
	)

	return server
}
