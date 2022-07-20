package steampipecloud

import (
	"context"

	openapi "github.com/turbot/steampipe-cloud-sdk-go"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func tableSteampipeCloudOrganization(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "steampipecloud_organization",
		Description: "Organizations include multiple users and can be used to share workspaces and connections.",
		List: &plugin.ListConfig{
			Hydrate: listOrganizations,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"handle"}),
			Hydrate:    getOrganization,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The unique identifier for a organization.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "handle",
				Description: "The handle name for the organization.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "display_name",
				Description: "The display name for the organization.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "avatar_url",
				Description: "The avatar URL of the organization.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "url",
				Description: "The URL of the organization.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_at",
				Description: "The organization's creation time.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "updated_at",
				Description: "The organization's last updated time.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "created_by",
				Description: "ID of the user who created the organization.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CreatedById"),
			},
			{
				Name:        "updated_by",
				Description: "ID of the user who last updated the organization.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("UpdatedById"),
			},
			{
				Name:        "version_id",
				Description: "The organization version ID.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromCamel(),
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

	// execute list call
	pagesLeft := true

	var resp openapi.ListActorOrgsResponse
	var listDetails func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error)

	for pagesLeft {
		if resp.NextToken != nil {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.Actors.ListOrgs(ctx).NextToken(*resp.NextToken).Limit(maxResults).Execute()
				return resp, err
			}
		} else {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.Actors.ListOrgs(ctx).Limit(maxResults).Execute()
				return resp, err
			}
		}

		response, err := plugin.RetryHydrate(ctx, d, h, listDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

		if err != nil {
			plugin.Logger(ctx).Error("listOrganizations", "list", err)
			return nil, err
		}

		result := response.(openapi.ListActorOrgsResponse)

		if result.HasItems() {
			for _, org := range *result.Items {
				d.StreamListItem(ctx, org.Org)

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

func getOrganization(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	handle := d.KeyColumnQuals["handle"].GetStringValue()

	// check if handle is empty
	if handle == "" {
		return nil, nil
	}

	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getOrganization", "connection_error", err)
		return nil, err
	}

	var resp openapi.Org

	getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		resp, _, err = svc.Orgs.Get(ctx, handle).Execute()
		return resp, err
	}

	response, err := plugin.RetryHydrate(ctx, d, h, getDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})
	if err != nil {
		plugin.Logger(ctx).Error("getOrganization", "get", err)
		return nil, err
	}

	if response == nil {
		return nil, nil
	}

	return response.(openapi.Org), nil
}
