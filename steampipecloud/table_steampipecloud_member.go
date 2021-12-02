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

func tableSteampipecloudMember(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "steampipecloud_member",
		Description: "Steampipecloud Member",
		List: &plugin.ListConfig{
			ParentHydrate: listOrganizations,
			Hydrate:       listMembers,
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
				Description: "The email id for the user.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "role",
				Description: "The role of the member.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_at",
				Description: "The member creation time.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "updated_at",
				Description: "The last time member was updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "version_id",
				Description: "The current version id for the member.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "user",
				Description: "The user details.",
				Type:        proto.ColumnType_JSON,
			},
		},
	}
}

//// LIST FUNCTION

func listMembers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org := h.Item.(openapi.TypesOrg)

	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listMembers", "connection_error", err)
		return nil, err
	}

	// execute ListAcceptedOrgMembers call
	pagesLeft := true
	var resp openapi.TypesListOrgUsersResponse
	var httpResp *http.Response

	for pagesLeft {
		b, err := retry.NewFibonacci(100 * time.Millisecond)
		if resp.NextToken != nil {
			err = retry.Do(ctx, retry.WithMaxRetries(10, b), func(ctx context.Context) error {
				resp, httpResp, err = svc.OrgMembersApi.ListAcceptedOrgMembers(context.Background(), org.Handle).NextToken(*resp.NextToken).Execute()
				// 429 too many request
				if httpResp.StatusCode == 429 {
					return retry.RetryableError(err)
				}
				return nil
			})
		} else {
			err = retry.Do(ctx, retry.WithMaxRetries(10, b), func(ctx context.Context) error {
				resp, httpResp, err = svc.OrgMembersApi.ListAcceptedOrgMembers(context.Background(), org.Handle).Execute()
				// 429 too many request
				if httpResp.StatusCode == 429 {
					return retry.RetryableError(err)
				}
				return nil
			})
		}

		if err != nil {
			plugin.Logger(ctx).Error("listAcceptedOrgMembers", "list", err)
			return nil, err
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

	// execute ListInvitedOrgMembers call
	pagesLeft = true

	for pagesLeft {
		b, err := retry.NewFibonacci(100 * time.Millisecond)
		if resp.NextToken != nil {
			err = retry.Do(ctx, retry.WithMaxRetries(10, b), func(ctx context.Context) error {
				resp, httpResp, err = svc.OrgMembersApi.ListInvitedOrgMembers(context.Background(), org.Handle).NextToken(*resp.NextToken).Execute()
				// 429 too many request
				if httpResp.StatusCode == 429 {
					return retry.RetryableError(err)
				}
				return nil
			})

		} else {
			err = retry.Do(ctx, retry.WithMaxRetries(10, b), func(ctx context.Context) error {
				resp, httpResp, err = svc.OrgMembersApi.ListInvitedOrgMembers(context.Background(), org.Handle).Execute()
				// 429 too many request
				if httpResp.StatusCode == 429 {
					return retry.RetryableError(err)
				}
				return nil
			})

		}

		if err != nil {
			plugin.Logger(ctx).Error("listInvitedOrgMembers", "list", err)
			return nil, err
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

	return nil, nil
}
