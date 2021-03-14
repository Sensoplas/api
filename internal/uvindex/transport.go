package uvindex

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Sensoplas/api/models"
)

func DecodeUVIRequestHTTP(ctx context.Context, req *http.Request) (interface{}, error) {
	var request UVIRequest
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return nil, models.NewGenericClientError(fmt.Sprintf("bad request, could not decode json body: %s", err.Error()), err)
	}
	return request, nil
}

func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
