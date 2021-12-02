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

func tableSteampipecloudOrganization(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "steampipecloud_org",
		Description: "Steampipecloud Organization",
		List: &plugin.ListConfig{
			Hydrate: listOrganizations,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "org_id",
				Description: "",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "org_handle",
				Description: "",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Org.Handle"),
			},
			{
				Name:        "status",
				Description: "",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "role",
				Description: "",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "user_id",
				Description: "",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "created_at",
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
			{
				Name:        "org",
				Description: "",
				Type:        proto.ColumnType_JSON,
			},
		},
	}
}

//// LIST FUNCTION

func listOrganizations(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listOrganizations", "connection_error", err)
		return nil, err
	}

	getUserIdentityCached := plugin.HydrateFunc(getUserIdentity).WithCache()
	commonData, err := getUserIdentityCached(ctx, d, h)
	user := commonData.(openapi.TypesUser)

	// execute list call
	pagesLeft := true

	var resp openapi.TypesListUserOrgsResponse
	var httpResp *http.Response

	for pagesLeft {
		b, err := retry.NewFibonacci(100 * time.Millisecond)
		if resp.NextToken != nil {
			err = retry.Do(ctx, retry.WithMaxRetries(10, b), func(ctx context.Context) error {
				resp, httpResp, err = svc.UsersApi.ListOrgUsers(context.Background(), user.Handle).NextToken(*resp.NextToken).Execute()
				// 429 too many request
				if httpResp.StatusCode == 429 {
					return retry.RetryableError(err)
				}
				return nil
			})
		} else {
			err = retry.Do(ctx, retry.WithMaxRetries(10, b), func(ctx context.Context) error {
				resp, httpResp, err = svc.UsersApi.ListOrgUsers(context.Background(), user.Handle).Execute()
				// 429 too many request
				if httpResp.StatusCode == 429 {
					return retry.RetryableError(err)
				}
				return nil
			})
		}

		if err != nil {
			plugin.Logger(ctx).Error("listOrganizations", "list", err)
			return nil, err
		}
		if resp.HasItems() {
			for _, org := range *resp.Items {
				d.StreamListItem(ctx, org)
			}
		}

		if resp.NextToken == nil {
			pagesLeft = false
		}
	}

	return nil, nil
}
