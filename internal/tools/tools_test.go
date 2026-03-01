package tools

import (
	"context"
	"os/exec"
	"testing"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

// ghAvailable reports whether the gh CLI is installed and authenticated.
func ghAvailable() bool {
	cmd := exec.Command("gh", "auth", "status")
	return cmd.Run() == nil
}

// callTool invokes a handler with the given string args map and returns the response.
func callTool(
	t *testing.T,
	handler func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error),
	args map[string]any,
) *pluginv1.ToolResponse {
	t.Helper()
	s, err := structpb.NewStruct(args)
	if err != nil {
		t.Fatalf("failed to build args struct: %v", err)
	}
	req := &pluginv1.ToolRequest{Arguments: s}
	resp, err := handler(context.Background(), req)
	if err != nil {
		t.Fatalf("handler returned unexpected error: %v", err)
	}
	return resp
}

// isError reports whether the response carries an error.
func isError(resp *pluginv1.ToolResponse) bool {
	return resp != nil && !resp.Success
}

// errorCode extracts the error_code field from a failed response.
func errorCode(resp *pluginv1.ToolResponse) string {
	return resp.GetErrorCode()
}

// ---------------------------------------------------------------------------
// devops_list_pipelines
// ---------------------------------------------------------------------------

func TestDevopsListPipelines_NoGh(t *testing.T) {
	if !ghAvailable() {
		t.Skip("gh CLI not authenticated")
	}
	resp := callTool(t, DevopsListPipelines(), map[string]any{})
	// Either a successful listing or a gh_error (no workflows in current dir) is acceptable.
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
}

// ---------------------------------------------------------------------------
// devops_trigger_pipeline
// ---------------------------------------------------------------------------

func TestDevopsTriggerPipeline_MissingWorkflow(t *testing.T) {
	resp := callTool(t, DevopsTriggerPipeline(), map[string]any{})
	if !isError(resp) {
		t.Fatalf("expected error response, got success")
	}
	if errorCode(resp) != "validation_error" {
		t.Fatalf("expected validation_error, got %q", errorCode(resp))
	}
}

func TestDevopsTriggerPipeline_NoGh(t *testing.T) {
	if !ghAvailable() {
		t.Skip("gh CLI not authenticated")
	}
	resp := callTool(t, DevopsTriggerPipeline(), map[string]any{
		"workflow": "ci.yml",
		"repo":     "nonexistent/repo",
	})
	if !isError(resp) {
		t.Fatalf("expected gh_error for nonexistent repo, got success")
	}
	if errorCode(resp) != "gh_error" {
		t.Fatalf("expected gh_error, got %q", errorCode(resp))
	}
}

// ---------------------------------------------------------------------------
// devops_pipeline_status
// ---------------------------------------------------------------------------

func TestDevopsPipelineStatus_MissingRunID(t *testing.T) {
	resp := callTool(t, DevopsPipelineStatus(), map[string]any{})
	if !isError(resp) {
		t.Fatalf("expected error response, got success")
	}
	if errorCode(resp) != "validation_error" {
		t.Fatalf("expected validation_error, got %q", errorCode(resp))
	}
}

func TestDevopsPipelineStatus_NoGh(t *testing.T) {
	if !ghAvailable() {
		t.Skip("gh CLI not authenticated")
	}
	resp := callTool(t, DevopsPipelineStatus(), map[string]any{
		"run_id": "0",
		"repo":   "nonexistent/repo",
	})
	if !isError(resp) {
		t.Fatalf("expected gh_error for nonexistent repo/run, got success")
	}
	if errorCode(resp) != "gh_error" {
		t.Fatalf("expected gh_error, got %q", errorCode(resp))
	}
}

// ---------------------------------------------------------------------------
// devops_pipeline_logs
// ---------------------------------------------------------------------------

func TestDevopsPipelineLogs_MissingRunID(t *testing.T) {
	resp := callTool(t, DevopsPipelineLogs(), map[string]any{})
	if !isError(resp) {
		t.Fatalf("expected error response, got success")
	}
	if errorCode(resp) != "validation_error" {
		t.Fatalf("expected validation_error, got %q", errorCode(resp))
	}
}

func TestDevopsPipelineLogs_NoGh(t *testing.T) {
	if !ghAvailable() {
		t.Skip("gh CLI not authenticated")
	}
	resp := callTool(t, DevopsPipelineLogs(), map[string]any{
		"run_id": "0",
		"repo":   "nonexistent/repo",
	})
	// gh run view --log on a nonexistent repo/run should return gh_error or success.
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
}

// ---------------------------------------------------------------------------
// devops_list_deployments
// ---------------------------------------------------------------------------

func TestDevopsListDeployments_NoGh(t *testing.T) {
	if !ghAvailable() {
		t.Skip("gh CLI not authenticated")
	}
	resp := callTool(t, DevopsListDeployments(), map[string]any{})
	// Either a successful listing or a gh_error is acceptable.
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
}

// ---------------------------------------------------------------------------
// devops_env_vars
// ---------------------------------------------------------------------------

func TestDevopsEnvVars_NoGh(t *testing.T) {
	if !ghAvailable() {
		t.Skip("gh CLI not authenticated")
	}
	resp := callTool(t, DevopsEnvVars(), map[string]any{})
	// Either a successful listing or a gh_error is acceptable.
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
}

// ---------------------------------------------------------------------------
// devops_deploy
// ---------------------------------------------------------------------------

func TestDevopsDeploy_MissingWorkflow(t *testing.T) {
	resp := callTool(t, DevopsDeploy(), map[string]any{})
	if !isError(resp) {
		t.Fatalf("expected error response, got success")
	}
	if errorCode(resp) != "validation_error" {
		t.Fatalf("expected validation_error, got %q", errorCode(resp))
	}
}

func TestDevopsDeploy_NoGh(t *testing.T) {
	if !ghAvailable() {
		t.Skip("gh CLI not authenticated")
	}
	resp := callTool(t, DevopsDeploy(), map[string]any{
		"workflow": "deploy.yml",
		"repo":     "nonexistent/repo",
	})
	if !isError(resp) {
		t.Fatalf("expected gh_error for nonexistent repo, got success")
	}
	if errorCode(resp) != "gh_error" {
		t.Fatalf("expected gh_error, got %q", errorCode(resp))
	}
}

// ---------------------------------------------------------------------------
// devops_rollback
// ---------------------------------------------------------------------------

func TestDevopsRollback_MissingRunID(t *testing.T) {
	resp := callTool(t, DevopsRollback(), map[string]any{})
	if !isError(resp) {
		t.Fatalf("expected error response, got success")
	}
	if errorCode(resp) != "validation_error" {
		t.Fatalf("expected validation_error, got %q", errorCode(resp))
	}
}

func TestDevopsRollback_NoGh(t *testing.T) {
	if !ghAvailable() {
		t.Skip("gh CLI not authenticated")
	}
	resp := callTool(t, DevopsRollback(), map[string]any{
		"run_id": "0",
		"repo":   "nonexistent/repo",
	})
	if !isError(resp) {
		t.Fatalf("expected gh_error for nonexistent repo/run, got success")
	}
	if errorCode(resp) != "gh_error" {
		t.Fatalf("expected gh_error, got %q", errorCode(resp))
	}
}
