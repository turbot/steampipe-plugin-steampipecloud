package steampipecloud

import (
	"context"
	"strings"

	openapi "github.com/turbot/steampipe-cloud-sdk-go"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type IdentityWorkspaceDetailsForWorkspaceMod struct {
	IdentityHandle  string `json:"identity_handle"`
	IdentityType    string `json:"identity_type"`
	WorkspaceHandle string `json:"workspace_handle"`
}

//// TABLE DEFINITION

func tableSteampipeCloudWorkspaceMod(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "steampipecloud_workspace_mod",
		Description: "A Steampipe mod is a portable, versioned collection of related Steampipe resources such as dashboards, benchmarks, queries, and controls.",
		List: &plugin.ListConfig{
			ParentHydrate: listWorkspaces,
			Hydrate:       listWorkspaceMods,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"identity_id", "workspace_id", "alias"}),
			Hydrate:    getWorkspaceMod,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The unique identifier for the workspace mod.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "identity_id",
				Description: "The unique identifier for the identity which contains the workspace.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "identity_handle",
				Description: "The handle of the identity which contains the workspace.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIdentityWorkspaceDetailsForWorkspaceMod,
			},
			{
				Name:        "identity_type",
				Description: "The type of identity, which can be 'user' or 'org'.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIdentityWorkspaceDetailsForWorkspaceMod,
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
				Hydrate:     getIdentityWorkspaceDetailsForWorkspaceMod,
			},
			{
				Name:        "constraint",
				Description: "Version constraint for the mod.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "alias",
				Description: "Short name used to identify the mod.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "installed_version",
				Description: "Version of the mod installed.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "state",
				Description: "State of the mod. Can be one of 'installing', 'installed' or 'error'.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "path",
				Description: "Full path name for the mod.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "details",
				Description: "Extra stored details about the mod.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "created_at",
				Description: "The time when the mod was installed.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "created_by_id",
				Description: "The unique identifier of the user who installed the mod.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "created_by",
				Description: "Information about the user who installed the mod.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "updated_at",
				Description: "The time when the mod was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "updated_by_id",
				Description: "The unique identifier of the user who last updated the mod.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "updated_by",
				Description: "Information about the user who last updated the mod.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "version_id",
				Description: "The current version ID of the mod record.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromCamel(),
			},
		},
	}
}

//// LIST FUNCTION

func listWorkspaceMods(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	workspace := h.Item.(*openapi.Workspace)

	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listWorkspaceMods", "connection_error", err)
		return nil, err
	}

	getUserIdentityCached := plugin.HydrateFunc(getUserIdentity).WithCache()
	commonData, err := getUserIdentityCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("listWorkspaceMods", "getUserIdentityCached", err)
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
		err = listUserWorkspaceMods(ctx, d, h, workspace.IdentityId, workspace.Id, svc, maxResults)
	} else {
		err = listOrgWorkspaceMods(ctx, d, h, workspace.IdentityId, workspace.Id, svc, maxResults)
	}

	if err != nil {
		plugin.Logger(ctx).Error("listWorkspaceMods", "list", err)
		return nil, err
	}
	return nil, nil
}

