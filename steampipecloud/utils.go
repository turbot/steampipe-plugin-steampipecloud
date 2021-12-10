package steampipecloud

import (
	"context"
	"strings"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
	openapi "github.com/turbot/steampipecloud-sdk-go"
)

func getUserIdentity(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	cacheKey := "GetUserIdentity"

	// if found in cache, return the result
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(openapi.TypesUser), nil
	}

	// get the service connection for the service
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("GetUserIdentity", "connection_error", err)
		return nil, err
	}

	var resp openapi.TypesUser

	getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		resp, _, err = svc.UsersApi.GetActor(ctx).Execute()
		return resp, err
	}

	response, err := plugin.RetryHydrate(ctx, d, h, getDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

	user := response.(openapi.TypesUser)

	if err != nil {
		plugin.Logger(ctx).Error("GetUserIdentity", "error", err)
		return nil, err
	}

	// save to extension cache
	d.ConnectionManager.Cache.Set(cacheKey, user)

	return user, nil
}

func setIdentityType(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("setIdentityType", "Value", d.Value)
	if d.Value == nil {
		return nil, nil
	}
	id := types.SafeString(d.Value)
	if strings.Contains(id, "o_") {
		return "org", nil
	} else if strings.Contains(id, "u_") {
		return "user", nil
	} else {
		return nil, nil
	}
}
