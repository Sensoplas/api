package trend

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/Sensoplas/api/models"
	"github.com/wcharczuk/go-chart/v2"
)

func DecodeExposureTrendRequestHTTP(ctx context.Context, req *http.Request) (interface{}, error) {
	// var request Request
	// if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
	// 	return nil, models.NewGenericClientError(fmt.Sprintf("bad request, could not decoded json body: %s", err.Error()), err)
	// }
	vals := req.URL.Query()

	width, we := strconv.Atoi(vals.Get("width"))
	height, he := strconv.Atoi(vals.Get("height"))
	if we != nil || he != nil {
		return nil, models.NewGenericClientError("need width and height specified in query param", errors.New("unspecified dimension"))
	}

	return Request{width, height}, nil
}

func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "image/svg")
	return response.(*Response).Plot.Render(chart.SVG, w)
}