func listUserWorkspaceMods(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, userHandle string, workspaceHandle string, svc *openapi.APIClient, maxResults int32) error {
	var err error

	// execute list call
	pagesLeft := true
	var resp openapi.ListWorkspaceModsResponse
	var listDetails func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error)

	for pagesLeft {
		if resp.NextToken != nil {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.UserWorkspaceMods.List(ctx, userHandle, workspaceHandle).NextToken(*resp.NextToken).Limit(maxResults).Execute()
				return resp, err
			}
		} else {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.UserWorkspaceMods.List(ctx, userHandle, workspaceHandle).Limit(maxResults).Execute()
				return resp, err
			}
		}

		response, err := plugin.RetryHydrate(ctx, d, h, listDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

		if err != nil {
			plugin.Logger(ctx).Error("listUserWorkspaceMods", "list", err)
			return err
		}

		result := response.(openapi.ListWorkspaceModsResponse)

		if result.HasItems() {
			for _, workspaceMod := range *result.Items {
				d.StreamListItem(ctx, workspaceMod)

				// Context can be cancelled due to manual cancellation or the limit has been hit
				if d.RowsRemaining(ctx) == 0 {
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

func listOrgWorkspaceMods(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, orgHandle string, workspaceHandle string, svc *openapi.APIClient, maxResults int32) error {
	var err error

	// execute list call
	pagesLeft := true
	var resp openapi.ListWorkspaceModsResponse
	var listDetails func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error)

	for pagesLeft {
		if resp.NextToken != nil {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.OrgWorkspaceMods.List(ctx, orgHandle, workspaceHandle).NextToken(*resp.NextToken).Limit(maxResults).Execute()
				return resp, err
			}
		} else {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.OrgWorkspaceMods.List(ctx, orgHandle, workspaceHandle).Limit(maxResults).Execute()
				return resp, err
			}
		}

		response, err := plugin.RetryHydrate(ctx, d, h, listDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

		if err != nil {
			plugin.Logger(ctx).Error("listUserWorkspaceMods", "list", err)
			return err
		}

		result := response.(openapi.ListWorkspaceModsResponse)

		if result.HasItems() {
			for _, workspaceMod := range *result.Items {
				d.StreamListItem(ctx, workspaceMod)

				// Context can be cancelled due to manual cancellation or the limit has been hit
				if d.RowsRemaining(ctx) == 0 {
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

//// GET FUNCTION

func getWorkspaceMod(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	identityId := d.EqualsQuals["identity_id"].GetStringValue()
	workspaceId := d.EqualsQuals["workspace_id"].GetStringValue()
	alias := d.EqualsQuals["alias"].GetStringValue()

	// check if identity or workspace or alias information is missing
	if identityId == "" || workspaceId == "" || alias == "" {
		return nil, nil
	}

	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getWorkspaceMod", "connection_error", err)
		return nil, err
	}

	getUserIdentityCached := plugin.HydrateFunc(getUserIdentity).WithCache()
	commonData, err := getUserIdentityCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("getWorkspaceMod", "getUserIdentityCached", err)
		return nil, err
	}

	user := commonData.(openapi.User)
	var resp interface{}

	if identityId == user.Id {
		resp, err = getUserWorkspaceMod(ctx, d, h, identityId, workspaceId, alias, svc)
	} else {
		resp, err = getOrgWorkspaceMod(ctx, d, h, identityId, workspaceId, alias, svc)
	}

	if err != nil {
		plugin.Logger(ctx).Error("getWorkspaceMod", "get", err)
		return nil, err
	}

	if resp == nil {
		return nil, nil
	}

	return resp.(openapi.WorkspaceMod), nil
}

func getUserWorkspaceMod(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, identityId string, workspaceId string, alias string, svc *openapi.APIClient) (interface{}, error) {
	var err error

	// execute get call
	var resp openapi.WorkspaceMod

	getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		resp, _, err = svc.UserWorkspaceMods.Get(ctx, identityId, workspaceId, alias).Execute()
		return resp, err
	}

	response, err := plugin.RetryHydrate(ctx, d, h, getDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

	workspaceMod := response.(openapi.WorkspaceMod)

	if err != nil {
		plugin.Logger(ctx).Error("getUserWorkspaceMod", "get", err)
		return nil, err
	}

	return workspaceMod, nil
}

func getOrgWorkspaceMod(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, identityId string, workspaceId string, alias string, svc *openapi.APIClient) (interface{}, error) {
	var err error

	// execute get call
	var resp openapi.WorkspaceMod

	getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		resp, _, err = svc.OrgWorkspaceMods.Get(ctx, identityId, workspaceId, alias).Execute()
		return resp, err
	}

	response, err := plugin.RetryHydrate(ctx, d, h, getDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

	workspaceMod := response.(openapi.WorkspaceMod)

	if err != nil {
		plugin.Logger(ctx).Error("getOrgWorkspaceMod", "get", err)
		return nil, err
	}

	return workspaceMod, nil
}

func getIdentityWorkspaceDetailsForWorkspaceMod(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getIdentityWorkspaceDetailsForWorkspaceMod", "connection_error", err)
		return nil, err
	}

	var identityWorkspaceDetails IdentityWorkspaceDetailsForWorkspaceMod
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
		plugin.Logger(ctx).Debug("getIdentityWorkspaceDetailsForWorkspaceMod", "Unknown Type", w)
	}

	identityId := h.Item.(openapi.WorkspaceMod).IdentityId
	workspaceId := h.Item.(openapi.WorkspaceMod).WorkspaceId
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
