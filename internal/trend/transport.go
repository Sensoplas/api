package trend

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Sensoplas/api/models"
)

func DecodeExposureTrendRequestHTTP(ctx context.Context, req *http.Request) (interface{}, error) {
	var request Request
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return nil, models.NewGenericClientError(fmt.Sprintf("bad request, could not decoded json body: %s", err.Error()), err)
	}
	return request, nil
}

func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
