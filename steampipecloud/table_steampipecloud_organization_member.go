package steampipecloud

import (
	"context"

	openapi "github.com/turbot/steampipe-cloud-sdk-go"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

type OrgDetails struct {
	OrgHandle string `json:"org_handle"`
}

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
			KeyColumns: plugin.AllColumns([]string{"org_handle", "user_handle"}),
			Hydrate:    getOrganizationMember,
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
				Name:        "org_handle",
				Description: "The handle of the organization.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getOrgDetails,
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
				Name:        "user",
				Description: "Information about the user.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "role",
				Description: "The role of the member.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "scope",
				Description: "The scope of the role. Will always be 'org'.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_at",
				Description: "The time when the member was invited.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "created_by_id",
				Description: "The unique identifier of the user who invited the member.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "created_by",
				Description: "Information about the user who invited the member.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "updated_at",
				Description: "The member's last update time.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "updated_by_id",
				Description: "The unique identifier of the user who last updated the member.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "updated_by",
				Description: "Information about the user who last updated the member.",
				Type:        proto.ColumnType_JSON,
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

	err := listOrgMembers(ctx, d, h, org.Handle, maxResults)
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
				resp, _, err = svc.OrgMembers.List(ctx, handle).NextToken(*resp.NextToken).Limit(maxResults).Execute()
				return resp, err
			}
		} else {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.OrgMembers.List(ctx, handle).Limit(maxResults).Execute()
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

func getOrganizationMember(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	orgHandle := d.KeyColumnQuals["org_handle"].GetStringValue()
	userhandle := d.KeyColumnQuals["user_handle"].GetStringValue()

	// check if handle or identityHandle is empty
	if orgHandle == "" || userhandle == "" {
		return nil, nil
	}

	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getOrganizationMember", "connection_error", err)
		return nil, err
	}

	var orgUser openapi.OrgUser

	getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		orgUser, _, err = svc.OrgMembers.Get(ctx, orgHandle, userhandle).Execute()
		return orgUser, err
	}

	response, err := plugin.RetryHydrate(ctx, d, h, getDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

	orgUser = response.(openapi.OrgUser)

	if err != nil {
		plugin.Logger(ctx).Error("getOrganizationMember", "get", err)
		return nil, err
	}

	return orgUser, nil
}

func getOrgDetails(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// get org details from hydrate data
	// org details reside in the parent item in this case
	switch o := h.ParentItem.(type) {
	case openapi.Org:
		return &OrgDetails{OrgHandle: h.ParentItem.(openapi.Org).Handle}, nil
	default:
		plugin.Logger(ctx).Debug("getOrgDetails", "Unknown Type", o)
	}

	// If we are in this section - it means that the org details are not present, so we query for the org
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getOrgDetails", "connection_error", err)
		return nil, err
	}

	getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		resp, _, err := svc.Orgs.Get(ctx, h.Item.(openapi.OrgUser).OrgId).Execute()
		return resp, err
	}

	response, _ := plugin.RetryHydrate(ctx, d, h, getDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

	return &OrgDetails{OrgHandle: response.(openapi.Org).Handle}, nil
}
