package helpers

import (
	"errors"

	"go.uber.org/dig"
	"go.uber.org/zap"
)

func FatalOnInvokeError(container *dig.Container, message string, err error) {
	invokeErr := container.Invoke(func(logger *zap.Logger) {
		logger.Fatal(message, zap.Error(err))
	})
	if invokeErr != nil {
		panic(errors.Join(err, invokeErr))
	}
}
