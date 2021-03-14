package trend

import (
	"context"
	"net/http"

	firebase "firebase.google.com/go/v4"
	"github.com/Sensoplas/api/internal/auth"
	"github.com/Sensoplas/api/models"
	"github.com/go-kit/kit/auth/jwt"
	httptransport "github.com/go-kit/kit/transport/http"
	"go.uber.org/zap"
)

func MakeHTTPHandler(svc Service, logger *zap.Logger, fbApp *firebase.App) *httptransport.Server {
	authClient, err := fbApp.Auth(context.Background())

	if err != nil {
		logger.Panic("start up panic: cannot configure firebase auth client", zap.Error(err))
		panic(err)
	}

	// nolint
	var service Service
	service = svc
	endpoint := MakeGetExposureTrendEndpoint(service)
	endpoint = auth.FirebaseIDTokenMiddleware(authClient, logger)(endpoint)

	server := httptransport.NewServer(
		endpoint,
		DecodeExposureTrendRequestHTTP,
		EncodeResponse,
		httptransport.ServerBefore(func(c context.Context, r *http.Request) context.Context {
			logger.Info("received request", zap.String("method", r.Method), zap.String("host", r.Host), zap.String("path", r.URL.Path))

			return c

		}, jwt.HTTPToContext()),
		httptransport.ServerErrorEncoder(models.EncodeErrorHTTP(logger)),
	)

	return server
}
