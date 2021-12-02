package steampipecloud

import (
	"context"

	openapi "github.com/turbot/steampipe-cloud-sdk-go"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableSteampipecloudUser(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "steampipecloud_user",
		Description: "Steampipecloud User",
		List: &plugin.ListConfig{
			Hydrate: getUser,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "display_name",
				Description: "",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "handle",
				Description: "",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "url",
				Description: "",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "avatar_url",
				Description: "",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_at",
				Description: "",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "email",
				Description: "",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "preview_access_mode",
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
				Type:        proto.ColumnType_STRING,
			},
		},
	}
}

//// LIST FUNCTION

func getUser(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	getUserIdentityCached := plugin.HydrateFunc(getUserIdentity).WithCache()
	commonData, err := getUserIdentityCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("getUser", "error", err)
		return nil, err
	}

	user := commonData.(openapi.TypesUser)

	d.StreamListItem(ctx, user)

	return nil, nil
}
