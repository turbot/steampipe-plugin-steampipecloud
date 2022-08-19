package steampipecloud

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
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
			"steampipecloud_audit_log":                     tableSteampipeCloudAuditLog(ctx),
			"steampipecloud_connection":                    tableSteampipeCloudConnection(ctx),
			"steampipecloud_organization_member":           tableSteampipeCloudOrganizationMember(ctx),
			"steampipecloud_organization":                  tableSteampipeCloudOrganization(ctx),
			"steampipecloud_organization_workspace_member": tableSteampipeCloudOrganizationWorkspaceMember(ctx),
			"steampipecloud_token":                         tableSteampipeCloudToken(ctx),
			"steampipecloud_user":                          tableSteampipeCloudUser(ctx),
			"steampipecloud_workspace":                     tableSteampipeCloudWorkspace(ctx),
			"steampipecloud_workspace_connection":          tableSteampipeCloudWorkspaceConnection(ctx),
			"steampipecloud_workspace_mod":                 tableSteampipeCloudWorkspaceMod(ctx),
			"steampipecloud_workspace_mod_variable":        tableSteampipeCloudWorkspaceModVariable(ctx),
			"steampipecloud_workspace_db_log":              tableSteampipeCloudWorkspaceDBLog(ctx),
		},
	}

	return p
}
