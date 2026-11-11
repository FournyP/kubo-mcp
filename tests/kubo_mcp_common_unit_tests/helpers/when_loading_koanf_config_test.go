package helpers_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp_common/helpers"
)

func TestWhenLoadingKoanfConfig(t *testing.T) {
	t.Run("Given the config file is missing", func(t *testing.T) {
		t.Run("Should load the env variables alone", func(t *testing.T) {
			t.Setenv("KUBO__API_URL", "http://env-kubo:5001")

			k, err := helpers.LoadKoanfConfig(filepath.Join(t.TempDir(), "missing.json"))

			require.NoError(t, err)
			assert.Equal(t, "http://env-kubo:5001", k.String("kubo.api_url"))
		})
	})

	t.Run("Given the config file is not valid JSON", func(t *testing.T) {
		t.Run("Should return a parse error", func(t *testing.T) {
			configFile := filepath.Join(t.TempDir(), "appsettings.json")
			require.NoError(t, os.WriteFile(configFile, []byte("{"), 0o600))

			_, err := helpers.LoadKoanfConfig(configFile)

			assert.Error(t, err)
		})
	})

	t.Run("Given a key is set in both the file and the env", func(t *testing.T) {
		t.Run("Should use the env value", func(t *testing.T) {
			t.Setenv("KUBO__API_URL", "http://env-kubo:5001")
			configFile := filepath.Join(t.TempDir(), "appsettings.json")
			content := `{"kubo":{"api_url":"http://file-kubo:5001","upload_timeout":"1m"}}`
			require.NoError(t, os.WriteFile(configFile, []byte(content), 0o600))

			k, err := helpers.LoadKoanfConfig(configFile)

			require.NoError(t, err)
			assert.Equal(t, "http://env-kubo:5001", k.String("kubo.api_url"))
			assert.Equal(t, "1m", k.String("kubo.upload_timeout"))
		})
	})
}
