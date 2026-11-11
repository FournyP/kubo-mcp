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

func TestWhenLoadingMcpSettings(t *testing.T) {
	t.Run("Given the MCP settings are not provided", func(t *testing.T) {
		t.Parallel()

		t.Run("Should fall back to exposing every tool", func(t *testing.T) {
			t.Parallel()

			s := settings.NewMcpSettings(zap.NewNop(), koanf.New("."))

			assert.Empty(t, s.EnabledTools)
		})
	})

	t.Run("Given the MCP settings are provided", func(t *testing.T) {
		t.Parallel()

		values := map[string]any{"mcp.enabled_tools": []any{"add_file", " get_file "}}

		t.Run("Should populate every field", func(t *testing.T) {
			t.Parallel()

			k := koanf.New(".")
			require.NoError(t, k.Load(confmap.Provider(values, "."), nil))

			s := settings.NewMcpSettings(zap.NewNop(), k)

			assert.Equal(t, []string{"add_file", "get_file"}, s.EnabledTools)
		})
	})

	t.Run("Given the MCP settings are provided through env variables", func(t *testing.T) {
		t.Run("Should populate every field from the env", func(t *testing.T) {
			t.Setenv("MCP__ENABLED_TOOLS", "add_file, get_file,")

			k, err := common_helpers.LoadKoanfConfig(filepath.Join(t.TempDir(), "missing.json"))
			require.NoError(t, err)

			s := settings.NewMcpSettings(zap.NewNop(), k)

			assert.Equal(t, []string{"add_file", "get_file"}, s.EnabledTools)
		})
	})
}
