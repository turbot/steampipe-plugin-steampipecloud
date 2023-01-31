package steampipecloud

import (
	"context"
	"strings"

	openapi "github.com/turbot/steampipe-cloud-sdk-go"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

type IdentityWorkspaceDetailsForPipeline struct {
	IdentityId      string `json:"identity_id"`
	IdentityHandle  string `json:"identity_handle"`
	IdentityType    string `json:"identity_type"`
	WorkspaceHandle string `json:"workspace_handle"`
}

//// TABLE DEFINITION

func tableSteampipeCloudWorkspacePipeline(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "steampipecloud_workspace_pipeline",
		Description: "Pipelines allow users to run different kinds of activities in steampipe cloud on a schedule.",
		List: &plugin.ListConfig{
			ParentHydrate: listWorkspaces,
			Hydrate:       listWorkspacePipelines,
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
			KeyColumns: plugin.AllColumns([]string{"identity_handle", "workspace_handle", "id"}),
			Hydrate:    getWorkspacePipeline,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The unique identifier for the pipeline.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "identity_id",
				Description: "The unique identifier of the identity to which the pipeline belongs to.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIdentityWorkspaceDetailsForPipeline,
			},
			{
				Name:        "identity_handle",
				Description: "The handle of the identity.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIdentityWorkspaceDetailsForPipeline,
			},
			{
				Name:        "identity_type",
				Description: "The type of identity, can be org/user.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIdentityWorkspaceDetailsForPipeline,
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
				Hydrate:     getIdentityWorkspaceDetailsForPipeline,
			},
			{
				Name:        "title",
				Description: "The title of the pipeline.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "frequency",
				Description: "The frequency at which the pipeline will be executed.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "pipeline",
				Description: "The name of the pipeline to be executed.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "args",
				Description: "Arguments to be passed to the pipeline.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags",
				Description: "The tags for the pipeline.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "last_process_id",
				Description: "The unique identifier of the last process that was executed for the pipeline.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "last_process",
				Description: "Information about the process that was last executed for the pipeline.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "created_at",
				Description: "The time when the pipeline was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "created_by_id",
				Description: "The unique identifier of the user who created the pipeline.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "created_by",
				Description: "Information about the user who created the pipeline.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "updated_at",
				Description: "The time when the pipeline was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "updated_by_id",
				Description: "The unique identifier of the user who last updated the pipeline.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "updated_by",
				Description: "Information about the user who last updated the pipeline.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "version_id",
				Description: "The current version ID for the pipeline.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromCamel(),
			},
		},
	}
}

//// LIST FUNCTION

func listWorkspacePipelines(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var workspace *openapi.Workspace
	switch w := h.Item.(type) {
	case openapi.Workspace:
		wo := h.Item.(openapi.Workspace)
		workspace = &wo
	case *openapi.Workspace:
		workspace = h.Item.(*openapi.Workspace)
	default:
		plugin.Logger(ctx).Error("listWorkspacePipelines", "unknown response type for workspace list parent hydrate call", w)
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
		plugin.Logger(ctx).Error("listWorkspacePipelines", "please pass any one of workspace_id or workspace_handle")
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
		err = listUserWorkspacePipelines(ctx, d, h, workspace.IdentityId, workspaceToPass, maxResults)
	} else {
		err = listOrgWorkspacePipelines(ctx, d, h, workspace.IdentityId, workspaceToPass, maxResults)
	}

	if err != nil {
		plugin.Logger(ctx).Error("listWorkspacePipelines", "error", err)
		return nil, err
	}

	return nil, nil
}

func listUserWorkspacePipelines(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, userHandle string, workspaceHandle string, maxResults int32) error {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listUserWorkspacePipelines", "connection_error", err)
		return err
	}

	pagesLeft := true
	var resp openapi.ListPipelinesResponse
	var listDetails func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error)

	for pagesLeft {
		if resp.NextToken != nil {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.UserWorkspacePipelines.List(ctx, userHandle, workspaceHandle).NextToken(*resp.NextToken).Limit(maxResults).Execute()
				return resp, err
			}
		} else {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.UserWorkspacePipelines.List(ctx, userHandle, workspaceHandle).Limit(maxResults).Execute()
				return resp, err
			}
		}

		response, err := plugin.RetryHydrate(ctx, d, h, listDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

		if err != nil {
			plugin.Logger(ctx).Error("listUserWorkspacePipelines", "list", err)
			return err
		}

		result := response.(openapi.ListPipelinesResponse)

		if result.HasItems() {
			for _, pipeline := range *result.Items {
				d.StreamListItem(ctx, pipeline)

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

func listOrgWorkspacePipelines(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, orgHandle string, workspaceHandle string, maxResults int32) error {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listOrgWorkspacePipelines", "connection_error", err)
		return err
	}

	pagesLeft := true
	var resp openapi.ListPipelinesResponse
	var listDetails func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error)

	for pagesLeft {
		if resp.NextToken != nil {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.OrgWorkspacePipelines.List(ctx, orgHandle, workspaceHandle).NextToken(*resp.NextToken).Limit(maxResults).Execute()
				return resp, err
			}
		} else {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.OrgWorkspacePipelines.List(ctx, orgHandle, workspaceHandle).Limit(maxResults).Execute()
				return resp, err
			}
		}

		response, err := plugin.RetryHydrate(ctx, d, h, listDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

		if err != nil {
			plugin.Logger(ctx).Error("listOrgWorkspacePipelines", "list", err)
			return err
		}

		result := response.(openapi.ListPipelinesResponse)

		if result.HasItems() {
			for _, pipeline := range *result.Items {
				d.StreamListItem(ctx, pipeline)

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

func getWorkspacePipeline(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	identityHandle := d.KeyColumnQuals["identity_handle"].GetStringValue()
	workspaceHandle := d.KeyColumnQuals["workspace_handle"].GetStringValue()
	pipelineId := d.KeyColumnQuals["id"].GetStringValue()

	// check if identityHandle or workspaceHandle or pipeline id is empty
	if identityHandle == "" || workspaceHandle == "" || pipelineId == "" {
		return nil, nil
	}

	getUserIdentityCached := plugin.HydrateFunc(getUserIdentity).WithCache()
	commonData, err := getUserIdentityCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("getWorkspacePipeline", "getUserIdentityCached", err)
		return nil, err
	}

	user := commonData.(openapi.User)
	var response interface{}
	if identityHandle == user.Handle {
		response, err = getUserWorkspacePipeline(ctx, d, h, identityHandle, workspaceHandle, pipelineId)
	} else {
		response, err = getOrgWorkspacePipeline(ctx, d, h, identityHandle, workspaceHandle, pipelineId)
	}

	if err != nil {
		plugin.Logger(ctx).Error("getWorkspacePipeline", "error", err)
		return nil, err
	}

	return response.(openapi.Pipeline), nil
}

func getUserWorkspacePipeline(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, userHandle, workspaceHandle, pipelineId string) (interface{}, error) {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getUserWorkspacePipeline", "connection_error", err)
		return nil, err
	}

	var pipeline openapi.Pipeline

	getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		pipeline, _, err = svc.UserWorkspacePipelines.Get(ctx, userHandle, workspaceHandle, pipelineId).Execute()
		return pipeline, err
	}

	response, err := plugin.RetryHydrate(ctx, d, h, getDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})
	if err != nil {
		plugin.Logger(ctx).Error("getUserWorkspacePipeline", "get", err)
		return nil, err
	}

	return response, nil
}

