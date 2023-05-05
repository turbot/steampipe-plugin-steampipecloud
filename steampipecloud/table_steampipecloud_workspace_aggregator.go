package steampipecloud

import (
	"context"
	"fmt"
	"strings"

	openapi "github.com/turbot/steampipe-cloud-sdk-go"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type IdentityWorkspaceDetailsForAggregator struct {
	IdentityHandle  string `json:"identity_handle"`
	IdentityType    string `json:"identity_type"`
	WorkspaceHandle string `json:"workspace_handle"`
}

//// TABLE DEFINITION

func tableSteampipeCloudWorkspaceAggregator(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "steampipecloud_workspace_aggregator",
		Description: "Aggregators allow users to define a collection of connections in a workspace.",
		DefaultIgnoreConfig: &plugin.IgnoreConfig{
			ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"404"}),
		},
		List: &plugin.ListConfig{
			ParentHydrate: listWorkspaces,
			Hydrate:       listWorkspaceAggregators,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "identity_handle",
					Require: plugin.Optional,
				},
				{
					Name:    "identity_id",
					Require: plugin.Optional,
				},
				{
					Name:    "workspace_handle",
					Require: plugin.Optional,
				},
				{
					Name:    "workspace_id",
					Require: plugin.Optional,
				},
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"identity_handle", "workspace_handle", "handle"}),
			Hydrate:    getWorkspaceAggregator,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The unique identifier for the aggregator.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "handle",
				Description: "The handle of the aggregator.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "identity_id",
				Description: "The unique identifier of the identity to which the aggregator belongs to.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "identity_handle",
				Description: "The handle of the identity.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIdentityWorkspaceDetailsForAggregator,
			},
			{
				Name:        "identity_type",
				Description: "The type of identity, can be org/user.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIdentityWorkspaceDetailsForAggregator,
			},
			{
				Name:        "workspace_id",
				Description: "The unique identifier for the workspace.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "workspace_handle",
				Description: "The handle of the workspace.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIdentityWorkspaceDetailsForAggregator,
			},
			{
				Name:        "type",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "plugin",
				Description: "The type of the aggregator connections.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "connections",
				Description: "The list of connections defined for the aggregator.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "created_at",
				Description: "The time when the aggregator was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "created_by_id",
				Description: "The unique identifier of the user who created the aggregator.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "created_by",
				Description: "Information about the user who created the aggregator.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "updated_at",
				Description: "The time when the aggregator was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "updated_by_id",
				Description: "The unique identifier of the user who last updated the aggregator.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "updated_by",
				Description: "Information about the user who last updated the aggregator.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "version_id",
				Description: "The current version ID for the aggregator.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromCamel(),
			},
		},
	}
}

//// LIST FUNCTION

func listWorkspaceAggregators(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var workspace *openapi.Workspace
	switch w := h.Item.(type) {
	case openapi.Workspace:
		wo := h.Item.(openapi.Workspace)
		workspace = &wo
	case *openapi.Workspace:
		workspace = h.Item.(*openapi.Workspace)
	default:
		plugin.Logger(ctx).Error("steampipecloud_workspace_aggregator.listWorkspaceAggregators", "unknown response type for workspace list parent hydrate call", w)
	}

	// If the requested number of items is less than the paging max limit
	// set the limit to that instead
	maxResults := int32(100)
	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(maxResults) {
			maxResults = int32(*limit)
		}
	}

	workspaceHandle := d.EqualsQualString("workspace_handle")
	workspaceId := d.EqualsQualString("workspace_id")
	var workspaceToPass string

	// Error out if both workspace_handle and workspace_id is passed
	if workspaceHandle != "" && workspaceId != "" {
		return nil, fmt.Errorf("please pass any one of workspace_id or workspace_handle")
	}
	// If either one has been passed, check whether either of the handle or the id matches with the workspace in context
	if workspaceHandle != "" || workspaceId != "" {
		if workspaceHandle == workspace.Handle {
			workspaceToPass = workspaceHandle
		} else if workspaceId == workspace.Id {
			workspaceToPass = workspaceId
		} else {
			return nil, nil
		}
	} else {
		// If neither is passed, we pass the context over to the call
		workspaceToPass = workspace.Id
	}

	var err error
	if strings.HasPrefix(workspace.IdentityId, "u_") {
		err = listUserWorkspaceAggregators(ctx, d, h, workspace.IdentityId, workspaceToPass, maxResults)
	} else {
		err = listOrgWorkspaceAggregators(ctx, d, h, workspace.IdentityId, workspaceToPass, maxResults)
	}

	if err != nil {
		plugin.Logger(ctx).Error("steampipecloud_workspace_aggregator.listWorkspaceAggregators", "query_error", err)
		return nil, err
	}

	return nil, nil
}

