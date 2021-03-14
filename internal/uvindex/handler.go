package uvindex

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"github.com/Sensoplas/api/internal/auth"
	"github.com/Sensoplas/api/models"
	"github.com/go-kit/kit/auth/jwt"
	httptransport "github.com/go-kit/kit/transport/http"
	"go.uber.org/zap"
)

func MakeHTTPHandler(svc Service, logger *zap.Logger) *httptransport.Server {
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		logger.Panic("start up panic: cannot configure firebase app", zap.Error(err))
		panic(err)
	}

	authClient, err := app.Auth(context.Background())

	if err != nil {
		logger.Panic("start up panic: cannot configure firebase auth client", zap.Error(err))
		panic(err)
	}

	service := svc
	endpoint := MakeUVIComputeEndpoint(service)
	endpoint = endpointLoggingMiddleware(logger)(endpoint)
	endpoint = auth.FirebaseIDTokenMiddleware(authClient, logger)(endpoint)

	server := httptransport.NewServer(
		endpoint,
		httpDecoderMiddleware(logger)(DecodeUVIRequestHTTP),
		EncodeResponse,
		httptransport.ServerBefore(jwt.HTTPToContext()),
		httptransport.ServerErrorEncoder(models.EncodeErrorHTTP(logger)),
	)

	return server
}
