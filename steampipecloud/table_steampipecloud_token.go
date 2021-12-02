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

func tableSteampipecloudToken(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "steampipecloud_token",
		Description: "Steampipecloud Token",
		List: &plugin.ListConfig{
			Hydrate: listTokens,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getToken,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "status",
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
				Name:        "last4",
				Description: "",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_at",
				Description: "",
				Type:        proto.ColumnType_TIMESTAMP,
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
		},
	}
}

//// LIST FUNCTION

func listTokens(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listTokens", "connection_error", err)
		return nil, err
	}
	getUserIdentityCached := plugin.HydrateFunc(getUserIdentity).WithCache()
	commonData, err := getUserIdentityCached(ctx, d, h)
	user := commonData.(openapi.TypesUser)

	// execute list call
	pagesLeft := true

	var resp openapi.TypesListTokensResponse
	var httpResp *http.Response

	for pagesLeft {
		b, err := retry.NewFibonacci(100 * time.Millisecond)
		if resp.NextToken != nil {
			err = retry.Do(ctx, retry.WithMaxRetries(10, b), func(ctx context.Context) error {
				resp, httpResp, err = svc.UserTokensApi.ListTokens(context.Background(), user.Handle).NextToken(*resp.NextToken).Execute()
				// 429 too many request
				if httpResp.StatusCode == 429 {
					return retry.RetryableError(err)
				}
				return nil
			})
		} else {
			err = retry.Do(ctx, retry.WithMaxRetries(10, b), func(ctx context.Context) error {
				resp, httpResp, err = svc.UserTokensApi.ListTokens(context.Background(), user.Handle).Execute()
				// 429 too many request
				if httpResp.StatusCode == 429 {
					return retry.RetryableError(err)
				}
				return nil
			})
		}

		if err != nil {
			plugin.Logger(ctx).Error("listTokens", "list", err)
			return nil, err
		}
		if resp.HasItems() {
			for _, token := range *resp.Items {
				d.StreamListItem(ctx, token)
			}
		}
		if resp.NextToken == nil {
			pagesLeft = false
		}
	}

	return nil, nil
}

func getToken(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getToken", "connection_error", err)
		return nil, err
	}
	id := d.KeyColumnQuals["id"].GetStringValue()

	// check if id is empty
	if id == "" {
		return nil, nil
	}

	getUserIdentityCached := plugin.HydrateFunc(getUserIdentity).WithCache()
	commonData, err := getUserIdentityCached(ctx, d, h)
	user := commonData.(openapi.TypesUser)

	var resp openapi.TypesToken
	var httpResp *http.Response

	// execute get call
	b, err := retry.NewFibonacci(100 * time.Millisecond)
	err = retry.Do(ctx, retry.WithMaxRetries(10, b), func(ctx context.Context) error {
		resp, httpResp, err = svc.UserTokensApi.GetToken(context.Background(), id, user.Handle).Execute()
		// 429 too many request
		if httpResp.StatusCode == 429 {
			return retry.RetryableError(err)
		}
		return nil
	})

	if err != nil {
		plugin.Logger(ctx).Error("getToken", "get", err)
		return nil, err
	}

	// 404 Not Found
	if httpResp.StatusCode == 404 {
		return nil, nil
	}

	return resp, nil
}
