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
	"github.com/FournyP/kubo-mcp/src/kubo_mcp_http_kubo/settings"
	"github.com/FournyP/kubo-mcp/tests/common_tests"
)

func TestWhenLoadingKuboSettings(t *testing.T) {
	t.Run("Given an API URL that is not absolute", func(t *testing.T) {
		t.Parallel()

		values := map[string]any{"kubo.api_url": "kubo:5001"}

		t.Run("Should stop at startup", func(t *testing.T) {
			t.Parallel()

			k := koanf.New(".")
			require.NoError(t, k.Load(confmap.Provider(values, "."), nil))

			assert.Panics(t, func() { settings.NewKuboSettings(common_tests.PanickingLogger(), k) })
		})
	})

	t.Run("Given an invalid request timeout", func(t *testing.T) {
		t.Parallel()

		values := map[string]any{"kubo.request_timeout": "30"}

		t.Run("Should stop at startup", func(t *testing.T) {
			t.Parallel()

			k := koanf.New(".")
			require.NoError(t, k.Load(confmap.Provider(values, "."), nil))

			assert.Panics(t, func() { settings.NewKuboSettings(common_tests.PanickingLogger(), k) })
		})
	})

	t.Run("Given the Kubo settings are not provided", func(t *testing.T) {
		t.Parallel()

		t.Run("Should fall back to the defaults", func(t *testing.T) {
			t.Parallel()

			s := settings.NewKuboSettings(zap.NewNop(), koanf.New("."))

			assert.Equal(t, settings.DefaultApiUrl, s.ApiUrl)
			assert.Empty(t, s.Authorization)
			assert.Equal(t, settings.DefaultRequestTimeout, s.RequestTimeout)
			assert.Equal(t, settings.DefaultUploadTimeout, s.UploadTimeout)
		})
	})

	t.Run("Given the Kubo settings are provided", func(t *testing.T) {
		t.Parallel()

		values := map[string]any{
			"kubo.api_url":         "http://kubo.railway.internal:5001/",
			"kubo.authorization":   "Bearer secret",
			"kubo.request_timeout": "15s",
			"kubo.upload_timeout":  "2m",
		}

		t.Run("Should populate every field", func(t *testing.T) {
			t.Parallel()

			k := koanf.New(".")
			require.NoError(t, k.Load(confmap.Provider(values, "."), nil))

			s := settings.NewKuboSettings(zap.NewNop(), k)

			assert.Equal(t, "http://kubo.railway.internal:5001", s.ApiUrl)
			assert.Equal(t, "Bearer secret", s.Authorization)
			assert.Equal(t, 15*time.Second, s.RequestTimeout)
			assert.Equal(t, 2*time.Minute, s.UploadTimeout)
		})
	})

	t.Run("Given the Kubo settings are provided through env variables", func(t *testing.T) {
		t.Run("Should populate every field from the env", func(t *testing.T) {
			t.Setenv("KUBO__API_URL", "http://env-kubo:5001")
			t.Setenv("KUBO__AUTHORIZATION", "Basic ZW52")
			t.Setenv("KUBO__REQUEST_TIMEOUT", "5s")
			t.Setenv("KUBO__UPLOAD_TIMEOUT", "1m")

			k, err := common_helpers.LoadKoanfConfig(filepath.Join(t.TempDir(), "missing.json"))
			require.NoError(t, err)

			s := settings.NewKuboSettings(zap.NewNop(), k)

			assert.Equal(t, "http://env-kubo:5001", s.ApiUrl)
			assert.Equal(t, "Basic ZW52", s.Authorization)
			assert.Equal(t, 5*time.Second, s.RequestTimeout)
			assert.Equal(t, time.Minute, s.UploadTimeout)
		})
	})
}
