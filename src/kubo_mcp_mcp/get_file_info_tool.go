package kubo_mcp_mcp

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp/use_cases"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp/use_cases/requests"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp_mcp/helpers"
)

const GetFileInfoToolName = "get_file_info"

type GetFileInfoArgs struct {
	Path string `json:"path" jsonschema:"CID, optionally followed by a sub-path (<cid>/images/logo.png)"`
}

type GetFileInfoTool struct {
	getFileInfo use_cases.GetFileInfoInterface
}

func NewGetFileInfoTool(getFileInfo use_cases.GetFileInfoInterface) *GetFileInfoTool {
	return &GetFileInfoTool{getFileInfo: getFileInfo}
}

func (t *GetFileInfoTool) Name() string {
	return GetFileInfoToolName
}

func (t *GetFileInfoTool) Register(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        GetFileInfoToolName,
		Description: "Describe a file or directory on IPFS: CID, type, size, block count and gateway URL.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, req *mcp.CallToolRequest, args GetFileInfoArgs) (*mcp.CallToolResult, any, error) {
		response, err := t.getFileInfo.Get(ctx, &requests.GetFileInfoRequest{
			PathRequest: requests.PathRequest{Path: args.Path},
		})
		if err != nil {
			return nil, nil, err
		}
		result, marshalErr := helpers.TextResult(response)
		return result, nil, marshalErr
	})
}
