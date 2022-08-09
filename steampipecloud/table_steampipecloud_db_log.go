package steampipecloud

import (
	"context"

	openapi "github.com/turbot/steampipe-cloud-sdk-go"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func tableSteampipeCloudWorkspaceDbLog(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "steampipecloud_workspace_db_log",
		Description: "Database logs records the underlying queries executed when a user executes a query.",
		List: &plugin.ListConfig{
			ParentHydrate: listWorkspaces,
			Hydrate:       listWorkspaceDbLogs,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The unique identifier for a db log.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "actor_id",
				Description: "The unique identifier for the user who executed the query.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "actor_handle",
				Description: "The handle of the user who executed the query.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "actor_display_name",
				Description: "The display name of the user who executed the query.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "actor_avatar_url",
				Description: "The avatar of the user who executed the query.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "workspace_id",
				Description: "The unique identifier of the workspace on which the query was executed.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "workspace_handle",
				Description: "The handle of the workspace on which the query was executed.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "duration",
				Description: "The duration of the query.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "query",
				Description: "The query that was executed in the workspace.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "log_timestamp",
				Description: "The time when the log got captured in postgres.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "created_at",
				Description: "The time when the db log record was generated.",
				Type:        proto.ColumnType_STRING,
			},
		},
	}
}

//// LIST FUNCTION

func listWorkspaceDbLogs(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get the workspace object from the parent hydrate
	workspace := h.Item.(*openapi.Workspace)

	// Create the connection
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listDbLogs", "connection_error", err)
		return nil, err
	}

	// Get cached user identity
	getUserIdentityCached := plugin.HydrateFunc(getUserIdentity).WithCache()
	commonData, err := getUserIdentityCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("listDbLogs", "getUserIdentityCached", err)
		return nil, err
	}

	// Extract the user object from the cached identity
	user := commonData.(openapi.User)

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

	// If we want to get the db logs for the user
	if user.Id == workspace.IdentityId {
		err = listUserWorkspaceDbLogs(ctx, d, h, svc, maxResults, user.Id, workspace.Id)
	} else {
		err = listOrgWorkspaceDbLogs(ctx, d, h, svc, maxResults, workspace.IdentityId, workspace.Id)
	}
	if err != nil {
		plugin.Logger(ctx).Error("listDbLogs", "error", err)
		return nil, err
	}

	if err != nil {
		plugin.Logger(ctx).Error("listDbLogs", "error", err)
		return nil, err
	}
	return nil, nil
}

func listUserWorkspaceDbLogs(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, svc *openapi.APIClient, maxResults int32, identityId, workspaceId string) error {
	var err error

	// execute list call
	pagesLeft := true
	var resp openapi.ListLogsResponse
	var listDetails func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error)

	for pagesLeft {
		if resp.NextToken != nil {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.UserWorkspaces.ListDBLogs(ctx, identityId, workspaceId).NextToken(*resp.NextToken).Limit(maxResults).Execute()
				return resp, err
			}
		} else {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.UserWorkspaces.ListDBLogs(ctx, identityId, workspaceId).Limit(maxResults).Execute()
				return resp, err
			}
		}

		response, err := plugin.RetryHydrate(ctx, d, h, listDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

		if err != nil {
			plugin.Logger(ctx).Error("listUserDbLogs", "list", err)
			return err
		}

		result := response.(openapi.ListLogsResponse)

		if result.HasItems() {
			for _, log := range *result.Items {
				d.StreamListItem(ctx, log)

				// Context can be cancelled due to manual cancellation or the limit has been hit
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return nil
				}
			}
		}
		if resp.NextToken == nil {
			pagesLeft = false
		} else {
			resp.NextToken = result.NextToken
		}
	}

	return nil
}

func listOrgWorkspaceDbLogs(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, svc *openapi.APIClient, maxResults int32, identityId, workspaceId string) error {
	var err error

	// execute list call
	pagesLeft := true
	var resp openapi.ListLogsResponse
	var listDetails func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error)

	for pagesLeft {
		if resp.NextToken != nil {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.OrgWorkspaces.ListDBLogs(ctx, identityId, workspaceId).NextToken(*resp.NextToken).Limit(maxResults).Execute()
				return resp, err
			}
		} else {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.OrgWorkspaces.ListDBLogs(ctx, identityId, workspaceId).Limit(maxResults).Execute()
				return resp, err
			}
		}

		response, err := plugin.RetryHydrate(ctx, d, h, listDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

		if err != nil {
			plugin.Logger(ctx).Error("listOrgDbLogs", "list", err)
			return err
		}

		result := response.(openapi.ListLogsResponse)

		if result.HasItems() {
			for _, log := range *result.Items {
				d.StreamListItem(ctx, log)

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
