package steampipecloud

import (
	"context"

	openapi "github.com/turbot/steampipe-cloud-sdk-go"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableSteampipeCloudUser(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "steampipecloud_user",
		Description: "Users can manage connections, organizations, and workspaces.",
		List: &plugin.ListConfig{
			Hydrate: getUser,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The unique identifier for the user.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "display_name",
				Description: "The display name for the user.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The user status.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "handle",
				Description: "The handle name of the user.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "url",
				Description: "The URL of the user.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "avatar_url",
				Description: "The avatar URL of the user.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_at",
				Description: "The creation time of the user.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "email",
				Description: "The email address of the user.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "preview_access_mode",
				Description: "The preview mode for the current user.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "version_id",
				Description: "The version ID of the user.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "updated_at",
				Description: "The user's last updated time.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
		},
	}
}

//// LIST FUNCTION

func getUser(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	getUserIdentityCached := plugin.HydrateFunc(getUserIdentity)
	commonData, err := getUserIdentityCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("getUser", "error", err)
		return nil, err
	}

	user := commonData.(openapi.User)

	d.StreamListItem(ctx, user)

	return nil, nil
}
