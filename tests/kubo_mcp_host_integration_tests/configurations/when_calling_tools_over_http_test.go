package configurations_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	common_configurations "github.com/FournyP/kubo-mcp/src/kubo_mcp_common/configurations"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp_host/configurations"
	"github.com/FournyP/kubo-mcp/tests/common_tests/container_builders"
)

// onePixelPng is the smallest valid PNG, so get_file returns image content.
const onePixelPng = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR4nGP4DwAAAQEABRjYTgAAAABJRU5ErkJggg=="

type WhenCallingToolsOverHttpTestingSuite struct {
	session *mcp.ClientSession
}

// WhenCallingToolsOverHttpBeforeEach wires the real host against Kubo and connects an MCP client to it.
func WhenCallingToolsOverHttpBeforeEach(t *testing.T, values map[string]any) *WhenCallingToolsOverHttpTestingSuite {
	values["kubo.api_url"] = container_builders.GetKuboContainerBuilder().ApiUrl()
	values["logging.level"] = "error"
	k := koanf.New(".")
	require.NoError(t, k.Load(confmap.Provider(values, "."), nil))

	container := configurations.ConfigureDI(k)
	common_configurations.ConfigureLogging(container, k)
	configurations.ConfigureServer(container)

	var handler http.Handler
	require.NoError(t, container.Invoke(func(httpHandler http.Handler) { handler = httpHandler }))
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	client := mcp.NewClient(&mcp.Implementation{Name: "integration-test", Version: "0"}, nil)
	session, err := client.Connect(
		context.Background(),
		&mcp.StreamableClientTransport{Endpoint: server.URL + "/mcp"},
		nil,
	)
	require.NoError(t, err)
	t.Cleanup(func() { _ = session.Close() })

	return &WhenCallingToolsOverHttpTestingSuite{session: session}
}

func (suite *WhenCallingToolsOverHttpTestingSuite) callTool(
	t *testing.T,
	name string,
	arguments map[string]any,
) *mcp.CallToolResult {
	t.Helper()

	result, err := suite.session.CallTool(context.Background(), &mcp.CallToolParams{Name: name, Arguments: arguments})
	require.NoError(t, err)
	return result
}

func TestWhenCallingToolsOverHttp(t *testing.T) {
	t.Parallel()

	t.Run("Given the enabled tools are restricted", func(t *testing.T) {
		t.Parallel()

		values := map[string]any{"mcp.enabled_tools": "add_file,get_file"}

		t.Run("Should list only the enabled tools", func(t *testing.T) {
			t.Parallel()

			suite := WhenCallingToolsOverHttpBeforeEach(t, values)

			tools, err := suite.session.ListTools(context.Background(), nil)

			require.NoError(t, err)
			names := make([]string, 0, len(tools.Tools))
			for _, tool := range tools.Tools {
				names = append(names, tool.Name)
			}
			assert.ElementsMatch(t, []string{"add_file", "get_file"}, names)
		})
	})

	t.Run("Given an invalid CID", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a tool error carrying the validation message", func(t *testing.T) {
			t.Parallel()

			suite := WhenCallingToolsOverHttpBeforeEach(t, map[string]any{})

			result := suite.callTool(t, "pin_file", map[string]any{"cid": "not-a-cid"})

			assert.True(t, result.IsError)
			assert.Equal(t, "validation error on field 'cid': cid is not a valid CID",
				result.Content[0].(*mcp.TextContent).Text)
		})
	})

	t.Run("Given an image uploaded with add_file", func(t *testing.T) {
		t.Parallel()

		t.Run("Should read it back with get_file as image content", func(t *testing.T) {
			t.Parallel()

			suite := WhenCallingToolsOverHttpBeforeEach(t, map[string]any{})

			added := suite.callTool(
				t,
				"add_file",
				map[string]any{"fileName": "pixel.png", "contentBase64": onePixelPng},
			)
			require.False(t, added.IsError)
			var addedFile struct {
				Cid        string `json:"cid"`
				GatewayUrl string `json:"gatewayUrl"`
			}
			require.NoError(t, json.Unmarshal([]byte(added.Content[0].(*mcp.TextContent).Text), &addedFile))

			read := suite.callTool(t, "get_file", map[string]any{"path": addedFile.Cid})

			require.False(t, read.IsError)
			require.Len(t, read.Content, 2)
			image := read.Content[1].(*mcp.ImageContent)
			expected, err := base64.StdEncoding.DecodeString(onePixelPng)
			require.NoError(t, err)
			assert.Equal(t, "image/png", image.MIMEType)
			assert.Equal(t, expected, image.Data)
			assert.Equal(t, "https://ipfs.io/ipfs/"+addedFile.Cid+"?filename=pixel.png", addedFile.GatewayUrl)
		})
	})
}
