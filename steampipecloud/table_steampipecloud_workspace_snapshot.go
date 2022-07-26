package steampipecloud

import (
	"context"
	"encoding/json"
	"strings"

	openapi "github.com/turbot/steampipe-cloud-sdk-go"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

type IdentityWorkspaceDetails struct {
	IdentityHandle  string `json:"identity_handle"`
	IdentityType    string `json:"identity_type"`
	WorkspaceHandle string `json:"workspace_handle"`
}

type SnapshotData struct {
	Data string `json:"data"`
}

//// TABLE DEFINITION

func tableSteampipeCloudWorkspaceSnapshot(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "steampipecloud_workspace_snapshot",
		Description: "Organization workspace members can collaborate and share connections and dashboards.",
		List: &plugin.ListConfig{
			ParentHydrate: listWorkspaces,
			Hydrate:       listWorkspaceSnapshots,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"identity_handle", "workspace_handle", "id"}),
			Hydrate:    getWorkspaceSnapshot,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The unique identifier for the member.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "identity_id",
				Description: "The unique identifier of the indentity to which the snapshot belongs to.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "identity_handle",
				Description: "The handle of the identity.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIdentityWorkspaceDetails,
			},
			{
				Name:        "identity_type",
				Description: "The type of identity, can be org/user.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIdentityWorkspaceDetails,
			},
			{
				Name:        "workspace_id",
				Description: "The id of the workspace.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "workspace_handle",
				Description: "The handle of the workspace.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIdentityWorkspaceDetails,
			},
			{
				Name:        "state",
				Description: "The current state of the snapshot.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "visibility",
				Description: "The visibility of the snapshot.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "dashboard_name",
				Description: "The mod-prefixed name of the dashboard this snapshot belongs to.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "schema_version",
				Description: "The schema version of the underlying snapshot.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "inputs",
				Description: "The inputs used for this snapshot.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "variables",
				Description: "The variables used for this snapshot.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "search_path",
				Description: "The search path used for this snapshot.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "start_time",
				Description: "The time the dashboard execution started.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "end_time",
				Description: "The time the dashboard execution ended.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "data",
				Description: "The data for snapshot.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getSnapshotData,
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

func listWorkspaceSnapshots(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	workspace := h.Item.(*openapi.Workspace)

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
	if strings.HasPrefix(workspace.IdentityId, "u_") {
		err = listUserWorkspaceSnapshots(ctx, d, h, workspace.IdentityId, workspace.Handle, maxResults)
	} else {
		err = listOrgWorkspaceSnapshots(ctx, d, h, workspace.IdentityId, workspace.Handle, maxResults)
	}

	if err != nil {
		plugin.Logger(ctx).Error("listWorkspaceSnapshots", "error", err)
		return nil, err
	}

	return nil, nil
}

func listUserWorkspaceSnapshots(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, userHandle string, workspaceHandle string, maxResults int32) error {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listUserWorkspaceSnapshots", "connection_error", err)
		return err
	}

	pagesLeft := true
	var resp openapi.ListWorkspaceSnapshotsResponse
	var listDetails func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error)

	for pagesLeft {
		if resp.NextToken != nil {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.UserWorkspaceSnapshots.List(ctx, userHandle, workspaceHandle).NextToken(*resp.NextToken).Limit(maxResults).Execute()
				return resp, err
			}
		} else {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.UserWorkspaceSnapshots.List(ctx, userHandle, workspaceHandle).Limit(maxResults).Execute()
				return resp, err
			}
		}

		response, err := plugin.RetryHydrate(ctx, d, h, listDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

		if err != nil {
			plugin.Logger(ctx).Error("listUserWorkspaceSnapshots", "list", err)
			return err
		}

		result := response.(openapi.ListWorkspaceSnapshotsResponse)

		if result.HasItems() {
			for _, snapshot := range *result.Items {
				d.StreamListItem(ctx, snapshot)

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

func listOrgWorkspaceSnapshots(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, orgHandle string, workspaceHandle string, maxResults int32) error {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listOrgWorkspaceSnapshots", "connection_error", err)
		return err
	}

	pagesLeft := true
	var resp openapi.ListWorkspaceSnapshotsResponse
	var listDetails func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error)

	for pagesLeft {
		if resp.NextToken != nil {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.OrgWorkspaceSnapshots.List(ctx, orgHandle, workspaceHandle).NextToken(*resp.NextToken).Limit(maxResults).Execute()
				return resp, err
			}
		} else {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.OrgWorkspaceSnapshots.List(ctx, orgHandle, workspaceHandle).Limit(maxResults).Execute()
				return resp, err
			}
		}

		response, err := plugin.RetryHydrate(ctx, d, h, listDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

		if err != nil {
			plugin.Logger(ctx).Error("listOrgWorkspaceSnapshots", "list", err)
			return err
		}

		result := response.(openapi.ListWorkspaceSnapshotsResponse)

		if result.HasItems() {
			for _, snapshot := range *result.Items {
				d.StreamListItem(ctx, snapshot)

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

func getWorkspaceSnapshot(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	identityHandle := d.KeyColumnQuals["identity_handle"].GetStringValue()
	workspaceHandle := d.KeyColumnQuals["workspace_handle"].GetStringValue()
	snapshotId := d.KeyColumnQuals["id"].GetStringValue()

	// check if identityHandle or workspaceHandle or snapshot id is empty
	if identityHandle == "" || workspaceHandle == "" || snapshotId == "" {
		return nil, nil
	}

	getUserIdentityCached := plugin.HydrateFunc(getUserIdentity).WithCache()
	commonData, err := getUserIdentityCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("getWorkspaceSnapshot", "getUserIdentityCached", err)
		return nil, err
	}

	user := commonData.(openapi.User)
	var response interface{}
	if identityHandle == user.Handle {
		response, err = getUserWorkspaceSnapshot(ctx, d, h, identityHandle, workspaceHandle, snapshotId)
	} else {
		response, err = getOrgWorkspaceSnapshot(ctx, d, h, identityHandle, workspaceHandle, snapshotId)
	}

	if err != nil {
		plugin.Logger(ctx).Error("getWorkspaceSnapshot", "error", err)
		return nil, err
	}

	return response.(openapi.WorkspaceSnapshot), nil
}

func getUserWorkspaceSnapshot(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, userHandle, workspaceHandle, snapshotId string) (interface{}, error) {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getUserWorkspaceSnapshot", "connection_error", err)
		return nil, err
	}

	var snapshot openapi.WorkspaceSnapshot

	getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		snapshot, _, err = svc.UserWorkspaceSnapshots.Get(ctx, userHandle, workspaceHandle, snapshotId).Execute()
		return snapshot, err
	}

	response, err := plugin.RetryHydrate(ctx, d, h, getDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})
	if err != nil {
		plugin.Logger(ctx).Error("getUserWorkspaceSnapshot", "get", err)
		return nil, err
	}

	return response, nil
}

