package kubo_mcp_mcp

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp/use_cases"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp_mcp/helpers"
)

const ListPinsToolName = "list_pins"

type ListPinsTool struct {
	listPins use_cases.ListPinsInterface
}

func NewListPinsTool(listPins use_cases.ListPinsInterface) *ListPinsTool {
	return &ListPinsTool{listPins: listPins}
}

func (t *ListPinsTool) Name() string {
	return ListPinsToolName
}

func (t *ListPinsTool) Register(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        ListPinsToolName,
		Description: "List the CIDs recursively pinned on the node.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, req *mcp.CallToolRequest, args struct{}) (*mcp.CallToolResult, any, error) {
		response, err := t.listPins.List(ctx)
		if err != nil {
			return nil, nil, err
		}
		result, marshalErr := helpers.TextResult(response)
		return result, nil, marshalErr
	})
}
