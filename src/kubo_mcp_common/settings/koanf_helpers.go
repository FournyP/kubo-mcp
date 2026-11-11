package settings

import (
	"strconv"
	"strings"
	"time"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

// Duration requires a unit: koanf's own Duration reads bare numbers as nanoseconds.
func Duration(logger *zap.Logger, k *koanf.Koanf, path string, defaultValue time.Duration) time.Duration {
	raw := strings.TrimSpace(k.String(path))
	if raw == "" {
		logger.Info(path+" is not set, using default value", zap.Duration("value", defaultValue))
		return defaultValue
	}

	value, err := time.ParseDuration(raw)
	if err != nil || value <= 0 {
		logger.Fatal(path+" must be a positive duration with a unit (e.g. \"30s\")",
			zap.String("value", raw),
			zap.Error(err),
		)
	}

	return value
}

func PositiveInt64(logger *zap.Logger, k *koanf.Koanf, path string, defaultValue int64) int64 {
	raw := strings.TrimSpace(k.String(path))
	if raw == "" {
		logger.Info(path+" is not set, using default value", zap.Int64("value", defaultValue))
		return defaultValue
	}

	value, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || value <= 0 {
		logger.Fatal(path+" must be a positive integer", zap.String("value", raw), zap.Error(err))
	}

	return value
}

func Bool(logger *zap.Logger, k *koanf.Koanf, path string, defaultValue bool) bool {
	raw := strings.TrimSpace(k.String(path))
	if raw == "" {
		return defaultValue
	}

	value, err := strconv.ParseBool(raw)
	if err != nil {
		logger.Fatal(path+" must be a boolean", zap.String("value", raw), zap.Error(err))
	}

	return value
}

// StringSlice reads a JSON array or a comma-separated env string, dropping empty entries.
func StringSlice(k *koanf.Koanf, path string) []string {
	raw := k.Strings(path)
	if len(raw) == 0 {
		raw = []string{k.String(path)}
	}

	values := make([]string, 0, len(raw))
	for _, entry := range raw {
		for _, value := range strings.Split(entry, ",") {
			if value = strings.TrimSpace(value); value != "" {
				values = append(values, value)
			}
		}
	}

	return values
}
