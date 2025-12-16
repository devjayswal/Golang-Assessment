package logger

import (
	"go.uber.org/zap"
)

func New() *zap.Logger {
	lz, _ := zap.NewProduction()
	return lz
}
