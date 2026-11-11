package settings_test

import (
	"path/filepath"
	"testing"

	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp/settings"
	common_helpers "github.com/FournyP/kubo-mcp/src/kubo_mcp_common/helpers"
	"github.com/FournyP/kubo-mcp/tests/common_tests"
)

func TestWhenLoadingIpfsSettings(t *testing.T) {
	t.Run("Given an invalid max upload size", func(t *testing.T) {
		t.Parallel()

		values := map[string]any{"ipfs.max_upload_size_bytes": "-1"}

		t.Run("Should stop at startup", func(t *testing.T) {
			t.Parallel()

			k := koanf.New(".")
			require.NoError(t, k.Load(confmap.Provider(values, "."), nil))

			assert.Panics(t, func() { settings.NewIpfsSettings(common_tests.PanickingLogger(), k) })
		})
	})

	t.Run("Given an invalid max read size", func(t *testing.T) {
		t.Parallel()

		values := map[string]any{"ipfs.max_read_size_bytes": "five"}

		t.Run("Should stop at startup", func(t *testing.T) {
			t.Parallel()

			k := koanf.New(".")
			require.NoError(t, k.Load(confmap.Provider(values, "."), nil))

			assert.Panics(t, func() { settings.NewIpfsSettings(common_tests.PanickingLogger(), k) })
		})
	})

	t.Run("Given a public gateway URL that is not absolute", func(t *testing.T) {
		t.Parallel()

		values := map[string]any{"ipfs.public_gateway_url": "gateway.example"}

		t.Run("Should stop at startup", func(t *testing.T) {
			t.Parallel()

			k := koanf.New(".")
			require.NoError(t, k.Load(confmap.Provider(values, "."), nil))

			assert.Panics(t, func() { settings.NewIpfsSettings(common_tests.PanickingLogger(), k) })
		})
	})

	t.Run("Given the IPFS settings are not provided", func(t *testing.T) {
		t.Parallel()

		t.Run("Should fall back to the defaults", func(t *testing.T) {
			t.Parallel()

			s := settings.NewIpfsSettings(zap.NewNop(), koanf.New("."))

			assert.Equal(t, settings.DefaultMaxUploadSizeBytes, s.MaxUploadSizeBytes)
			assert.Equal(t, settings.DefaultMaxReadSizeBytes, s.MaxReadSizeBytes)
			assert.Equal(t, settings.DefaultPublicGatewayUrl, s.PublicGatewayUrl)
		})
	})

	t.Run("Given the IPFS settings are provided", func(t *testing.T) {
		t.Parallel()

		values := map[string]any{
			"ipfs.max_upload_size_bytes": 2048,
			"ipfs.max_read_size_bytes":   1024,
			"ipfs.public_gateway_url":    "https://gateway.example/",
		}

		t.Run("Should populate every field", func(t *testing.T) {
			t.Parallel()

			k := koanf.New(".")
			require.NoError(t, k.Load(confmap.Provider(values, "."), nil))

			s := settings.NewIpfsSettings(zap.NewNop(), k)

			assert.Equal(t, int64(2048), s.MaxUploadSizeBytes)
			assert.Equal(t, int64(1024), s.MaxReadSizeBytes)
			assert.Equal(t, "https://gateway.example", s.PublicGatewayUrl)
		})
	})

	t.Run("Given the IPFS settings are provided through env variables", func(t *testing.T) {
		t.Run("Should populate every field from the env", func(t *testing.T) {
			t.Setenv("IPFS__MAX_UPLOAD_SIZE_BYTES", "4096")
			t.Setenv("IPFS__MAX_READ_SIZE_BYTES", "512")
			t.Setenv("IPFS__PUBLIC_GATEWAY_URL", "https://env.gateway.example")

			k, err := common_helpers.LoadKoanfConfig(filepath.Join(t.TempDir(), "missing.json"))
			require.NoError(t, err)

			s := settings.NewIpfsSettings(zap.NewNop(), k)

			assert.Equal(t, int64(4096), s.MaxUploadSizeBytes)
			assert.Equal(t, int64(512), s.MaxReadSizeBytes)
			assert.Equal(t, "https://env.gateway.example", s.PublicGatewayUrl)
		})
	})
}
