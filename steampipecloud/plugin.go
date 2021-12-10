package steampipecloud

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

const pluginName = "steampipe-plugin-steampipecloud"

// Plugin creates this (steampipecloud) plugin
func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name:             pluginName,
		DefaultTransform: transform.FromGo(),
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		TableMap: map[string]*plugin.Table{
			"steampipecloud_audit_log":            tableSteampipecloudAuditLog(ctx),
			"steampipecloud_connection":           tableSteampipecloudConnection(ctx),
			"steampipecloud_member":               tableSteampipecloudMember(ctx),
			"steampipecloud_org":                  tableSteampipecloudOrganization(ctx),
			"steampipecloud_token":                tableSteampipecloudToken(ctx),
			"steampipecloud_user":                 tableSteampipecloudUser(ctx),
			"steampipecloud_workspace":            tableSteampipecloudWorkspace(ctx),
			"steampipecloud_workspace_connection": tableSteampipecloudWorkspaceConnection(ctx),
		},
	}

	return p
}
