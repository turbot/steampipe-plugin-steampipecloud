package steampipecloud

import (
	"context"
	"errors"

	openapi "github.com/turbot/steampipe-cloud-sdk-go"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableSteampipeCloudOrganizationMember(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "steampipecloud_organization_member",
		Description: "SteampipeCloud Organization Member",
		List: &plugin.ListConfig{
			ParentHydrate: listOrganizations,
			Hydrate:       listOrganizationMembers,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "status",
					Require: plugin.Optional,
				},
			},
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

func listOrganizationMembers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org := h.Item.(*openapi.TypesOrg)

	status := d.KeyColumnQuals["status"].GetStringValue()

	var err error
	if status == "" {
		err = listInvitedOrgMembers(ctx, d, h, org.Handle)
		if err != nil {
			plugin.Logger(ctx).Error("listInvitedOrgMembers", "error", err)
			return nil, err
		}
		err = listAcceptedOrgMembers(ctx, d, h, org.Handle)
	} else if status == "invited" {
		err = listInvitedOrgMembers(ctx, d, h, org.Handle)
	} else if status == "accepted" {
		err = listAcceptedOrgMembers(ctx, d, h, org.Handle)
	} else {
		return nil, errors.New("possible values are: invited and accepted")
	}

	if err != nil {
		plugin.Logger(ctx).Error("listOrganizationMembers", "list", err)
		return nil, err
	}
	return nil, nil
}

func listAcceptedOrgMembers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, handle string) error {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listAcceptedOrgMembers", "connection_error", err)
		return err
	}

	pagesLeft := true
	var resp openapi.TypesListOrgUsersResponse
	var listDetails func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error)

	for pagesLeft {
		if resp.NextToken != nil {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.OrgMembers.ListAccepted(context.Background(), handle).NextToken(*resp.NextToken).Execute()
				return resp, err
			}
		} else {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.OrgMembers.ListAccepted(context.Background(), handle).Execute()
				return resp, err
			}
		}

		response, err := plugin.RetryHydrate(ctx, d, h, listDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

		if err != nil {
			plugin.Logger(ctx).Error("listAcceptedOrgMembers", "list", err)
			return err
		}

		result := response.(openapi.TypesListOrgUsersResponse)

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

func listInvitedOrgMembers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, handle string) error {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listInvitedOrgMembers", "connection_error", err)
		return err
	}

	pagesLeft := true
	var resp openapi.TypesListOrgUsersResponse
	var listDetails func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error)

	for pagesLeft {
		if resp.NextToken != nil {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.OrgMembers.ListInvited(context.Background(), handle).NextToken(*resp.NextToken).Execute()
				return resp, err
			}
		} else {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.OrgMembers.ListInvited(context.Background(), handle).Execute()
				return resp, err
			}
		}

		response, err := plugin.RetryHydrate(ctx, d, h, listDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

		if err != nil {
			plugin.Logger(ctx).Error("listInvitedOrgMembers", "list", err)
			return err
		}

		result := response.(openapi.TypesListOrgUsersResponse)

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
