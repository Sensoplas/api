package uvindex

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
)

type UVIRequest struct {
	Data SensorData `json:"data"`
}

type UVIResponse struct {
	Index float32 `json:"prediction"`
	Err   string  `json:"err,omitempty"`
}

func MakeUVIComputeEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(UVIRequest)
		if !ok {
			return nil, errors.New("request of unexpected type")
		}

		pred, err := svc.Compute(ctx, req.Data)
		if err != nil {
			return UVIResponse{Err: err.Error()}, nil
		}

		return UVIResponse{pred, ""}, nil
	}
}
