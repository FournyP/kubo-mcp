package settings_test

import (
	"path/filepath"
	"testing"

	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	common_helpers "github.com/FournyP/kubo-mcp/src/kubo_mcp_common/helpers"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp_host/settings"
	"github.com/FournyP/kubo-mcp/tests/common_tests"
)

func TestWhenLoadingHostSettings(t *testing.T) {
	t.Run("Given an out of range port", func(t *testing.T) {
		t.Parallel()

		values := map[string]any{"host.port": 70000}

		t.Run("Should stop at startup", func(t *testing.T) {
			t.Parallel()

			k := koanf.New(".")
			require.NoError(t, k.Load(confmap.Provider(values, "."), nil))

			assert.Panics(t, func() { settings.NewHostSettings(common_tests.PanickingLogger(), k) })
		})
	})

	t.Run("Given the host settings are not provided", func(t *testing.T) {
		t.Parallel()

		t.Run("Should fall back to the default port", func(t *testing.T) {
			t.Parallel()

			s := settings.NewHostSettings(zap.NewNop(), koanf.New("."))

			assert.Equal(t, settings.DefaultPort, s.Port)
		})
	})

	t.Run("Given the host settings are provided", func(t *testing.T) {
		t.Parallel()

		values := map[string]any{"host.port": 8080}

		t.Run("Should use the provided value", func(t *testing.T) {
			t.Parallel()

			k := koanf.New(".")
			require.NoError(t, k.Load(confmap.Provider(values, "."), nil))

			s := settings.NewHostSettings(zap.NewNop(), k)

			assert.Equal(t, 8080, s.Port)
		})
	})

	t.Run("Given the host settings are provided through env variables", func(t *testing.T) {
		t.Run("Should use the env value", func(t *testing.T) {
			t.Setenv("HOST__PORT", "9000")

			k, err := common_helpers.LoadKoanfConfig(filepath.Join(t.TempDir(), "missing.json"))
			require.NoError(t, err)

			s := settings.NewHostSettings(zap.NewNop(), k)

			assert.Equal(t, 9000, s.Port)
		})
	})
}