func getOrgWorkspaceSnapshot(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, orgHandle, workspaceHandle, snapshotId string) (interface{}, error) {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getOrgWorkspaceSnapshot", "connection_error", err)
		return nil, err
	}

	var snapshot openapi.WorkspaceSnapshot

	getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		snapshot, _, err = svc.OrgWorkspaceSnapshots.Get(ctx, orgHandle, workspaceHandle, snapshotId).Execute()
		return snapshot, err
	}

	response, err := plugin.RetryHydrate(ctx, d, h, getDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})
	if err != nil {
		plugin.Logger(ctx).Error("getOrgWorkspaceSnapshot", "get", err)
		return nil, err
	}

	return response, nil
}

func getIdentityWorkspaceDetails(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getIdentityWorkspaceDetails", "connection_error", err)
		return nil, err
	}

	var identityWorkspaceDetails IdentityWorkspaceDetails
	// get workspace details from hydrate data
	// workspace details reside in the parent item in this case
	switch w := h.ParentItem.(type) {
	case openapi.Workspace:
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
	default:
		plugin.Logger(ctx).Debug("getIdentityWorkspaceDetails", "Unknown Type", w)
	}

	identityId := h.Item.(openapi.WorkspaceSnapshot).IdentityId
	workspaceId := h.Item.(openapi.WorkspaceSnapshot).WorkspaceId
	getIdentityDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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
	_, _ = plugin.RetryHydrate(ctx, d, h, getIdentityDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

	getWorkspaceDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		if strings.HasPrefix(identityId, "u_") {
			resp, _, err := svc.UserWorkspaces.Get(ctx, identityId, workspaceId).Execute()
			identityWorkspaceDetails.WorkspaceHandle = resp.Handle
			return nil, err
		} else {
			resp, _, err := svc.OrgWorkspaces.Get(ctx, identityId, workspaceId).Execute()
			identityWorkspaceDetails.WorkspaceHandle = resp.Handle
			return nil, err
		}
	}
	_, _ = plugin.RetryHydrate(ctx, d, h, getWorkspaceDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

	return identityWorkspaceDetails, nil
}

func getSnapshotData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getSnapshotData", "connection_error", err)
		return nil, err
	}

	var snapshotData SnapshotData
	workspaceSnapshot := h.Item.(openapi.WorkspaceSnapshot)
	var response openapi.WorkspaceSnapshotData
	getSnapshotData := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		if strings.HasPrefix(workspaceSnapshot.IdentityId, "u_") {
			response, _, err = svc.UserWorkspaceSnapshots.Download(ctx, workspaceSnapshot.IdentityId, workspaceSnapshot.WorkspaceId, workspaceSnapshot.Id, "json").Execute()
			if err != nil {
				return nil, err
			}

		} else {
			response, _, err = svc.OrgWorkspaceSnapshots.Download(ctx, workspaceSnapshot.IdentityId, workspaceSnapshot.WorkspaceId, workspaceSnapshot.Id, "json").Execute()
			if err != nil {
				return nil, err
			}
		}
		byteArr, _ := json.Marshal(response)
		snapshotData.Data = string(byteArr)
		return nil, nil
	}
	_, err = plugin.RetryHydrate(ctx, d, h, getSnapshotData, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

	return snapshotData, nil
}
