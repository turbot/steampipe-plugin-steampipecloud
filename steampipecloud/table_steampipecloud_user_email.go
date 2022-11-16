package steampipecloud

import (
	"context"

	openapi "github.com/turbot/steampipe-cloud-sdk-go"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableSteampipeCloudUserEmail(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "steampipecloud_user_email",
		Description: "User Email table allows users to manage their emails.",
		List: &plugin.ListConfig{
			Hydrate: listUserEmails,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The unique identifier for the user email.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "email",
				Description: "The email ID.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The current status of the email. Can be one of pending / verified",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_at",
				Description: "The time when the user email was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "version_id",
				Description: "The version ID of the user email.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromCamel(),
			},
		},
	}
}

func listUserEmails(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listUserEmails", "connection_error", err)
		return nil, err
	}

	// declare variables
	pagesLeft := true
	var resp openapi.ListUserEmailsResponse
	var listDetails func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error)
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

	// Get Cached user
	getUserIdentityCached := plugin.HydrateFunc(getUserIdentity)
	commonData, err := getUserIdentityCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("listUserEmails", "error", err)
		return nil, err
	}
	user := commonData.(openapi.User)

	for pagesLeft {
		if resp.NextToken != nil {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.Users.ListEmails(ctx, user.Handle).NextToken(*resp.NextToken).Limit(maxResults).Execute()
				return resp, err
			}
		} else {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.Users.ListEmails(ctx, user.Handle).Limit(maxResults).Execute()
				return resp, err
			}
		}

		response, err := plugin.RetryHydrate(ctx, d, h, listDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

		if err != nil {
			plugin.Logger(ctx).Error("listUserEmails", "list", err)
			return nil, err
		}

		result := response.(openapi.ListUserEmailsResponse)

		if result.HasItems() {
			for _, userEmail := range *result.Items {
				d.StreamListItem(ctx, userEmail)

				// Context can be cancelled due to manual cancellation or the limit has been hit
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return nil, nil
				}
			}
		}
		if result.NextToken == nil {
			pagesLeft = false
		} else {
			resp.NextToken = result.NextToken
		}
	}

	return nil, nil
}
