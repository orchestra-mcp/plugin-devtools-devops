package tools

import (
	"context"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/plugin-devtools-devops/internal/gh"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// DevopsListPipelinesSchema returns the JSON Schema for the devops_list_pipelines tool.
func DevopsListPipelinesSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"repo": map[string]any{
				"type":        "string",
				"description": "Repository in OWNER/REPO format. Uses current directory repo if omitted.",
			},
		},
	})
	return s
}

// DevopsListPipelines returns a tool handler that lists GitHub Actions workflows.
func DevopsListPipelines() func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		args := []string{"workflow", "list", "--json", "name,id,state"}

		repo := helpers.GetString(req.Arguments, "repo")
		if repo != "" {
			args = append(args, "--repo", repo)
		}

		out, err := gh.Run(ctx, args...)
		if err != nil {
			return helpers.ErrorResult("gh_error", err.Error()), nil
		}
		return helpers.TextResult(out), nil
	}
}
