package steampipecloud

import (
	"context"

	openapi "github.com/turbot/steampipe-cloud-sdk-go"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func tableSteampipeCloudWorkspaceModVariable(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "steampipecloud_workspace_mod_variable",
		Description: "Variables are module level objects that allow you to pass values to your module at runtime.",
		List: &plugin.ListConfig{
			ParentHydrate: listWorkspaces,
			Hydrate:       listWorkspaceModVariables,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name: "workspace_id",
				},
				{
					Name: "mod_alias",
				},
			},
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The unique identifier for the workspace mod variable.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "workspace_id",
				Description: "The identifier of the workspace to which the variable belongs.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("workspace_id"),
			},
			{
				Name:        "mod_alias",
				Description: "The alias of the mod to which the variable belongs to.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "description",
				Description: "Description of the variable.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "name",
				Description: "Name of the variable.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "value_default",
				Description: "Default Value of the variable.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "value_setting",
				Description: "An explicit setting defined for the variable.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "value",
				Description: "Winning Value of the variable.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "type",
				Description: "Type of value expected by the variable.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "created_at",
				Description: "Time when the mod variable was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "created_by_id",
				Description: "Unique identifier of the user who created the setting.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "created_by",
				Description: "Information about the user who created the Setting.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "updated_at",
				Description: "Time when the mod variable was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "updated_by_id",
				Description: "Unique identifier of the user who last updated the setting.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "updated_by",
				Description: "Information about the user who updated the Setting.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "version_id",
				Description: "The current version ID of the variable.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromCamel(),
			},
		},
	}
}

//// LIST FUNCTION

func listWorkspaceModVariables(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	workspaceId := d.KeyColumnQuals["workspace_id"].GetStringValue()
	modId := d.KeyColumnQuals["mod_alias"].GetStringValue()

	// If key qual columns are not mentioned, exit
	if workspaceId == "" || modId == "" {
		return nil, nil
	}

	workspace := h.Item.(*openapi.Workspace)

	if workspace.Id != workspaceId {
		return nil, nil
	}

	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listWorkspaceModVariables", "connection_error", err)
		return nil, err
	}

	getUserIdentityCached := plugin.HydrateFunc(getUserIdentity).WithCache()
	commonData, err := getUserIdentityCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("listWorkspaceModVariables", "getUserIdentityCached", err)
		return nil, err
	}

	user := commonData.(openapi.User)

	// If the requested number of items is less than the paging max limit
	// set the limit to that instead
	maxResults := int32(100)
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(maxResults) {
			if *limit < 1 {
				maxResults = int32(1)
			} else {
				maxResults = int32(*limit)
			}
		}
	}

	if workspace.IdentityId == user.Id {
		err = listUserWorkspaceModVariables(ctx, d, h, workspace.IdentityId, workspace.Id, modId, svc, maxResults)
	} else {
		err = listOrgWorkspaceModVariables(ctx, d, h, workspace.IdentityId, workspace.Id, modId, svc, maxResults)
	}

	if err != nil {
		plugin.Logger(ctx).Error("listWorkspaceModVariables", "list", err)
		return nil, err
	}
	return nil, nil
}

func listUserWorkspaceModVariables(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, userHandle string, workspaceHandle string, modAlias string, svc *openapi.APIClient, maxResults int32) error {
	var err error

	// execute list call
	pagesLeft := true
	var resp openapi.ListWorkspaceModVariablesResponse
	var listDetails func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error)

	for pagesLeft {
		if resp.NextToken != nil {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.UserWorkspaceModVariables.List(ctx, userHandle, workspaceHandle, modAlias).NextToken(*resp.NextToken).Limit(maxResults).Execute()
				return resp, err
			}
		} else {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.UserWorkspaceModVariables.List(ctx, userHandle, workspaceHandle, modAlias).Limit(maxResults).Execute()
				return resp, err
			}
		}

		response, err := plugin.RetryHydrate(ctx, d, h, listDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

		if err != nil {
			plugin.Logger(ctx).Error("listUserWorkspaceModVariables", "list", err)
			return err
		}

		result := response.(openapi.ListWorkspaceModVariablesResponse)

		if result.HasItems() {
			for _, workspaceMod := range *result.Items {
				d.StreamListItem(ctx, workspaceMod)

				// Context can be cancelled due to manual cancellation or the limit has been hit
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return nil
				}
			}
		}
		if result.NextToken == nil {
			pagesLeft = false
		} else {
			resp.NextToken = result.NextToken
		}
	}

	return nil
}

func listOrgWorkspaceModVariables(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, orgHandle string, workspaceHandle string, modAlias string, svc *openapi.APIClient, maxResults int32) error {
	var err error

	// execute list call
	pagesLeft := true
	var resp openapi.ListWorkspaceModVariablesResponse
	var listDetails func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error)

	for pagesLeft {
		if resp.NextToken != nil {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.OrgWorkspaceModVariables.List(ctx, orgHandle, workspaceHandle, modAlias).NextToken(*resp.NextToken).Limit(maxResults).Execute()
				return resp, err
			}
		} else {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.OrgWorkspaceModVariables.List(ctx, orgHandle, workspaceHandle, modAlias).Limit(maxResults).Execute()
				return resp, err
			}
		}

		response, err := plugin.RetryHydrate(ctx, d, h, listDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

		if err != nil {
			plugin.Logger(ctx).Error("listUserWorkspaceModVariables", "list", err)
			return err
		}

		result := response.(openapi.ListWorkspaceModVariablesResponse)

		if result.HasItems() {
			for _, workspaceMod := range *result.Items {
				d.StreamListItem(ctx, workspaceMod)

				// Context can be cancelled due to manual cancellation or the limit has been hit
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return nil
				}
			}
		}
		if result.NextToken == nil {
			pagesLeft = false
		} else {
			resp.NextToken = result.NextToken
		}
	}

	return nil
}
