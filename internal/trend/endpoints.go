package trend

import (
	"context"
	"errors"
	"time"

	"github.com/Sensoplas/api/internal/cont"
	"github.com/Sensoplas/api/models"
	"github.com/go-kit/kit/endpoint"
	chart "github.com/wcharczuk/go-chart/v2"
)

type Request struct {
	width  int
	height int
}

type Response struct {
	Plot *chart.Chart
	Err  error
}

func MakeGetExposureTrendEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		token, err := cont.GetIDTokenClaims(ctx)
		if err != nil {
			return nil, models.NewUnauthenticatedError(errors.New("authentication token is not found, please authenticate with bearer token"))
		}

		req := request.(Request)

		timeSeries, err := svc.GetTimeSeries(ctx, token.Subject, time.Now().Add(-7*24*time.Hour).Unix(), time.Now().Unix(), int(req.width), int(req.height))
		if err != nil {
			return nil, err
		}
		return &Response{timeSeries, err}, nil
	}
}
