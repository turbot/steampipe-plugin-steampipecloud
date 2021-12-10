package steampipecloud

import (
	"context"

	openapi "github.com/turbot/steampipe-cloud-sdk-go"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableSteampipeCloudConnection(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "steampipecloud_connection",
		Description: "SteampipeCloud Connection",
		List: &plugin.ListConfig{
			Hydrate: listConnections,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "identity_handle",
					Require: plugin.Optional,
				},
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"handle", "identity_handle"}),
			Hydrate:    getConnection,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The unique identifier for the connection.",
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
				Name:        "identity_handle",
				Description: "The handle name for an identity where the action has been performed.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getConnectionIdentityHandle,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "identity_type",
				Description: "The unique identifier for an identity where the action has been performed.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("IdentityId").Transform(setIdentityType),
			},
			{
				Name:        "handle",
				Description: "The handle name for the connection.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "plugin",
				Description: "The plugin name for the connection.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_at",
				Description: "The connection created time.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "type",
				Description: "The connection type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "version_id",
				Description: "The current version id for the connection.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "updated_at",
				Description: "The connection updated time.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "config",
				Description: "The connection config details.",
				Type:        proto.ColumnType_JSON,
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

func listConnections(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	getUserIdentityCached := plugin.HydrateFunc(getUserIdentity).WithCache()
	commonData, err := getUserIdentityCached(ctx, d, h)
	user := commonData.(openapi.TypesUser)

	handle := d.KeyColumnQuals["identity_handle"].GetStringValue()

	if handle == "" {
		err = listActorConnections(ctx, d, h)
	} else if handle == user.Handle {
		err = listUserConnections(ctx, d, h, handle)
	} else {
		err = listOrgConnections(ctx, d, h, handle)
	}

	if err != nil {
		plugin.Logger(ctx).Error("listConnections", "list", err)
		return nil, err
	}
	return nil, nil
}

func listOrgConnections(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, handle string) error {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listOrgConnections", "connection_error", err)
		return err
	}

	// execute list call
	pagesLeft := true
	var resp openapi.TypesListConnectionsResponse
	var listDetails func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error)

	for pagesLeft {
		if resp.NextToken != nil {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.OrgConnectionsApi.ListOrgConnections(context.Background(), handle).NextToken(*resp.NextToken).Execute()
				return resp, err
			}
		} else {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.OrgConnectionsApi.ListOrgConnections(context.Background(), handle).Execute()
				return resp, err
			}
		}

		response, err := plugin.RetryHydrate(ctx, d, h, listDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

		if err != nil {
			plugin.Logger(ctx).Error("listOrgConnections", "list", err)
			return err
		}

		result := response.(openapi.TypesListConnectionsResponse)

		for _, connection := range *result.Items {
			d.StreamListItem(ctx, connection)
		}
		if result.NextToken == nil {
			pagesLeft = false
		} else {
			resp.NextToken = result.NextToken
		}
	}

	return nil
}

func listUserConnections(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, handle string) error {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listUserConnections", "connection_error", err)
		return err
	}

	// execute list call
	pagesLeft := true
	var resp openapi.TypesListConnectionsResponse
	var listDetails func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error)

	for pagesLeft {
		if resp.NextToken != nil {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.UserConnectionsApi.ListUserConnections(context.Background(), handle).NextToken(*resp.NextToken).Execute()
				return resp, err
			}
		} else {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.UserConnectionsApi.ListUserConnections(context.Background(), handle).Execute()
				return resp, err
			}
		}

		response, err := plugin.RetryHydrate(ctx, d, h, listDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

		if err != nil {
			plugin.Logger(ctx).Error("listUserConnections", "list", err)
			return err
		}

		result := response.(openapi.TypesListConnectionsResponse)

		for _, connection := range *result.Items {
			d.StreamListItem(ctx, connection)
		}
		if result.NextToken == nil {
			pagesLeft = false
		} else {
			resp.NextToken = result.NextToken
		}
	}

	return nil
}

func listActorConnections(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) error {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listActorConnections", "connection_error", err)
		return err
	}

	// execute list call
	pagesLeft := true

	var resp openapi.TypesListConnectionsResponse
	var listDetails func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error)

	for pagesLeft {
		if resp.NextToken != nil {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.UserConnectionsApi.ListActorConnections(context.Background()).NextToken(*resp.NextToken).Execute()
				return resp, err
			}
		} else {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.UserConnectionsApi.ListActorConnections(context.Background()).Execute()
				return resp, err
			}
		}

		response, err := plugin.RetryHydrate(ctx, d, h, listDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

		if err != nil {
			plugin.Logger(ctx).Error("listActorConnections", "list", err)
			return err
		}

		result := response.(openapi.TypesListConnectionsResponse)

		if result.HasItems() {
			for _, connection := range *result.Items {
				d.StreamListItem(ctx, connection)
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

func getConnection(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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
		resp, err = getUserConnection(ctx, d, h, identityHandle, handle)
	} else {
		resp, err = getOrgConnection(ctx, d, h, identityHandle, handle)
	}

	if err != nil {
		plugin.Logger(ctx).Error("getConnection", "get", err)
		return nil, err
	}

	if resp == nil {
		return nil, nil
	}

	return resp.(openapi.TypesConnection), nil
}

func getOrgConnection(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, identityHandle string, handle string) (interface{}, error) {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getOrgConnection", "connection_error", err)
		return nil, err
	}

	// execute get call
	var resp openapi.TypesConnection

	getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		resp, _, err = svc.OrgConnectionsApi.GetOrgConnection(context.Background(), identityHandle, handle).Execute()
		return resp, err
	}

	response, err := plugin.RetryHydrate(ctx, d, h, getDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

	connection := response.(openapi.TypesConnection)

	if err != nil {
		plugin.Logger(ctx).Error("getOrgConnection", "get", err)
		return nil, err
	}

	return connection, nil
}

func getUserConnection(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, identityHandle string, handle string) (interface{}, error) {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getUserConnection", "connection_error", err)
		return nil, err
	}

	// execute get call
	var resp openapi.TypesConnection

	getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		resp, _, err = svc.UserConnectionsApi.GetUserConnection(context.Background(), identityHandle, handle).Execute()
		return resp, err
	}

	response, err := plugin.RetryHydrate(ctx, d, h, getDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

	connection := response.(openapi.TypesConnection)

	if err != nil {
		plugin.Logger(ctx).Error("getUserConnection", "get", err)
		return nil, err
	}

	return connection, nil
}

func getConnectionIdentityHandle(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	connection := h.Item.(openapi.TypesConnection)
	handle := d.KeyColumnQuals["identity_handle"].GetStringValue()

	if handle == "" {
		return connection.Identity.Handle, nil
	}

	return handle, nil
}
