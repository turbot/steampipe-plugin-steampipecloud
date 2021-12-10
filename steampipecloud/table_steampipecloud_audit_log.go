package steampipecloud

import (
	"context"

	openapi "github.com/turbot/steampipe-cloud-sdk-go"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableSteampipeCloudAuditLog(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "steampipecloud_audit_log",
		Description: "SteampipeCloud Audit Log",
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
				Description: "The unique identifier for an audit log.",
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
			},
			{
				Name:        "action_type",
				Description: "The action performed on the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "actor_avatar_url",
				Description: "The avatar of an actor who has performed the action.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "actor_display_name",
				Description: "The display name of an actor.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "actor_handle",
				Description: "The handle name of an actor.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "actor_id",
				Description: "The unique identifier of an actor.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "actor_ip",
				Description: "The IP address of the actor.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_at",
				Description: "The time when the action was performed.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "target_handle",
				Description: "The handle name of the entity where the action has been performed.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "target_id",
				Description: "The unique identifier of the entity where the action has been performed.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "data",
				Description: "The data which has been modified on the entity.",
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
	var listDetails func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error)

	for pagesLeft {
		if resp.NextToken != nil {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.OrgsApi.ListOrgAuditLogs(context.Background(), handle).NextToken(*resp.NextToken).Execute()
				return resp, err
			}
		} else {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.OrgsApi.ListOrgAuditLogs(context.Background(), handle).Execute()
				return resp, err
			}
		}

		response, err := plugin.RetryHydrate(ctx, d, h, listDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

		if err != nil {
			plugin.Logger(ctx).Error("listOrgAuditLogs", "list", err)
			return err
		}

		result := response.(openapi.TypesListAuditLogsResponse)

		if result.HasItems() {
			for _, log := range *result.Items {
				d.StreamListItem(ctx, log)
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
	var listDetails func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error)

	for pagesLeft {
		if resp.NextToken != nil {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.UsersApi.ListUserAuditLogs(context.Background(), handle).NextToken(*resp.NextToken).Execute()
				return resp, err
			}
		} else {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.UsersApi.ListUserAuditLogs(context.Background(), handle).Execute()
				return resp, err
			}
		}

		response, err := plugin.RetryHydrate(ctx, d, h, listDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

		if err != nil {
			plugin.Logger(ctx).Error("listUserAuditLogs", "list", err)
			return err
		}

		result := response.(openapi.TypesListAuditLogsResponse)

		if result.HasItems() {
			for _, log := range *result.Items {
				d.StreamListItem(ctx, log)
			}
		}
		if resp.NextToken == nil {
			pagesLeft = false
		} else {
			resp.NextToken = result.NextToken
		}
	}

	return nil
}
