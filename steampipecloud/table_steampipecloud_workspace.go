package steampipecloud

import (
	"context"
	"strings"

	openapi "github.com/turbot/steampipe-cloud-sdk-go"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func tableSteampipeCloudWorkspace(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "steampipecloud_workspace",
		Description: "Workspaces provide a bounded context for managing and securing Steampipe resources.",
		List: &plugin.ListConfig{
			Hydrate: listWorkspaces,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "identity_handle",
					Require: plugin.Optional,
				},
				{
					Name:    "identity_id",
					Require: plugin.Optional,
				},
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"handle", "identity_handle"}),
			Hydrate:    getWorkspace,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The unique identifier for the workspace.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "workspace_state",
				Description: "The current workspace state.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "handle",
				Description: "The handle name for the workspace.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "database_name",
				Description: "The database name for the workspace.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_at",
				Description: "The creation time of the workspace.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "hive",
				Description: "The database hive for this workspace.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "identity_id",
				Description: "The unique identifier for an identity where the workspace has been created.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "identity_handle",
				Description: "The handle name for an identity where the workspace has been created.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Identity.Handle"),
			},
			{
				Name:        "identity_type",
				Description: "The type of identity, which can be 'user' or 'org'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Identity.Type"),
			},
			{
				Name:        "version_id",
				Description: "The version ID of the workspace.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "updated_at",
				Description: "The workspace's last updated time.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "identity",
				Description: "Information about the identity.",
				Type:        proto.ColumnType_JSON,
			},
		},
	}
}

//// LIST FUNCTION

func listWorkspaces(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listWorkspaces", "connection_error", err)
		return nil, err
	}

	getUserIdentityCached := plugin.HydrateFunc(getUserIdentity).WithCache()
	commonData, err := getUserIdentityCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("listWorkspaces", "getUserIdentityCached", err)
		return nil, err
	}

	user := commonData.(openapi.User)

	identityHandle := d.KeyColumnQuals["identity_handle"].GetStringValue()
	identityId := d.KeyColumnQuals["identity_id"].GetStringValue()

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

	if identityHandle == "" && identityId == "" {
		err = listActorWorkspaces(ctx, d, h, svc, maxResults)
	} else if identityId != "" && strings.HasPrefix(identityId, "u_") {
		err = listUserWorkspaces(ctx, d, h, identityId, svc, maxResults)
	} else if identityId != "" && strings.HasPrefix(identityId, "o_") {
		err = listOrgWorkspaces(ctx, d, h, identityId, svc, maxResults)
	} else if identityHandle == user.Handle {
		err = listUserWorkspaces(ctx, d, h, identityHandle, svc, maxResults)
	} else {
		err = listOrgWorkspaces(ctx, d, h, identityHandle, svc, maxResults)
	}

	if err != nil {
		plugin.Logger(ctx).Error("listWorkspaces", "list", err)
		return nil, err
	}
	return nil, nil
}

