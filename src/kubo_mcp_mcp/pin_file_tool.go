package kubo_mcp_mcp

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp/use_cases"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp/use_cases/requests"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp_mcp/helpers"
)

const PinFileToolName = "pin_file"

type PinFileArgs struct {
	Cid string `json:"cid" jsonschema:"CID to pin on the node"`
}

type PinFileTool struct {
	pinFile use_cases.PinFileInterface
}

func NewPinFileTool(pinFile use_cases.PinFileInterface) *PinFileTool {
	return &PinFileTool{pinFile: pinFile}
}

func (t *PinFileTool) Name() string {
	return PinFileToolName
}

func (t *PinFileTool) Register(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        PinFileToolName,
		Description: "Pin an existing CID on the node so it is kept and served. Fetches the content from the network if needed.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args PinFileArgs) (*mcp.CallToolResult, any, error) {
		response, err := t.pinFile.Pin(ctx, &requests.PinFileRequest{CidRequest: requests.CidRequest{Cid: args.Cid}})
		if err != nil {
			return nil, nil, err
		}
		result, marshalErr := helpers.TextResult(response)
		return result, nil, marshalErr
	})
}
