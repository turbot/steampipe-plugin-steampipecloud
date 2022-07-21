package steampipecloud

import (
	"context"
	"strings"

	openapi "github.com/turbot/steampipe-cloud-sdk-go"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

type OrgWorkspaceDetails struct {
	OrgHandle       string `json:"org_handle"`
	WorkspaceHandle string `json:"workspace_handle"`
}

//// TABLE DEFINITION

func tableSteampipeCloudOrganizationWorkspaceMember(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "steampipecloud_organization_workspace_member",
		Description: "Organization workspace members can collaborate and share connections and dashboards.",
		List: &plugin.ListConfig{
			ParentHydrate: listWorkspaces,
			Hydrate:       listOrganizationWorkspaceMembers,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"org_handle", "workspace_handle", "user_handle"}),
			Hydrate:    getOrganizationWorkspaceMember,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The unique identifier for the member.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "org_id",
				Description: "The unique identifier for the organization.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "org_handle",
				Description: "The handle of the organization.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getOrgWorkspaceDetails,
			},
			{
				Name:        "workspace_id",
				Description: "The unique identifier for the workspace.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "workspace_handle",
				Description: "The handle of the workspace.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getOrgWorkspaceDetails,
			},
			{
				Name:        "status",
				Description: "The member current status.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "user_id",
				Description: "The unique identifier for the user.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "user_handle",
				Description: "The handle name for the user.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "email",
				Description: "The email address for the member.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "role",
				Description: "The role of the member.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "scope",
				Description: "The scope of the role. Can be one of org / workspace.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_at",
				Description: "The member's creation time.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "updated_at",
				Description: "The member's last update time.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "created_by",
				Description: "ID of the user who invited the member.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CreatedById"),
			},
			{
				Name:        "updated_by",
				Description: "ID of the user who last updated the member.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("UpdatedById"),
			},
			{
				Name:        "version_id",
				Description: "The current version ID for the member.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromCamel(),
			},
		},
	}
}

//// LIST FUNCTION

func listOrganizationWorkspaceMembers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	workspace := h.Item.(*openapi.Workspace)

	if strings.HasPrefix(workspace.IdentityId, "u_") {
		return nil, nil
	}

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

	err := listOrgWorkspaceMembers(ctx, d, h, workspace.IdentityId, workspace.Handle, maxResults)
	if err != nil {
		plugin.Logger(ctx).Error("listOrganizationWorkspaceMembers", "error", err)
		return nil, err
	}

	if err != nil {
		plugin.Logger(ctx).Error("listOrganizationWorkspaceMembers", "error", err)
		return nil, err
	}
	return nil, nil
}

func listOrgWorkspaceMembers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, orgHandle string, workspaceHandle string, maxResults int32) error {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listOrgWorkspaceMembers", "connection_error", err)
		return err
	}

	pagesLeft := true
	var resp openapi.ListOrgWorkspaceUsersResponse
	var listDetails func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error)

	for pagesLeft {
		if resp.NextToken != nil {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.OrgWorkspaceMembers.List(ctx, orgHandle, workspaceHandle).NextToken(*resp.NextToken).Limit(maxResults).Execute()
				return resp, err
			}
		} else {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.OrgWorkspaceMembers.List(ctx, orgHandle, workspaceHandle).Limit(maxResults).Execute()
				return resp, err
			}
		}

		response, err := plugin.RetryHydrate(ctx, d, h, listDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

		if err != nil {
			plugin.Logger(ctx).Error("listOrgWorkspaceMembers", "list", err)
			return err
		}

		result := response.(openapi.ListOrgWorkspaceUsersResponse)

		if result.HasItems() {
			for _, member := range *result.Items {
				d.StreamListItem(ctx, member)

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

func getOrganizationWorkspaceMember(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	orgHandle := d.KeyColumnQuals["org_handle"].GetStringValue()
	workspaceHandle := d.KeyColumnQuals["workspace_handle"].GetStringValue()
	userhandle := d.KeyColumnQuals["user_handle"].GetStringValue()

	// check if handle or identityHandle is empty
	if orgHandle == "" || workspaceHandle == "" || userhandle == "" {
		return nil, nil
	}

	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getOrganizationWorkspaceMember", "connection_error", err)
		return nil, err
	}

	var orgWorkspaceUser openapi.OrgWorkspaceUser

	getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		orgWorkspaceUser, _, err = svc.OrgWorkspaceMembers.Get(ctx, orgHandle, workspaceHandle, userhandle).Execute()
		return orgWorkspaceUser, err
	}

	response, err := plugin.RetryHydrate(ctx, d, h, getDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

	orgWorkspaceUser = response.(openapi.OrgWorkspaceUser)

	if err != nil {
		plugin.Logger(ctx).Error("getOrganizationWorkspaceMember", "get", err)
		return nil, err
	}

	return orgWorkspaceUser, nil
}

func getOrgWorkspaceDetails(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getOrgDetails", "connection_error", err)
		return nil, err
	}

	// get workspace details from hydrate data
	// workspace details reside in the parent item in this case
	switch w := h.ParentItem.(type) {
	case openapi.Workspace:
		getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
			resp, _, err := svc.Orgs.Get(ctx, h.ParentItem.(openapi.Workspace).IdentityId).Execute()
			return resp, err
		}
		response, _ := plugin.RetryHydrate(ctx, d, h, getDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})
		return &OrgWorkspaceDetails{OrgHandle: response.(openapi.Org).Handle, WorkspaceHandle: w.Handle}, nil
	default:
		plugin.Logger(ctx).Debug("getOrgDetails", "Unknown Type", w)
	}

	getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		resp, _, err := svc.Orgs.Get(ctx, h.Item.(openapi.OrgWorkspaceUser).OrgId).Execute()
		return resp, err
	}
	response, _ := plugin.RetryHydrate(ctx, d, h, getDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

	return &OrgWorkspaceDetails{OrgHandle: response.(openapi.Org).Handle, WorkspaceHandle: h.Item.(openapi.OrgWorkspaceUser).WorkspaceHandle}, nil
}
