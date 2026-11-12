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
	"github.com/FournyP/kubo-mcp/src/kubo_mcp_mcp/settings"
)

func TestWhenLoadingAuthSettings(t *testing.T) {
	t.Run("Given the auth settings are not provided", func(t *testing.T) {
		t.Parallel()

		t.Run("Should fall back to no API key", func(t *testing.T) {
			t.Parallel()

			s := settings.NewAuthSettings(zap.NewNop(), koanf.New("."))

			assert.Empty(t, s.ApiKeys)
		})
	})

	t.Run("Given the auth settings are provided", func(t *testing.T) {
		t.Parallel()

		values := map[string]any{"auth.api_keys": []any{"first-key", "second-key"}}

		t.Run("Should populate every field", func(t *testing.T) {
			t.Parallel()

			k := koanf.New(".")
			require.NoError(t, k.Load(confmap.Provider(values, "."), nil))

			s := settings.NewAuthSettings(zap.NewNop(), k)

			assert.Equal(t, []string{"first-key", "second-key"}, s.ApiKeys)
		})
	})

	t.Run("Given the auth settings are provided through env variables", func(t *testing.T) {
		t.Run("Should populate every field from the env", func(t *testing.T) {
			t.Setenv("AUTH__API_KEYS", "first-key,second-key")

			k, err := common_helpers.LoadKoanfConfig(filepath.Join(t.TempDir(), "missing.json"))
			require.NoError(t, err)

			s := settings.NewAuthSettings(zap.NewNop(), k)

			assert.Equal(t, []string{"first-key", "second-key"}, s.ApiKeys)
		})
	})
}
