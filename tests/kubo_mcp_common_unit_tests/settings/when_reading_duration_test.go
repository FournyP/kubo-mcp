package settings_test

import (
	"testing"
	"time"

	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp_common/settings"
	"github.com/FournyP/kubo-mcp/tests/common_tests"
)

func TestWhenReadingDuration(t *testing.T) {
	t.Parallel()

	t.Run("Given a duration without a unit", func(t *testing.T) {
		t.Parallel()

		values := map[string]any{"section.timeout": "30"}

		t.Run("Should stop at startup", func(t *testing.T) {
			t.Parallel()

			k := koanf.New(".")
			require.NoError(t, k.Load(confmap.Provider(values, "."), nil))

			assert.Panics(t, func() {
				settings.Duration(common_tests.PanickingLogger(), k, "section.timeout", time.Second)
			})
		})
	})

	t.Run("Given a negative duration", func(t *testing.T) {
		t.Parallel()

		values := map[string]any{"section.timeout": "-5s"}

		t.Run("Should stop at startup", func(t *testing.T) {
			t.Parallel()

			k := koanf.New(".")
			require.NoError(t, k.Load(confmap.Provider(values, "."), nil))

			assert.Panics(t, func() {
				settings.Duration(common_tests.PanickingLogger(), k, "section.timeout", time.Second)
			})
		})
	})

	t.Run("Given the duration is not provided", func(t *testing.T) {
		t.Parallel()

		t.Run("Should fall back to the default value", func(t *testing.T) {
			t.Parallel()

			assert.Equal(
				t,
				time.Second,
				settings.Duration(zap.NewNop(), koanf.New("."), "section.timeout", time.Second),
			)
		})
	})

	t.Run("Given the duration is provided", func(t *testing.T) {
		t.Parallel()

		values := map[string]any{"section.timeout": "90s"}

		t.Run("Should use the provided value", func(t *testing.T) {
			t.Parallel()

			k := koanf.New(".")
			require.NoError(t, k.Load(confmap.Provider(values, "."), nil))

			assert.Equal(t, 90*time.Second, settings.Duration(zap.NewNop(), k, "section.timeout", time.Second))
		})
	})
}
