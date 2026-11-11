package kubo_mcp_mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp/use_cases"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp/use_cases/requests"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp/use_cases/responses"
)

const GetFileToolName = "get_file"

type GetFileArgs struct {
	Path string `json:"path" jsonschema:"CID, optionally followed by a sub-path (<cid>/images/logo.png)"`
}

type GetFileTool struct {
	getFile use_cases.GetFileInterface
}

func NewGetFileTool(getFile use_cases.GetFileInterface) *GetFileTool {
	return &GetFileTool{getFile: getFile}
}

func (t *GetFileTool) Name() string {
	return GetFileToolName
}

func (t *GetFileTool) Register(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        GetFileToolName,
		Description: "Read a file from IPFS. Images come back as image content, text as text, anything else as base64.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, req *mcp.CallToolRequest, args GetFileArgs) (*mcp.CallToolResult, any, error) {
		response, err := t.getFile.Get(ctx, &requests.GetFileRequest{
			PathRequest: requests.PathRequest{Path: args.Path},
		})
		if err != nil {
			return nil, nil, err
		}
		result, marshalErr := t.toResult(response)
		return result, nil, marshalErr
	})
}

func (t *GetFileTool) toResult(response *responses.GetFileResponse) (*mcp.CallToolResult, error) {
	metadata, err := json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal tool result: %w", err)
	}

	var content mcp.Content
	switch {
	case strings.HasPrefix(response.MimeType, "image/"):
		content = &mcp.ImageContent{Data: response.Content, MIMEType: response.MimeType}
	case utf8.Valid(response.Content):
		content = &mcp.TextContent{Text: string(response.Content)}
	default:
		content = &mcp.EmbeddedResource{Resource: &mcp.ResourceContents{
			URI:      "ipfs://" + response.Path,
			MIMEType: response.MimeType,
			Blob:     response.Content,
		}}
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: string(metadata)}, content},
	}, nil
}
