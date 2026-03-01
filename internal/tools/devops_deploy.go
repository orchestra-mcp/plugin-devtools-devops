package tools

import (
	"context"
	"fmt"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/plugin-devtools-devops/internal/gh"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// DevopsDeploySchema returns the JSON Schema for the devops_deploy tool.
func DevopsDeploySchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"workflow": map[string]any{
				"type":        "string",
				"description": "Workflow file name (e.g. deploy.yml) or workflow ID to trigger as deployment",
			},
			"repo": map[string]any{
				"type":        "string",
				"description": "Repository in OWNER/REPO format. Uses current directory repo if omitted.",
			},
			"ref": map[string]any{
				"type":        "string",
				"description": "Branch or tag to deploy (default: the repo default branch)",
			},
			"environment": map[string]any{
				"type":        "string",
				"description": "Target environment (e.g. production, staging)",
			},
		},
		"required": []any{"workflow"},
	})
	return s
}

// DevopsDeploy returns a tool handler that triggers a deployment via GitHub Actions workflow dispatch.
func DevopsDeploy() func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
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
			environment := helpers.GetString(req.Arguments, "environment")
			if environment == "" {
				environment = "default"
			}
			out = fmt.Sprintf("Deployment triggered successfully for environment: %s", environment)
		}
		return helpers.TextResult(out), nil
	}
}
