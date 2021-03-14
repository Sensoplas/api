package logging

import "go.uber.org/zap"

const (
	layerLoggerKey = "layer"
	opLoggerKey    = "op"
)

func WithLayerName(l *zap.Logger, name string) *zap.Logger {
	return l.With(zap.String(layerLoggerKey, name))
}

func WithOperationName(l *zap.Logger, name string) *zap.Logger {
	return l.With(zap.String(opLoggerKey, name))
}
