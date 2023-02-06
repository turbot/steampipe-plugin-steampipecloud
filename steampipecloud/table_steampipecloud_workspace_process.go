package steampipecloud

import (
	"context"
	"fmt"
	"strings"
	"time"

	openapi "github.com/turbot/steampipe-cloud-sdk-go"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

type IdentityWorkspaceDetailsForProcess struct {
	IdentityHandle  string `json:"identity_handle"`
	IdentityType    string `json:"identity_type"`
	WorkspaceHandle string `json:"workspace_handle"`
}

//// TABLE DEFINITION

func tableSteampipeCloudWorkspaceProcess(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "steampipecloud_workspace_process",
		Description: "Allows to track various processes for a workspace of an identity in Steampipe Cloud.",
		List: &plugin.ListConfig{
			ParentHydrate: listWorkspaces,
			Hydrate:       listWorkspaceProcesses,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:      "created_at",
					Require:   plugin.Optional,
					Operators: []string{">", ">=", "=", "<", "<="},
				},
				{
					Name:      "id",
					Require:   plugin.Optional,
					Operators: []string{"=", "<>"},
				},
				{
					Name:    "identity_handle",
					Require: plugin.Optional,
				},
				{
					Name:    "identity_id",
					Require: plugin.Optional,
				},
				{
					Name:      "pipeline_id",
					Require:   plugin.Optional,
					Operators: []string{"=", "<>"},
				},
				{
					Name:       "query_where",
					Require:    plugin.Optional,
					CacheMatch: "exact",
				},
				{
					Name:      "state",
					Require:   plugin.Optional,
					Operators: []string{"=", "<>"},
				},
				{
					Name:      "type",
					Require:   plugin.Optional,
					Operators: []string{"=", "<>"},
				},
				{
					Name:      "updated_at",
					Require:   plugin.Optional,
					Operators: []string{">", ">=", "=", "<", "<="},
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
			KeyColumns: plugin.AllColumns([]string{"identity_handle", "workspace_handle", "id"}),
			Hydrate:    getWorkspaceProcess,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The unique identifier for the process.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "identity_id",
				Description: "The unique identifier of the identity to which the process belongs to.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "identity_handle",
				Description: "The handle of the identity.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIdentityWorkspaceDetailsForWorkspaceProcess,
			},
			{
				Name:        "identity_type",
				Description: "The type of identity, can be org/user.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIdentityWorkspaceDetailsForWorkspaceProcess,
			},
			{
				Name:        "workspace_id",
				Description: "The unique identifier of the workspace to which the process belongs to.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "workspace_handle",
				Description: "The handle of the workspace.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIdentityWorkspaceDetailsForWorkspaceProcess,
			},
			{
				Name:        "pipeline_id",
				Description: "The unique identifier for the pipeline if a process is for a pipeline run/execution.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "type",
				Description: "The type of action executed by the process.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "The current state of the process.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "query_where",
				Description: "The query where expression to filter workspace processes.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("query_where"),
			},
			{
				Name:        "created_at",
				Description: "The time when the process was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "created_by_id",
				Description: "The unique identifier of the user who created the process.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "created_by",
				Description: "Information about the user who created the process.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "updated_at",
				Description: "The time when the process was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "updated_by_id",
				Description: "The unique identifier of the user who last updated the process.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "updated_by",
				Description: "Information about the user who last updated the process.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "version_id",
				Description: "The current version ID for the process.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromCamel(),
			},
		},
	}
}

//// LIST FUNCTION

