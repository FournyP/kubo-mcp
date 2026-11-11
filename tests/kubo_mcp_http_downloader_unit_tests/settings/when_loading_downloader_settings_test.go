package settings_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	common_helpers "github.com/FournyP/kubo-mcp/src/kubo_mcp_common/helpers"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp_http_downloader/settings"
	"github.com/FournyP/kubo-mcp/tests/common_tests"
)

func TestWhenLoadingDownloaderSettings(t *testing.T) {
	t.Run("Given an invalid private networks flag", func(t *testing.T) {
		t.Parallel()

		values := map[string]any{"downloader.allow_private_networks": "maybe"}

		t.Run("Should stop at startup", func(t *testing.T) {
			t.Parallel()

			k := koanf.New(".")
			require.NoError(t, k.Load(confmap.Provider(values, "."), nil))

			assert.Panics(t, func() { settings.NewDownloaderSettings(common_tests.PanickingLogger(), k) })
		})
	})

	t.Run("Given the downloader settings are not provided", func(t *testing.T) {
		t.Parallel()

		t.Run("Should fall back to the defaults", func(t *testing.T) {
			t.Parallel()

			s := settings.NewDownloaderSettings(zap.NewNop(), koanf.New("."))

			assert.Equal(t, settings.DefaultRequestTimeout, s.RequestTimeout)
			assert.False(t, s.AllowPrivateNetworks)
		})
	})

	t.Run("Given the downloader settings are provided", func(t *testing.T) {
		t.Parallel()

		values := map[string]any{
			"downloader.request_timeout":        "45s",
			"downloader.allow_private_networks": true,
		}

		t.Run("Should populate every field", func(t *testing.T) {
			t.Parallel()

			k := koanf.New(".")
			require.NoError(t, k.Load(confmap.Provider(values, "."), nil))

			s := settings.NewDownloaderSettings(zap.NewNop(), k)

			assert.Equal(t, 45*time.Second, s.RequestTimeout)
			assert.True(t, s.AllowPrivateNetworks)
		})
	})

	t.Run("Given the downloader settings are provided through env variables", func(t *testing.T) {
		t.Run("Should populate every field from the env", func(t *testing.T) {
			t.Setenv("DOWNLOADER__REQUEST_TIMEOUT", "10s")
			t.Setenv("DOWNLOADER__ALLOW_PRIVATE_NETWORKS", "true")

			k, err := common_helpers.LoadKoanfConfig(filepath.Join(t.TempDir(), "missing.json"))
			require.NoError(t, err)

			s := settings.NewDownloaderSettings(zap.NewNop(), k)

			assert.Equal(t, 10*time.Second, s.RequestTimeout)
			assert.True(t, s.AllowPrivateNetworks)
		})
	})
}
