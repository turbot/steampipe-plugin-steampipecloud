package steampipecloud

import (
	"context"

	openapi "github.com/turbot/steampipe-cloud-sdk-go"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableSteampipeCloudUserPreferences(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "steampipecloud_user_preferences",
		Description: "User Preferences represents various preferences settings for a user e.g. email settings.",
		List: &plugin.ListConfig{
			Hydrate: getUserPreferences,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The unique identifier for the user preferences.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "communication_community_updates",
				Description: "Is the user subscribed to receiving community update emails.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "communication_product_updates",
				Description: "Is the user subscribed to receiving product update emails.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "communication_tips_and_tricks",
				Description: "Is the user subscribed to receiving tips and tricks emails.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_at",
				Description: "The time when the user preferences was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "updated_at",
				Description: "The time when any of the user preferences was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "version_id",
				Description: "The version ID of the user preferences.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromCamel(),
			},
		},
	}
}

func getUserPreferences(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getUserPreferences", "connection_error", err)
		return nil, err
	}

	// Get Cached user
	getUserIdentityCached := plugin.HydrateFunc(getUserIdentity)
	commonData, err := getUserIdentityCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("getUserPreferences", "error", err)
		return nil, err
	}
	user := commonData.(openapi.User)

	// Function to fetch the user preferences
	getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		resp, _, err := svc.Users.GetPreferences(ctx, user.Handle).Execute()
		return resp, err
	}

	// Execute function to fetch the user preferences
	response, err := plugin.RetryHydrate(ctx, d, h, getDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})
	if err != nil {
		plugin.Logger(ctx).Error("getUserPreferences", "error", err)
		return nil, err
	}

	// Extract user preferences from the response object
	userPreferences := response.(openapi.UserPreferences)

	// Push the preferences object to the list stream
	d.StreamListItem(ctx, userPreferences)

	return nil, nil
}