func listWorkspaceProcesses(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var workspace *openapi.Workspace
	switch w := h.Item.(type) {
	case openapi.Workspace:
		wo := h.Item.(openapi.Workspace)
		workspace = &wo
	case *openapi.Workspace:
		workspace = h.Item.(*openapi.Workspace)
	default:
		plugin.Logger(ctx).Error("listWorkspaceProcesses", "unknown response type for workspace list parent hydrate call", w)
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

	workspaceHandle := d.KeyColumnQuals["workspace_handle"].GetStringValue()
	workspaceId := d.KeyColumnQuals["workspace_id"].GetStringValue()
	var workspaceToPass string

	// Error out if both workspace_handle and workspace_id is passed
	if workspaceHandle != "" && workspaceId != "" {
		plugin.Logger(ctx).Error("listWorkspaceProcesses", "please pass any one of workspace_id or workspace_handle")
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
		err = listUserWorkspaceProcesses(ctx, d, h, workspace.IdentityId, workspaceToPass, maxResults)
	} else {
		err = listOrgWorkspaceProcesses(ctx, d, h, workspace.IdentityId, workspaceToPass, maxResults)
	}

	if err != nil {
		plugin.Logger(ctx).Error("listWorkspaceProcesses", "error", err)
		return nil, err
	}

	return nil, nil
}

func listUserWorkspaceProcesses(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, userHandle, workspaceHandle string, maxResults int32) error {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listUserWorkspaceProcesses", "connection_error", err)
		return err
	}

	var filter string
	// collect all clauses passed as quals except for "query_where", "identity_id", "identity_handle", "workspace_id", "workspace_handle"
	var clauses []string
	for _, keyQual := range d.Table.List.KeyColumns {
		filterQual := d.Quals[keyQual.Name]
		if filterQual == nil || keyQual.Name == "query_where" || keyQual.Name == "identity_id" || keyQual.Name == "identity_handle" || keyQual.Name == "workspace_id" || keyQual.Name == "workspace_handle" {
			continue
		}
		for _, qual := range filterQual.Quals {
			if qual.Value != nil {
				var value string
				if keyQual.Name == "created_at" || keyQual.Name == "updated_at" {
					t := time.Unix(qual.Value.GetTimestampValue().Seconds, int64(qual.Value.GetTimestampValue().Nanos)).UTC()
					value = t.Format("2006-01-02 15:04:05.00000")
				} else {
					value = qual.Value.GetStringValue()
				}
				switch qual.Operator {
				case "=":
					clauses = append(clauses, fmt.Sprintf(`%s = '%s'`, keyQual.Name, value))
				case "<>":
					clauses = append(clauses, fmt.Sprintf(`%s <> '%s'`, keyQual.Name, value))
				case ">":
					clauses = append(clauses, fmt.Sprintf(`%s > '%s'`, keyQual.Name, value))
				case ">=":
					clauses = append(clauses, fmt.Sprintf(`%s >= '%s'`, keyQual.Name, value))
				case "<":
					clauses = append(clauses, fmt.Sprintf(`%s < '%s'`, keyQual.Name, value))
				case "<=":
					clauses = append(clauses, fmt.Sprintf(`%s <= '%s'`, keyQual.Name, value))
				}
			}
		}
	}

	// Frame the filter string by joining the collected quals by "and"
	filter = strings.Join(clauses, " and ")

	// Check if a query_where qual has been passed and add it to the filter string if yes
	if d.KeyColumnQuals["query_where"] != nil {
		if len(filter) >= 1 {
			filter = filter + " and " + d.KeyColumnQuals["query_where"].GetStringValue()
		} else {
			filter = d.KeyColumnQuals["query_where"].GetStringValue()
		}
	}

	pagesLeft := true
	var resp openapi.ListProcessesResponse
	var listDetails func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error)

	for pagesLeft {
		if resp.NextToken != nil {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.UserWorkspaceProcesses.List(ctx, userHandle, workspaceHandle).Where(filter).NextToken(*resp.NextToken).Limit(maxResults).Execute()
				return resp, err
			}
		} else {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.UserWorkspaceProcesses.List(ctx, userHandle, workspaceHandle).Where(filter).Limit(maxResults).Execute()
				return resp, err
			}
		}

		response, err := plugin.RetryHydrate(ctx, d, h, listDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

		if err != nil {
			plugin.Logger(ctx).Error("listUserWorkspaceProcesses", "list", err)
			return err
		}

		result := response.(openapi.ListProcessesResponse)

		if result.HasItems() {
			for _, process := range *result.Items {
				d.StreamListItem(ctx, process)

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

func listOrgWorkspaceProcesses(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, orgHandle, workspaceHandle string, maxResults int32) error {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listOrgWorkspaceProcesses", "connection_error", err)
		return err
	}

	var filter string
	// collect all clauses passed as quals except for "query_where", "identity_id", "identity_handle", "workspace_id", "workspace_handle"
	var clauses []string
	for _, keyQual := range d.Table.List.KeyColumns {
		filterQual := d.Quals[keyQual.Name]
		if filterQual == nil || keyQual.Name == "query_where" || keyQual.Name == "identity_id" || keyQual.Name == "identity_handle" || keyQual.Name == "workspace_id" || keyQual.Name == "workspace_handle" {
			continue
		}
		for _, qual := range filterQual.Quals {
			if qual.Value != nil {
				var value string
				if keyQual.Name == "created_at" || keyQual.Name == "updated_at" {
					t := time.Unix(qual.Value.GetTimestampValue().Seconds, int64(qual.Value.GetTimestampValue().Nanos)).UTC()
					value = t.Format("2006-01-02 15:04:05.00000")
				} else {
					value = qual.Value.GetStringValue()
				}
				switch qual.Operator {
				case "=":
					clauses = append(clauses, fmt.Sprintf(`%s = '%s'`, keyQual.Name, value))
				case "<>":
					clauses = append(clauses, fmt.Sprintf(`%s <> '%s'`, keyQual.Name, value))
				case ">":
					clauses = append(clauses, fmt.Sprintf(`%s > '%s'`, keyQual.Name, value))
				case ">=":
					clauses = append(clauses, fmt.Sprintf(`%s >= '%s'`, keyQual.Name, value))
				case "<":
					clauses = append(clauses, fmt.Sprintf(`%s < '%s'`, keyQual.Name, value))
				case "<=":
					clauses = append(clauses, fmt.Sprintf(`%s <= '%s'`, keyQual.Name, value))
				}
			}
		}
	}

	// Frame the filter string by joining the collected quals by "and"
	filter = strings.Join(clauses, " and ")

	// Check if a query_where qual has been passed and add it to the filter string if yes
	if d.KeyColumnQuals["query_where"] != nil {
		if len(filter) >= 1 {
			filter = filter + " and " + d.KeyColumnQuals["query_where"].GetStringValue()
		} else {
			filter = d.KeyColumnQuals["query_where"].GetStringValue()
		}
	}

	pagesLeft := true
	var resp openapi.ListProcessesResponse
	var listDetails func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error)

	for pagesLeft {
		if resp.NextToken != nil {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.OrgWorkspaceProcesses.List(ctx, orgHandle, workspaceHandle).Where(filter).NextToken(*resp.NextToken).Limit(maxResults).Execute()
				return resp, err
			}
		} else {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.OrgWorkspaceProcesses.List(ctx, orgHandle, workspaceHandle).Where(filter).Limit(maxResults).Execute()
				return resp, err
			}
		}

		response, err := plugin.RetryHydrate(ctx, d, h, listDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

		if err != nil {
			plugin.Logger(ctx).Error("listOrgWorkspaceProcesses", "list", err)
			return err
		}

		result := response.(openapi.ListProcessesResponse)

		if result.HasItems() {
			for _, process := range *result.Items {
				d.StreamListItem(ctx, process)

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

func getWorkspaceProcess(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	identityHandle := d.KeyColumnQuals["identity_handle"].GetStringValue()
	workspaceHandle := d.KeyColumnQuals["workspace_handle"].GetStringValue()
	processId := d.KeyColumnQuals["id"].GetStringValue()

	// check if identityHandle or workspaceHandle or pipeline id is empty
	if identityHandle == "" || workspaceHandle == "" || processId == "" {
		return nil, nil
	}

	getUserIdentityCached := plugin.HydrateFunc(getUserIdentity).WithCache()
	commonData, err := getUserIdentityCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("getWorkspaceProcess", "getUserIdentityCached", err)
		return nil, err
	}

	user := commonData.(openapi.User)
	var response interface{}
	if identityHandle == user.Handle {
		response, err = getUserWorkspaceProcess(ctx, d, h, identityHandle, workspaceHandle, processId)
	} else {
		response, err = getOrgWorkspaceProcess(ctx, d, h, identityHandle, workspaceHandle, processId)
	}

	if err != nil {
		plugin.Logger(ctx).Error("getWorkspaceProcess", "error", err)
		return nil, err
	}

	return response.(openapi.SpProcess), nil
}

func getUserWorkspaceProcess(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, userHandle, workspaceHandle, processId string) (interface{}, error) {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getUserWorkspaceProcess", "connection_error", err)
		return nil, err
	}

	var process openapi.SpProcess

	getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		process, _, err = svc.UserWorkspaceProcesses.Get(ctx, userHandle, workspaceHandle, processId).Execute()
		return process, err
	}

	response, err := plugin.RetryHydrate(ctx, d, h, getDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})
	if err != nil {
		plugin.Logger(ctx).Error("getUserWorkspaceProcess", "get", err)
		return nil, err
	}

	return response, nil
}

