package steampipecloud

import (
	"context"
	"strings"

	openapi "github.com/turbot/steampipe-cloud-sdk-go"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

type IdentityWorkspaceDetailsForWorkspaceConn struct {
	IdentityHandle  string `json:"identity_handle"`
	IdentityType    string `json:"identity_type"`
	WorkspaceHandle string `json:"workspace_handle"`
}

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
				Name:        "identity_id",
				Description: "The unique identifier of the identity.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "identity_handle",
				Description: "The handle of the identity.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIdentityWorkspaceDetailsForWorkspaceConn,
			},
			{
				Name:        "identity_type",
				Description: "The type of identity. Can be one of 'user' or 'org'",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIdentityWorkspaceDetailsForWorkspaceConn,
			},
			{
				Name:        "workspace_id",
				Description: "The unique identifier for the workspace.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "workspace_handle",
				Description: "The handle for the workspace.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIdentityWorkspaceDetailsForWorkspaceConn,
			},
			{
				Name:        "connection_id",
				Description: "The unique identifier for the connection.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "connection_handle",
				Description: "The handle for the connection.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Connection.Handle"),
			},
			{
				Name:        "connection",
				Description: "Additional information about the connection.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "created_at",
				Description: "The time when the connection was added to the workspace.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "created_by_id",
				Description: "The unique identifier of the user who added the connection to the workspace.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "created_by",
				Description: "Information about the user who added the connection to the workspace.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "updated_at",
				Description: "The time when the association was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "updated_by_id",
				Description: "The unique identifier of the user who last updated the association.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "updated_by",
				Description: "Information about the user who last updated the association.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "version_id",
				Description: "The current version ID for the association.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromCamel(),
			},
		},
	}
}

//// LIST FUNCTION

func listWorkspaceConnections(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	workspace := h.Item.(*openapi.Workspace)

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

	if workspace.IdentityId == user.Id {
		err = listUserWorkspaceConnectionAssociations(ctx, d, h, user.Handle, workspace.Handle, svc, maxResults)
	} else {
		err = listOrgWorkspaceConnectionAssociations(ctx, d, h, workspace.IdentityId, workspace.Handle, svc, maxResults)
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
				resp, _, err = svc.UserWorkspaceConnectionAssociations.List(ctx, userHandle, workspaceHandle).NextToken(*resp.NextToken).Limit(maxResults).Execute()
				return resp, err
			}
		} else {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.UserWorkspaceConnectionAssociations.List(ctx, userHandle, workspaceHandle).Limit(maxResults).Execute()
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
				resp, _, err = svc.OrgWorkspaceConnectionAssociations.List(ctx, orgHandle, workspaceHandle).NextToken(*resp.NextToken).Limit(maxResults).Execute()
				return resp, err
			}
		} else {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.OrgWorkspaceConnectionAssociations.List(ctx, orgHandle, workspaceHandle).Limit(maxResults).Execute()
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

func getIdentityWorkspaceDetailsForWorkspaceConn(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getIdentityWorkspaceDetailsForWorkspaceConn", "connection_error", err)
		return nil, err
	}

	var identityWorkspaceDetails IdentityWorkspaceDetailsForWorkspaceConn
	// get workspace details from hydrate data
	// workspace details reside in the parent item in this case
	switch w := h.ParentItem.(type) {
	case openapi.Workspace:
		identityId := h.ParentItem.(openapi.Workspace).IdentityId
		identityWorkspaceDetails.WorkspaceHandle = h.ParentItem.(openapi.Workspace).Handle
		getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
			if strings.HasPrefix(identityId, "u_") {
				resp, _, err := svc.Users.Get(ctx, identityId).Execute()
				identityWorkspaceDetails.IdentityType = "user"
				identityWorkspaceDetails.IdentityHandle = resp.Handle
				return nil, err
			} else {
				resp, _, err := svc.Orgs.Get(ctx, identityId).Execute()
				identityWorkspaceDetails.IdentityType = "org"
				identityWorkspaceDetails.IdentityHandle = resp.Handle
				return nil, err
			}
		}
		_, _ = plugin.RetryHydrate(ctx, d, h, getDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})
		return identityWorkspaceDetails, nil
	default:
		plugin.Logger(ctx).Debug("getIdentityWorkspaceDetailsForWorkspaceConn", "Unknown Type", w)
	}

	identityId := h.Item.(openapi.WorkspaceConn).IdentityId
	workspaceId := h.Item.(openapi.WorkspaceConn).WorkspaceId
	getIdentityDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		if strings.HasPrefix(identityId, "u_") {
			resp, _, err := svc.Users.Get(ctx, identityId).Execute()
			identityWorkspaceDetails.IdentityType = "user"
			identityWorkspaceDetails.IdentityHandle = resp.Handle
			return nil, err
		} else {
			resp, _, err := svc.Orgs.Get(ctx, identityId).Execute()
			identityWorkspaceDetails.IdentityType = "org"
			identityWorkspaceDetails.IdentityHandle = resp.Handle
			return nil, err
		}
	}
	_, _ = plugin.RetryHydrate(ctx, d, h, getIdentityDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

	getWorkspaceDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		if strings.HasPrefix(identityId, "u_") {
			resp, _, err := svc.UserWorkspaces.Get(ctx, identityId, workspaceId).Execute()
			identityWorkspaceDetails.WorkspaceHandle = resp.Handle
			return nil, err
		} else {
			resp, _, err := svc.OrgWorkspaces.Get(ctx, identityId, workspaceId).Execute()
			identityWorkspaceDetails.WorkspaceHandle = resp.Handle
			return nil, err
		}
	}
	_, _ = plugin.RetryHydrate(ctx, d, h, getWorkspaceDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

	return identityWorkspaceDetails, nil
}
