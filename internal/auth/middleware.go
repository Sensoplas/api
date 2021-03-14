package auth

import (
	"context"
	"errors"

	"firebase.google.com/go/v4/auth"
	"github.com/Sensoplas/api/internal/cont"
	"github.com/Sensoplas/api/models"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	"go.uber.org/zap"
)

func FirebaseIDTokenMiddleware(client *auth.Client, logger *zap.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			tokenString, ok := ctx.Value(jwt.JWTTokenContextKey).(string)
			if !ok {
				return nil, models.NewUnauthenticatedError(errors.New("token not found"))
			}

			token, err := client.VerifyIDToken(ctx, tokenString)
			if err != nil {
				return nil, models.NewUnauthenticatedError(err)
			}

			return next(cont.WithIDTokenClaims(ctx, token), request)
		}
	}
}
