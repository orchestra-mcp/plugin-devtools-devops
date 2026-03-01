package internal

import (
	"github.com/orchestra-mcp/sdk-go/plugin"
	"github.com/orchestra-mcp/plugin-devtools-devops/internal/tools"
)

// ToolsPlugin registers all DevOps tools.
type ToolsPlugin struct{}

// RegisterTools registers all 8 DevOps tools with the plugin builder.
func (tp *ToolsPlugin) RegisterTools(builder *plugin.PluginBuilder) {
	builder.RegisterTool("devops_list_pipelines",
		"List GitHub Actions workflows for a repository",
		tools.DevopsListPipelinesSchema(), tools.DevopsListPipelines())

	builder.RegisterTool("devops_trigger_pipeline",
		"Trigger a GitHub Actions workflow dispatch event",
		tools.DevopsTriggerPipelineSchema(), tools.DevopsTriggerPipeline())

	builder.RegisterTool("devops_pipeline_status",
		"Get the status and conclusion of a GitHub Actions workflow run",
		tools.DevopsPipelineStatusSchema(), tools.DevopsPipelineStatus())

	builder.RegisterTool("devops_pipeline_logs",
		"Fetch the logs for a GitHub Actions workflow run",
		tools.DevopsPipelineLogsSchema(), tools.DevopsPipelineLogs())

	builder.RegisterTool("devops_list_deployments",
		"List deployments or releases for a GitHub repository",
		tools.DevopsListDeploymentsSchema(), tools.DevopsListDeployments())

	builder.RegisterTool("devops_env_vars",
		"List GitHub repository secrets (names only — values are never exposed)",
		tools.DevopsEnvVarsSchema(), tools.DevopsEnvVars())

	builder.RegisterTool("devops_deploy",
		"Trigger a deployment via GitHub Actions workflow dispatch",
		tools.DevopsDeploySchema(), tools.DevopsDeploy())

	builder.RegisterTool("devops_rollback",
		"Roll back by re-running a previous GitHub Actions workflow run",
		tools.DevopsRollbackSchema(), tools.DevopsRollback())
}
