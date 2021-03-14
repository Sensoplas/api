package trend

import (
	"context"
	"errors"

	"github.com/Sensoplas/api/internal/cont"
	"github.com/Sensoplas/api/models"
	"github.com/go-kit/kit/endpoint"
)

type Request struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}

type Response struct {
	TimeSeries []ExposureEntry `json:"data"`
	Err        string          `json:"err,omitempty"`
}

func MakeGetExposureTrendEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		token, err := cont.GetIDTokenClaims(ctx)
		if err != nil {
			return nil, models.NewUnauthenticatedError(errors.New("authentication token is not found, please authenticate with bearer token"))
		}
		req, ok := request.(Request)
		if !ok {
			return nil, errors.New("request of unexpected type")
		}
		timeSeries, err := svc.GetTimeSeries(ctx, token.Subject, req.Start, req.End)
		if err != nil {
			return nil, err
		}
		return &Response{timeSeries, ""}, nil
	}
}
