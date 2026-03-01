package tools

import (
	"context"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/plugin-devtools-devops/internal/gh"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// DevopsEnvVarsSchema returns the JSON Schema for the devops_env_vars tool.
func DevopsEnvVarsSchema() *structpb.Struct {
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

// DevopsEnvVars returns a tool handler that lists GitHub repository secrets.
// Note: secret values are never returned by the GitHub API — only names are listed.
func DevopsEnvVars() func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		args := []string{"secret", "list"}

		repo := helpers.GetString(req.Arguments, "repo")
		if repo != "" {
			args = append(args, "--repo", repo)
		}

		out, err := gh.Run(ctx, args...)
		if err != nil {
			return helpers.ErrorResult("gh_error", err.Error()), nil
		}
		if out == "" {
			out = "No secrets found."
		}
		return helpers.TextResult(out), nil
	}
}