func getOrgWorkspacePipeline(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, orgHandle, workspaceHandle, pipelineId string) (interface{}, error) {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getOrgWorkspacePipeline", "connection_error", err)
		return nil, err
	}

	var pipeline openapi.Pipeline

	getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		pipeline, _, err = svc.OrgWorkspacePipelines.Get(ctx, orgHandle, workspaceHandle, pipelineId).Execute()
		return pipeline, err
	}

	response, err := plugin.RetryHydrate(ctx, d, h, getDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})
	if err != nil {
		plugin.Logger(ctx).Error("getOrgWorkspacePipeline", "get", err)
		return nil, err
	}

	return response, nil
}

func getIdentityWorkspaceDetailsForPipeline(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getIdentityWorkspaceDetails", "connection_error", err)
		return nil, err
	}

	var identityWorkspaceDetails IdentityWorkspaceDetailsForPipeline
	// get workspace details from hydrate data
	// workspace details reside in the parent item in this case
	switch w := h.ParentItem.(type) {
	case openapi.Workspace:
		plugin.Logger(ctx).Debug("getIdentityWorkspaceDetails", "openapi.Workspace")
		identityId := h.ParentItem.(openapi.Workspace).IdentityId
		identityWorkspaceDetails.WorkspaceHandle = h.ParentItem.(openapi.Workspace).Handle
		getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
			if strings.HasPrefix(identityId, "u_") {
				resp, _, err := svc.Users.Get(ctx, identityId).Execute()
				identityWorkspaceDetails.IdentityId = resp.Id
				identityWorkspaceDetails.IdentityType = "user"
				identityWorkspaceDetails.IdentityHandle = resp.Handle
				return nil, err
			} else {
				resp, _, err := svc.Orgs.Get(ctx, identityId).Execute()
				identityWorkspaceDetails.IdentityId = resp.Id
				identityWorkspaceDetails.IdentityType = "org"
				identityWorkspaceDetails.IdentityHandle = resp.Handle
				return nil, err
			}
		}
		_, _ = plugin.RetryHydrate(ctx, d, h, getDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})
		return identityWorkspaceDetails, nil
	case *openapi.Workspace:
		plugin.Logger(ctx).Debug("getIdentityWorkspaceDetails", "*openapi.Workspace")
		identityId := h.ParentItem.(*openapi.Workspace).IdentityId
		identityWorkspaceDetails.WorkspaceHandle = h.ParentItem.(*openapi.Workspace).Handle
		getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
			if strings.HasPrefix(identityId, "u_") {
				resp, _, err := svc.Users.Get(ctx, identityId).Execute()
				identityWorkspaceDetails.IdentityId = resp.Id
				identityWorkspaceDetails.IdentityType = "user"
				identityWorkspaceDetails.IdentityHandle = resp.Handle
				return nil, err
			} else {
				resp, _, err := svc.Orgs.Get(ctx, identityId).Execute()
				identityWorkspaceDetails.IdentityId = resp.Id
				identityWorkspaceDetails.IdentityType = "org"
				identityWorkspaceDetails.IdentityHandle = resp.Handle
				return nil, err
			}
		}
		_, _ = plugin.RetryHydrate(ctx, d, h, getDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})
		plugin.Logger(ctx).Debug("getIdentityWorkspaceDetails", "identityWorkspaceDetails", identityWorkspaceDetails)
		return identityWorkspaceDetails, nil
	default:
		plugin.Logger(ctx).Debug("getIdentityWorkspaceDetails", "Unknown Type", w)
	}
	return identityWorkspaceDetails, nil
}
