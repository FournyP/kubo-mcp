package kubo_mcp_mcp

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp/use_cases"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp/use_cases/requests"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp_mcp/helpers"
)

const UnpinFileToolName = "unpin_file"

type UnpinFileArgs struct {
	Cid string `json:"cid" jsonschema:"CID to unpin from the node"`
}

type UnpinFileTool struct {
	unpinFile use_cases.UnpinFileInterface
}

func NewUnpinFileTool(unpinFile use_cases.UnpinFileInterface) *UnpinFileTool {
	return &UnpinFileTool{unpinFile: unpinFile}
}

func (t *UnpinFileTool) Name() string {
	return UnpinFileToolName
}

func (t *UnpinFileTool) Register(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        UnpinFileToolName,
		Description: "Unpin a CID from the node. The content stays until the node garbage-collects it.",
		Annotations: &mcp.ToolAnnotations{DestructiveHint: new(true)},
	}, func(ctx context.Context, req *mcp.CallToolRequest, args UnpinFileArgs) (*mcp.CallToolResult, any, error) {
		response, err := t.unpinFile.Unpin(ctx, &requests.UnpinFileRequest{
			CidRequest: requests.CidRequest{Cid: args.Cid},
		})
		if err != nil {
			return nil, nil, err
		}
		result, marshalErr := helpers.TextResult(response)
		return result, nil, marshalErr
	})
}
