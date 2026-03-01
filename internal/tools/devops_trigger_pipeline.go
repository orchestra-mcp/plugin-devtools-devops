package tools

import (
	"context"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/plugin-devtools-devops/internal/gh"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// DevopsTriggerPipelineSchema returns the JSON Schema for the devops_trigger_pipeline tool.
func DevopsTriggerPipelineSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"workflow": map[string]any{
				"type":        "string",
				"description": "Workflow file name (e.g. ci.yml) or workflow ID to trigger",
			},
			"repo": map[string]any{
				"type":        "string",
				"description": "Repository in OWNER/REPO format. Uses current directory repo if omitted.",
			},
			"ref": map[string]any{
				"type":        "string",
				"description": "Branch or tag ref to run the workflow on (default: the repo default branch)",
			},
		},
		"required": []any{"workflow"},
	})
	return s
}

// DevopsTriggerPipeline returns a tool handler that triggers a GitHub Actions workflow.
func DevopsTriggerPipeline() func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "workflow"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}

		workflow := helpers.GetString(req.Arguments, "workflow")
		args := []string{"workflow", "run", workflow}

		repo := helpers.GetString(req.Arguments, "repo")
		if repo != "" {
			args = append(args, "--repo", repo)
		}

		ref := helpers.GetString(req.Arguments, "ref")
		if ref != "" {
			args = append(args, "--ref", ref)
		}

		out, err := gh.Run(ctx, args...)
		if err != nil {
			return helpers.ErrorResult("gh_error", err.Error()), nil
		}
		if out == "" {
			out = "Workflow dispatch triggered successfully."
		}
		return helpers.TextResult(out), nil
	}
}
