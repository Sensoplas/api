package models

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"go.uber.org/zap"
)

type HTTPErrorResponse struct {
	Err string `json:"error"`
}

func EncodeErrorHTTP(logger *zap.Logger) func(ctx context.Context, err error, w http.ResponseWriter) {
	responseErrorLogger := logger.With(zap.String("co", "http-error-encoder"))
	return func(ctx context.Context, err error, w http.ResponseWriter) {
		// check if is explicit server error
		var serverError ServerError
		if errors.As(err, &serverError) {
			code := func() int {
				switch serverError.Type() {
				case NotImplemented:
					return http.StatusNotImplemented
				case ServiceUnavailable:
					return http.StatusServiceUnavailable
				default:
					return http.StatusInternalServerError
				}
			}()

			responseErrorLogger.Error("encountered server error, dispatching response", zap.Int("status", code), zap.Error(err))

			body := &HTTPErrorResponse{serverError.ClientMessage()}
			w.WriteHeader(code)
			w.Header().Set("Content-Type", "application-json")
			encoderErr := json.NewEncoder(w).Encode(body)
			if encoderErr != nil {
				panic(encoderErr)
			}
			return
		}

		// check if is explicit client error
		var clientError ClientError
		if errors.As(err, &clientError) {
			code := func() int {
				switch clientError.Type() {
				case Unauthorized:
					return http.StatusUnauthorized
				case Forbidden:
					return http.StatusForbidden
				case NotFound:
					return http.StatusNotFound
				case Conflict:
					return http.StatusConflict
				default:
					return http.StatusBadRequest
				}
			}()

			responseErrorLogger.Warn("encountered client error, dispatching response", zap.Int("status", code), zap.Error(err))

			body := &HTTPErrorResponse{clientError.ClientMessage()}
			w.WriteHeader(code)
			w.Header().Set("Content-Type", "application-json")
			encoderErr := json.NewEncoder(w).Encode(body)
			if encoderErr != nil {
				panic(encoderErr)
			}
			return
		}

		responseErrorLogger.Error("encountered error of unknown type, dispatching response", zap.Int("status", http.StatusInternalServerError), zap.Error(err))

		body := &HTTPErrorResponse{"Internal Server Error"}
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application-json")
		encoderErr := json.NewEncoder(w).Encode(body)
		if encoderErr != nil {
			panic(encoderErr)
		}
	}
}
