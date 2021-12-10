package steampipecloud

import (
	"context"

	openapi "github.com/turbot/steampipecloud-sdk-go"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableSteampipecloudWorkspaceConnection(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "steampipecloud_workspace_connection",
		Description: "Steampipecloud Workspace Connection",
		List: &plugin.ListConfig{
			ParentHydrate: listWorkspaces,
			Hydrate:       listWorkspaceConnections,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The unique identifier for the workspace connection association.",
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
				Description: "The creation time of the workspace connection association.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "version_id",
				Description: "The creation version id of the workspace connection association.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "updated_at",
				Description: "The last updated time of the workspace connection association.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "connection",
				Description: "The connection details of the workspace connection association.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "workspace",
				Description: "The workspace details of the workspace connection association.",
				Type:        proto.ColumnType_JSON,
			},
		},
	}
}

//// LIST FUNCTION

func listWorkspaceConnections(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	workspace := h.Item.(openapi.TypesWorkspace)

	getUserIdentityCached := plugin.HydrateFunc(getUserIdentity).WithCache()
	commonData, err := getUserIdentityCached(ctx, d, h)
	user := commonData.(openapi.TypesUser)

	if workspace.Identity.Handle == user.Handle {
		err = listUserWorkspaceConnectionAssociations(ctx, d, h, user.Handle, workspace.Handle)
	} else {
		err = listOrgWorkspaceConnectionAssociations(ctx, d, h, workspace.Identity.Handle, workspace.Handle)
	}

	if err != nil {
		plugin.Logger(ctx).Error("listWorkspaceConnections", "list", err)
		return nil, err
	}
	return nil, nil
}

func listUserWorkspaceConnectionAssociations(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, userHandle string, workspaceHandle string) error {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listUserWorkspaceConnectionAssociations", "connection_error", err)
		return err
	}

	// execute list call
	pagesLeft := true
	var resp openapi.TypesListWorkspaceConnResponse
	var listDetails func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error)

	for pagesLeft {
		if resp.NextToken != nil {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.UserWorkspaceConnectionAssociationsApi.ListUserWorkspaceConnectionAssociations(context.Background(), userHandle, workspaceHandle).NextToken(*resp.NextToken).Execute()
				return resp, err
			}
		} else {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.UserWorkspaceConnectionAssociationsApi.ListUserWorkspaceConnectionAssociations(context.Background(), userHandle, workspaceHandle).Execute()
				return resp, err
			}
		}

		response, err := plugin.RetryHydrate(ctx, d, h, listDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

		if err != nil {
			plugin.Logger(ctx).Error("listUserWorkspaceConnectionAssociations", "list", err)
			return err
		}

		result := response.(openapi.TypesListWorkspaceConnResponse)

		if result.HasItems() {
			for _, workspaceConn := range *result.Items {
				d.StreamListItem(ctx, workspaceConn)
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

func listOrgWorkspaceConnectionAssociations(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, orgHandle string, workspaceHandle string) error {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listOrgWorkspaceConnectionAssociations", "connection_error", err)
		return err
	}

	// execute list call
	pagesLeft := true
	var resp openapi.TypesListWorkspaceConnResponse
	var listDetails func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error)

	for pagesLeft {
		if resp.NextToken != nil {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.OrgWorkspaceConnectionAssociationsApi.ListOrgWorkspaceConnectionAssociations(context.Background(), orgHandle, workspaceHandle).NextToken(*resp.NextToken).Execute()
				return resp, err
			}
		} else {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.OrgWorkspaceConnectionAssociationsApi.ListOrgWorkspaceConnectionAssociations(context.Background(), orgHandle, workspaceHandle).Execute()
				return resp, err
			}
		}

		response, err := plugin.RetryHydrate(ctx, d, h, listDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

		if err != nil {
			plugin.Logger(ctx).Error("listOrgWorkspaceConnectionAssociations", "list", err)
			return err
		}

		result := response.(openapi.TypesListWorkspaceConnResponse)

		if result.HasItems() {
			for _, workspaceConn := range *result.Items {
				d.StreamListItem(ctx, workspaceConn)
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
