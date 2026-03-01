package tools

import (
	"context"
	"fmt"
	"strings"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/plugin-devtools-devops/internal/gh"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// DevopsListDeploymentsSchema returns the JSON Schema for the devops_list_deployments tool.
func DevopsListDeploymentsSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"repo": map[string]any{
				"type":        "string",
				"description": "Repository in OWNER/REPO format. Uses current directory repo if omitted.",
			},
			"environment": map[string]any{
				"type":        "string",
				"description": "Filter deployments by environment name (e.g. production, staging).",
			},
		},
	})
	return s
}

// DevopsListDeployments returns a tool handler that lists GitHub deployments or releases.
func DevopsListDeployments() func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		repo := helpers.GetString(req.Arguments, "repo")
		environment := helpers.GetString(req.Arguments, "environment")

		// Use the GitHub API to list deployments when an environment filter is given
		// or when we need structured deployment data.
		// Fall back to gh release list for a simpler listing.
		if environment != "" || repo != "" {
			// Build the API path. If repo is not given we rely on the inferred repo.
			var apiPath string
			if repo != "" {
				apiPath = fmt.Sprintf("/repos/%s/deployments", repo)
			} else {
				// Let gh infer the repo from git remote.
				apiPath = "/repos/{owner}/{repo}/deployments"
			}

			args := []string{"api", apiPath}
			if environment != "" {
				args = append(args, "--field", fmt.Sprintf("environment=%s", environment))
			}

			out, err := gh.Run(ctx, args...)
			if err != nil {
				// Fall back to release list on API error.
				return fallbackReleaseList(ctx, repo)
			}
			return helpers.TextResult(out), nil
		}

		return fallbackReleaseList(ctx, repo)
	}
}

func fallbackReleaseList(ctx context.Context, repo string) (*pluginv1.ToolResponse, error) {
	args := []string{"release", "list"}
	if repo != "" {
		args = append(args, "--repo", repo)
	}
	out, err := gh.Run(ctx, args...)
	if err != nil {
		return helpers.ErrorResult("gh_error", err.Error()), nil
	}
	if strings.TrimSpace(out) == "" {
		out = "No releases found."
	}
	return helpers.TextResult(out), nil
}
