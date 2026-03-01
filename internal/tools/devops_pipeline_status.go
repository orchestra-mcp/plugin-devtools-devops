package tools

import (
	"context"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/plugin-devtools-devops/internal/gh"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// DevopsPipelineStatusSchema returns the JSON Schema for the devops_pipeline_status tool.
func DevopsPipelineStatusSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"run_id": map[string]any{
				"type":        "string",
				"description": "GitHub Actions run ID",
			},
			"repo": map[string]any{
				"type":        "string",
				"description": "Repository in OWNER/REPO format. Uses current directory repo if omitted.",
			},
		},
		"required": []any{"run_id"},
	})
	return s
}

// DevopsPipelineStatus returns a tool handler that checks the status of a workflow run.
func DevopsPipelineStatus() func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "run_id"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}

		runID := helpers.GetString(req.Arguments, "run_id")
		args := []string{"run", "view", runID, "--json", "status,conclusion,name,createdAt,updatedAt"}

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
