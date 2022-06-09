package steampipecloud

import (
	"context"

	openapi "github.com/turbot/steampipe-cloud-sdk-go"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func tableSteampipeCloudWorkspaceConnection(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "steampipecloud_workspace_connection",
		Description: "Workspace connections are the associations between workspaces and connections.",
		List: &plugin.ListConfig{
			ParentHydrate: listWorkspaces,
			Hydrate:       listWorkspaceConnections,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The unique identifier for the association.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "connection_id",
				Description: "The unique identifier for the connection.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "workspace_id",
				Description: "The unique identifier for the workspace.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "identity_id",
				Description: "The unique identifier for an identity where the action has been performed.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "created_at",
				Description: "The creation time of the association.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "version_id",
				Description: "The version ID of the association.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "updated_at",
				Description: "The association's last updated time.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
		},
	}
}

//// LIST FUNCTION

func listWorkspaceConnections(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	workspace := h.Item.(openapi.Workspace)

	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listWorkspaceConnections", "connection_error", err)
		return nil, err
	}

	getUserIdentityCached := plugin.HydrateFunc(getUserIdentity).WithCache()
	commonData, err := getUserIdentityCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("listWorkspaceConnections", "getUserIdentityCached", err)
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

	if workspace.Identity.Handle == user.Handle {
		err = listUserWorkspaceConnectionAssociations(ctx, d, h, user.Handle, workspace.Handle, svc, maxResults)
	} else {
		err = listOrgWorkspaceConnectionAssociations(ctx, d, h, workspace.Identity.Handle, workspace.Handle, svc, maxResults)
	}

	if err != nil {
		plugin.Logger(ctx).Error("listWorkspaceConnections", "list", err)
		return nil, err
	}
	return nil, nil
}

func listUserWorkspaceConnectionAssociations(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, userHandle string, workspaceHandle string, svc *openapi.APIClient, maxResults int32) error {
	var err error

	// execute list call
	pagesLeft := true
	var resp openapi.ListWorkspaceConnResponse
	var listDetails func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error)

	for pagesLeft {
		if resp.NextToken != nil {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.UserWorkspaceConnectionAssociations.List(context.Background(), userHandle, workspaceHandle).NextToken(*resp.NextToken).Limit(maxResults).Execute()
				return resp, err
			}
		} else {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.UserWorkspaceConnectionAssociations.List(context.Background(), userHandle, workspaceHandle).Limit(maxResults).Execute()
				return resp, err
			}
		}

		response, err := plugin.RetryHydrate(ctx, d, h, listDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

		if err != nil {
			plugin.Logger(ctx).Error("listUserWorkspaceConnectionAssociations", "list", err)
			return err
		}

		result := response.(openapi.ListWorkspaceConnResponse)

		if result.HasItems() {
			for _, workspaceConn := range *result.Items {
				d.StreamListItem(ctx, workspaceConn)

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

func listOrgWorkspaceConnectionAssociations(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, orgHandle string, workspaceHandle string, svc *openapi.APIClient, maxResults int32) error {
	var err error

	// execute list call
	pagesLeft := true
	var resp openapi.ListWorkspaceConnResponse
	var listDetails func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error)

	for pagesLeft {
		if resp.NextToken != nil {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.OrgWorkspaceConnectionAssociations.List(context.Background(), orgHandle, workspaceHandle).NextToken(*resp.NextToken).Limit(maxResults).Execute()
				return resp, err
			}
		} else {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.OrgWorkspaceConnectionAssociations.List(context.Background(), orgHandle, workspaceHandle).Limit(maxResults).Execute()
				return resp, err
			}
		}

		response, err := plugin.RetryHydrate(ctx, d, h, listDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

		if err != nil {
			plugin.Logger(ctx).Error("listOrgWorkspaceConnectionAssociations", "list", err)
			return err
		}

		result := response.(openapi.ListWorkspaceConnResponse)

		if result.HasItems() {
			for _, workspaceConn := range *result.Items {
				d.StreamListItem(ctx, workspaceConn)

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