func listUserWorkspaces(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, handle string, svc *openapi.APIClient, maxResults int32) error {
	var err error

	// execute list call
	pagesLeft := true

	var resp openapi.ListWorkspacesResponse
	var listDetails func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error)

	for pagesLeft {
		if resp.NextToken != nil {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.UserWorkspaces.List(context.Background(), handle).NextToken(*resp.NextToken).Limit(maxResults).Execute()
				return resp, err
			}
		} else {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.UserWorkspaces.List(context.Background(), handle).Limit(maxResults).Execute()
				return resp, err
			}
		}

		response, err := plugin.RetryHydrate(ctx, d, h, listDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

		if err != nil {
			plugin.Logger(ctx).Error("listUserWorkspaces", "list", err)
			return err
		}

		result := response.(openapi.ListWorkspacesResponse)

		if result.HasItems() {
			for _, workspace := range *result.Items {
				d.StreamListItem(ctx, workspace)

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

func listOrgWorkspaces(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, handle string, svc *openapi.APIClient, maxResults int32) error {
	var err error

	// execute list call
	pagesLeft := true

	var resp openapi.ListWorkspacesResponse
	var listDetails func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error)

	for pagesLeft {
		if resp.NextToken != nil {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.OrgWorkspaces.List(context.Background(), handle).NextToken(*resp.NextToken).Limit(maxResults).Execute()
				return resp, err
			}
		} else {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.OrgWorkspaces.List(context.Background(), handle).Limit(maxResults).Execute()
				return resp, err
			}
		}

		response, err := plugin.RetryHydrate(ctx, d, h, listDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

		if err != nil {
			plugin.Logger(ctx).Error("listOrgWorkspaces", "list", err)
			return err
		}

		result := response.(openapi.ListWorkspacesResponse)

		if result.HasItems() {
			for _, workspace := range *result.Items {
				d.StreamListItem(ctx, workspace)

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

func listActorWorkspaces(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, svc *openapi.APIClient, maxResults int32) error {
	var err error

	// execute list call
	pagesLeft := true

	var resp openapi.ListWorkspacesResponse
	var listDetails func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error)

	for pagesLeft {
		if resp.NextToken != nil {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.Actors.ListWorkspaces(context.Background()).NextToken(*resp.NextToken).Limit(maxResults).Execute()
				return resp, err
			}
		} else {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.Actors.ListWorkspaces(context.Background()).Limit(maxResults).Execute()
				return resp, err
			}
		}

		response, err := plugin.RetryHydrate(ctx, d, h, listDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

		if err != nil {
			plugin.Logger(ctx).Error("listActorWorkspaces", "list", err)
			return err
		}

		result := response.(openapi.ListWorkspacesResponse)

		if result.HasItems() {
			for _, workspace := range *result.Items {
				d.StreamListItem(ctx, workspace)

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

//// HYDRATE FUNCTIONS

func getWorkspace(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	identityHandle := d.KeyColumnQuals["identity_handle"].GetStringValue()
	handle := d.KeyColumnQuals["handle"].GetStringValue()

	// check if handle or identityHandle is empty
	if identityHandle == "" || handle == "" {
		return nil, nil
	}

	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getWorkspace", "connection_error", err)
		return nil, err
	}

	getUserIdentityCached := plugin.HydrateFunc(getUserIdentity).WithCache()
	commonData, err := getUserIdentityCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("getWorkspace", "getUserIdentityCached", err)
		return nil, err
	}

	user := commonData.(openapi.User)
	var resp interface{}

	if identityHandle == user.Handle {
		resp, err = getUserWorkspace(ctx, d, h, identityHandle, handle, svc)
	} else {
		resp, err = getOrgWorkspace(ctx, d, h, identityHandle, handle, svc)
	}

	if err != nil {
		plugin.Logger(ctx).Error("getWorkspace", "get", err)
		return nil, err
	}

	if resp == nil {
		return nil, nil
	}

	return resp.(openapi.Workspace), nil
}

func getOrgWorkspace(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, identityHandle string, handle string, svc *openapi.APIClient) (interface{}, error) {
	var err error

	// execute get call
	var resp openapi.Workspace

	getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		resp, _, err = svc.OrgWorkspaces.Get(context.Background(), identityHandle, handle).Execute()
		return resp, err
	}

	response, err := plugin.RetryHydrate(ctx, d, h, getDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

	workspace := response.(openapi.Workspace)

	if err != nil {
		plugin.Logger(ctx).Error("getOrgWorkspace", "get", err)
		return nil, err
	}

	return workspace, nil
}

func getUserWorkspace(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, identityHandle string, handle string, svc *openapi.APIClient) (interface{}, error) {
	var err error

	// execute get call
	var resp openapi.Workspace

	getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		resp, _, err = svc.UserWorkspaces.Get(context.Background(), identityHandle, handle).Execute()
		return resp, err
	}

	response, err := plugin.RetryHydrate(ctx, d, h, getDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

	workspace := response.(openapi.Workspace)

	if err != nil {
		plugin.Logger(ctx).Error("getOrgWorkspace", "get", err)
		return nil, err
	}

	return workspace, nil
}
