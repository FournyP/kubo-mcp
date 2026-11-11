package common_tests

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// PanickingLogger turns logger.Fatal into a panic.
func PanickingLogger() *zap.Logger {
	return zap.NewNop().WithOptions(zap.OnFatal(zapcore.WriteThenPanic))
}
