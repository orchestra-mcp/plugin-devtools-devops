package tools

import (
	"context"
	"fmt"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/plugin-devtools-devops/internal/gh"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// DevopsRollbackSchema returns the JSON Schema for the devops_rollback tool.
func DevopsRollbackSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"run_id": map[string]any{
				"type":        "string",
				"description": "GitHub Actions run ID to re-run for rollback",
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

// DevopsRollback returns a tool handler that triggers a rollback by re-running a previous workflow run.
func DevopsRollback() func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "run_id"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}

		runID := helpers.GetString(req.Arguments, "run_id")
		args := []string{"run", "rerun", runID}

		repo := helpers.GetString(req.Arguments, "repo")
		if repo != "" {
			args = append(args, "--repo", repo)
		}

		_, err := gh.Run(ctx, args...)
		if err != nil {
			return helpers.ErrorResult("gh_error", err.Error()), nil
		}
		return helpers.TextResult(fmt.Sprintf("Rollback initiated by re-running workflow run %s", runID)), nil
	}
}
