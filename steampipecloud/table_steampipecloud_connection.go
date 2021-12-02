package steampipecloud

import (
	"context"
	"net/http"
	"time"

	"github.com/sethvargo/go-retry"
	openapi "github.com/turbot/steampipe-cloud-sdk-go"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableSteampipecloudConnection(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "steampipecloud_connection",
		Description: "Steampipecloud Connection",
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
				Description: "",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "identity_id",
				Description: "",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "identity_handle",
				Description: "",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getConnectionIdentityHandle,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "handle",
				Description: "",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "plugin",
				Description: "",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_at",
				Description: "",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "type",
				Description: "",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "version_id",
				Description: "",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "updated_at",
				Description: "",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "config",
				Description: "",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "identity",
				Description: "",
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
	var httpResp *http.Response

	for pagesLeft {
		b, err := retry.NewFibonacci(100 * time.Millisecond)
		if resp.NextToken != nil {
			err = retry.Do(ctx, retry.WithMaxRetries(10, b), func(ctx context.Context) error {
				resp, httpResp, err = svc.OrgConnectionsApi.ListOrgConnections(context.Background(), handle).NextToken(*resp.NextToken).Execute()
				// 429 too many request
				if httpResp.StatusCode == 429 {
					return retry.RetryableError(err)
				}
				return nil
			})
		} else {
			err = retry.Do(ctx, retry.WithMaxRetries(10, b), func(ctx context.Context) error {
				resp, httpResp, err = svc.OrgConnectionsApi.ListOrgConnections(context.Background(), handle).Execute()
				// 429 too many request
				if httpResp.StatusCode == 429 {
					return retry.RetryableError(err)
				}
				return nil
			})
		}

		if err != nil {
			plugin.Logger(ctx).Error("listOrgConnections", "list", err)
			return err
		}

		for _, connection := range *resp.Items {
			d.StreamListItem(ctx, connection)
		}
		if resp.NextToken == nil {
			pagesLeft = false
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
	var httpResp *http.Response

	for pagesLeft {
		b, err := retry.NewFibonacci(100 * time.Millisecond)
		if resp.NextToken != nil {
			err = retry.Do(ctx, retry.WithMaxRetries(10, b), func(ctx context.Context) error {
				resp, httpResp, err = svc.UserConnectionsApi.ListUserConnections(context.Background(), handle).NextToken(*resp.NextToken).Execute()
				// 429 too many request
				if httpResp.StatusCode == 429 {
					return retry.RetryableError(err)
				}
				return nil
			})
		} else {
			err = retry.Do(ctx, retry.WithMaxRetries(10, b), func(ctx context.Context) error {
				resp, httpResp, err = svc.UserConnectionsApi.ListUserConnections(context.Background(), handle).Execute()
				// 429 too many request
				if httpResp.StatusCode == 429 {
					return retry.RetryableError(err)
				}
				return nil
			})
		}

		if err != nil {
			plugin.Logger(ctx).Error("listUserConnections", "list", err)
			return err
		}

		for _, connection := range *resp.Items {
			d.StreamListItem(ctx, connection)
		}
		if resp.NextToken == nil {
			pagesLeft = false
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
	var httpResp *http.Response

	for pagesLeft {
		b, err := retry.NewFibonacci(100 * time.Millisecond)
		if resp.NextToken != nil {
			err = retry.Do(ctx, retry.WithMaxRetries(10, b), func(ctx context.Context) error {
				resp, httpResp, err = svc.UserConnectionsApi.ListActorConnections(context.Background()).NextToken(*resp.NextToken).Execute()
				// 429 too many request
				if httpResp.StatusCode == 429 {
					return retry.RetryableError(err)
				}
				return nil
			})
		} else {
			err = retry.Do(ctx, retry.WithMaxRetries(10, b), func(ctx context.Context) error {
				resp, httpResp, err = svc.UserConnectionsApi.ListActorConnections(context.Background()).Execute()
				// 429 too many request
				if httpResp.StatusCode == 429 {
					return retry.RetryableError(err)
				}
				return nil
			})
		}

		if err != nil {
			plugin.Logger(ctx).Error("listActorConnections", "list", err)
			return err
		}

		if resp.HasItems() {
			for _, connection := range *resp.Items {
				d.StreamListItem(ctx, connection)
			}
		}
		if resp.NextToken == nil {
			pagesLeft = false
		}
	}

	return nil
}

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
	var httpResp *http.Response

	b, err := retry.NewFibonacci(100 * time.Millisecond)
	err = retry.Do(ctx, retry.WithMaxRetries(10, b), func(ctx context.Context) error {
		resp, httpResp, err = svc.OrgConnectionsApi.GetOrgConnection(context.Background(), identityHandle, handle).Execute()
		// 429 too many request
		if httpResp.StatusCode == 429 {
			return retry.RetryableError(err)
		}
		return nil
	})

	if err != nil {
		plugin.Logger(ctx).Error("getOrgConnection", "get", err)
		return nil, err
	}

	// 404 Not Found
	if httpResp.StatusCode == 404 {
		return nil, nil
	}

	return resp, nil
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
	var httpResp *http.Response

	b, err := retry.NewFibonacci(100 * time.Millisecond)
	err = retry.Do(ctx, retry.WithMaxRetries(10, b), func(ctx context.Context) error {
		resp, httpResp, err = svc.UserConnectionsApi.GetUserConnection(context.Background(), identityHandle, handle).Execute()
		// 429 too many request
		if httpResp.StatusCode == 429 {
			return retry.RetryableError(err)
		}
		return nil
	})

	if err != nil {
		plugin.Logger(ctx).Error("getUserConnection", "get", err)
		return nil, err
	}

	// 404 Not Found
	if httpResp.StatusCode == 404 {
		return nil, nil
	}

	return resp, nil
}

func getConnectionIdentityHandle(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	connection := h.Item.(openapi.TypesConnection)
	handle := d.KeyColumnQuals["identity_handle"].GetStringValue()

	if handle == "" {
		return connection.Identity.Handle, nil
	}

	return handle, nil
}
