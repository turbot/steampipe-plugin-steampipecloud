package steampipecloud

import (
	"context"

	openapi "github.com/turbot/steampipe-cloud-sdk-go"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableSteampipecloudWorkspace(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "steampipecloud_workspace",
		Description: "Steampipecloud Workspace",
		List: &plugin.ListConfig{
			Hydrate: listWorkspaces,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "identity_handle",
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
				Description: "The handle name for the connection.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "database_name",
				Description: "The database name for the connection.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_at",
				Description: "The creation time of the connection.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "hive",
				Description: "The database hive for this workspace.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "identity_id",
				Description: "The unique identifier for an identity where the action has been performed.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "identity_handle",
				Description: "The handle name for an identity where the action has been performed.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getWorkspaceIdentityHandle,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "identity_type",
				Description: "The unique identifier for an identity where the action has been performed.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("IdentityId").Transform(setIdentityType),
			},
			{
				Name:        "version_id",
				Description: "The current version id of the workspace.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "updated_at",
				Description: "The last updated time of the workspace.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "identity",
				Description: "The identity where the action has been performed.",
				Type:        proto.ColumnType_JSON,
			},
		},
	}
}

//// LIST FUNCTION

func listWorkspaces(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	getUserIdentityCached := plugin.HydrateFunc(getUserIdentity).WithCache()
	commonData, err := getUserIdentityCached(ctx, d, h)
	user := commonData.(openapi.TypesUser)

	handle := d.KeyColumnQuals["identity_handle"].GetStringValue()

	if handle == "" {
		err = listActorWorkspaces(ctx, d, h)
	} else if handle == user.Handle {
		err = listUserWorkspaces(ctx, d, h, handle)
	} else {
		err = listOrgWorkspaces(ctx, d, h, handle)
	}

	if err != nil {
		plugin.Logger(ctx).Error("listWorkspaces", "list", err)
		return nil, err
	}
	return nil, nil
}

func listUserWorkspaces(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, handle string) error {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listUserWorkspaces", "connection_error", err)
		return err
	}

	// execute list call
	pagesLeft := true

	var resp openapi.TypesListWorkspacesResponse
	var listDetails func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error)

	for pagesLeft {
		if resp.NextToken != nil {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.UserWorkspacesApi.ListUserWorkspaces(context.Background(), handle).NextToken(*resp.NextToken).Execute()
				return resp, err
			}
		} else {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.UserWorkspacesApi.ListUserWorkspaces(context.Background(), handle).Execute()
				return resp, err
			}
		}

		response, err := plugin.RetryHydrate(ctx, d, h, listDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

		if err != nil {
			plugin.Logger(ctx).Error("listUserWorkspaces", "list", err)
			return err
		}

		result := response.(openapi.TypesListWorkspacesResponse)

		if result.HasItems() {
			for _, workspace := range *result.Items {
				d.StreamListItem(ctx, workspace)
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

func listOrgWorkspaces(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, handle string) error {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listOrgWorkspaces", "connection_error", err)
		return err
	}

	// execute list call
	pagesLeft := true

	var resp openapi.TypesListWorkspacesResponse
	var listDetails func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error)

	for pagesLeft {
		if resp.NextToken != nil {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.OrgWorkspacesApi.ListOrgWorkspaces(context.Background(), handle).NextToken(*resp.NextToken).Execute()
				return resp, err
			}
		} else {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.OrgWorkspacesApi.ListOrgWorkspaces(context.Background(), handle).Execute()
				return resp, err
			}
		}

		response, err := plugin.RetryHydrate(ctx, d, h, listDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

		if err != nil {
			plugin.Logger(ctx).Error("listOrgWorkspaces", "list", err)
			return err
		}

		result := response.(openapi.TypesListWorkspacesResponse)

		if result.HasItems() {
			for _, workspace := range *result.Items {
				d.StreamListItem(ctx, workspace)
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

func listActorWorkspaces(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) error {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listActorWorkspaces", "connection_error", err)
		return err
	}

	// execute list call
	pagesLeft := true

	var resp openapi.TypesListWorkspacesResponse
	var listDetails func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error)

	for pagesLeft {
		if resp.NextToken != nil {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.UserWorkspacesApi.ListActorWorkspaces(context.Background()).NextToken(*resp.NextToken).Execute()
				return resp, err
			}
		} else {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.UserWorkspacesApi.ListActorWorkspaces(context.Background()).Execute()
				return resp, err
			}
		}

		response, err := plugin.RetryHydrate(ctx, d, h, listDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

		if err != nil {
			plugin.Logger(ctx).Error("listActorWorkspaces", "list", err)
			return err
		}

		result := response.(openapi.TypesListWorkspacesResponse)

		if result.HasItems() {
			for _, workspace := range *result.Items {
				d.StreamListItem(ctx, workspace)
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
	getUserIdentityCached := plugin.HydrateFunc(getUserIdentity).WithCache()
	commonData, err := getUserIdentityCached(ctx, d, h)
	user := commonData.(openapi.TypesUser)
	var resp interface{}

	if identityHandle == user.Handle {
		resp, err = getUserWorkspace(ctx, d, h, identityHandle, handle)
	} else {
		resp, err = getOrgWorkspace(ctx, d, h, identityHandle, handle)
	}

	if err != nil {
		plugin.Logger(ctx).Error("getWorkspace", "get", err)
		return nil, err
	}

	if resp == nil {
		return nil, nil
	}

	return resp.(openapi.TypesWorkspace), nil
}

func getOrgWorkspace(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, identityHandle string, handle string) (interface{}, error) {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getOrgWorkspace", "connection_error", err)
		return nil, err
	}

	// execute get call
	var resp openapi.TypesWorkspace

	getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		resp, _, err = svc.OrgWorkspacesApi.GetOrgWorkspace(context.Background(), identityHandle, handle).Execute()
		return resp, err
	}

	response, err := plugin.RetryHydrate(ctx, d, h, getDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

	workspace := response.(openapi.TypesWorkspace)

	if err != nil {
		plugin.Logger(ctx).Error("getOrgWorkspace", "get", err)
		return nil, err
	}

	return workspace, nil
}

func getUserWorkspace(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, identityHandle string, handle string) (interface{}, error) {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getUserWorkspace", "connection_error", err)
		return nil, err
	}

	// execute get call
	var resp openapi.TypesWorkspace

	getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		resp, _, err = svc.UserWorkspacesApi.GetUserWorkspace(context.Background(), identityHandle, handle).Execute()
		return resp, err
	}

	response, err := plugin.RetryHydrate(ctx, d, h, getDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

	workspace := response.(openapi.TypesWorkspace)

	if err != nil {
		plugin.Logger(ctx).Error("getOrgWorkspace", "get", err)
		return nil, err
	}

	return workspace, nil
}

func getWorkspaceIdentityHandle(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	workspace := h.Item.(openapi.TypesWorkspace)
	handle := d.KeyColumnQuals["identity_handle"].GetStringValue()

	if handle == "" {
		return workspace.Identity.Handle, nil
	}

	return handle, nil
}
