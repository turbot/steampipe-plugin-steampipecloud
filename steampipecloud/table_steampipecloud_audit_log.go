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

func tableSteampipecloudAuditLog(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "steampipecloud_audit_log",
		Description: "Steampipecloud Audit Log",
		List: &plugin.ListConfig{
			Hydrate: listAuditLogs,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "identity_handle",
					Require: plugin.Required,
				},
			},
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
			},
			{
				Name:        "action_type",
				Description: "",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "actor_avatar_url",
				Description: "",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "actor_display_name",
				Description: "",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "actor_handle",
				Description: "",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "actor_id",
				Description: "",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "actor_ip",
				Description: "",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_at",
				Description: "",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "target_handle",
				Description: "",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "target_id",
				Description: "",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "data",
				Description: "",
				Type:        proto.ColumnType_JSON,
			},
		},
	}
}

//// LIST FUNCTION

func listAuditLogs(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	getUserIdentityCached := plugin.HydrateFunc(getUserIdentity).WithCache()
	commonData, err := getUserIdentityCached(ctx, d, h)
	user := commonData.(openapi.TypesUser)

	handle := d.KeyColumnQuals["identity_handle"].GetStringValue()

	if handle == user.Handle {
		err = listUserAuditLogs(ctx, d, h, handle)
	} else {
		err = listOrgAuditLogs(ctx, d, h, handle)
	}

	if err != nil {
		plugin.Logger(ctx).Error("listAuditLogs", "list", err)
		return nil, err
	}
	return nil, nil
}

func listOrgAuditLogs(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, handle string) error {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listOrgAuditLogs", "connection_error", err)
		return err
	}

	// execute list call
	pagesLeft := true
	var resp openapi.TypesListAuditLogsResponse
	var httpResp *http.Response

	for pagesLeft {
		b, err := retry.NewFibonacci(100 * time.Millisecond)
		if resp.NextToken != nil {
			err = retry.Do(ctx, retry.WithMaxRetries(10, b), func(ctx context.Context) error {
				resp, httpResp, err = svc.OrgsApi.ListOrgAuditLogs(context.Background(), handle).NextToken(*resp.NextToken).Execute()
				// 429 too many request
				if httpResp.StatusCode == 429 {
					return retry.RetryableError(err)
				}
				return nil
			})
		} else {
			err = retry.Do(ctx, retry.WithMaxRetries(10, b), func(ctx context.Context) error {
				resp, httpResp, err = svc.OrgsApi.ListOrgAuditLogs(context.Background(), handle).Execute()
				// 429 too many request
				if httpResp.StatusCode == 429 {
					return retry.RetryableError(err)
				}
				return nil
			})
		}

		if err != nil {
			plugin.Logger(ctx).Error("listOrgAuditLogs", "list", err)
			return err
		}

		if resp.HasItems() {
			for _, log := range *resp.Items {
				d.StreamListItem(ctx, log)
			}
		}
		if resp.NextToken == nil {
			pagesLeft = false
		}
	}

	return nil
}

func listUserAuditLogs(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, handle string) error {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listUserAuditLogs", "connection_error", err)
		return err
	}

	// execute list call
	pagesLeft := true
	var resp openapi.TypesListAuditLogsResponse
	var httpResp *http.Response

	for pagesLeft {
		b, err := retry.NewFibonacci(100 * time.Millisecond)
		if resp.NextToken != nil {
			err = retry.Do(ctx, retry.WithMaxRetries(10, b), func(ctx context.Context) error {
				resp, httpResp, err = svc.UsersApi.ListUserAuditLogs(context.Background(), handle).NextToken(*resp.NextToken).Execute()
				// 429 too many request
				if httpResp.StatusCode == 429 {
					return retry.RetryableError(err)
				}
				return nil
			})

		} else {
			err = retry.Do(ctx, retry.WithMaxRetries(10, b), func(ctx context.Context) error {
				resp, httpResp, err = svc.UsersApi.ListUserAuditLogs(context.Background(), handle).Execute()
				// 429 too many request
				if httpResp.StatusCode == 429 {
					return retry.RetryableError(err)
				}
				return nil
			})

		}

		if err != nil {
			plugin.Logger(ctx).Error("listUserAuditLogs", "list", err)
			return err
		}

		if resp.HasItems() {
			for _, log := range *resp.Items {
				d.StreamListItem(ctx, log)
			}
		}
		if resp.NextToken == nil {
			pagesLeft = false
		}
	}

	return nil
}
