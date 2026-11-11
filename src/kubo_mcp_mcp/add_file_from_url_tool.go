package kubo_mcp_mcp

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp/use_cases"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp/use_cases/requests"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp_mcp/helpers"
)

const AddFileFromUrlToolName = "add_file_from_url"

type AddFileFromUrlArgs struct {
	Url      string `json:"url"                jsonschema:"public http(s) URL of the file to upload"`
	FileName string `json:"fileName,omitempty" jsonschema:"file name override; defaults to the name served by the URL"`
}

type AddFileFromUrlTool struct {
	addFileFromUrl use_cases.AddFileFromUrlInterface
}

func NewAddFileFromUrlTool(addFileFromUrl use_cases.AddFileFromUrlInterface) *AddFileFromUrlTool {
	return &AddFileFromUrlTool{addFileFromUrl: addFileFromUrl}
}

func (t *AddFileFromUrlTool) Name() string {
	return AddFileFromUrlToolName
}

func (t *AddFileFromUrlTool) Register(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        AddFileFromUrlToolName,
		Description: "Download a file from a public URL, upload it to IPFS and pin it. Returns its CID, ipfs:// URI and a gateway URL.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args AddFileFromUrlArgs) (*mcp.CallToolResult, any, error) {
		response, err := t.addFileFromUrl.Add(ctx, &requests.AddFileFromUrlRequest{
			FileNameRequest: requests.FileNameRequest{FileName: args.FileName},
			Url:             args.Url,
		})
		if err != nil {
			return nil, nil, err
		}
		result, marshalErr := helpers.TextResult(response)
		return result, nil, marshalErr
	})
}