func listUserWorkspaceAggregators(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, userHandle string, workspaceHandle string, maxResults int32) error {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("steampipecloud_workspace_aggregator.listUserWorkspaceAggregators", "connection_error", err)
		return err
	}

	pagesLeft := true
	var resp openapi.ListWorkspaceAggregatorsResponse
	var listDetails func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error)

	for pagesLeft {
		if resp.NextToken != nil {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.UserWorkspaceAggregators.List(ctx, userHandle, workspaceHandle).NextToken(*resp.NextToken).Limit(maxResults).Execute()
				return resp, err
			}
		} else {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.UserWorkspaceAggregators.List(ctx, userHandle, workspaceHandle).Limit(maxResults).Execute()
				return resp, err
			}
		}

		response, err := plugin.RetryHydrate(ctx, d, h, listDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

		if err != nil {
			plugin.Logger(ctx).Error("steampipecloud_workspace_aggregator.listUserWorkspaceAggregators", "query_error", err)
			return err
		}

		result := response.(openapi.ListWorkspaceAggregatorsResponse)

		if result.HasItems() {
			for _, aggregator := range *result.Items {
				d.StreamListItem(ctx, aggregator)

				// Context can be cancelled due to manual cancellation or the limit has been hit
				if d.RowsRemaining(ctx) == 0 {
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

func listOrgWorkspaceAggregators(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, orgHandle string, workspaceHandle string, maxResults int32) error {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("steampipecloud_workspace_aggregator.listOrgWorkspaceAggregators", "connection_error", err)
		return err
	}

	pagesLeft := true
	var resp openapi.ListWorkspaceAggregatorsResponse
	var listDetails func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error)

	for pagesLeft {
		if resp.NextToken != nil {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.OrgWorkspaceAggregators.List(ctx, orgHandle, workspaceHandle).NextToken(*resp.NextToken).Limit(maxResults).Execute()
				return resp, err
			}
		} else {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.OrgWorkspaceAggregators.List(ctx, orgHandle, workspaceHandle).Limit(maxResults).Execute()
				return resp, err
			}
		}

		response, err := plugin.RetryHydrate(ctx, d, h, listDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

		if err != nil {
			plugin.Logger(ctx).Error("steampipecloud_workspace_aggregator.listOrgWorkspaceAggregators", "query_error", err)
			return err
		}

		result := response.(openapi.ListWorkspaceAggregatorsResponse)

		if result.HasItems() {
			for _, aggregator := range *result.Items {
				d.StreamListItem(ctx, aggregator)

				// Context can be cancelled due to manual cancellation or the limit has been hit
				if d.RowsRemaining(ctx) == 0 {
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

func getWorkspaceAggregator(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	identityHandle := d.EqualsQualString("identity_handle")
	workspaceHandle := d.EqualsQualString("workspace_handle")
	aggregatorHandle := d.EqualsQualString("handle")

	// check if identityHandle or workspaceHandle or aggregator id is empty
	if identityHandle == "" || workspaceHandle == "" || aggregatorHandle == "" {
		return nil, nil
	}

	getUserIdentityCached := plugin.HydrateFunc(getUserIdentity).WithCache()
	commonData, err := getUserIdentityCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("getWorkspaceAggregator", "getUserIdentityCached", err)
		return nil, err
	}

	user := commonData.(openapi.User)
	var response interface{}
	if identityHandle == user.Handle {
		response, err = getUserWorkspaceAggregator(ctx, d, h, identityHandle, workspaceHandle, aggregatorHandle)
	} else {
		response, err = getOrgWorkspaceAggregator(ctx, d, h, identityHandle, workspaceHandle, aggregatorHandle)
	}

	if err != nil {
		plugin.Logger(ctx).Error("getWorkspaceAggregator", "error", err)
		return nil, err
	}

	return response.(openapi.WorkspaceAggregator), nil
}

func getUserWorkspaceAggregator(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, userHandle, workspaceHandle, aggregatorHandle string) (interface{}, error) {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("steampipecloud_workspace_aggregator.getUserWorkspaceAggregator", "connection_error", err)
		return nil, err
	}

	var aggregator openapi.WorkspaceAggregator

	getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		aggregator, _, err = svc.UserWorkspaceAggregators.Get(ctx, userHandle, workspaceHandle, aggregatorHandle).Execute()
		return aggregator, err
	}

	response, err := plugin.RetryHydrate(ctx, d, h, getDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})
	if err != nil {
		plugin.Logger(ctx).Error("steampipecloud_workspace_aggregator.getUserWorkspaceAggregator", "retry_error", err)
		return nil, err
	}

	aggregator = response.(openapi.WorkspaceAggregator)

	return aggregator, nil
}

func getOrgWorkspaceAggregator(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, orgHandle, workspaceHandle, aggregatorHandle string) (interface{}, error) {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getOrgWorkspaceAggregator", "connection_error", err)
		return nil, err
	}

	var aggregator openapi.WorkspaceAggregator

	getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		aggregator, _, err = svc.OrgWorkspaceAggregators.Get(ctx, orgHandle, workspaceHandle, aggregatorHandle).Execute()
		return aggregator, err
	}

	response, err := plugin.RetryHydrate(ctx, d, h, getDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})
	if err != nil {
		plugin.Logger(ctx).Error("getOrgWorkspaceAggregator", "get", err)
		return nil, err
	}

	aggregator = response.(openapi.WorkspaceAggregator)

	return aggregator, nil
}

