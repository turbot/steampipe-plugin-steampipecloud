package steampipecloud

import (
	"context"
	"fmt"
	"strings"

	openapi "github.com/turbot/steampipe-cloud-sdk-go"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

type IdentityDetailsForProcess struct {
	IdentityHandle string `json:"identity_handle"`
	IdentityType   string `json:"identity_type"`
}

//// TABLE DEFINITION

func tableSteampipeCloudProcess(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "steampipecloud_process",
		Description: "Allows to track various processes for an identity in Steampipe Cloud.",
		List: &plugin.ListConfig{
			Hydrate: listIdentityProcesses,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "identity_handle",
					Require: plugin.Optional,
				},
				{
					Name:    "identity_id",
					Require: plugin.Optional,
				},
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"identity_handle", "id"}),
			Hydrate:    getIdentityProcess,
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
				Hydrate:     getIdentityDetailsForProcess,
			},
			{
				Name:        "identity_type",
				Description: "The type of identity, can be org/user.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIdentityDetailsForProcess,
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

func listIdentityProcesses(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var user openapi.User
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

	identityHandle := d.KeyColumnQuals["identity_handle"].GetStringValue()
	identityId := d.KeyColumnQuals["identity_id"].GetStringValue()
	var identityToPass string

	getUserIdentityCached := plugin.HydrateFunc(getUserIdentity).WithCache()
	commonData, err := getUserIdentityCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("steampipecloud_process.listIdentityProcesses", "getUserIdentityCached", err)
		return nil, err
	}

	user = commonData.(openapi.User)

	// Error out if both identity_handle and identity_id is passed
	if identityHandle != "" && identityId != "" {
		return nil, fmt.Errorf("please pass any one of identity_handle or identity_id")
	}
	if identityHandle != "" {
		identityToPass = identityHandle
	} else if identityId != "" {
		identityToPass = identityId
	} else {
		identityToPass = user.Id
	}

	if strings.HasPrefix(identityToPass, "u_") || identityToPass == user.Handle {
		err = listUserProcesses(ctx, d, h, identityToPass, maxResults)
	} else if strings.HasPrefix(identityToPass, "o_") || identityToPass != user.Handle {
		err = listOrgProcesses(ctx, d, h, identityToPass, maxResults)
	}

	if err != nil {
		plugin.Logger(ctx).Error("steampipecloud_process.listIdentityProcesses", "query_error", err)
		return nil, err
	}

	return nil, nil
}

func listUserProcesses(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, userHandle string, maxResults int32) error {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("steampipecloud_process.listUserProcesses", "connection_error", err)
		return err
	}

	pagesLeft := true
	var resp openapi.ListProcessesResponse
	var listDetails func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error)

	for pagesLeft {
		if resp.NextToken != nil {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.UserProcesses.List(ctx, userHandle).NextToken(*resp.NextToken).Limit(maxResults).Execute()
				return resp, err
			}
		} else {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.UserProcesses.List(ctx, userHandle).Limit(maxResults).Execute()
				return resp, err
			}
		}

		response, err := plugin.RetryHydrate(ctx, d, h, listDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

		if err != nil {
			plugin.Logger(ctx).Error("steampipecloud_process.listUserProcesses", "query_error", err)
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

func listOrgProcesses(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, orgHandle string, maxResults int32) error {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("steampipecloud_process.listOrgProcesses", "connection_error", err)
		return err
	}

	pagesLeft := true
	var resp openapi.ListProcessesResponse
	var listDetails func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error)

	for pagesLeft {
		if resp.NextToken != nil {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.OrgProcesses.List(ctx, orgHandle).NextToken(*resp.NextToken).Limit(maxResults).Execute()
				return resp, err
			}
		} else {
			listDetails = func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
				resp, _, err = svc.OrgProcesses.List(ctx, orgHandle).Limit(maxResults).Execute()
				return resp, err
			}
		}

		response, err := plugin.RetryHydrate(ctx, d, h, listDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

		if err != nil {
			plugin.Logger(ctx).Error("steampipecloud_process.listOrgProcesses", "query_error", err)
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

func getIdentityProcess(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	identityHandle := d.KeyColumnQuals["identity_handle"].GetStringValue()
	processId := d.KeyColumnQuals["id"].GetStringValue()

	// check if identityHandle or workspaceHandle or process id is empty
	if identityHandle == "" || processId == "" {
		return nil, nil
	}

	getUserIdentityCached := plugin.HydrateFunc(getUserIdentity).WithCache()
	commonData, err := getUserIdentityCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("steampipecloud_process.getIdentityProcess", "connection_error", err)
		return nil, err
	}

	user := commonData.(openapi.User)
	var response interface{}
	if identityHandle == user.Handle {
		response, err = getUserProcess(ctx, d, h, identityHandle, processId)
	} else {
		response, err = getOrgProcess(ctx, d, h, identityHandle, processId)
	}

	if err != nil {
		plugin.Logger(ctx).Error("steampipecloud_process.getIdentityProcess", "query_error", err)
		return nil, err
	}

	return response.(openapi.SpProcess), nil
}

func getUserProcess(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, userHandle, processId string) (interface{}, error) {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("steampipecloud_process.getUserProcess", "connection_error", err)
		return nil, err
	}

	var process openapi.SpProcess

	getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		process, _, err = svc.UserProcesses.Get(ctx, userHandle, processId).Execute()
		return process, err
	}

	response, err := plugin.RetryHydrate(ctx, d, h, getDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})
	if err != nil {
		plugin.Logger(ctx).Error("steampipecloud_process.getUserProcess", "query_error", err)
		return nil, err
	}

	return response, nil
}

func getOrgProcess(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, orgHandle, processId string) (interface{}, error) {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("steampipecloud_process.getOrgProcess", "connection_error", err)
		return nil, err
	}

	var process openapi.SpProcess

	getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		process, _, err = svc.OrgProcesses.Get(ctx, orgHandle, processId).Execute()
		return process, err
	}

	response, err := plugin.RetryHydrate(ctx, d, h, getDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})
	if err != nil {
		plugin.Logger(ctx).Error("steampipecloud_process.getOrgProcess", "query_error", err)
		return nil, err
	}

	return response, nil
}

func getIdentityDetailsForProcess(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("steampipecloud_process.getIdentityDetailsForProcess", "connection_error", err)
		return nil, err
	}

	getUserIdentityCached := plugin.HydrateFunc(getUserIdentity).WithCache()
	commonData, err := getUserIdentityCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("steampipecloud_process.getIdentityDetailsForProcess", "getUserIdentityCached", err)
		return nil, err
	}

	user := commonData.(openapi.User)

	var identityDetails IdentityDetailsForProcess
	process := h.Item.(openapi.SpProcess)
	plugin.Logger(ctx).Info("getIdentityDetailsForProcess", "process Item", process)
	if *process.IdentityId == user.Id {
		identityDetails.IdentityHandle = user.Handle
		identityDetails.IdentityType = "user"
	} else {
		org, _, err := svc.Orgs.Get(ctx, *process.IdentityId).Execute()
		if err != nil {
			plugin.Logger(ctx).Error("steampipecloud_process.getIdentityDetailsForProcess", "query_error", err)
			return nil, err
		}
		identityDetails.IdentityHandle = org.Handle
		identityDetails.IdentityType = "org"
	}
	return identityDetails, nil
}
