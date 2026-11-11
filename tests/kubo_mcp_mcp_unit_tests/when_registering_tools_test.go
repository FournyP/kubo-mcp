package kubo_mcp_mcp_test

import (
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp_mcp"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp_mcp/settings"
	"github.com/FournyP/kubo-mcp/tests/kubo_mcp_mcp_mocks"
)

type WhenRegisteringToolsTestingSuite struct {
	sut *kubo_mcp_mcp.ToolRegistrar

	mockAddFileTool *kubo_mcp_mcp_mocks.MockToolInterface
	mockGetFileTool *kubo_mcp_mcp_mocks.MockToolInterface
}

func WhenRegisteringToolsBeforeEach(t *testing.T, enabledTools []string) *WhenRegisteringToolsTestingSuite {
	mockController := gomock.NewController(t)
	mockAddFileTool := kubo_mcp_mcp_mocks.NewMockToolInterface(mockController)
	mockGetFileTool := kubo_mcp_mcp_mocks.NewMockToolInterface(mockController)
	// Name is a constant getter, not the behavior under test.
	mockAddFileTool.EXPECT().Name().Return("add_file").AnyTimes()
	mockGetFileTool.EXPECT().Name().Return("get_file").AnyTimes()

	return &WhenRegisteringToolsTestingSuite{
		sut: kubo_mcp_mcp.NewToolRegistrar(
			[]kubo_mcp_mcp.ToolInterface{mockAddFileTool, mockGetFileTool},
			&settings.McpSettings{EnabledTools: enabledTools},
		),
		mockAddFileTool: mockAddFileTool,
		mockGetFileTool: mockGetFileTool,
	}
}

func newTestServer() *mcp.Server {
	return mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0"}, nil)
}

func TestWhenRegisteringTools(t *testing.T) {
	t.Parallel()

	t.Run("Given an unknown enabled tool", func(t *testing.T) {
		t.Parallel()

		server := newTestServer()

		initSuite := func(suite *WhenRegisteringToolsTestingSuite) {
			suite.mockAddFileTool.EXPECT().Register(gomock.Any()).Times(0)
			suite.mockGetFileTool.EXPECT().Register(gomock.Any()).Times(0)
		}

		t.Run("Should return an error without registering any tool", func(t *testing.T) {
			t.Parallel()

			suite := WhenRegisteringToolsBeforeEach(t, []string{"add_file", "delete_everything"})
			initSuite(suite)

			err := suite.sut.Register(server)

			assert.EqualError(t, err, `mcp.enabled_tools: unknown tool "delete_everything"`)
		})
	})

	t.Run("Given no enabled tools are configured", func(t *testing.T) {
		t.Parallel()

		server := newTestServer()

		initSuite := func(suite *WhenRegisteringToolsTestingSuite) {
			suite.mockAddFileTool.EXPECT().Register(server).Times(1)
			suite.mockGetFileTool.EXPECT().Register(server).Times(1)
		}

		t.Run("Should register every tool", func(t *testing.T) {
			t.Parallel()

			suite := WhenRegisteringToolsBeforeEach(t, nil)
			initSuite(suite)

			assert.NoError(t, suite.sut.Register(server))
		})
	})

	t.Run("Given a subset of tools is enabled", func(t *testing.T) {
		t.Parallel()

		server := newTestServer()

		initSuite := func(suite *WhenRegisteringToolsTestingSuite) {
			suite.mockAddFileTool.EXPECT().Register(gomock.Any()).Times(0)
			suite.mockGetFileTool.EXPECT().Register(server).Times(1)
		}

		t.Run("Should register only the enabled tools, once each", func(t *testing.T) {
			t.Parallel()

			suite := WhenRegisteringToolsBeforeEach(t, []string{"get_file", "get_file"})
			initSuite(suite)

			assert.NoError(t, suite.sut.Register(server))
		})
	})
}
