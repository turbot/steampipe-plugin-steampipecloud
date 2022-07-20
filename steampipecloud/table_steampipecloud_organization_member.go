package steampipecloud

import (
	"context"

	openapi "github.com/turbot/steampipe-cloud-sdk-go"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func tableSteampipeCloudOrganizationMember(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "steampipecloud_organization_member",
		Description: "Organization members can collaborate and share workspaces and connections.",
		List: &plugin.ListConfig{
			ParentHydrate: listOrganizations,
			Hydrate:       listOrganizationMembers,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"handle"}),
			Hydrate:    getOrganization,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The unique identifier for the member.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "org_id",
				Description: "The unique identifier for the organization.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "status",
				Description: "The member current status.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "user_id",
				Description: "The unique identifier for the user.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "user_handle",
				Description: "The handle name for the user.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "email",
				Description: "The email address for the member.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "role",
				Description: "The role of the member.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_at",
				Description: "The member's creation time.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "updated_at",
				Description: "The member's last update time.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "created_by",
				Description: "ID of the user who invited the member.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CreatedById"),
			},
			{
				Name:        "updated_by",
				Description: "ID of the user who last updated the member.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("UpdatedById"),
			},
			{
				Name:        "version_id",
				Description: "The current version ID for the member.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromCamel(),
			},
		},
	}
}

//// LIST FUNCTION

func listOrganizationMembers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org := h.Item.(openapi.Org)

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

	var err error
	err = listOrgMembers(ctx, d, h, org.Handle, maxResults)
	if err != nil {
		plugin.Logger(ctx).Error("listOrganizationMembers", "error", err)
		return nil, err
	}

	if err != nil {
		plugin.Logger(ctx).Error("listOrganizationMembers", "error", err)
		return nil, err
	}
	return nil, nil
}

func listOrgMembers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, handle string, maxResults int32) error {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listOrgMembers", "connection_error", err)
		return err
	}

	pagesLeft := true
	var resp openapi.ListOrgUsersResponse
	var listDetails func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error)

	for pagesLeft {
		if resp.NextToken != nil {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.OrgMembers.List(context.Background(), handle).NextToken(*resp.NextToken).Limit(maxResults).Execute()
				return resp, err
			}
		} else {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.OrgMembers.List(context.Background(), handle).Limit(maxResults).Execute()
				return resp, err
			}
		}

		response, err := plugin.RetryHydrate(ctx, d, h, listDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

		if err != nil {
			plugin.Logger(ctx).Error("listOrgMembers", "list", err)
			return err
		}

		result := response.(openapi.ListOrgUsersResponse)

		if result.HasItems() {
			for _, member := range *result.Items {
				d.StreamListItem(ctx, member)

				// Context can be cancelled due to manual cancellation or the limit has been hit
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return nil
				}
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
