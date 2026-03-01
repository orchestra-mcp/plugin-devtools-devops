package devtoolsdevops

import (
	"github.com/orchestra-mcp/plugin-devtools-devops/internal"
	"github.com/orchestra-mcp/sdk-go/plugin"
)

// Register adds all DevOps tools to the builder.
func Register(builder *plugin.PluginBuilder) {
	tp := &internal.ToolsPlugin{}
	tp.RegisterTools(builder)
}
