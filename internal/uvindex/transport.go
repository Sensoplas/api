package uvindex

import (
	"context"
	"encoding/json"
	"net/http"
)

func DecodeUVIRequestHTTP(ctx context.Context, req *http.Request) (interface{}, error) {
	var request UVIRequest
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