func getOrgWorkspaceProcess(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, orgHandle, workspaceHandle, processId string) (interface{}, error) {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getOrgWorkspaceProcess", "connection_error", err)
		return nil, err
	}

	var process openapi.SpProcess

	getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		process, _, err = svc.OrgWorkspaceProcesses.Get(ctx, orgHandle, workspaceHandle, processId).Execute()
		return process, err
	}

	response, err := plugin.RetryHydrate(ctx, d, h, getDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})
	if err != nil {
		plugin.Logger(ctx).Error("getOrgWorkspaceProcess", "get", err)
		return nil, err
	}

	return response, nil
}

func getIdentityWorkspaceDetailsForWorkspaceProcess(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getIdentityWorkspaceDetailsForWorkspaceProcess", "connection_error", err)
		return nil, err
	}

	var identityWorkspaceDetails IdentityWorkspaceDetailsForProcess
	// get workspace details from hydrate data
	// workspace details reside in the parent item in this case
	switch w := h.ParentItem.(type) {
	case openapi.Workspace:
		plugin.Logger(ctx).Debug("getIdentityWorkspaceDetailsForWorkspaceProcess", "openapi.Workspace")
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
		plugin.Logger(ctx).Debug("getIdentityWorkspaceDetailsForWorkspaceProcess", "*openapi.Workspace")
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
		plugin.Logger(ctx).Debug("getIdentityWorkspaceDetailsForWorkspaceProcess", "identityWorkspaceDetails", identityWorkspaceDetails)
		return &identityWorkspaceDetails, nil
	default:
		plugin.Logger(ctx).Debug("getIdentityWorkspaceDetailsForWorkspaceProcess", "Unknown Type", w)
	}
	return &identityWorkspaceDetails, nil
}
