package kubo_mcp_mcp

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp/use_cases"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp/use_cases/requests"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp_mcp/helpers"
)

const AddFileToolName = "add_file"

type AddFileArgs struct {
	FileName      string `json:"fileName"      jsonschema:"file name with extension, e.g. logo.png"`
	ContentBase64 string `json:"contentBase64" jsonschema:"file content as base64 or as a data URI"`
}

type AddFileTool struct {
	addFile use_cases.AddFileInterface
}

func NewAddFileTool(addFile use_cases.AddFileInterface) *AddFileTool {
	return &AddFileTool{addFile: addFile}
}

func (t *AddFileTool) Name() string {
	return AddFileToolName
}

func (t *AddFileTool) Register(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        AddFileToolName,
		Description: "Upload a file (e.g. an image) to IPFS from base64 content and pin it. Returns its CID, ipfs:// URI and a gateway URL. Prefer add_file_from_url for large files.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args AddFileArgs) (*mcp.CallToolResult, any, error) {
		response, err := t.addFile.Add(ctx, &requests.AddFileRequest{
			FileNameRequest: requests.FileNameRequest{FileName: args.FileName},
			ContentBase64:   args.ContentBase64,
		})
		if err != nil {
			return nil, nil, err
		}
		result, marshalErr := helpers.TextResult(response)
		return result, nil, marshalErr
	})
}