func getIdentityWorkspaceDetailsForAggregator(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("steampipecloud_workspace_aggregator.getIdentityWorkspaceDetailsForAggregator", "connection_error", err)
		return nil, err
	}

	var identityWorkspaceDetails IdentityWorkspaceDetailsForAggregator
	// get workspace details from hydrate data
	// workspace details reside in the parent item in this case
	switch w := h.ParentItem.(type) {
	case openapi.Workspace:
		plugin.Logger(ctx).Debug("getIdentityWorkspaceDetailsForAggregator", "openapi.Workspace")
		identityId := h.ParentItem.(openapi.Workspace).IdentityId
		identityWorkspaceDetails.WorkspaceHandle = h.ParentItem.(openapi.Workspace).Handle
		getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
			if strings.HasPrefix(identityId, "u_") {
				resp, _, err := svc.Users.Get(ctx, identityId).Execute()
				identityWorkspaceDetails.IdentityType = "user"
				identityWorkspaceDetails.IdentityHandle = resp.Handle
				return nil, err
			} else {
				resp, _, err := svc.Orgs.Get(ctx, identityId).Execute()
				identityWorkspaceDetails.IdentityType = "org"
				identityWorkspaceDetails.IdentityHandle = resp.Handle
				return nil, err
			}
		}
		_, _ = plugin.RetryHydrate(ctx, d, h, getDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})
		return identityWorkspaceDetails, nil
	case *openapi.Workspace:
		plugin.Logger(ctx).Debug("getIdentityWorkspaceDetailsForAggregator", "*openapi.Workspace")
		identityId := h.ParentItem.(*openapi.Workspace).IdentityId
		identityWorkspaceDetails.WorkspaceHandle = h.ParentItem.(*openapi.Workspace).Handle
		getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
			if strings.HasPrefix(identityId, "u_") {
				resp, _, err := svc.Users.Get(ctx, identityId).Execute()
				identityWorkspaceDetails.IdentityType = "user"
				identityWorkspaceDetails.IdentityHandle = resp.Handle
				return nil, err
			} else {
				resp, _, err := svc.Orgs.Get(ctx, identityId).Execute()
				identityWorkspaceDetails.IdentityType = "org"
				identityWorkspaceDetails.IdentityHandle = resp.Handle
				return nil, err
			}
		}
		_, _ = plugin.RetryHydrate(ctx, d, h, getDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})
		plugin.Logger(ctx).Debug("getIdentityWorkspaceDetailsForAggregator", "identityWorkspaceDetails", identityWorkspaceDetails)
		return &identityWorkspaceDetails, nil
	default:
		// The default case will come up when the parent items does not exist which happens when a get call is executed instead of a list
		// In a get call the parent hydrate is not executed and hence the parent item will not exist
		plugin.Logger(ctx).Debug("getIdentityWorkspaceDetailsForAggregator", "Unknown Type", w)
		identityId := h.Item.(openapi.WorkspaceAggregator).IdentityId
		workspaceId := h.Item.(openapi.WorkspaceAggregator).WorkspaceId
		getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
			if strings.HasPrefix(identityId, "u_") {
				user, _, err := svc.Users.Get(ctx, identityId).Execute()
				if err != nil {
					return nil, err
				}
				identityWorkspaceDetails.IdentityType = "user"
				identityWorkspaceDetails.IdentityHandle = user.Handle
				workspace, _, err := svc.UserWorkspaces.Get(ctx, identityId, workspaceId).Execute()
				if err != nil {
					return nil, err
				}
				identityWorkspaceDetails.WorkspaceHandle = workspace.Handle
			} else {
				org, _, err := svc.Orgs.Get(ctx, identityId).Execute()
				if err != nil {
					return nil, err
				}
				identityWorkspaceDetails.IdentityType = "org"
				identityWorkspaceDetails.IdentityHandle = org.Handle
				workspace, _, err := svc.UserWorkspaces.Get(ctx, identityId, workspaceId).Execute()
				if err != nil {
					return nil, err
				}
				identityWorkspaceDetails.WorkspaceHandle = workspace.Handle
			}
			return nil, nil
		}
		_, _ = plugin.RetryHydrate(ctx, d, h, getDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})
		plugin.Logger(ctx).Debug("getIdentityWorkspaceDetailsForAggregator", "identityWorkspaceDetails", identityWorkspaceDetails)
		return &identityWorkspaceDetails, err
	}
}
